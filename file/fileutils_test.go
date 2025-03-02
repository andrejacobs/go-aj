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
	"os"
	"path/filepath"
	"testing"

	"github.com/andrejacobs/go-aj/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExists(t *testing.T) {
	d, err := makeValidDir()
	defer os.RemoveAll(d)
	require.NoError(t, err)
	require.DirExists(t, d)
	exists, err := file.PathExists(d)
	require.NoError(t, err)
	require.True(t, exists)
	exists, err = file.DirExists(d)
	require.NoError(t, err)
	require.True(t, exists)
	exists, err = file.FileExists(d)
	require.NoError(t, err)
	require.False(t, exists)

	f, err := makeValidFile()
	defer os.Remove(f)
	require.NoError(t, err)
	require.FileExists(t, f)
	exists, err = file.PathExists(f)
	require.NoError(t, err)
	require.True(t, exists)
	exists, err = file.DirExists(f)
	require.NoError(t, err)
	require.False(t, exists)
	exists, err = file.FileExists(f)
	require.NoError(t, err)
	require.True(t, exists)

	d, err = makeInvalidDir()
	require.NoError(t, err)
	require.NoDirExists(t, d)
	exists, err = file.PathExists(d)
	require.NoError(t, err)
	require.False(t, exists)
	exists, err = file.DirExists(d)
	require.NoError(t, err)
	require.False(t, exists)
	exists, err = file.FileExists(d)
	require.NoError(t, err)
	require.False(t, exists)

	f, err = makeInvalidFile()
	require.NoError(t, err)
	require.NoFileExists(t, f)
	exists, err = file.PathExists(f)
	require.NoError(t, err)
	require.False(t, exists)
	exists, err = file.DirExists(f)
	require.NoError(t, err)
	require.False(t, exists)
	exists, err = file.FileExists(f)
	require.NoError(t, err)
	require.False(t, exists)
}

func TestAbsPaths(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(cwd)

	tempDir, err := makeValidDir()
	defer os.RemoveAll(tempDir)
	require.NoError(t, err)
	require.NoError(t, os.Chdir(tempDir))

	// Without validation
	paths, err := file.AbsPaths([]string{"dirOne"}, false)
	require.NoError(t, err)
	expected, err := filepath.Abs("dirOne")
	require.NoError(t, err)
	assert.Contains(t, paths, expected)

	// With validation
	_, err = file.AbsPaths([]string{"dirOne"}, true)
	require.Error(t, err)
}

func TestReplaceExt(t *testing.T) {
	assert.Equal(t, "/a/b/c.md", file.ReplaceExt("/a/b/c.txt", ".md"))
	assert.Equal(t, "/a/b/c.md", file.ReplaceExt("/a/b/c", ".md"))
	assert.Equal(t, "/a/b/cmd", file.ReplaceExt("/a/b/c.txt", "md"))
}

func TestRemoveIfExists(t *testing.T) {
	f, err := os.CreateTemp("", "delme")
	require.NoError(t, err)
	require.NoError(t, os.Remove(f.Name()))

	assert.NoError(t, file.RemoveIfExists(f.Name()))
}

//-----------------------------------------------------------------------------

func makeValidDir() (string, error) {
	return os.MkdirTemp("", "unit-tests")
}

func makeValidFile() (string, error) {
	f, err := os.CreateTemp("", "unit-tests")
	if err != nil {
		return "", err
	}
	defer f.Close()
	return f.Name(), nil
}

func makeInvalidDir() (string, error) {
	p, err := makeValidDir()
	if err != nil {
		return "", err
	}
	if err := os.RemoveAll(p); err != nil {
		return "", err
	}
	return p, nil
}

func makeInvalidFile() (string, error) {
	p, err := makeValidFile()
	if err != nil {
		return "", err
	}
	if err := os.Remove(p); err != nil {
		return "", err
	}
	return p, nil
}
