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
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/andrejacobs/go-aj/matches"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegexScannerFailToCompile(t *testing.T) {
	r := &matches.RegexScanner{}
	assert.Error(t, r.Add("fail", "a(b", nil))
}

func TestRegexScanner(t *testing.T) {
	input := `The quick
brown fox
jumped over
the lazy
dog!

alpha: 42
bravo 007 delta
charlie
bravo 7 delta
echo
`
	r := &matches.RegexScanner{}
	r.Add("one", "\\bquick\\b", nil)
	r.Add("two", "fox$", nil)

	var jumpedLine string
	var jumpedLineNumber int
	r.Add("three", "^jumped\\b", func(key string, line string, lineNumber int, matches []string) error {
		jumpedLine = line
		jumpedLineNumber = lineNumber
		return nil
	})

	r.Add("four", "ox|ov", func(key string, line string, lineNumber int, matches []string) error {
		assert.Len(t, matches, 1)
		if matches[0] == "ox" || matches[0] == "ov" {
			return nil
		}
		assert.Fail(t, "should not have matched: "+line)
		return fmt.Errorf("test failed")
	})

	r.Add("five", "(?i)DOG", nil)
	r.Add("no-match", "zebra", nil)

	r.Add("capture", "bravo\\s+(\\d+)\\s+delta", nil)

	result, err := r.Process(strings.NewReader(input))
	require.NoError(t, err)
	assert.Len(t, result, 6)
	assert.Equal(t, "jumped over", jumpedLine)
	assert.Equal(t, 2, jumpedLineNumber)
	assert.Equal(t, []string{"dog"}, result["five"])

	_, exists := result["no-match"]
	assert.False(t, exists)

	assert.Len(t, result["capture"], 2)
	assert.Equal(t, "bravo 7 delta", result["capture"][0])
	assert.Equal(t, "7", result["capture"][1])
}

func TestRegexScannerWriteToOut(t *testing.T) {
	input := `The quick brown
	fox jumped`

	r := &matches.RegexScanner{}
	r.Add("one", "\\bquick\\b", nil)

	buf := bytes.Buffer{}

	r.SetOut(&buf)
	_, err := r.Process(strings.NewReader(input))
	require.NoError(t, err)

	assert.Equal(t, input+"\n", buf.String())
}
