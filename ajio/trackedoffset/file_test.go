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
	"os"
	"path/filepath"
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

	assert.Equal(t, uint64(0), tracker.Offset())

	tracker.SetOffset(0x414A)
	assert.Equal(t, uint64(0x414A), tracker.Offset())

	assert.NotNil(t, tracker.File())
	assert.Equal(t, tempFile, tracker.File().Name())
	assert.NotNil(t, tracker.Reader())
	assert.NotNil(t, tracker.Writer())
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

	tracker, err := trackedoffset.Open(tempFile)
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

	_, err = tracker.ReadByte()
	require.NoError(t, err)
	assert.Equal(t, uint64(7), tracker.Offset())
}

func TestFileWrite(t *testing.T) {
	fileSize := 10
	tempFile, err := random.CreateTempFile("", "unit-testing", 10)
	require.NoError(t, err)
	defer os.Remove(tempFile)

	tracker, err := trackedoffset.OpenFile(tempFile, os.O_RDWR, 0)
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

	require.NoError(t, tracker.WriteByte(0x41))
	assert.Equal(t, uint64(7), tracker.Offset())

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

func TestWriteRead(t *testing.T) {
	tempFile := filepath.Join(os.TempDir(), "unit-testing")
	_ = os.Remove(tempFile)
	defer os.Remove(tempFile)

	writer, err := trackedoffset.Create(tempFile)
	require.NoError(t, err)
	defer writer.Close()

	expected := []byte("The quick brown fox jumped over the lazy dog!")
	wc, err := writer.Write(expected)
	require.NoError(t, err)
	assert.Equal(t, len(expected), wc)

	require.NoError(t, writer.WriteByte(0x41))
	require.NoError(t, writer.WriteByte(0x4A))
	assert.Equal(t, uint64(len(expected)+2), writer.Offset())

	require.NoError(t, writer.Flush())

	// Validate

	reader, err := trackedoffset.Open(tempFile)
	require.NoError(t, err)
	defer writer.Close()

	buffer := make([]byte, len(expected))
	rc, err := reader.Read(buffer)
	require.NoError(t, err)
	assert.Equal(t, len(expected), rc)
	assert.Equal(t, expected, buffer)

	b, err := reader.ReadByte()
	require.NoError(t, err)
	assert.Equal(t, byte(0x41), b)

	b, err = reader.ReadByte()
	require.NoError(t, err)
	assert.Equal(t, byte(0x4A), b)

	assert.Equal(t, uint64(len(expected)+2), reader.Offset())
}
