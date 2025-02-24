package random

import (
	"encoding/binary"
	"math/rand"
	"strings"
	"time"

	crand "crypto/rand"
)

// -----------------------------------------------------------------------------
// Amazing! Someone went through a number of implementations and benchmarked it
// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
// I am using the RandStringBytesMaskImprSrcSB version
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

// String produces a string of length n that contains random characters.
// Characters are chosen from the following set: abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ
func String(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}

//-----------------------------------------------------------------------------

// Int returns a random integer between the minimum and maximum.
func Int(min int, max int) int {
	return rand.Intn(max-min+1) + min
}

// Read 4 bytes from the secure random number generator and convert it to an uint32
func SecureUint32() (uint32, error) {
	var result uint32
	err := binary.Read(crand.Reader, binary.LittleEndian, &result)
	if err != nil {
		return 0, err
	}
	return result, nil
}
