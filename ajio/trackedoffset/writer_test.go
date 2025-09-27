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
