package ajio_test

import (
	"bufio"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/andrejacobs/go-aj/ajio"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//-----------------------------------------------------------------------------

func TestOffsetTrackedReader(t *testing.T) {
	text := "The quick brown fox jumped over the lazy dog!"
	sr := strings.NewReader(text)
	br := bufio.NewReader(sr)

	tr := ajio.NewOffsetTrackedReader(br, 0)
	assert.Equal(t, int64(0), tr.Offset())

	buffer := make([]byte, 4)
	for i := 0; i < len(text)/4; i++ {
		_, err := tr.Read(buffer)
		require.NoError(t, err)
		assert.Equal(t, int64((i+1)*4), tr.Offset())
	}
}

//-----------------------------------------------------------------------------

func TestMultiByteOffsetTrackedReader(t *testing.T) {
	text := "The quick brown fox jumped over the lazy dog!"
	sr := strings.NewReader(text)
	br := bufio.NewReader(sr)

	tr := ajio.NewMultiByteOffsetTrackedReader(br, 0)
	assert.Equal(t, int64(0), tr.Offset())

	buffer := make([]byte, 4)
	for i := 0; i < len(text)/4; i++ {
		_, err := tr.Read(buffer)
		require.NoError(t, err)
		assert.Equal(t, int64((i+1)*4), tr.Offset())
	}

	sr = strings.NewReader(text)
	br = bufio.NewReader(sr)
	for i := 0; i < 4; i++ {
		br.ReadByte()
	}

	tr = ajio.NewMultiByteOffsetTrackedReader(br, 4)
	assert.Equal(t, int64(4), tr.Offset())

	b, err := tr.ReadByte()
	require.NoError(t, err)
	assert.Equal(t, byte('q'), b)

	b, err = tr.ReadByte()
	require.NoError(t, err)
	assert.Equal(t, byte('u'), b)

	assert.Equal(t, int64(6), tr.Offset())
}

//-----------------------------------------------------------------------------

func TestNewFileOffsetTracker(t *testing.T) {
	tempFile, err := createTempFile(10)
	require.NoError(t, err)
	defer os.Remove(tempFile)

	f, err := os.Open(tempFile)
	require.NoError(t, err)
	defer f.Close()

	tracker, err := ajio.NewFileOffsetTracker(f)
	require.NoError(t, err)
	offset := tracker.Offset()
	assert.Equal(t, int64(0), offset)
}

func TestNewFileOffsetTrackerWithExistingOffset(t *testing.T) {
	tempFile, err := createTempFile(10)
	require.NoError(t, err)
	defer os.Remove(tempFile)

	f, err := os.Open(tempFile)
	require.NoError(t, err)
	defer f.Close()
	expectedOffset, err := f.Seek(5, io.SeekStart)
	require.NoError(t, err)

	tracker, err := ajio.NewFileOffsetTracker(f)
	require.NoError(t, err)
	offset := tracker.Offset()
	assert.Equal(t, expectedOffset, offset)
}

func TestFileOffsetTrackerSeek(t *testing.T) {
	tempFile, err := createTempFile(10)
	require.NoError(t, err)
	defer os.Remove(tempFile)

	f, err := os.Open(tempFile)
	require.NoError(t, err)
	defer f.Close()

	tracker, err := ajio.NewFileOffsetTracker(f)
	require.NoError(t, err)

	expectedOffset, err := tracker.Seek(4, io.SeekStart)
	require.NoError(t, err)
	assert.Equal(t, expectedOffset, tracker.Offset())
	assert.Equal(t, int64(4), expectedOffset)

	actualOffset, err := f.Seek(0, io.SeekCurrent)
	require.NoError(t, err)
	assert.Equal(t, actualOffset, expectedOffset)

	expectedOffset, err = tracker.Seek(-2, io.SeekCurrent)
	require.NoError(t, err)
	assert.Equal(t, expectedOffset, tracker.Offset())
	assert.Equal(t, int64(2), expectedOffset)

	actualOffset, err = f.Seek(0, io.SeekCurrent)
	require.NoError(t, err)
	assert.Equal(t, actualOffset, expectedOffset)

	expectedOffset, err = tracker.Seek(-3, io.SeekEnd)
	require.NoError(t, err)
	assert.Equal(t, expectedOffset, tracker.Offset())
	assert.Equal(t, int64(7), expectedOffset)

	actualOffset, err = f.Seek(0, io.SeekCurrent)
	require.NoError(t, err)
	assert.Equal(t, actualOffset, expectedOffset)
}

func TestFileOffsetTrackerRead(t *testing.T) {
	fileSize := 10
	tempFile, err := createTempFile(int64(fileSize))
	require.NoError(t, err)
	defer os.Remove(tempFile)

	f, err := os.Open(tempFile)
	require.NoError(t, err)
	defer f.Close()

	tracker, err := ajio.NewFileOffsetTracker(f)
	require.NoError(t, err)

	buffer := make([]byte, 2)

	for i := 0; i < fileSize; i += len(buffer) {
		rc, err := tracker.Read(buffer)
		require.NoError(t, err)
		assert.Equal(t, 2, rc)
		assert.Equal(t, int64(i+2), tracker.Offset())
	}

	_, err = tracker.Seek(4, io.SeekStart)
	require.NoError(t, err)
	_, err = tracker.Read(buffer)
	require.NoError(t, err)
	assert.Equal(t, int64(6), tracker.Offset())
}

func TestFileOffsetTrackerWrite(t *testing.T) {
	fileSize := 10
	tempFile, err := createTempFile(int64(fileSize))
	require.NoError(t, err)
	defer os.Remove(tempFile)

	f, err := os.OpenFile(tempFile, os.O_RDWR, 0)
	require.NoError(t, err)
	defer f.Close()

	tracker, err := ajio.NewFileOffsetTracker(f)
	require.NoError(t, err)

	buffer := make([]byte, 2)

	for i := 0; i < fileSize; i += len(buffer) {
		wc, err := tracker.Write(buffer)
		require.NoError(t, err)
		assert.Equal(t, 2, wc)
		assert.Equal(t, int64(i+2), tracker.Offset())
	}

	_, err = tracker.Seek(4, io.SeekStart)
	require.NoError(t, err)
	_, err = tracker.Write(buffer)
	require.NoError(t, err)
	assert.Equal(t, int64(6), tracker.Offset())
}

func TestFileOffsetTrackerSyncOffset(t *testing.T) {
	tempFile, err := createTempFile(10)
	require.NoError(t, err)
	defer os.Remove(tempFile)

	f, err := os.Open(tempFile)
	require.NoError(t, err)
	defer f.Close()

	tracker, err := ajio.NewFileOffsetTracker(f)
	require.NoError(t, err)

	// Out of sync
	actualOffset, err := f.Seek(2, io.SeekStart)
	require.NoError(t, err)
	assert.NotEqual(t, actualOffset, tracker.Offset())

	// Back in sync
	tracker.SyncOffset()
	actualOffset, err = f.Seek(0, io.SeekCurrent)
	require.NoError(t, err)
	assert.Equal(t, actualOffset, tracker.Offset())
}

func TestFileOffsetTrackerReaderAtWriterAt(t *testing.T) {
	fileSize := 10
	tempFile, err := createTempFile(int64(fileSize))
	require.NoError(t, err)
	defer os.Remove(tempFile)

	f, err := os.OpenFile(tempFile, os.O_RDWR, 0)
	require.NoError(t, err)
	defer f.Close()

	tracker, err := ajio.NewFileOffsetTracker(f)
	require.NoError(t, err)

	expected := []byte{0x41, 0x4a}
	wc, err := tracker.WriteAt(expected, 4)
	require.NoError(t, err)
	assert.Equal(t, 2, wc)
	assert.Equal(t, int64(6), tracker.Offset())

	buffer := make([]byte, 2)
	rc, err := tracker.ReadAt(buffer, 4)
	require.NoError(t, err)
	assert.Equal(t, 2, rc)
	assert.Equal(t, int64(6), tracker.Offset())

	assert.Equal(t, expected, buffer)
}
