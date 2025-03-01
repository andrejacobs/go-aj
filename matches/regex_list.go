package matches

import (
	"fmt"
	"regexp"
)

// A list of compiled regular expressions that can be used to match things.
type RegexList struct {
	compiled []*regexp.Regexp
}

// Create a new RegexList that compiles the given regular expressions.
func NewRegexList(expressions []string) (*RegexList, error) {
	l := &RegexList{}
	err := l.compile(expressions)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func (l *RegexList) compile(expressions []string) error {
	for i, exp := range expressions {
		r, err := regexp.Compile(exp)
		if err != nil {
			return &RegexListCompileErr{
				Input: exp,
				Index: i,
				Err:   err,
			}
		}
		l.compiled = append(l.compiled, r)
	}

	return nil
}

// Returns true if the needle matches any of the compiled regular expressions.
func (l *RegexList) MatchesAny(needle string) bool {
	return matchesAnyRegexp(l.compiled, needle)
}

// Returns true if the needle matches all of the compiled regular expressions.
func (l *RegexList) MatchesAll(needle string) bool {
	return matchesAllRegexp(l.compiled, needle)
}

// Returns the slice of needles that matched any of the compiled regular expressions.
func (l *RegexList) Matches(needles []string) []string {
	return matchesRegexp(l.compiled, needles)
}

type RegexListCompileErr struct {
	Input string
	Index int
	Err   error
}

func (e *RegexListCompileErr) Error() string {
	return fmt.Sprintf("the regular expression at index [%d] %q is not valid. %v", e.Index, e.Input, e.Err)
}

//-----------------------------------------------------------------------------

func matchesAnyRegexp(expressions []*regexp.Regexp, needle string) bool {
	for _, re := range expressions {
		if re.MatchString(needle) {
			return true
		}
	}

	return false
}

func matchesAllRegexp(expressions []*regexp.Regexp, needle string) bool {
	for _, re := range expressions {
		if !re.MatchString(needle) {
			return false
		}
	}

	return true
}

func matchesRegexp(expressions []*regexp.Regexp, needles []string) []string {
	var result []string
	for _, needle := range needles {
		found := false
		for _, re := range expressions {
			if re.MatchString(needle) {
				found = true
				break
			}
		}
		if found {
			result = append(result, needle)
		}
	}

	return result
}
