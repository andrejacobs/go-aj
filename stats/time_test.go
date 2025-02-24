package stats_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/andrejacobs/go-micropkg/stats"
	"github.com/stretchr/testify/assert"
)

func TestPrintTimeTaken(t *testing.T) {
	var buf bytes.Buffer
	stats.PrintTimeTaken(&buf, "test", time.Now(), time.Now().Add(time.Minute))
	assert.True(t, strings.HasPrefix(buf.String(), "test took: 1m0.0"))
}

func TestMeasureElapsedTime(t *testing.T) {
	var buf bytes.Buffer
	stats.MeasureElapsedTime(&buf, "test", time.Now().Add(-1*time.Minute))
	assert.True(t, strings.HasPrefix(buf.String(), "test took: 1m0.0"))
}
