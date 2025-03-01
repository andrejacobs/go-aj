package ajio_test

import (
	"io"
	"os"
	"testing"

	"github.com/andrejacobs/go-aj/ajio"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMultiByteReaderSeeker(t *testing.T) {
	f, err := os.CreateTemp("", "MultiByteReaderSeeker")
	require.NoError(t, err)
	defer f.Close()

	_, err = f.WriteString("The quick brown fox jumped over the lazy dog!")
	require.NoError(t, err)
	_, err = f.Seek(0, io.SeekStart)
	require.NoError(t, err)

	rd := ajio.NewMultiByteReaderSeeker(f)

	buf := make([]byte, 10)
	_, err = rd.Read(buf)
	require.NoError(t, err)
	assert.Equal(t, []byte("The quick "), buf)

	b, err := rd.ReadByte()
	require.NoError(t, err)
	assert.Equal(t, byte('b'), b)

	rd.Seek(20, io.SeekStart)
	buf = buf[:6]
	_, err = rd.Read(buf)
	require.NoError(t, err)
	assert.Equal(t, []byte("jumped"), buf)
}
