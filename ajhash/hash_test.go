package ajhash_test

import (
	"crypto"
	"testing"

	"github.com/andrejacobs/go-aj/ajhash"
	"github.com/stretchr/testify/assert"
)

func TestHashAssumptions(t *testing.T) {
	assert.Equal(t, crypto.SHA1.Size(), ajhash.AlgoSHA1.Size())
	assert.Equal(t, crypto.SHA256.Size(), ajhash.AlgoSHA256.Size())
	assert.Equal(t, crypto.SHA512.Size(), ajhash.AlgoSHA512.Size())

	assert.Equal(t, ajhash.AlgoSHA256, ajhash.DefaultAlgo)

	assert.Equal(t, "SHA-1", ajhash.AlgoSHA1.String())
	assert.Equal(t, "SHA-256", ajhash.AlgoSHA256.String())
	assert.Equal(t, "SHA-512", ajhash.AlgoSHA512.String())

	// shasum -a 1 /dev/null
	assert.Equal(t, "da39a3ee5e6b4b0d3255bfef95601890afd80709", ajhash.AlgoSHA1.HashedStringForZeroBytes())
	// shasum -a 256 /dev/null
	assert.Equal(t, "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", ajhash.AlgoSHA256.HashedStringForZeroBytes())
	// shasum -a 512 /dev/null
	assert.Equal(t, "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e", ajhash.AlgoSHA512.HashedStringForZeroBytes())

	invalid := ajhash.Algo(42)
	assert.Equal(t, "unknown", invalid.String())
	assert.Panics(t, func() { invalid.Size() })
	assert.Equal(t, "", invalid.HashedStringForZeroBytes())
	assert.Panics(t, func() { invalid.Hasher() })
}

func TestAllZeroBytes(t *testing.T) {
	zeroes := make([]byte, 10)
	notZeroes := make([]byte, 10)
	notZeroes[7] = 0x41

	assert.True(t, ajhash.AllZeroBytes(zeroes))
	assert.False(t, ajhash.AllZeroBytes(notZeroes))
}
