// Copyright (c) 2025 Andre Jacobs
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package ajio_test

import (
	"bufio"
	"io"
	"math"
	"os"
	"strings"
	"testing"

	"github.com/andrejacobs/go-aj/ajio"
	"github.com/andrejacobs/go-aj/ajmath"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//-----------------------------------------------------------------------------

func TestTrackedOffsetReader(t *testing.T) {
	text := "The quick brown fox jumped over the lazy dog!"
	sr := strings.NewReader(text)
	br := bufio.NewReader(sr)

	baseOffset := uint64(42)
	tr := ajio.NewTrackedOffsetReader(br, baseOffset)
	assert.Equal(t, baseOffset, tr.Offset())

	buffer := make([]byte, 4)
	for i := 0; i < len(text)/4; i++ {
		_, err := tr.Read(buffer)
		require.NoError(t, err)
		assert.Equal(t, baseOffset+uint64((i+1)*4), tr.Offset())
	}
}

func TestTrackedOffsetReaderOverflow(t *testing.T) {
	text := "The quick brown fox jumped over the lazy dog!"
	sr := strings.NewReader(text)
	br := bufio.NewReader(sr)

	baseOffset := uint64(math.MaxUint64 - 2)
	tr := ajio.NewTrackedOffsetReader(br, uint64(baseOffset))
	assert.Equal(t, baseOffset, tr.Offset())

	buffer := make([]byte, 4)
	_, err := tr.Read(buffer)
	assert.ErrorIs(t, err, ajmath.ErrIntegerOverflow)
}

//-----------------------------------------------------------------------------

func TestTrackedOffsetWriter(t *testing.T) {
	baseOffset := uint64(42)
	tw := ajio.NewTrackedOffsetWriter(io.Discard, baseOffset)
	assert.Equal(t, baseOffset, tw.Offset())

	data := []byte("The quick brown fox jumped over the lazy dog!")
	c, err := tw.Write(data)
	require.NoError(t, err)
	assert.Equal(t, len(data), c)
	assert.Equal(t, baseOffset+uint64(len(data)), tw.Offset())
}

func TestTrackedOffsetWriterOverflow(t *testing.T) {
	baseOffset := uint64(math.MaxUint64 - 4)
	tw := ajio.NewTrackedOffsetWriter(io.Discard, baseOffset)
	assert.Equal(t, baseOffset, tw.Offset())

	data := []byte("The quick brown fox jumped over the lazy dog!")
	_, err := tw.Write(data)
	assert.ErrorIs(t, err, ajmath.ErrIntegerOverflow)
}

//-----------------------------------------------------------------------------

func TestMultiByteTrackedOffsetReader(t *testing.T) {
	text := "The quick brown fox jumped over the lazy dog!"
	sr := strings.NewReader(text)
	br := bufio.NewReader(sr)

	tr := ajio.NewTrackedOffsetReaderMultiByte(br, 0)
	assert.Equal(t, uint64(0), tr.Offset())

	buffer := make([]byte, 4)
	for i := 0; i < len(text)/4; i++ {
		_, err := tr.Read(buffer)
		require.NoError(t, err)
		assert.Equal(t, uint64((i+1)*4), tr.Offset())
	}

	sr = strings.NewReader(text)
	br = bufio.NewReader(sr)
	for i := 0; i < 4; i++ {
		br.ReadByte()
	}

	tr = ajio.NewTrackedOffsetReaderMultiByte(br, 4)
	assert.Equal(t, uint64(4), tr.Offset())

	b, err := tr.ReadByte()
	require.NoError(t, err)
	assert.Equal(t, byte('q'), b)

	b, err = tr.ReadByte()
	require.NoError(t, err)
	assert.Equal(t, byte('u'), b)

	assert.Equal(t, uint64(6), tr.Offset())
}

func TestMultiByteTrackedOffsetReaderOverflow(t *testing.T) {
	text := "The quick brown fox jumped over the lazy dog!"
	sr := strings.NewReader(text)
	br := bufio.NewReader(sr)

	baseOffset := uint64(math.MaxUint64 - 2)
	tr := ajio.NewTrackedOffsetReaderMultiByte(br, baseOffset)
	assert.Equal(t, baseOffset, tr.Offset())

	buffer := make([]byte, 4)
	_, err := tr.Read(buffer)
	assert.ErrorIs(t, err, ajmath.ErrIntegerOverflow)

	_, err = tr.ReadByte()
	assert.NoError(t, err)
	_, err = tr.ReadByte()
	assert.NoError(t, err)

	_, err = tr.ReadByte()
	assert.ErrorIs(t, err, ajmath.ErrIntegerOverflow)
}

//-----------------------------------------------------------------------------

func TestNewTrackedOffsetFile(t *testing.T) {
	tempFile, err := createTempFile(10)
	require.NoError(t, err)
	defer os.Remove(tempFile)

	f, err := os.Open(tempFile)
	require.NoError(t, err)
	defer f.Close()

	tracker, err := ajio.NewTrackedOffsetFile(f)
	require.NoError(t, err)
	offset := tracker.Offset()
	assert.Equal(t, uint64(0), offset)
}

func TestNewTrackedOffsetFileWithExistingOffset(t *testing.T) {
	tempFile, err := createTempFile(10)
	require.NoError(t, err)
	defer os.Remove(tempFile)

	f, err := os.Open(tempFile)
	require.NoError(t, err)
	defer f.Close()
	expectedOffset, err := f.Seek(5, io.SeekStart)
	require.NoError(t, err)

	tracker, err := ajio.NewTrackedOffsetFile(f)
	require.NoError(t, err)
	offset := tracker.Offset()
	assert.Equal(t, uint64(expectedOffset), offset)
}

func TestFileTrackedOffsetSeek(t *testing.T) {
	tempFile, err := createTempFile(10)
	require.NoError(t, err)
	defer os.Remove(tempFile)

	f, err := os.Open(tempFile)
	require.NoError(t, err)
	defer f.Close()

	tracker, err := ajio.NewTrackedOffsetFile(f)
	require.NoError(t, err)

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

func TestFileTrackedOffsetRead(t *testing.T) {
	fileSize := 10
	tempFile, err := createTempFile(int64(fileSize))
	require.NoError(t, err)
	defer os.Remove(tempFile)

	f, err := os.Open(tempFile)
	require.NoError(t, err)
	defer f.Close()

	tracker, err := ajio.NewTrackedOffsetFile(f)
	require.NoError(t, err)

	buffer := make([]byte, 2)

	for i := 0; i < fileSize; i += len(buffer) {
		rc, err := tracker.Read(buffer)
		require.NoError(t, err)
		assert.Equal(t, 2, rc)
		assert.Equal(t, uint64(i+2), tracker.Offset())
	}

	_, err = tracker.Seek(4, io.SeekStart)
	require.NoError(t, err)
	_, err = tracker.Read(buffer)
	require.NoError(t, err)
	assert.Equal(t, uint64(6), tracker.Offset())
}

func TestFileTrackedOffsetWrite(t *testing.T) {
	fileSize := 10
	tempFile, err := createTempFile(int64(fileSize))
	require.NoError(t, err)
	defer os.Remove(tempFile)

	f, err := os.OpenFile(tempFile, os.O_RDWR, 0)
	require.NoError(t, err)
	defer f.Close()

	tracker, err := ajio.NewTrackedOffsetFile(f)
	require.NoError(t, err)

	buffer := make([]byte, 2)

	for i := 0; i < fileSize; i += len(buffer) {
		wc, err := tracker.Write(buffer)
		require.NoError(t, err)
		assert.Equal(t, 2, wc)
		assert.Equal(t, uint64(i+2), tracker.Offset())
	}

	_, err = tracker.Seek(4, io.SeekStart)
	require.NoError(t, err)
	_, err = tracker.Write(buffer)
	require.NoError(t, err)
	assert.Equal(t, uint64(6), tracker.Offset())
}

func TestFileTrackedOffsetSyncOffset(t *testing.T) {
	tempFile, err := createTempFile(10)
	require.NoError(t, err)
	defer os.Remove(tempFile)

	f, err := os.Open(tempFile)
	require.NoError(t, err)
	defer f.Close()

	tracker, err := ajio.NewTrackedOffsetFile(f)
	require.NoError(t, err)

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

func TestFileTrackedOffsetReaderAtWriterAt(t *testing.T) {
	fileSize := 10
	tempFile, err := createTempFile(int64(fileSize))
	require.NoError(t, err)
	defer os.Remove(tempFile)

	f, err := os.OpenFile(tempFile, os.O_RDWR, 0)
	require.NoError(t, err)
	defer f.Close()

	tracker, err := ajio.NewTrackedOffsetFile(f)
	require.NoError(t, err)

	expected := []byte{0x41, 0x4a}
	wc, err := tracker.WriteAt(expected, 4)
	require.NoError(t, err)
	assert.Equal(t, 2, wc)
	assert.Equal(t, uint64(6), tracker.Offset())

	buffer := make([]byte, 2)
	rc, err := tracker.ReadAt(buffer, 4)
	require.NoError(t, err)
	assert.Equal(t, 2, rc)
	assert.Equal(t, uint64(6), tracker.Offset())

	assert.Equal(t, expected, buffer)
}
