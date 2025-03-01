package stats

import (
	"fmt"
	"io"
	"time"
)

// Used to time how long a function takes to execute and print to the writer.
// Example usage: defer stats.MeasureTimeTaken(os.Stdout, "reading entire database and building a cache", time.Now()).
func MeasureElapsedTime(w io.Writer, name string, start time.Time) {
	elapsed := time.Since(start)
	fmt.Fprintf(w, "%s took: %s\n", name, elapsed)
}

// Print the elapsed time between two times.
// Use this version when defer is not possible.
func PrintTimeTaken(w io.Writer, name string, start time.Time, end time.Time) {
	elapsed := end.Sub(start)
	fmt.Fprintf(w, "%s took: %s\n", name, elapsed)
}
