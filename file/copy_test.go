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
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/andrejacobs/go-aj/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopyFile(t *testing.T) {
	expected := "The quick brown fox jumped over the lazy dog!"
	src, err := os.CreateTemp("", "unit-test-source")
	require.NoError(t, err)
	defer os.Remove(src.Name())
	_, err = src.WriteString(expected)
	require.NoError(t, err)
	require.NoError(t, src.Close())

	destPath := filepath.Join(t.TempDir(), "unit-test-dest")
	defer os.Remove(destPath)
	wc, err := file.CopyFile(context.Background(), src.Name(), destPath)
	require.NoError(t, err)
	assert.Equal(t, int64(len(expected)), wc)

	dest, err := os.Open(destPath)
	require.NoError(t, err)
	defer dest.Close()

	data, err := io.ReadAll(dest)
	require.NoError(t, err)

	assert.Equal(t, expected, string(data))
}

func TestCopyFileN(t *testing.T) {
	expected := "The quick brown fox jumped over the lazy dog!"
	src, err := os.CreateTemp("", "unit-test-source")
	require.NoError(t, err)
	defer os.Remove(src.Name())
	_, err = src.WriteString(expected)
	require.NoError(t, err)
	require.NoError(t, src.Close())

	destPath := filepath.Join(t.TempDir(), "unit-test-dest")
	defer os.Remove(destPath)
	wc, err := file.CopyFileN(context.Background(), src.Name(), destPath, 9)
	require.NoError(t, err)
	assert.Equal(t, int64(9), wc)

	dest, err := os.Open(destPath)
	require.NoError(t, err)
	defer dest.Close()

	data, err := io.ReadAll(dest)
	require.NoError(t, err)

	assert.Equal(t, "The quick", string(data))
}
