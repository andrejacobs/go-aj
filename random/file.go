package random

import (
	"crypto/rand"
	"io"
	"os"
)

// Create a file and fill it with random bytes.
// NOTE: This will override any existing file.
// path The path of the file to be created.
// size The number of random bytes to write to the file.
func CreateFile(path string, size int64) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.CopyN(f, rand.Reader, size)
	if err != nil {
		return err
	}

	return nil
}

// Create a temporary file and fill it with random bytes.
// NOTE: This will override any existing file.
// See os.CreateTemp for details on dir and pattern.
// size The number of random bytes to write to the file.
// Returns the path to the file that was created.
func CreateTempFile(dir, pattern string, size int64) (string, error) {
	f, err := os.CreateTemp(dir, pattern)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = io.CopyN(f, rand.Reader, size)
	if err != nil {
		return "", err
	}

	return f.Name(), nil
}
