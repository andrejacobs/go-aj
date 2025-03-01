package ajio

import (
	"io"
	"os"
)

// Keep track of the offset within an io.Reader source
type TrackedOffsetReader interface {
	io.Reader

	// Return the current offset
	Offset() int64
}

// Keep track of the offset within an io.Writer source
type TrackedOffsetWriter interface {
	io.Writer

	// Return the current offset
	Offset() int64
}

// TrackedOffset keeps track of the current offset within file like objects without requiring to make calls to Seek
type TrackedOffset interface {
	io.Reader
	io.Writer
	io.Seeker

	io.ReaderAt
	io.WriterAt

	// Return the current offset
	Offset() int64

	// Ensure the current offset matches the underlying source. In the case of a file like object this should involve a call to Seek.
	SyncOffset() error
}

// MultiByteTrackedOffsetReader is the same as TrackedOffsetReader but with the io.ByteReader added
type MultiByteTrackedOffsetReader interface {
	TrackedOffsetReader
	io.ByteReader
}

//-----------------------------------------------------------------------------

type reader struct {
	rd     io.Reader
	offset int64
}

// Create a new TrackedOffsetReader that will keep track of the offset within the source io.Reader object.
func NewTrackedOffsetReader(rd io.Reader, baseOffset int64) TrackedOffsetReader {
	t := &reader{
		rd:     rd,
		offset: int64(baseOffset),
	}
	return t
}

// Reader implementation
func (t *reader) Read(p []byte) (int, error) {
	n, err := t.rd.Read(p)
	if err != nil {
		return n, err
	}

	t.offset += int64(n)

	return n, nil
}

// TrackedOffsetReader implementation
func (t *reader) Offset() int64 {
	return t.offset
}

//-----------------------------------------------------------------------------

type writer struct {
	wd     io.Writer
	offset int64
}

// Create a new TrackedOffsetWriter that will keep track of the offset within the source io.Writer object.
func NewTrackedOffsetWriter(wd io.Writer, baseOffset int64) TrackedOffsetWriter {
	t := &writer{
		wd:     wd,
		offset: int64(baseOffset),
	}
	return t
}

// Writer implementation
func (t *writer) Write(p []byte) (int, error) {
	n, err := t.wd.Write(p)
	if err != nil {
		return n, err
	}

	t.offset += int64(n)

	return n, nil
}

// TrackedOffsetWriter implementation
func (t *writer) Offset() int64 {
	return t.offset
}

//-----------------------------------------------------------------------------

// Wrap os.File to keep track of the current offset without needing to make constant calls to Seek which involves syscall Lseek
type fileTrackedOffset struct {
	f      *os.File
	offset int64
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

// Reader implementation
func (t *fileTrackedOffset) Read(p []byte) (int, error) {
	n, err := t.f.Read(p)
	if err != nil {
		return n, err
	}

	t.offset += int64(n)

	return n, nil
}

// Writer implementation
func (t *fileTrackedOffset) Write(p []byte) (int, error) {
	n, err := t.f.Write(p)
	if err != nil {
		return n, err
	}

	t.offset += int64(n)

	return n, nil
}

// Seeker implementation
func (t *fileTrackedOffset) Seek(offset int64, whence int) (int64, error) {
	newOffset, err := t.f.Seek(offset, whence)
	if err != nil {
		return newOffset, err
	}
	t.offset = newOffset
	return newOffset, err
}

// ReaderAt implementation
func (t *fileTrackedOffset) ReadAt(p []byte, off int64) (int, error) {
	n, err := t.f.ReadAt(p, off)
	if err != nil {
		return n, err
	}

	t.offset = off + int64(n)

	return n, nil
}

// WriterAt implementation
func (t *fileTrackedOffset) WriteAt(p []byte, off int64) (int, error) {
	n, err := t.f.WriteAt(p, off)
	if err != nil {
		return n, err
	}

	t.offset = off + int64(n)

	return n, nil
}

//-----------------------------------------------------------------------------

// TrackedOffset implementation
func (t *fileTrackedOffset) Offset() int64 {
	return t.offset
}

// Ensure the tracker's offset and the file's actual offsets are the same.
// This will make a call to file.Seek
func (t *fileTrackedOffset) SyncOffset() error {
	offset, err := t.f.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}

	t.offset = offset
	return nil
}

//-----------------------------------------------------------------------------

type wrappedMBReader struct {
	rd     MultiByteReader
	offset int64
}

// Create a new bufio.Reader that also supports being able to do io.Seeker
func NewTrackedOffsetReaderMultiByte(rd MultiByteReader, baseOffset int64) MultiByteTrackedOffsetReader {
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

	w.offset += int64(n)

	return n, nil
}

func (w *wrappedMBReader) ReadByte() (byte, error) {
	b, err := w.rd.ReadByte()
	if err != nil {
		return b, err
	}
	w.offset++
	return b, nil
}

func (w *wrappedMBReader) Offset() int64 {
	return w.offset
}
