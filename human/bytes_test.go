package human_test

import (
	"testing"

	"github.com/andrejacobs/go-micropkg/human"
)

func TestBytes(t *testing.T) {
	testCases := []struct {
		desc string
		size uint64
		exp  string
	}{
		{desc: "Bytes(0)", size: 0, exp: "0 B"},
		{desc: "Bytes(1)", size: 1, exp: "1 B"},
		{desc: "Bytes(803)", size: 803, exp: "803 B"},
		{desc: "Bytes(999)", size: 999, exp: "999 B"},
		{desc: "Bytes(1011)", size: 1011, exp: "1.0 kB"},
		{desc: "Bytes(1024)", size: 1024, exp: "1.0 kB"},
		{desc: "Bytes(9999)", size: 9999, exp: "10 kB"},
		{desc: "Bytes(1M - 1)", size: MByte - 1, exp: "1000 kB"},
		{desc: "Bytes(1M)", size: 1024 * 1024, exp: "1.0 MB"},
		{desc: "Bytes(1GB - 1K)", size: GByte - KByte, exp: "1000 MB"},
		{desc: "Bytes(1GB)", size: GByte, exp: "1.0 GB"},
		{desc: "Bytes(1TB - 1M)", size: TByte - MByte, exp: "1000 GB"},
		{desc: "Bytes(10MB)", size: 9999 * 1000, exp: "10 MB"},
		{desc: "Bytes(1TB)", size: TByte, exp: "1.0 TB"},
		{desc: "Bytes(1PB - 1T)", size: PByte - TByte, exp: "999 TB"},
		{desc: "Bytes(1PB)", size: PByte, exp: "1.0 PB"},
		{desc: "Bytes(1EB - 1T)", size: EByte - PByte, exp: "999 PB"},
		{desc: "Bytes(1EB)", size: EByte, exp: "1.0 EB"},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			result := human.Bytes(tC.size)
			if result != tC.exp {
				t.Errorf("%v: expected '%v', but got '%v'", tC.desc, tC.exp, result)
			}
		})
	}
}

const (
	IByte = 1
	KByte = IByte * 1000
	MByte = KByte * 1000
	GByte = MByte * 1000
	TByte = GByte * 1000
	PByte = TByte * 1000
	EByte = PByte * 1000
)
