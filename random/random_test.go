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

package random_test

import (
	"testing"

	"github.com/andrejacobs/go-aj/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt(t *testing.T) {
	for i := 0; i < 100; i++ {
		x := random.Int(10, 42)
		assert.GreaterOrEqual(t, x, 10)
		assert.LessOrEqual(t, x, 42)
	}
}

func TestString(t *testing.T) {
	for i := 0; i < 100; i++ {
		s := random.String(i)
		assert.Len(t, s, i)
	}
}

func TestSecureUint32(t *testing.T) {

	seen := make(map[uint32]struct{})

	for i := 0; i < 100; i++ {
		r, err := random.SecureUint32()
		require.NoError(t, err)

		_, exists := seen[r]
		assert.False(t, exists)
		seen[r] = struct{}{}
	}

}
