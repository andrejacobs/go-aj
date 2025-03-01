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
	"testing"

	"github.com/andrejacobs/go-aj/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsDirEntryEqual(t *testing.T) {
	a, err := os.ReadDir(tempDir)
	require.NoError(t, err)

	b, err := os.ReadDir(tempDir)
	require.NoError(t, err)

	assert.True(t, file.IsDirEntryEqual(a[0], b[0]))
}

func TestIsDirEntryWithInfoEqual(t *testing.T) {
	a, err := os.ReadDir(tempDir)
	require.NoError(t, err)

	b, err := os.ReadDir(tempDir)
	require.NoError(t, err)

	equal, err := file.IsDirEntryWithInfoEqual(a[0], b[0])
	require.NoError(t, err)
	assert.True(t, equal)
}

func TestReadDirUnsorted(t *testing.T) {
	a, err := os.ReadDir(tempDir)
	require.NoError(t, err)

	b, err := file.ReadDirUnsorted(tempDir)
	require.NoError(t, err)
	file.SortDirEntries(b)

	require.Equal(t, len(a), len(b))
	for i := 0; i < len(a); i++ {
		assert.True(t, file.IsDirEntryEqual(a[i], b[i]))
	}
}
