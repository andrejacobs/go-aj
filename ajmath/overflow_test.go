package ajmath_test

import (
	"math"
	"testing"

	"github.com/andrejacobs/go-aj/ajmath"
	"github.com/stretchr/testify/assert"
)

func TestAdd32(t *testing.T) {
	v, err := ajmath.Add32(42, 42)
	assert.NoError(t, err)
	assert.Equal(t, uint32(84), v)

	v, err = ajmath.Add32(42, uint32(math.MaxUint32))
	assert.ErrorIs(t, err, ajmath.ErrIntegerOverflow)
	assert.Equal(t, uint32(0), v)
}

func TestAdd64(t *testing.T) {
	v, err := ajmath.Add64(42, 42)
	assert.NoError(t, err)
	assert.Equal(t, uint64(84), v)

	v, err = ajmath.Add64(42, uint64(math.MaxUint64))
	assert.ErrorIs(t, err, ajmath.ErrIntegerOverflow)
	assert.Equal(t, uint64(0), v)
}

func TestSub32(t *testing.T) {
	v, err := ajmath.Sub32(42, 42)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0), v)

	v, err = ajmath.Sub32(42, 45)
	assert.ErrorIs(t, err, ajmath.ErrIntegerUnderflow)
	assert.Equal(t, uint32(0), v)
}

func TestSub64(t *testing.T) {
	v, err := ajmath.Sub64(42, 42)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), v)

	v, err = ajmath.Sub64(42, 45)
	assert.ErrorIs(t, err, ajmath.ErrIntegerUnderflow)
	assert.Equal(t, uint64(0), v)
}
