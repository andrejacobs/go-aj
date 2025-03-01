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
	"bufio"
	"io"
)

// Need to be able to do io.Reader.Read but also io.ByteReader.ReadByte and guess
// what os.File doesn't support reading just a single byte.
// bufio.NewReader implements this interface.
type MultiByteReader interface {
	// MultiByte for the lack of a better name
	io.Reader
	io.ByteReader
}

// Need to do MultiByteReader but also support seeking.
type MultiByteReaderSeeker interface {
	MultiByteReader
	io.Seeker
}

//-----------------------------------------------------------------------------

type wrappedBufIOReadSeeker struct {
	rs io.ReadSeeker
	br *bufio.Reader
}

// Create a new bufio.Reader that also supports being able to do io.Seeker.
func NewMultiByteReaderSeeker(rs io.ReadSeeker) MultiByteReaderSeeker {
	return &wrappedBufIOReadSeeker{
		rs: rs,
		br: bufio.NewReader(rs),
	}
}

// Create a new bufio.Reader with size, that also supports being able to do io.Seeker.
func NewMultiByteReaderSeekerSize(rs io.ReadSeeker, size int) MultiByteReaderSeeker {
	return &wrappedBufIOReadSeeker{
		rs: rs,
		br: bufio.NewReaderSize(rs, size),
	}
}

func (w *wrappedBufIOReadSeeker) Read(p []byte) (n int, err error) {
	return w.br.Read(p)
}

func (w *wrappedBufIOReadSeeker) ReadByte() (byte, error) {
	return w.br.ReadByte()
}

func (w *wrappedBufIOReadSeeker) Seek(offset int64, whence int) (int64, error) {
	newOffset, err := w.rs.Seek(offset, whence)
	if err != nil {
		return newOffset, err
	}
	w.br.Reset(w.rs)
	return newOffset, nil
}
