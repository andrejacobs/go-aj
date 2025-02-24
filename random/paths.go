package random

// Provide utility functions for creating random file paths. Mainly used in unit-testing.

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

// Generate a path consisting of random depth (subdirectories) between min and max
// minDirs, maxDirs: random range between the minimum  and maximum amount of subdirectories to create
// minNameLen, maxNameLen: random range of length of characters used to generate each random subdirectory's name.
// The function will always return the base + range(min, max) paths.
func Path(base string, minDirs int, maxDirs int, minNameLen int, maxNameLen int) string {
	sb := strings.Builder{}
	count := Int(minDirs, maxDirs)
	minNameLen = max(1, minNameLen)
	for depth := 0; depth < count; depth++ {
		sb.WriteString(String(Int(minNameLen, maxNameLen)))
		if depth < (count - 1) {
			sb.WriteRune(os.PathSeparator)
		}
	}
	return path.Join(base, sb.String())
}

// Generate a slice of random paths
// count: is the number of random paths to create and return
func Paths(base string, count int, min int, max int, minNameLen int, maxNameLen int) []string {
	paths := make([]string, 0, count)
	for i := 0; i < count; i++ {
		paths = append(paths, Path(base, min, max, minNameLen, maxNameLen))
	}
	return paths
}

// Generate random files inside the specified directory
// Files will be created using data copied from the crypto random generator.
// dir: is the parent directory
// minFile: minimum number of files to create
// maxFile: maximum number of files to create
// minSize: the minimum size in bytes of a file.
// maxSize: the maximum size in bytes of a file.
// maxTotalSize: the maximum number of bytes to be used for all files being created.
// Return the total number of bytes written.
func CreateFiles(dir string,
	minFiles int, maxFiles int,
	minSize uint64, maxSize uint64,
	maxTotalSize uint64) (uint64, error) {

	currentTotalSize := uint64(0)

	for i := 0; i < Int(minFiles, maxFiles); i++ {
		path := path.Join(dir, fmt.Sprintf("%s-%d", String(Int(1, 16)), i))
		if currentTotalSize < maxTotalSize {
			amount := min(int64(Int(0, int(maxSize))), int64(maxTotalSize-currentTotalSize))
			wc, err := CreateFileWithSize(path, uint64(amount))
			if err != nil {
				return currentTotalSize, err
			}
			currentTotalSize += uint64(wc)
			if currentTotalSize >= maxTotalSize {
				break
			}
		}
	}

	return currentTotalSize, nil
}

// Create a file with the exact size in bytes, by copying bytes from the cryptographically secure random number generator.
func CreateFileWithSize(path string, size uint64) (uint64, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	wc, err := io.CopyN(f, rand.Reader, int64(size))
	return uint64(wc), err
}
