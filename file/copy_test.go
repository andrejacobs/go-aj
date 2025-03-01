package file_test

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/andrejacobs/go-aj/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopyFile(t *testing.T) {
	expected := "The quick brown fox jumped over the lazy dog!"
	src, err := os.CreateTemp("", "unit-test-source")
	require.NoError(t, err)
	defer os.Remove(src.Name())
	_, err = src.WriteString(expected)
	require.NoError(t, err)
	require.NoError(t, src.Close())

	destPath := filepath.Join(os.TempDir(), "unit-test-dest")
	defer os.Remove(destPath)
	wc, err := file.CopyFile(context.Background(), src.Name(), destPath)
	require.NoError(t, err)
	assert.Equal(t, int64(len(expected)), wc)

	dest, err := os.Open(destPath)
	require.NoError(t, err)
	defer dest.Close()

	data, err := io.ReadAll(dest)
	require.NoError(t, err)

	assert.Equal(t, expected, string(data))
}

func TestCopyFileN(t *testing.T) {
	expected := "The quick brown fox jumped over the lazy dog!"
	src, err := os.CreateTemp("", "unit-test-source")
	require.NoError(t, err)
	defer os.Remove(src.Name())
	_, err = src.WriteString(expected)
	require.NoError(t, err)
	require.NoError(t, src.Close())

	destPath := filepath.Join(os.TempDir(), "unit-test-dest")
	defer os.Remove(destPath)
	wc, err := file.CopyFileN(context.Background(), src.Name(), destPath, 9)
	require.NoError(t, err)
	assert.Equal(t, int64(9), wc)

	dest, err := os.Open(destPath)
	require.NoError(t, err)
	defer dest.Close()

	data, err := io.ReadAll(dest)
	require.NoError(t, err)

	assert.Equal(t, "The quick", string(data))
}
