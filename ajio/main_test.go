package ajio_test

import (
	"crypto/rand"
	"io"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

//-----------------------------------------------------------------------------

func createTempFile(size int64) (string, error) {
	f, err := os.CreateTemp("", "unit-testing")
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = io.CopyN(f, rand.Reader, size)
	if err != nil {
		return "", err
	}

	return f.Name(), nil
}
