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

//go:build amd64 || arm64
// +build amd64 arm64

package safe_test

import (
	"math"
	"testing"

	"github.com/andrejacobs/go-aj/ajmath/safe"
	"github.com/stretchr/testify/assert"
)

func TestIntToInt32_on_64bit(t *testing.T) {
	v, err := safe.IntToInt32(math.MinInt32 - 1)
	assert.ErrorIs(t, err, safe.ErrIntegerUnderflow)
	assert.Equal(t, int32(0), v)

	v, err = safe.IntToInt32(math.MaxInt32 + 1)
	assert.ErrorIs(t, err, safe.ErrIntegerOverflow)
	assert.Equal(t, int32(0), v)
}

func TestUintToUint32_on_64bit(t *testing.T) {
	v, err := safe.UintToUint32(math.MaxUint32 + 1)
	assert.ErrorIs(t, err, safe.ErrIntegerOverflow)
	assert.Equal(t, uint32(0), v)
}

func TestUint32ToInt_on_64bit(t *testing.T) {
	v, err := safe.Uint32ToInt(math.MaxUint32)
	assert.NoError(t, err)
	assert.Equal(t, math.MaxUint32, v)
}
