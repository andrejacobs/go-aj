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

package concurrency

import (
	"context"
)

// Consume from the 'in' channel and produce the same value to all of the output channels.
func Fanout[T any](ctx context.Context, in <-chan T, outs ...chan T) {
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case data, ok := <-in:
			if !ok {
				break loop
			}
			for _, out := range outs {
				out <- data
			}
		}
	}

	for _, out := range outs {
		close(out)
	}
}

// Consume from the 'in' channel and produce the a transformed value to the output channels.
// Meaning consume T and produce V.
func TransformedFanout[T any, V any](ctx context.Context,
	transformer func(in T) V,
	in <-chan T, outs ...chan V) {
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case data, ok := <-in:
			if !ok {
				break loop
			}
			for _, out := range outs {
				out <- transformer(data)
			}
		}
	}

	for _, out := range outs {
		close(out)
	}
}
