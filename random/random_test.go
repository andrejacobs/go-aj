package random_test

import (
	"testing"

	"github.com/andrejacobs/go-micropkg/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt(t *testing.T) {
	for i := 0; i < 100; i++ {
		x := random.Int(10, 42)
		assert.GreaterOrEqual(t, x, 10)
		assert.LessOrEqual(t, x, 42)
	}
}

func TestString(t *testing.T) {
	for i := 0; i < 100; i++ {
		s := random.String(i)
		assert.Len(t, s, i)
	}
}

func TestSecureUint32(t *testing.T) {

	seen := make(map[uint32]struct{})

	for i := 0; i < 100; i++ {
		r, err := random.SecureUint32()
		require.NoError(t, err)

		_, exists := seen[r]
		assert.False(t, exists)
		seen[r] = struct{}{}
	}

}
