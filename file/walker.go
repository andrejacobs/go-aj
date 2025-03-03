package file

import (
	"io/fs"
	"path/filepath"
)

// Walker is used to walk a file hierarchy.
type Walker struct {
	DirExcluder  MatchPathFn // Determine which directories should not be walked
	FileExcluder MatchPathFn // Determine which files should not be walked
}

// Create a new Walker.
//
// By default all files and directories found will be walked and not be excluded.
func NewWalker() *Walker {
	return &Walker{
		DirExcluder:  NeverMatch,
		FileExcluder: NeverMatch,
	}
}

// Set the excluder used to determine which directories should not be walked.
func (w *Walker) SetDirExcluder(excluder MatchPathFn) *Walker {
	w.DirExcluder = excluder
	return w
}

// Set the excluder used to determine which files should not be walked.
func (w *Walker) SetFileExcluder(excluder MatchPathFn) *Walker {
	w.FileExcluder = excluder
	return w
}

// Walk walks the file tree rooted at root, calling fn for each file or
// directory in the tree, including root.
//
// All errors that arise visiting files and directories are filtered by fn:
// see the [fs.WalkDirFunc] documentation for details.
//
// The files are walked in lexical order, which makes the output deterministic
// but requires Walk to read an entire directory into memory before proceeding
// to walk that directory.
//
// Walk does not follow symbolic links.
//
// Walk calls fn with paths that use the separator character appropriate
// for the operating system.
//
// Walk uses [fs.WalkDir] for implementation.
//
// For each directory that is found, the DirExcluder will be called to determine
// if the path should not be walked.
//
// For each file that is found, the FileExcluder will be called to determine
// if the path should not be walked.
func (w *Walker) Walk(root string, fn fs.WalkDirFunc) error {

	rErr := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		//AJ### do filtering
		fnErr := fn(path, d, err)
		return fnErr
	})

	return rErr
}

//-----------------------------------------------------------------------------
// Excluders

// MatchPathFn determines if the path matches a criteria and if so returns true.
type MatchPathFn func(path string) (bool, error)

// MatchPathMiddleware specifies the function signature for a wrapping MatchPathFn into a call chain.
// Similiarly to how http middleware works in popular frameworks.
type MatchPathMiddleware func(next MatchPathFn) MatchPathFn

// NeverMatch is a MatchPathFn that will always return false.
func NeverMatch(path string) (bool, error) {
	return false, nil
}

// func DefaultDirExcluder(next MatchPathFn) MatchPathFn {
// 	fmt.Println("DefaultDirExcluder init")
// 	return func(path string) (bool, error) {
// 		fmt.Println("DefaultDirExcluder called")
// 		return next(path)
// 	}
// }
