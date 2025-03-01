package file_test

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/andrejacobs/go-aj/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMD5(t *testing.T) {
	tempFile, err := makeHashFile()
	require.NoError(t, err)
	defer os.Remove(tempFile)

	hash, _, err := file.HashMD5(context.Background(), tempFile, nil)
	require.NoError(t, err)

	assert.Equal(t, expectedMD5, fmt.Sprintf("%x", hash))
}

func TestSHA1(t *testing.T) {
	tempFile, err := makeHashFile()
	require.NoError(t, err)
	defer os.Remove(tempFile)

	hash, _, err := file.HashSHA1(context.Background(), tempFile, nil)
	require.NoError(t, err)

	assert.Equal(t, expectedSHA1, fmt.Sprintf("%x", hash))
}

func TestSHA256(t *testing.T) {
	tempFile, err := makeHashFile()
	require.NoError(t, err)
	defer os.Remove(tempFile)

	hash, _, err := file.HashSHA256(context.Background(), tempFile, nil)
	require.NoError(t, err)

	assert.Equal(t, expectedSHA256, fmt.Sprintf("%x", hash))
}

func TestSHA512(t *testing.T) {
	tempFile, err := makeHashFile()
	require.NoError(t, err)
	defer os.Remove(tempFile)

	hash, _, err := file.HashSHA512(context.Background(), tempFile, nil)
	require.NoError(t, err)

	assert.Equal(t, expectedSHA512, fmt.Sprintf("%x", hash))
}

func TestCancel(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()

	_, _, err := file.HashFromReader(ctx, rand.Reader, md5.New(), nil)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestOptionalWriter(t *testing.T) {
	expected := "The quick brown fox jumped over the lazy dog!"
	rd := strings.NewReader(expected)
	w := bytes.Buffer{}
	_, wcount, err := file.HashFromReader(context.Background(), rd, sha256.New(), &w)
	require.NoError(t, err)
	assert.Equal(t, uint64(len(expected)), wcount)

	result := make([]byte, len(expected))
	_, err = io.ReadFull(&w, result)
	require.NoError(t, err)
	assert.Equal(t, expected, string(result))
}

//-----------------------------------------------------------------------------

func makeHashFile() (string, error) {
	f, err := os.CreateTemp("", "unit-test-hashfile")
	if err != nil {
		return "", err
	}
	defer f.Close()

	f.WriteString("The quick brown fox jumped over the lazy dog!")

	return f.Name(), nil
}

const (
	expectedMD5    = "efc05c070367008abb4388b189ac2b1e"
	expectedSHA1   = "98f77361955d663bf42d2c828c929d507edc3613"
	expectedSHA256 = "2d2a94f4aebb2aaa87da022f344b14ed4d49843838cd2511b42065b6d661564f"
	expectedSHA512 = "cf12f85cfeada9999644c8f73e3b258a44d363506eea7f105e7c93304f80abdd51d3c5107b799a3bd149683f588ff3948be8e5bc697d40e6437785a69ab896dc"
)
