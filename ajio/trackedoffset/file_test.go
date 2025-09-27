package trackedoffset_test

import (
	"io"
	"os"
	"testing"

	"github.com/andrejacobs/go-aj/ajio/trackedoffset"
	"github.com/andrejacobs/go-aj/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFile(t *testing.T) {
	tempFile, err := random.CreateTempFile("", "unit-testing", 10)
	require.NoError(t, err)
	defer os.Remove(tempFile)

	f, err := os.Open(tempFile)
	require.NoError(t, err)

	tracker, err := trackedoffset.NewFile(f)
	require.NoError(t, err)
	defer tracker.Close()

	assert.Equal(t, f.Name(), tracker.Name())
	_, err = tracker.Stat()
	assert.NoError(t, err)

	offset := tracker.Offset()
	assert.Equal(t, uint64(0), offset)
}

func TestNewFileWithExistingOffset(t *testing.T) {
	tempFile, err := random.CreateTempFile("", "unit-testing", 10)
	require.NoError(t, err)
	defer os.Remove(tempFile)

	f, err := os.Open(tempFile)
	require.NoError(t, err)
	expectedOffset, err := f.Seek(5, io.SeekStart)
	require.NoError(t, err)

	tracker, err := trackedoffset.NewFile(f)
	require.NoError(t, err)
	defer tracker.Close()
	offset := tracker.Offset()
	assert.Equal(t, uint64(expectedOffset), offset)
}

func TestFileSeek(t *testing.T) {
	tempFile, err := random.CreateTempFile("", "unit-testing", 10)
	require.NoError(t, err)
	defer os.Remove(tempFile)

	f, err := os.Open(tempFile)
	require.NoError(t, err)

	tracker, err := trackedoffset.NewFile(f)
	require.NoError(t, err)
	defer tracker.Close()

	expectedOffset, err := tracker.Seek(4, io.SeekStart)
	require.NoError(t, err)
	assert.Equal(t, uint64(expectedOffset), tracker.Offset())
	assert.Equal(t, int64(4), expectedOffset)

	actualOffset, err := f.Seek(0, io.SeekCurrent)
	require.NoError(t, err)
	assert.Equal(t, actualOffset, expectedOffset)

	expectedOffset, err = tracker.Seek(-2, io.SeekCurrent)
	require.NoError(t, err)
	assert.Equal(t, uint64(expectedOffset), tracker.Offset())
	assert.Equal(t, int64(2), expectedOffset)

	actualOffset, err = f.Seek(0, io.SeekCurrent)
	require.NoError(t, err)
	assert.Equal(t, actualOffset, expectedOffset)

	expectedOffset, err = tracker.Seek(-3, io.SeekEnd)
	require.NoError(t, err)
	assert.Equal(t, uint64(expectedOffset), tracker.Offset())
	assert.Equal(t, int64(7), expectedOffset)

	actualOffset, err = f.Seek(0, io.SeekCurrent)
	require.NoError(t, err)
	assert.Equal(t, actualOffset, expectedOffset)
}

func TestFileRead(t *testing.T) {
	fileSize := 10
	tempFile, err := random.CreateTempFile("", "unit-testing", 10)
	require.NoError(t, err)
	defer os.Remove(tempFile)

	f, err := os.Open(tempFile)
	require.NoError(t, err)

	tracker, err := trackedoffset.NewFile(f)
	require.NoError(t, err)
	defer tracker.Close()

	buffer := make([]byte, 2)

	for i := 0; i < fileSize; i += len(buffer) {
		rc, err := tracker.Read(buffer)
		require.NoError(t, err)
		assert.Equal(t, 2, rc)
		assert.Equal(t, uint64(i+2), tracker.Offset())
	}

	_, err = tracker.Seek(4, io.SeekStart)
	require.NoError(t, err)
	tracker.ResetReadBuffer()

	_, err = tracker.Read(buffer)
	require.NoError(t, err)
	assert.Equal(t, uint64(6), tracker.Offset())
}

func TestFileWrite(t *testing.T) {
	fileSize := 10
	tempFile, err := random.CreateTempFile("", "unit-testing", 10)
	require.NoError(t, err)
	defer os.Remove(tempFile)

	f, err := os.OpenFile(tempFile, os.O_RDWR, 0)
	require.NoError(t, err)

	tracker, err := trackedoffset.NewFile(f)
	require.NoError(t, err)
	defer tracker.Close()

	buffer := make([]byte, 2)

	for i := 0; i < fileSize; i += len(buffer) {
		wc, err := tracker.Write(buffer)
		require.NoError(t, err)
		assert.Equal(t, 2, wc)
		assert.Equal(t, uint64(i+2), tracker.Offset())
	}

	_, err = tracker.Seek(4, io.SeekStart)
	require.NoError(t, err)
	tracker.ResetWriteBuffer()

	_, err = tracker.Write(buffer)
	require.NoError(t, err)
	assert.Equal(t, uint64(6), tracker.Offset())

	err = tracker.Flush()
	assert.NoError(t, err)
	err = tracker.Sync()
	assert.NoError(t, err)
}

func TestFileSyncOffset(t *testing.T) {
	tempFile, err := random.CreateTempFile("", "unit-testing", 10)
	require.NoError(t, err)
	defer os.Remove(tempFile)

	f, err := os.Open(tempFile)
	require.NoError(t, err)

	tracker, err := trackedoffset.NewFile(f)
	require.NoError(t, err)
	defer tracker.Close()

	// Out of sync
	actualOffset, err := f.Seek(2, io.SeekStart)
	require.NoError(t, err)
	assert.NotEqual(t, actualOffset, tracker.Offset())

	// Back in sync
	tracker.SyncOffset()
	actualOffset, err = f.Seek(0, io.SeekCurrent)
	require.NoError(t, err)
	assert.Equal(t, uint64(actualOffset), tracker.Offset())
}
