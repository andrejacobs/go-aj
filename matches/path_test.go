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
