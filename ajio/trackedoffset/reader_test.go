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
