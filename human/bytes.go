package human

// I want to minimize 3rd party dependencies on my own go-micropkg packages
// and I use the following functions from Dustin Sallings great go-humanize package.
// So I decided to copy the functions from https://github.com/dustin/go-humanize/blob/master/bytes.go

import (
	"fmt"
	"math"
)

// Bytes produces a human readable representation of an SI size.
// Taken from: https://github.com/dustin/go-humanize/blob/master/bytes.go
func Bytes(s uint64) string {
	sizes := []string{"B", "kB", "MB", "GB", "TB", "PB", "EB"}
	return humanateBytes(s, 1000, sizes)
}

func humanateBytes(s uint64, base float64, sizes []string) string {
	if s < 10 {
		return fmt.Sprintf("%d B", s)
	}
	e := math.Floor(logn(float64(s), base))
	suffix := sizes[int(e)]
	val := math.Floor(float64(s)/math.Pow(base, e)*10+0.5) / 10
	f := "%.0f %s"
	if val < 10 {
		f = "%.1f %s"
	}

	return fmt.Sprintf(f, val, suffix)
}

func logn(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}
