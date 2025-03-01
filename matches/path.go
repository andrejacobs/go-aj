package matches

import (
	"fmt"
	"path/filepath"
)

// PathMatcher is used to determine if a file system path matches
type PathMatcher interface {
	// Match checks if the path matches and returns true if it does
	Match(path string) (bool, error)
}

//-----------------------------------------------------------------------------
// RegexPathMatcher

// RegexPathMatcher will match a file system path against a set of regular expressions
type RegexPathMatcher struct {
	regexList *RegexList
}

// Create a new RegexPathMatcher using the regular expression patterns.
func NewRegexPathMatcher(expressions []string) (*RegexPathMatcher, error) {
	regexList, err := NewRegexList(expressions)
	if err != nil {
		return nil, fmt.Errorf("failed to create the RegexPathMatcher. %v", err)
	}

	matcher := RegexPathMatcher{
		regexList: regexList,
	}

	return &matcher, nil
}

func (r *RegexPathMatcher) Match(path string) (bool, error) {
	matched := r.regexList.MatchesAny(path)
	return matched, nil
}

//-----------------------------------------------------------------------------
// ShellPatternPathMatcher

// ShellPatternPathMatcher will match a file system path against a set of shell patterns.
// See https://pkg.go.dev/path/filepath#Match for details.
type ShellPatternPathMatcher struct {
	patterns []string
}

// Create a new ShellPatternPathMatcher using the shell patterns.
func NewShellPatternPathMatcher(patterns []string) *ShellPatternPathMatcher {
	matcher := ShellPatternPathMatcher{
		patterns: patterns,
	}
	return &matcher
}

func (s *ShellPatternPathMatcher) Match(path string) (bool, error) {
	matched := false

	for _, pattern := range s.patterns {
		var err error
		matched, err = filepath.Match(pattern, path)
		if err != nil {
			return matched, err
		}
		if matched {
			break
		}
	}

	return matched, nil
}
