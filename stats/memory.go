package stats

import (
	"fmt"
	"io"
	"runtime"

	"github.com/andrejacobs/go-micropkg/human"
)

// Get the memory usage stats
func GetMemoryUsage() runtime.MemStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m
}

// Print out the memory usage
func PrintMemoryUsage(w io.Writer, msg string, m runtime.MemStats) {
	// Taken from: https://golangcode.com/print-the-current-memory-usage/
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	if msg != "" {
		fmt.Fprint(w, msg)
	}
	fmt.Fprintf(w, "Alloc = %s", human.Bytes(m.Alloc))
	fmt.Fprintf(w, "\tTotalAlloc = %s", human.Bytes(m.TotalAlloc))
	fmt.Fprintf(w, "\tSys = %s", human.Bytes(m.Sys))
	fmt.Fprintf(w, "\tNumGC = %v\n", m.NumGC)
}
