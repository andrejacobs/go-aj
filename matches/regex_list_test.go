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

func TestRegexListMatchesAny(t *testing.T) {
	l, err := matches.NewRegexList([]string{`\bHe`, `\bworld\b`, `\d+`})
	require.NoError(t, err)

	assert.True(t, l.MatchesAny("\t\tHello\t\n"))
	assert.True(t, l.MatchesAny(" the world!\t\n"))
	assert.True(t, l.MatchesAny(" The 9 tailed fox was 42 "))
	assert.False(t, l.MatchesAny("The quick brown fox"))
}

func TestRegexListMatchesAll(t *testing.T) {
	l, err := matches.NewRegexList([]string{`\bHe`, `\d{1,2}`})
	require.NoError(t, err)

	assert.True(t, l.MatchesAll("\t Hello42\t"))
	assert.True(t, l.MatchesAll("\t He is 42\t"))
	assert.True(t, l.MatchesAll("\t He said hello to 4242\t"))
	assert.False(t, l.MatchesAll("\t the answer is 42\t"))
	assert.False(t, l.MatchesAll("\t Hello \t"))
}

func TestRegexListMatches(t *testing.T) {
	l, err := matches.NewRegexList([]string{`\bHe`, `\d{1,2}`})
	require.NoError(t, err)

	items := []string{"No matches", "\tHello\t", "The year 2042"}
	found := l.Matches(items)
	assert.Contains(t, found, items[1])
	assert.Contains(t, found, items[2])
	assert.NotContains(t, found, items[0])
}

func TestRegexListCompileError(t *testing.T) {
	_, err := matches.NewRegexList([]string{`\bHe`, `\d{1,2}`, `\Knotvalid`})
	require.Error(t, err)

	assert.IsType(t, &matches.RegexListCompileErr{}, err)
	compErr, ok := err.(*matches.RegexListCompileErr)
	require.True(t, ok)
	assert.Equal(t, 2, compErr.Index)
	assert.Equal(t, `\Knotvalid`, compErr.Input)
}
