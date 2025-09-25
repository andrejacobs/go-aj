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

package ajio

import (
	"io"
	"os"

	"github.com/andrejacobs/go-aj/ajmath"
)

// Keep track of the offset within an io.Reader source.
type TrackedOffsetReader interface {
	io.Reader

	// Return the current offset
	Offset() uint64
}

// Keep track of the offset within an io.Writer source.
type TrackedOffsetWriter interface {
	io.Writer

	// Return the current offset
	Offset() uint64
}

// TrackedOffset keeps track of the current offset within file like objects without requiring to make calls to Seek.
type TrackedOffset interface {
	io.Reader
	io.Writer
	io.Seeker

	io.ReaderAt
	io.WriterAt

	// Return the current offset
	Offset() uint64

	// Ensure the current offset matches the underlying source. In the case of a file like object this should involve a call to Seek.
	SyncOffset() error
}

// MultiByteTrackedOffsetReader is the same as TrackedOffsetReader but with the io.ByteReader added.
type MultiByteTrackedOffsetReader interface {
	TrackedOffsetReader
	io.ByteReader
}

//-----------------------------------------------------------------------------
// TrackedOffsetReader

type reader struct {
	rd     io.Reader
	offset uint64
}

// Create a new TrackedOffsetReader that will keep track of the offset within the source io.Reader object.
func NewTrackedOffsetReader(rd io.Reader, baseOffset uint64) TrackedOffsetReader {
	t := &reader{
		rd:     rd,
		offset: baseOffset,
	}
	return t
}

// Reader implementation.
func (t *reader) Read(p []byte) (int, error) {
	n, err := t.rd.Read(p)
	if err != nil {
		return n, err
	}

	newOffset, err := ajmath.Add64(t.offset, uint64(n))
	if err != nil {
		return 0, err
	}
	t.offset = newOffset

	return n, nil
}

// TrackedOffsetReader implementation.
func (t *reader) Offset() uint64 {
	return t.offset
}

//-----------------------------------------------------------------------------
// TrackedOffsetWriter

type writer struct {
	wd     io.Writer
	offset uint64
}

// Create a new TrackedOffsetWriter that will keep track of the offset within the source io.Writer object.
func NewTrackedOffsetWriter(wd io.Writer, baseOffset uint64) TrackedOffsetWriter {
	t := &writer{
		wd:     wd,
		offset: baseOffset,
	}
	return t
}

// Writer implementation.
func (t *writer) Write(p []byte) (int, error) {
	n, err := t.wd.Write(p)
	if err != nil {
		return n, err
	}

	newOffset, err := ajmath.Add64(t.offset, uint64(n))
	if err != nil {
		return 0, err
	}
	t.offset = newOffset

	return n, nil
}

// TrackedOffsetWriter implementation.
func (t *writer) Offset() uint64 {
	return t.offset
}

//-----------------------------------------------------------------------------
// TrackedOffset file

// Wrap os.File to keep track of the current offset without needing to make constant calls to Seek which involves syscall Lseek.
type fileTrackedOffset struct {
	f      *os.File
	offset uint64
}

// Create a new TrackedOffset that will keep track of the file's offset.
// NOTE: An initital Seek will be called on the file to establish the current offset.
func NewTrackedOffsetFile(f *os.File) (TrackedOffset, error) {
	t := &fileTrackedOffset{
		f: f,
	}

	if err := t.SyncOffset(); err != nil {
		return nil, err
	}

	return t, nil
}

// Reader implementation.
func (t *fileTrackedOffset) Read(p []byte) (int, error) {
	n, err := t.f.Read(p)
	if err != nil {
		return n, err
	}

	newOffset, err := ajmath.Add64(t.offset, uint64(n))
	if err != nil {
		return 0, err
	}
	t.offset = newOffset

	return n, nil
}

// Writer implementation.
func (t *fileTrackedOffset) Write(p []byte) (int, error) {
	n, err := t.f.Write(p)
	if err != nil {
		return n, err
	}

	newOffset, err := ajmath.Add64(t.offset, uint64(n))
	if err != nil {
		return 0, err
	}
	t.offset = newOffset

	return n, nil
}

// Seeker implementation.
func (t *fileTrackedOffset) Seek(offset int64, whence int) (int64, error) {
	newOffset, err := t.f.Seek(offset, whence)
	if err != nil {
		return newOffset, err
	}

	t.offset, err = ajmath.Int64ToUint64(newOffset)
	if err != nil {
		return newOffset, err
	}

	return newOffset, nil
}

// ReaderAt implementation.
func (t *fileTrackedOffset) ReadAt(p []byte, off int64) (int, error) {
	n, err := t.f.ReadAt(p, off)
	if err != nil {
		return n, err
	}

	newOffset, err := ajmath.Add64(uint64(off), uint64(n))
	if err != nil {
		return 0, err
	}
	t.offset = newOffset

	return n, nil
}

// WriterAt implementation.
func (t *fileTrackedOffset) WriteAt(p []byte, off int64) (int, error) {
	n, err := t.f.WriteAt(p, off)
	if err != nil {
		return n, err
	}

	newOffset, err := ajmath.Add64(uint64(off), uint64(n))
	if err != nil {
		return 0, err
	}
	t.offset = newOffset

	return n, nil
}

//-----------------------------------------------------------------------------

// TrackedOffset implementation.
func (t *fileTrackedOffset) Offset() uint64 {
	return t.offset
}

// Ensure the tracker's offset and the file's actual offsets are the same.
// This will make a call to file.Seek.
func (t *fileTrackedOffset) SyncOffset() error {
	offset, err := t.f.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}

	t.offset, err = ajmath.Int64ToUint64(offset)
	if err != nil {
		return err
	}

	return nil
}

//-----------------------------------------------------------------------------
// MultiByteTrackedOffsetReader

type wrappedMBReader struct {
	rd     MultiByteReader
	offset uint64
}

// Create a new bufio.Reader that also supports being able to do io.Seeker.
func NewTrackedOffsetReaderMultiByte(rd MultiByteReader, baseOffset uint64) MultiByteTrackedOffsetReader {
	return &wrappedMBReader{
		rd:     rd,
		offset: baseOffset,
	}
}

func (w *wrappedMBReader) Read(p []byte) (int, error) {
	n, err := w.rd.Read(p)
	if err != nil {
		return n, err
	}

	newOffset, err := ajmath.Add64(w.offset, uint64(n))
	if err != nil {
		return 0, err
	}
	w.offset = newOffset

	return n, nil
}

func (w *wrappedMBReader) ReadByte() (byte, error) {
	b, err := w.rd.ReadByte()
	if err != nil {
		return b, err
	}

	newOffset, err := ajmath.Add64(w.offset, uint64(1))
	if err != nil {
		return 0, err
	}
	w.offset = newOffset

	return b, nil
}

func (w *wrappedMBReader) Offset() uint64 {
	return w.offset
}

//-----------------------------------------------------------------------------
