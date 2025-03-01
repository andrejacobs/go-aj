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

package matches_test

import (
	"testing"

	"github.com/andrejacobs/go-aj/matches"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegexPathMatcher(t *testing.T) {
	r, err := matches.NewRegexPathMatcher([]string{"foo.txt", "^/proc", "\\.DS_Store"})
	require.NoError(t, err)

	m, err := r.Match("/dev/foo.txt")
	require.NoError(t, err)
	assert.True(t, m)

	m, err = r.Match("/dev/foo.txt/wtf")
	require.NoError(t, err)
	assert.True(t, m)

	m, err = r.Match("/dev/bar.txt")
	require.NoError(t, err)
	assert.False(t, m)

	m, err = r.Match("/proc/something")
	require.NoError(t, err)
	assert.True(t, m)

	m, err = r.Match("/dev/proc/something")
	require.NoError(t, err)
	assert.False(t, m)

	m, err = r.Match("abc/.DS_Store")
	require.NoError(t, err)
	assert.True(t, m)
}

func TestShellPatternPathMatcher(t *testing.T) {
	s := matches.NewShellPatternPathMatcher([]string{"foo*", "/proc/*", "bar.?", "*/bar.txt"})

	m, err := s.Match("foo.txt")
	require.NoError(t, err)
	assert.True(t, m)

	m, err = s.Match("/bar/foo.txt")
	require.NoError(t, err)
	assert.False(t, m)

	m, err = s.Match("/proc/something")
	require.NoError(t, err)
	assert.True(t, m)

	m, err = s.Match("bar.1")
	require.NoError(t, err)
	assert.True(t, m)

	m, err = s.Match("bar.23")
	require.NoError(t, err)
	assert.False(t, m)

	m, err = s.Match("abc/bar.txt")
	require.NoError(t, err)
	assert.True(t, m)
}
