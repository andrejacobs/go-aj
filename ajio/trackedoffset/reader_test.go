// Copyright (c) 2025 Andre Jacobs
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package trackedoffset_test

import (
	"bufio"
	"math"
	"strings"
	"testing"

	"github.com/andrejacobs/go-aj/ajio/trackedoffset"
	"github.com/andrejacobs/go-aj/ajmath"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReader(t *testing.T) {
	text := "The quick brown fox jumped over the lazy dog!"
	sr := strings.NewReader(text)
	br := bufio.NewReader(sr)

	baseOffset := uint64(42)
	tr := trackedoffset.NewReader(br, baseOffset)
	assert.Equal(t, baseOffset, tr.Offset())

	buffer := make([]byte, 4)
	for i := 0; i < len(text)/4; i++ {
		_, err := tr.Read(buffer)
		require.NoError(t, err)
		assert.Equal(t, baseOffset+uint64((i+1)*4), tr.Offset())
	}
}

func TestReaderResetOffset(t *testing.T) {
	text := "The quick brown fox jumped over the lazy dog!"
	sr := strings.NewReader(text)
	br := bufio.NewReader(sr)

	baseOffset := uint64(42)
	tr := trackedoffset.NewReader(br, baseOffset)
	assert.Equal(t, baseOffset, tr.Offset())

	buffer := make([]byte, 4)
	_, err := tr.Read(buffer)
	assert.NoError(t, err)
	assert.Equal(t, baseOffset+4, tr.Offset())

	tr.ResetOffset(200)
	_, err = tr.Read(buffer)
	assert.NoError(t, err)
	assert.Equal(t, uint64(204), tr.Offset())
}

func TestReaderOverflow(t *testing.T) {
	text := "The quick brown fox jumped over the lazy dog!"
	sr := strings.NewReader(text)
	br := bufio.NewReader(sr)

	baseOffset := uint64(math.MaxUint64 - 2)
	tr := trackedoffset.NewReader(br, uint64(baseOffset))
	assert.Equal(t, baseOffset, tr.Offset())

	buffer := make([]byte, 4)
	_, err := tr.Read(buffer)
	assert.ErrorIs(t, err, ajmath.ErrIntegerOverflow)
}
