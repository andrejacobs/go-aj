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

package safe_test

import (
	"math"
	"strconv"
	"testing"

	"github.com/andrejacobs/go-aj/ajmath/safe"
	"github.com/stretchr/testify/assert"
)

func TestAdd32(t *testing.T) {
	v, err := safe.Add32(42, 42)
	assert.NoError(t, err)
	assert.Equal(t, uint32(84), v)

	v, err = safe.Add32(42, uint32(math.MaxUint32))
	assert.ErrorIs(t, err, safe.ErrIntegerOverflow)
	assert.Equal(t, uint32(0), v)
}

func TestAdd64(t *testing.T) {
	v, err := safe.Add64(42, 42)
	assert.NoError(t, err)
	assert.Equal(t, uint64(84), v)

	v, err = safe.Add64(42, uint64(math.MaxUint64))
	assert.ErrorIs(t, err, safe.ErrIntegerOverflow)
	assert.Equal(t, uint64(0), v)
}

func TestSub32(t *testing.T) {
	v, err := safe.Sub32(42, 42)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0), v)

	v, err = safe.Sub32(42, 45)
	assert.ErrorIs(t, err, safe.ErrIntegerUnderflow)
	assert.Equal(t, uint32(0), v)
}

func TestSub64(t *testing.T) {
	v, err := safe.Sub64(42, 42)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), v)

	v, err = safe.Sub64(42, 45)
	assert.ErrorIs(t, err, safe.ErrIntegerUnderflow)
	assert.Equal(t, uint64(0), v)
}

//-----------------------------------------------------------------------------
// Casting

func TestInt8ToUint8(t *testing.T) {
	v, err := safe.Int8ToUint8(0)
	assert.NoError(t, err)
	assert.Equal(t, uint8(0), v)

	v, err = safe.Int8ToUint8(math.MaxInt8)
	assert.NoError(t, err)
	assert.Equal(t, uint8(math.MaxInt8), v)

	v, err = safe.Int8ToUint8(-42)
	assert.ErrorIs(t, err, safe.ErrIntegerUnderflow)
	assert.Equal(t, uint8(0), v)
}

func TestUint8ToInt8(t *testing.T) {
	v, err := safe.Uint8ToInt8(0)
	assert.NoError(t, err)
	assert.Equal(t, int8(0), v)

	v, err = safe.Uint8ToInt8(math.MaxInt8)
	assert.NoError(t, err)
	assert.Equal(t, int8(math.MaxInt8), v)

	v, err = safe.Uint8ToInt8(math.MaxInt8 + 1)
	assert.ErrorIs(t, err, safe.ErrIntegerOverflow)
	assert.Equal(t, int8(0), v)
}

func TestInt16ToUint16(t *testing.T) {
	v, err := safe.Int16ToUint16(0)
	assert.NoError(t, err)
	assert.Equal(t, uint16(0), v)

	v, err = safe.Int16ToUint16(math.MaxInt16)
	assert.NoError(t, err)
	assert.Equal(t, uint16(math.MaxInt16), v)

	v, err = safe.Int16ToUint16(-42)
	assert.ErrorIs(t, err, safe.ErrIntegerUnderflow)
	assert.Equal(t, uint16(0), v)
}

func TestUint16ToInt16(t *testing.T) {
	v, err := safe.Uint16ToInt16(0)
	assert.NoError(t, err)
	assert.Equal(t, int16(0), v)

	v, err = safe.Uint16ToInt16(math.MaxInt16)
	assert.NoError(t, err)
	assert.Equal(t, int16(math.MaxInt16), v)

	v, err = safe.Uint16ToInt16(math.MaxInt16 + 1)
	assert.ErrorIs(t, err, safe.ErrIntegerOverflow)
	assert.Equal(t, int16(0), v)
}

func TestInt32ToUint32(t *testing.T) {
	v, err := safe.Int32ToUint32(0)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0), v)

	v, err = safe.Int32ToUint32(math.MaxInt32)
	assert.NoError(t, err)
	assert.Equal(t, uint32(math.MaxInt32), v)

	v, err = safe.Int32ToUint32(-42)
	assert.ErrorIs(t, err, safe.ErrIntegerUnderflow)
	assert.Equal(t, uint32(0), v)
}

func TestUint32ToInt32(t *testing.T) {
	v, err := safe.Uint32ToInt32(0)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), v)

	v, err = safe.Uint32ToInt32(math.MaxInt32)
	assert.NoError(t, err)
	assert.Equal(t, int32(math.MaxInt32), v)

	v, err = safe.Uint32ToInt32(math.MaxInt32 + 1)
	assert.ErrorIs(t, err, safe.ErrIntegerOverflow)
	assert.Equal(t, int32(0), v)
}

