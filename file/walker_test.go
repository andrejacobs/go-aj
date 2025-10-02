// Copyright (c) 2025 Andre Jacobs
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package file_test

import (
	"fmt"
	"io/fs"
	"os"
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
	var fn fs.WalkDirFunc = func(path string, d fs.DirEntry, err error) error {
		// fmt.Printf("%q\n", path)
		result = append(result, path)
		return nil
	}

	w := file.NewWalker()
	err = w.Walk(tempDir, fn)
	require.NoError(t, err)

	slices.Sort(result)
	assert.ElementsMatch(t, expected, result)
}

func TestWalkerIncludeDirs(t *testing.T) {
	expected := make([]string, 0, 10)
	err := filepath.WalkDir(tempDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() && (path != tempDir) {
			if d.Name() != "g" {
				return fs.SkipDir
			}
		}
		expected = append(expected, path)
		return nil
	})
	require.NoError(t, err)
	slices.Sort(expected)

	result := make([]string, 0, 10)
	var fn fs.WalkDirFunc = func(path string, d fs.DirEntry, err error) error {
		// fmt.Printf("%q\n", path)
		result = append(result, path)
		return nil
	}

	w := file.NewWalker()
	w.DirIncluder = (func(path string, d fs.DirEntry) (bool, error) {
		if d.IsDir() && d.Name() == "g" {
			return true, nil
		}
		return false, nil
	})
	err = w.Walk(tempDir, fn)
	require.NoError(t, err)

	slices.Sort(result)
	assert.ElementsMatch(t, expected, result)
}

func TestWalkerExcludeDirs(t *testing.T) {
	expected := make([]string, 0, 10)
	err := filepath.WalkDir(tempDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() && d.Name() == "d" {
			return fs.SkipDir
		}
		expected = append(expected, path)
		return nil
	})
	require.NoError(t, err)
	slices.Sort(expected)

	result := make([]string, 0, 10)
	var fn fs.WalkDirFunc = func(path string, d fs.DirEntry, err error) error {
		// fmt.Printf("%q\n", path)
		result = append(result, path)
		return nil
	}

	w := file.NewWalker()
	w.DirExcluder = (func(path string, d fs.DirEntry) (bool, error) {
		if d.IsDir() && d.Name() == "d" {
			return true, nil
		}
		return false, nil
	})
	err = w.Walk(tempDir, fn)
	require.NoError(t, err)

	slices.Sort(result)
	assert.ElementsMatch(t, expected, result)
}

func TestWalkerIncludeFiles(t *testing.T) {
	expected := make([]string, 0, 10)
	err := filepath.WalkDir(tempDir, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && (d.Name() != "a" && (d.Name() != "e")) {
			return nil
		}
		expected = append(expected, path)
		return nil
	})
	require.NoError(t, err)
	slices.Sort(expected)

	result := make([]string, 0, 10)
	var fn fs.WalkDirFunc = func(path string, d fs.DirEntry, err error) error {
		// fmt.Printf("%q\n", path)
		result = append(result, path)
		return nil
	}

	w := file.NewWalker()
	w.FileIncluder = func(path string, d fs.DirEntry) (bool, error) {
		if !d.IsDir() && (d.Name() == "a" || (d.Name() == "e")) {
			return true, nil
		}
		return false, nil
	}
	err = w.Walk(tempDir, fn)
	require.NoError(t, err)

	slices.Sort(result)
	assert.ElementsMatch(t, expected, result)
}

func TestWalkerExcludeFiles(t *testing.T) {
	expected := make([]string, 0, 10)
	err := filepath.WalkDir(tempDir, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && (d.Name() == "b" || (d.Name() == "e")) {
			return nil
		}
		expected = append(expected, path)
		return nil
	})
	require.NoError(t, err)
	slices.Sort(expected)

	result := make([]string, 0, 10)
	var fn fs.WalkDirFunc = func(path string, d fs.DirEntry, err error) error {
		// fmt.Printf("%q\n", path)
		result = append(result, path)
		return nil
	}

	w := file.NewWalker()
	w.FileExcluder = func(path string, d fs.DirEntry) (bool, error) {
		if !d.IsDir() && (d.Name() == "b" || (d.Name() == "e")) {
			return true, nil
		}
		return false, nil
	}
	err = w.Walk(tempDir, fn)
	require.NoError(t, err)

	slices.Sort(result)
	assert.ElementsMatch(t, expected, result)
}

func TestWalkerExcludeFilesAndMiddleware(t *testing.T) {
	expected := make([]string, 0, 10)
	err := filepath.WalkDir(tempDir, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && (d.Name() == "b" || (d.Name() == "e") || (d.Name() == ".DS_Store")) {
			return nil
		}
		expected = append(expected, path)
		return nil
	})
	require.NoError(t, err)
	slices.Sort(expected)

	result := make([]string, 0, 10)
	var fn fs.WalkDirFunc = func(path string, d fs.DirEntry, err error) error {
		// fmt.Printf("%q\n", path)
		result = append(result, path)
		return nil
	}

	w := file.NewWalker()
	w.FileExcluder = file.MatchAppleDSStore(
		func(path string, d fs.DirEntry) (bool, error) {
			if !d.IsDir() && (d.Name() == "b" || (d.Name() == "e")) {
				return true, nil
			}
			return false, nil
		})

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

func TestWalkerPassesReceivedError(t *testing.T) {
	var fn fs.WalkDirFunc = func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		return nil
	}

	w := file.NewWalker()
	err := w.Walk("/does-not-exist", fn)
	var expErr *fs.PathError
	require.ErrorAs(t, err, &expErr)
}

func TestWalkerExpandsUsersHomeDir(t *testing.T) {
	var fn fs.WalkDirFunc = func(path string, d fs.DirEntry, rcvErr error) error {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		require.Equal(t, fmt.Sprintf("%s/does-not-exist", home), path)

		if rcvErr != nil {
			return rcvErr
		}
		return nil
	}

	w := file.NewWalker()
	err := w.Walk("~/does-not-exist", fn)
	var expErr *fs.PathError
	require.ErrorAs(t, err, &expErr)
}
