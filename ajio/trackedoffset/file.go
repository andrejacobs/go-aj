package trackedoffset

import (
	"bufio"
	"io"
	"os"

	"github.com/andrejacobs/go-aj/ajmath"
)

// File wraps an os.File and keeps track of the current offset without requiring constant calls to Seek which involves syscall Lseek to be made.
// Reading and Writing is buffered by using the bufio package.
// Implements the following interfaces: io.Reader, io.Writer, io.Seeker.
type File struct {
	of     *os.File
	reader *bufio.Reader
	writer *bufio.Writer
	offset uint64
}

// Create a new File.
func NewFile(of *os.File) (*File, error) {
	f := &File{
		of:     of,
		reader: bufio.NewReader(of),
		writer: bufio.NewWriter(of),
	}

	if err := f.SyncOffset(); err != nil {
		return nil, err
	}

	return f, nil
}

// Close the file and release resources.
func (f *File) Close() error {
	err := f.of.Close()
	f.of = nil
	f.reader = nil
	f.writer = nil
	return err
}

// Path of the file.
func (f *File) Name() string {
	return f.of.Name()
}

// Return information describing the file.
func (f *File) Stat() (os.FileInfo, error) {
	return f.of.Stat()
}

// Reader implementation.
func (f *File) Read(p []byte) (int, error) {
	n, err := f.reader.Read(p)
	if err != nil {
		return n, err
	}

	newOffset, err := ajmath.Add64(f.offset, uint64(n))
	if err != nil {
		return 0, err
	}
	f.offset = newOffset

	return n, nil
}

// Writer implementation.
func (f *File) Write(p []byte) (int, error) {
	n, err := f.writer.Write(p)
	if err != nil {
		return n, err
	}

	newOffset, err := ajmath.Add64(f.offset, uint64(n))
	if err != nil {
		return 0, err
	}
	f.offset = newOffset

	return n, nil
}

// Seeker implementation.
// It is recommened that you ResetReadBuffer or ResetWriteBuffer.
func (f *File) Seek(offset int64, whence int) (int64, error) {
	newOffset, err := f.of.Seek(offset, whence)
	if err != nil {
		return newOffset, err
	}

	f.offset, err = ajmath.Int64ToUint64(newOffset)
	if err != nil {
		return newOffset, err
	}

	return newOffset, nil
}

// Return the current offset in bytes.
func (f *File) Offset() uint64 {
	return f.offset
}

// Ensure the File's offset and the underlying os.File's actual offsets are the same.
// This will make a call to file.Seek.
func (f *File) SyncOffset() error {
	offset, err := f.of.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}

	f.offset, err = ajmath.Int64ToUint64(offset)
	if err != nil {
		return err
	}

	return nil
}

// Ensure any unwritten data in the write buffer is written to the actual file.
func (f *File) Flush() error {
	return f.writer.Flush()
}

// Sync commits the current contents of the file to stable storage.
func (f *File) Sync() error {
	return f.of.Sync()
}

// Discard any buffered data that has not yet been read.
// Ensure you call this if you have changed the current offset using Seek.
func (f *File) ResetReadBuffer() {
	f.reader.Reset(f.of)
}

// Discard any buffered data that has not yet been written.
// Ensure you call this if you have changed the current offset using Seek.
func (f *File) ResetWriteBuffer() {
	f.writer.Reset(f.of)
}

//TODO: Flush writer buffer and maybe do a file sync
