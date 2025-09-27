package trackedoffset

import (
	"io"

	"github.com/andrejacobs/go-aj/ajmath"
)

// Reader keeps track of the offset within an io.Reader source.
type Reader struct {
	rd     io.Reader
	offset uint64
}

// Create a new Reader that will keep track of the offset within the source io.Reader.
// baseOffset is the known starting offset.
func NewReader(rd io.Reader, baseOffset uint64) *Reader {
	r := &Reader{
		rd:     rd,
		offset: baseOffset,
	}
	return r
}

// Reader implementation.
func (r *Reader) Read(p []byte) (int, error) {
	n, err := r.rd.Read(p)
	if err != nil {
		return n, err
	}

	newOffset, err := ajmath.Add64(r.offset, uint64(n))
	if err != nil {
		return 0, err
	}
	r.offset = newOffset

	return n, nil
}

// Return the current offset in bytes.
func (r *Reader) Offset() uint64 {
	return r.offset
}

// Set the known offset in bytes.
func (r *Reader) ResetOffset(offset uint64) {
	r.offset = offset
}
