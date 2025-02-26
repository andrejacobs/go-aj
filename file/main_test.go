package file_test

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
)

var tempDir string

func TestMain(m *testing.M) {
	var err error
	tempDir, err = makeFileTree()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	defer os.RemoveAll(tempDir)

	os.Exit(m.Run())
}

//-----------------------------------------------------------------------------

// Create a temporary directory with a couple of files for testing
func makeFileTree() (string, error) {
	tempDir, err := os.MkdirTemp("", "unit-testing")
	if err != nil {
		return "", err
	}

	if err := makeFile(filepath.Join(tempDir, "a"), 10); err != nil {
		return "", err
	}

	if err := makeFile(filepath.Join(tempDir, "b"), 20); err != nil {
		return "", err
	}

	if err := makeFile(filepath.Join(tempDir, "c"), 30); err != nil {
		return "", err
	}

	subDir := filepath.Join(tempDir, "d")
	if err := os.Mkdir(subDir, 0744); err != nil {
		return "", err
	}

	if err := makeFile(filepath.Join(subDir, "e"), 10); err != nil {
		return "", err
	}

	if err := makeFile(filepath.Join(subDir, "f"), 20); err != nil {
		return "", err
	}

	return tempDir, nil
}

// Create a file with the specified size
func makeFile(path string, size int64) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
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
