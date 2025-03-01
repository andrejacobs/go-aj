package file_test

import (
	"testing"

	"github.com/andrejacobs/go-aj/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculateDirSizeShallow(t *testing.T) {
	size, _, err := file.CalculateDirSizeShallow(tempDir)
	require.NoError(t, err)
	assert.Equal(t, int64(60), size)
}

func TestCalculateSize(t *testing.T) {
	result, err := file.CalculateSize(tempDir)
	require.NoError(t, err)

	assert.Equal(t, 2, result.Dirs)
	assert.Equal(t, 5, result.Files)
	assert.Equal(t, uint64(90), result.TotalSize)
}
