package file_test

import (
	"os"
	"testing"

	"github.com/andrejacobs/go-micropkg/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsDirEntryEqual(t *testing.T) {
	a, err := os.ReadDir(tempDir)
	require.NoError(t, err)

	b, err := os.ReadDir(tempDir)
	require.NoError(t, err)

	assert.True(t, file.IsDirEntryEqual(a[0], b[0]))
}

func TestIsDirEntryWithInfoEqual(t *testing.T) {
	a, err := os.ReadDir(tempDir)
	require.NoError(t, err)

	b, err := os.ReadDir(tempDir)
	require.NoError(t, err)

	equal, err := file.IsDirEntryWithInfoEqual(a[0], b[0])
	require.NoError(t, err)
	assert.True(t, equal)
}

func TestReadDirUnsorted(t *testing.T) {
	a, err := os.ReadDir(tempDir)
	require.NoError(t, err)

	b, err := file.ReadDirUnsorted(tempDir)
	require.NoError(t, err)
	file.SortDirEntries(b)

	require.Equal(t, len(a), len(b))
	for i := 0; i < len(a); i++ {
		assert.True(t, file.IsDirEntryEqual(a[i], b[i]))
	}
}
