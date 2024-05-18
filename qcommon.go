// qcommon.go - shared functions for the various q types
//
// (c) 2014 Sudhi Herle <sw-at-herle.net>
//
// Placed in the Public Domain
// This software does not come with any express or implied
// warranty; it is provided "as is". No claim  is made to its
// suitability for any purpose.

package utils

import (
	"fmt"
	"math/bits"
)

func nextpow2[T ~uint | ~uint16 | ~uint32 | ~uint64](z T) T {
	i := bits.LeadingZeros64(uint64(z))
	v := uint64(1) << (64 - i)
	return T(v)
}

func qlen(rd, wr, mask uint64) int {
	switch {
	case rd == wr:
		return 0
	case rd < wr:
		return int(wr - rd)
	default: // wrapped around
		return int((mask + 1) - rd + wr)
	}
}

func qempty(rd, wr, mask uint64) bool {
	return rd == wr
}

func qfull(rd, wr, mask uint64) bool {
	return rd == ((1 + wr) & mask)
}

func qrepr(rd, wr, mask uint64) string {
	full := qfull(rd, wr, mask)
	empty := qempty(rd, wr, mask)
	n := qlen(rd, wr, mask)

	var p string = ""
	if full {
		p = "[FULL] "
	} else if empty {
		p = "[EMPTY] "
	}

	return fmt.Sprintf("%scap=%d len=%d wr=%d rd=%d",
		p, mask, n, wr, rd)
}
