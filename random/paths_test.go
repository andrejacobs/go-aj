package random_test

import (
	"os"
	"strings"
	"testing"

	"github.com/andrejacobs/go-aj/file"
	"github.com/andrejacobs/go-aj/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPath(t *testing.T) {
	expectedPrefix := "dir1"

	assert.True(t, strings.HasPrefix(random.Path(expectedPrefix, 4, 10, 1, 10), expectedPrefix))

	parts := strings.Split(random.Path(expectedPrefix, 1, 1, 2, 8), string(os.PathSeparator))
	assert.Equal(t, len(parts), 2)
	assert.Equal(t, parts[0], expectedPrefix)

	parts = strings.Split(random.Path(expectedPrefix, 3, 3, 0, 20), string(os.PathSeparator))
	assert.Equal(t, len(parts), 4)
	assert.Equal(t, parts[0], expectedPrefix)

	// minNameLen = 0 (should at least use 1 character for name)
	parts = strings.Split(random.Path(expectedPrefix, 1, 1, 0, 2), string(os.PathSeparator))
	assert.Equal(t, len(parts), 2)
	assert.Equal(t, parts[0], expectedPrefix)
	assert.True(t, len(parts[1]) > 0)

	// minNameLen = maxNameLen = 4
	parts = strings.Split(random.Path(expectedPrefix, 1, 1, 4, 4), string(os.PathSeparator))
	assert.Equal(t, len(parts), 2)
	assert.Equal(t, parts[0], expectedPrefix)
	assert.Len(t, parts[1], 4)
}

func TestPaths(t *testing.T) {
	expectedPrefix := "dir1"
	expectedCount := 10

	paths := random.Paths(expectedPrefix, expectedCount, 4, 8, 1, 10)
	assert.Equal(t, len(paths), expectedCount)

	for _, v := range paths {
		assert.True(t, strings.HasPrefix(v, expectedPrefix))
	}
}

func TestCreateFiles(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "unit-testing")
	defer os.RemoveAll(tempDir)
	require.NoError(t, err)

	minFiles := 4
	maxFiles := 10
	minSize := uint64(4)
	maxSize := uint64(20)
	maxTotalSize := uint64(100)
	wc, err := random.CreateFiles(tempDir, minFiles, maxFiles, minSize, maxSize, maxTotalSize)
	require.NoError(t, err)
	assert.LessOrEqual(t, wc, maxTotalSize)

	totalSize, _, err := file.CalculateDirSizeShallow(tempDir)
	require.NoError(t, err)
	assert.LessOrEqual(t, uint64(totalSize), maxTotalSize)
}