func TestInt64ToUint64(t *testing.T) {
	v, err := safe.Int64ToUint64(0)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), v)

	v, err = safe.Int64ToUint64(math.MaxInt64)
	assert.NoError(t, err)
	assert.Equal(t, uint64(math.MaxInt64), v)

	v, err = safe.Int64ToUint64(-42)
	assert.ErrorIs(t, err, safe.ErrIntegerUnderflow)
	assert.Equal(t, uint64(0), v)
}

func TestUint64ToInt64(t *testing.T) {
	v, err := safe.Uint64ToInt64(0)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), v)

	v, err = safe.Uint64ToInt64(math.MaxInt64)
	assert.NoError(t, err)
	assert.Equal(t, int64(math.MaxInt64), v)

	v, err = safe.Uint64ToInt64(math.MaxInt64 + 1)
	assert.ErrorIs(t, err, safe.ErrIntegerOverflow)
	assert.Equal(t, int64(0), v)
}

//-----------------------------------------------------------------------------
// Downcasting

func TestUint64ToUint32(t *testing.T) {
	v, err := safe.Uint64ToUint32(0)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0), v)

	v, err = safe.Uint64ToUint32(4242)
	assert.NoError(t, err)
	assert.Equal(t, uint32(4242), v)

	v, err = safe.Uint64ToUint32(math.MaxUint32)
	assert.NoError(t, err)
	assert.Equal(t, uint32(math.MaxUint32), v)

	v, err = safe.Uint64ToUint32(math.MaxUint32 + 1)
	assert.ErrorIs(t, err, safe.ErrIntegerOverflow)
	assert.Equal(t, uint32(0), v)
}

func TestInt64ToInt32(t *testing.T) {
	v, err := safe.Int64ToInt32(0)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), v)

	v, err = safe.Int64ToInt32(4242)
	assert.NoError(t, err)
	assert.Equal(t, int32(4242), v)

	v, err = safe.Int64ToInt32(math.MaxInt32)
	assert.NoError(t, err)
	assert.Equal(t, int32(math.MaxInt32), v)

	v, err = safe.Int64ToInt32(math.MaxInt32 + 1)
	assert.ErrorIs(t, err, safe.ErrIntegerOverflow)
	assert.Equal(t, int32(0), v)
}

func TestInt64ToUint32(t *testing.T) {
	v, err := safe.Int64ToUint32(0)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0), v)

	v, err = safe.Int64ToUint32(4242)
	assert.NoError(t, err)
	assert.Equal(t, uint32(4242), v)

	v, err = safe.Int64ToUint32(-42)
	assert.ErrorIs(t, err, safe.ErrIntegerUnderflow)
	assert.Equal(t, uint32(0), v)

	v, err = safe.Int64ToUint32(math.MaxUint32 + 1)
	assert.ErrorIs(t, err, safe.ErrIntegerOverflow)
	assert.Equal(t, uint32(0), v)
}

func TestUint64ToInt32(t *testing.T) {
	v, err := safe.Uint64ToInt32(0)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), v)

	v, err = safe.Uint64ToInt32(4242)
	assert.NoError(t, err)
	assert.Equal(t, int32(4242), v)

	v, err = safe.Uint64ToInt32(math.MaxInt32 + 1)
	assert.ErrorIs(t, err, safe.ErrIntegerOverflow)
	assert.Equal(t, int32(0), v)
}

func TestIntToInt8(t *testing.T) {
	v, err := safe.IntToInt8(0)
	assert.NoError(t, err)
	assert.Equal(t, int8(0), v)

	v, err = safe.IntToInt8(42)
	assert.NoError(t, err)
	assert.Equal(t, int8(42), v)

	v, err = safe.IntToInt8(math.MinInt8)
	assert.NoError(t, err)
	assert.Equal(t, int8(math.MinInt8), v)

	v, err = safe.IntToInt8(math.MaxInt8)
	assert.NoError(t, err)
	assert.Equal(t, int8(math.MaxInt8), v)

	v, err = safe.IntToInt8(math.MinInt8 - 1)
	assert.ErrorIs(t, err, safe.ErrIntegerUnderflow)
	assert.Equal(t, int8(0), v)

	v, err = safe.IntToInt8(math.MaxInt8 + 1)
	assert.ErrorIs(t, err, safe.ErrIntegerOverflow)
	assert.Equal(t, int8(0), v)
}

