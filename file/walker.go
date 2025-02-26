package file

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sync"
	"syscall"
)

type WalkDirFunc func(path string, entry fs.DirEntry, fileInfo fs.FileInfo) error

// Provides a Concurrent file tree walker
type ConcurrentWalker struct {
	root     string      // The root path to start walking from
	walkFunc WalkDirFunc // The function that will be called concurrently with new file or directory info.

	errorWriter io.Writer // Writer used to report errors. Default is os.Stderr

	workerSemaCh chan struct{}  // Counting semaphore used to restrict the maximum number of workers
	pathInfoCh   chan pathInfo  // Channel on which found paths are published
	errorCh      chan error     // Channel on which any encountered errors are published on
	doneCh       chan struct{}  // Channel on which a caller will be notified when the work is done
	wg           sync.WaitGroup // WaitGroup used to wait until all work is finished

	hadErrors bool // Will be true if any errors were encountered during a walk.
	errorMu   sync.RWMutex

	inProgress bool // True while not done

	dirExcluder  PathExcluder
	fileExcluder PathExcluder
}

// Create a new concurrent walker
func NewConcurrentWalker() *ConcurrentWalker {
	return &ConcurrentWalker{
		errorWriter:  os.Stderr,
		dirExcluder:  &DefaultDirExcluder{},
		fileExcluder: &DefaultFileExcluder{},
	}
}

// Set the io.Writer to use for reporting any errors.
func (w *ConcurrentWalker) SetErrorWriter(writer io.Writer) *ConcurrentWalker {
	w.errorWriter = writer
	return w
}

// Return true if any errors were encountered during a walk.
func (w *ConcurrentWalker) HadErrors() bool {
	w.errorMu.RLock()
	defer w.errorMu.RUnlock()
	return w.hadErrors
}

// Set the excluder used to exclude directories from being walked
func (w *ConcurrentWalker) SetDirExcluder(excluder PathExcluder) *ConcurrentWalker {
	w.dirExcluder = excluder
	return w
}

// Set the excluder used to exclude files from being walked
func (w *ConcurrentWalker) SetFileExcluder(excluder PathExcluder) *ConcurrentWalker {
	w.fileExcluder = excluder
	return w
}

// Start walking the file tree recursively and concurrently.
// This function does not block the caller.
// Returns the channel on which to block until work is finished or an error that was encountered before
// walking recursively.
// Also returns the cancel function to invoke to stop the operation.
// ctx: Parent context.
// root: The path to start walking from.
// walkFunc The function to be called when new info is received.
func (w *ConcurrentWalker) StartWalking(ctx context.Context, root string, walkFunc WalkDirFunc) (<-chan struct{}, context.CancelFunc, error) {
	if w.inProgress {
		return nil, nil, fmt.Errorf("the previous walk has not yet been finished")
	}

	ctx, cancel := context.WithCancel(ctx)

	w.root = root
	w.walkFunc = walkFunc
	w.inProgress = true
	w.hadErrors = false

	// Get info about the root path
	info, err := os.Lstat(w.root)
	if err != nil {
		return nil, cancel, err
	}
	rootEntry := startDirEntry{info}

	// Determine the max number of open files allowed
	var rlimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlimit); err != nil {
		return nil, cancel, err
	}

	// Create a counting semaphore to restrict the number of active workers
	maxActiveWorkers := min(uint64(1024), rlimit.Cur)
	w.workerSemaCh = make(chan struct{}, maxActiveWorkers)

	w.pathInfoCh = make(chan pathInfo, maxActiveWorkers)
	w.errorCh = make(chan error, 100)
	w.doneCh = make(chan struct{})
	w.wg = sync.WaitGroup{}

	// Go: Start walking recursively
	w.wg.Add(1)
	go w.walk(w.root, &rootEntry)

	// Go: wait for all work to be done and then close the pathInfo channel to signal all work is done
	go func() {
		w.wg.Wait()
		close(w.pathInfoCh)
	}()

	// Go: Process
	go w.process(ctx)

	return w.doneCh, cancel, nil
}

// Block the caller and wait till the process is finished.
func (w *ConcurrentWalker) Wait() {
	<-w.doneCh
}

