package ajmath

import (
	"errors"
	"math/bits"
)

var (
	// An integer overflow occurred.
	ErrIntegerOverflow = errors.New("integer overflow occurred")

	// An integer underflow occurred.
	ErrIntegerUnderflow = errors.New("integer underflow occurred")
)

// Add two unsigned 32bit integers.
// Returns [ErrIntegerOverflow] if an overflow occurred.
func Add32(x, y uint32) (uint32, error) {
	new, overflowed := bits.Add32(x, y, 0)
	if overflowed > 0 {
		return 0, ErrIntegerOverflow
	}
	return new, nil
}

// Add two unsigned 64bit integers.
// Returns [ErrIntegerOverflow] if an overflow occurred.
func Add64(x, y uint64) (uint64, error) {
	new, overflowed := bits.Add64(x, y, 0)
	if overflowed > 0 {
		return 0, ErrIntegerOverflow
	}
	return new, nil
}

// Subtract two unsigned 32bit integers.
// Returns [ErrIntegerUnderflow] if an underflow occurred.
func Sub32(x, y uint32) (uint32, error) {
	new, underflowed := bits.Sub32(x, y, 0)
	if underflowed > 0 {
		return 0, ErrIntegerUnderflow
	}
	return new, nil
}

// Subtract two unsigned 64bit integers.
// Returns [ErrIntegerUnderflow] if an underflow occurred.
func Sub64(x, y uint64) (uint64, error) {
	new, underflowed := bits.Sub64(x, y, 0)
	if underflowed > 0 {
		return 0, ErrIntegerUnderflow
	}
	return new, nil
}
