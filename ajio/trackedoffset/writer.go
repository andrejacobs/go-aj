package trackedoffset

import (
	"io"

	"github.com/andrejacobs/go-aj/ajmath"
)

// Writer keeps track of the offset within an io.Writer source.
type Writer struct {
	wd     io.Writer
	offset uint64
}

// Create a new Writer that will keep track of the offset within the source io.Writer.
func NewWriter(wd io.Writer, baseOffset uint64) *Writer {
	w := &Writer{
		wd:     wd,
		offset: baseOffset,
	}
	return w
}

// Writer implementation.
func (w *Writer) Write(p []byte) (int, error) {
	n, err := w.wd.Write(p)
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

// Return the current offset in bytes.
func (w *Writer) Offset() uint64 {
	return w.offset
}

// Set the known offset in bytes.
func (w *Writer) ResetOffset(offset uint64) {
	w.offset = offset
}