func TestIntToInt16(t *testing.T) {
	v, err := safe.IntToInt16(0)
	assert.NoError(t, err)
	assert.Equal(t, int16(0), v)

	v, err = safe.IntToInt16(42)
	assert.NoError(t, err)
	assert.Equal(t, int16(42), v)

	v, err = safe.IntToInt16(math.MinInt16)
	assert.NoError(t, err)
	assert.Equal(t, int16(math.MinInt16), v)

	v, err = safe.IntToInt16(math.MaxInt16)
	assert.NoError(t, err)
	assert.Equal(t, int16(math.MaxInt16), v)

	v, err = safe.IntToInt16(math.MinInt16 - 1)
	assert.ErrorIs(t, err, safe.ErrIntegerUnderflow)
	assert.Equal(t, int16(0), v)

	v, err = safe.IntToInt16(math.MaxInt16 + 1)
	assert.ErrorIs(t, err, safe.ErrIntegerOverflow)
	assert.Equal(t, int16(0), v)
}

func TestIntToInt32(t *testing.T) {
	v, err := safe.IntToInt32(0)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), v)

	v, err = safe.IntToInt32(42)
	assert.NoError(t, err)
	assert.Equal(t, int32(42), v)

	v, err = safe.IntToInt32(math.MinInt32)
	assert.NoError(t, err)
	assert.Equal(t, int32(math.MinInt32), v)

	v, err = safe.IntToInt32(math.MaxInt32)
	assert.NoError(t, err)
	assert.Equal(t, int32(math.MaxInt32), v)
}

func TestUintToUint8(t *testing.T) {
	v, err := safe.UintToUint8(0)
	assert.NoError(t, err)
	assert.Equal(t, uint8(0), v)

	v, err = safe.UintToUint8(42)
	assert.NoError(t, err)
	assert.Equal(t, uint8(42), v)

	v, err = safe.UintToUint8(math.MaxUint8)
	assert.NoError(t, err)
	assert.Equal(t, uint8(math.MaxUint8), v)

	v, err = safe.UintToUint8(math.MaxUint8 + 1)
	assert.ErrorIs(t, err, safe.ErrIntegerOverflow)
	assert.Equal(t, uint8(0), v)
}

func TestUintToUint16(t *testing.T) {
	v, err := safe.UintToUint16(0)
	assert.NoError(t, err)
	assert.Equal(t, uint16(0), v)

	v, err = safe.UintToUint16(42)
	assert.NoError(t, err)
	assert.Equal(t, uint16(42), v)

	v, err = safe.UintToUint16(math.MaxUint16)
	assert.NoError(t, err)
	assert.Equal(t, uint16(math.MaxUint16), v)

	v, err = safe.UintToUint16(math.MaxUint16 + 1)
	assert.ErrorIs(t, err, safe.ErrIntegerOverflow)
	assert.Equal(t, uint16(0), v)
}

func TestUintToUint32(t *testing.T) {
	v, err := safe.UintToUint32(0)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0), v)

	v, err = safe.UintToUint32(42)
	assert.NoError(t, err)
	assert.Equal(t, uint32(42), v)

	v, err = safe.UintToUint32(math.MaxUint32)
	assert.NoError(t, err)
	assert.Equal(t, uint32(math.MaxUint32), v)
}

//-----------------------------------------------------------------------------
// Platform safety

func TestIntSize(t *testing.T) {
	assert.Equal(t, strconv.IntSize, safe.IntSize)
}

func TestUint32ToInt(t *testing.T) {
	// 32 bit machine
	if safe.IntSize == 32 {
		v, err := safe.Uint32ToInt(math.MaxUint32)
		assert.ErrorIs(t, err, safe.ErrIntegerOverflow)
		assert.Equal(t, 0, v)
	}
}

func TestUint64ToInt(t *testing.T) {
	if safe.IntSize == 32 {
		// 32 bit machine
		v, err := safe.Uint64ToInt(math.MaxUint32)
		assert.ErrorIs(t, err, safe.ErrIntegerOverflow)
		assert.Equal(t, 0, v)

		v, err = safe.Uint64ToInt(math.MaxUint64)
		assert.ErrorIs(t, err, safe.ErrIntegerOverflow)
		assert.Equal(t, 0, v)
	} else {
		// 64 bit machine
		v, err := safe.Uint64ToInt(math.MaxUint64)
		assert.ErrorIs(t, err, safe.ErrIntegerOverflow)
		assert.Equal(t, 0, v)
	}
}