// The main processing loop
func (w *ConcurrentWalker) process(ctx context.Context) {
loop:
	for {
		select {
		case <-ctx.Done():
			break loop

		// Process new found path
		case info, ok := <-w.pathInfoCh:
			if !ok {
				break loop
			}
			// Notify the walk function about the new path info found
			w.walkFunc(info.path, info.entry, info.fileInfo)

		// Process error
		case err := <-w.errorCh:
			w.errorMu.Lock()
			w.hadErrors = true
			fmt.Fprintf(w.errorWriter, "%v\n", err)
			w.errorMu.Unlock()
		}
	}

	w.signalDone()
}

// Signal to the caller (StartWalking) that we are done
func (w *ConcurrentWalker) signalDone() {
	w.inProgress = false
	close(w.doneCh)
}

// Recursive walker
func (w *ConcurrentWalker) walk(dir string, entry fs.DirEntry) {
	defer w.wg.Done()

	fileInfo, err := entry.Info()
	if err != nil {
		w.errorCh <- fmt.Errorf("failure getting info for a path %q: %v", dir, err)
		return
	}

	// Check if directory needs to be excluded
	if entry.IsDir() {
		exclude, err := w.dirExcluder.Match(dir)
		if err != nil {
			w.errorCh <- fmt.Errorf("failure checking if directory needs to be excluded %q: %v", dir, err)
			return
		}
		if exclude {
			return
		}
	} else {
		// Check if the file should be excluded
		exclude, err := w.fileExcluder.Match(dir)
		if err != nil {
			w.errorCh <- fmt.Errorf("failure checking if file needs to be excluded %q: %v", dir, err)
			return
		}
		if exclude {
			return
		}
	}

	w.pathInfoCh <- pathInfo{
		path:     dir,
		entry:    entry,
		fileInfo: fileInfo,
	}

	if !entry.IsDir() {
		return
	}

	// Walk children
	children, err := workerReadDir(w.workerSemaCh, dir)
	if err != nil {
		w.errorCh <- err
		return
	}

	for _, child := range children {
		childDir := path.Join(dir, child.Name())
		w.wg.Add(1)
		go w.walk(childDir, child)
	}
}

//-----------------------------------------------------------------------------

type pathInfo struct {
	path     string
	entry    fs.DirEntry
	fileInfo fs.FileInfo
}

// Used as a bridge from the Lstat result (fs.FileInfo) and the expected fs.DirEntry
// when starting to walk recursively
type startDirEntry struct {
	info fs.FileInfo
}

func (d *startDirEntry) Name() string               { return d.info.Name() }
func (d *startDirEntry) IsDir() bool                { return d.info.IsDir() }
func (d *startDirEntry) Type() fs.FileMode          { return d.info.Mode().Type() }
func (d *startDirEntry) Info() (fs.FileInfo, error) { return d.info, nil }

//-----------------------------------------------------------------------------

// workerReadDir is a version of os.ReadDir that works with a counting semaphore
// to ensure we can restrict the number of workers opening files
// Note this also doesn't sort the dir entries to save on a couple of allocs
func workerReadDir(workerSemaCh chan struct{}, dirname string) ([]fs.DirEntry, error) {
	// Acquire a worker slot
	workerSemaCh <- struct{}{}
	defer func() { <-workerSemaCh }()

	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	dirs, err := f.ReadDir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	// Ignoring sorting (this saves a couple of allocs, not so much of a speed improvement)
	// sort.Slice(dirs, func(i, j int) bool { return dirs[i].Name() < dirs[j].Name() })
	return dirs, nil
}

//-----------------------------------------------------------------------------
// PathMatcher (same as the one from the matches package)

// PathMatcher is used to determine if a file system path matches
type PathMatcher interface {
	// Match checks if the path matches and returns true if it does
	Match(path string) (bool, error)
}

//-----------------------------------------------------------------------------
// Excluders

// PathExcluder is used to determine if a path should be excluded from being walked
type PathExcluder interface {
	PathMatcher
}

// The default directory excluder
type DefaultDirExcluder struct {
	// See the walker_[platform].go specific files for matching rules
}

// The default file excluder
type DefaultFileExcluder struct {
}

func (e *DefaultFileExcluder) Match(path string) (bool, error) {
	filename := filepath.Base(path)
	// Exclude macOS stupid files on every platform
	if filename == ".DS_Store" {
		return true, nil
	}
	// Don't exclude
	return false, nil
}
