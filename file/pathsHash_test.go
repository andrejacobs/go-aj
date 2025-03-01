package file_test

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"testing"

	"github.com/andrejacobs/go-aj/file"
	"github.com/andrejacobs/go-aj/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculatePathHashConsistently(t *testing.T) {
	path := "/var/lib/ajfs"
	expected := "4e04b4b5415e5bef7e6c12736bb8b76f2ccb2751"
	sum := file.CalculatePathHash(path)
	require.Equal(t, expected, fmt.Sprintf("%x", sum))
}

func TestCalculatePathsHash(t *testing.T) {

	h1, err := file.CalculatePathsHash([]string{"/var", "/etc"})
	require.NoError(t, err)
	assert.NotEmpty(t, h1)

	h2, err := file.CalculatePathsHash([]string{"/etc", "/var"})
	require.NoError(t, err)
	assert.Equal(t, h1, h2)

	h3, err := file.CalculatePathsHash([]string{"/var", "/etc/aj"})
	require.NoError(t, err)
	assert.NotEqual(t, h1, h3)

	h4, err := file.CalculatePathsHash([]string{"/VAR", "/ETC"})
	require.NoError(t, err)
	assert.NotEqual(t, h1, h4)
}

func TestCalculatePathsHashConsistently(t *testing.T) {
	path := "/var/lib/ajfdb"
	expected := "397fb319d489c79c942221a055f298d06c24e95b"
	sum1, err := file.CalculatePathsHash([]string{path})
	require.NoError(t, err)
	require.Equal(t, expected, fmt.Sprintf("%x", sum1))
}

//-----------------------------------------------------------------------------

// Benchmark various hashing algorithms given a path
func BenchmarkHashingPaths(b *testing.B) {
	paths := random.Paths("/", 1000, 2, 100, 8, 16)

	// The result are quite interesting.
	// On the Intel machine (my Linux server) SHA1 is the fastest and SHA256 is the slowest
	// On my M2 Mac SHA256 is the fastest followed by SHA1
	// OpenSSL's benchmarks confirm the same: openssl speed md5 sha1 sha256

	b.Run("md5", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, path := range paths {
				md5.Sum([]byte(path))
			}
		}
	})

	b.Run("sha1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, path := range paths {
				sha1.Sum([]byte(path))
			}
		}
	})

	b.Run("sha256", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, path := range paths {
				sha256.Sum256([]byte(path))
			}
		}
	})

	b.Run("sha384", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, path := range paths {
				sha512.Sum384([]byte(path))
			}
		}
	})

	b.Run("sha512", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, path := range paths {
				sha512.Sum512([]byte(path))
			}
		}
	})
}
