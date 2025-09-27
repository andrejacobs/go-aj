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
	"io"
	"math"
	"testing"

	"github.com/andrejacobs/go-aj/ajio/trackedoffset"
	"github.com/andrejacobs/go-aj/ajmath"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriter(t *testing.T) {
	baseOffset := uint64(42)
	tw := trackedoffset.NewWriter(io.Discard, baseOffset)
	assert.Equal(t, baseOffset, tw.Offset())

	data := []byte("The quick brown fox jumped over the lazy dog!")
	c, err := tw.Write(data)
	require.NoError(t, err)
	assert.Equal(t, len(data), c)
	assert.Equal(t, baseOffset+uint64(len(data)), tw.Offset())
}

func TestWriterResetOffset(t *testing.T) {
	baseOffset := uint64(42)
	tw := trackedoffset.NewWriter(io.Discard, baseOffset)
	assert.Equal(t, baseOffset, tw.Offset())

	data := []byte("Moo said the cow")
	c, err := tw.Write(data)
	assert.NoError(t, err)
	assert.Equal(t, len(data), c)
	assert.Equal(t, baseOffset+uint64(len(data)), tw.Offset())

	data = []byte("The quick brown fox jumped over the lazy dog!")
	baseOffset = 200
	tw.ResetOffset(baseOffset)
	c, err = tw.Write(data)
	require.NoError(t, err)
	assert.Equal(t, len(data), c)
	assert.Equal(t, baseOffset+uint64(len(data)), tw.Offset())
}

func TestWriterOverflow(t *testing.T) {
	baseOffset := uint64(math.MaxUint64 - 4)
	tw := trackedoffset.NewWriter(io.Discard, baseOffset)
	assert.Equal(t, baseOffset, tw.Offset())

	data := []byte("The quick brown fox jumped over the lazy dog!")
	_, err := tw.Write(data)
	assert.ErrorIs(t, err, ajmath.ErrIntegerOverflow)
}
