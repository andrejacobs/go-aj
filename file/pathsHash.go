package file

import (
	"bytes"
	"crypto/sha1" // #nosec G505 -- SHA1 is not used for cryptography
	"sort"
)

const (
	PathHashSize = sha1.Size
)

type PathHash [PathHashSize]byte

// Calculate the unique hash for a path.
func CalculatePathHash(path string) PathHash {
	return sha1.Sum([]byte(path))
}

// Calculate the unique hash for a given slice of file paths.
func CalculatePathsHash(paths []string) (PathHash, error) {
	// Using sha1 since I need a hash that is consistent (maphash is great but requires to store the seed value)
	// sha1 turns out to be faster on the Intel CPU I intend to mainly run this code on
	// sha256 is slightly faster on my M2 Macbook
	// To test: openssl speed md5 sha1 sha256
	sorted := append([]string{}, paths...)
	sort.Strings(sorted)

	var buf bytes.Buffer
	for _, p := range sorted {
		if _, err := buf.WriteString(p); err != nil {
			return PathHash{}, err
		}
	}

	return sha1.Sum(buf.Bytes()), nil
}
