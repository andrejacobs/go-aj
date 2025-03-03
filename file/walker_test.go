package file_test

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"slices"
	"testing"

	"github.com/andrejacobs/go-aj/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWalkerDefaults(t *testing.T) {
	expected, err := expectedFilepathWalk(tempDir)
	require.NoError(t, err)

	result := make([]string, 0, 10)
	var fn fs.WalkDirFunc
	fn = func(path string, d fs.DirEntry, err error) error {
		fmt.Printf("%q\n", path)
		result = append(result, path)
		return nil
	}

	w := file.NewWalker()
	err = w.Walk(tempDir, fn)
	require.NoError(t, err)

	slices.Sort(result)
	assert.ElementsMatch(t, expected, result)
}

func expectedFilepathWalk(path string) ([]string, error) {
	expected := make([]string, 0, 10)
	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		expected = append(expected, path)
		return nil
	})
	if err != nil {
		return expected, err
	}

	slices.Sort(expected)
	return expected, nil
}
