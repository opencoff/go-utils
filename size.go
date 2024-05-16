// size.go -- Parse strings with a size suffix
//
// (c) 2016 Sudhi Herle <sudhi@herle.net>
//
// Licensing Terms: GPLv2
//
// If you need a commercial license for this work, please contact
// the author.
//
// This software does not come with any express or implied
// warranty; it is provided "as is". No claim  is made to its
// suitability for any purpose.
package utils

import (
	"fmt"
	"math/bits"
	"strconv"
	"strings"
)

const (
	_kB uint64 = 1 << 10
	_MB uint64 = 1 << 20
	_GB uint64 = 1 << 30
	_TB uint64 = 1 << 40
	_PB uint64 = 1 << 50
	_EB uint64 = 1 << 60
)

var multmap = map[string]uint64{
	"":  1,
	"B": 1,
	"k": _kB,
	"K": _kB,
	"M": _MB,
	"G": _GB,
	"T": _TB,
	"P": _PB,
	"E": _EB,
}

var orderedSizes = [...]struct {
	mult uint64
	suff string
}{
	{_EB, "EB"},
	{_PB, "PB"},
	{_TB, "TB"},
	{_GB, "GB"},
	{_MB, "MB"},
	{_kB, "kB"},
}

// Parse a string that has a size suffix (one of k, M, G, T, P, E).
// The suffix denotes multiples of 1024.
// e.g., "32k", "2M"
func ParseSize(s string) (uint64, error) {
	var suff string

	i := strings.LastIndexAny(s, "BkKMGTPE")
	if i > 0 {
		suff = s[i : i+1]
		s = s[:i]
	}

	if s == "" {
		return 0, nil
	}

	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}

	if m, ok := multmap[suff]; ok {
		// guard against overflow
		hi, lo := bits.Mul64(v, m)
		if hi > 0 {
			return 0, fmt.Errorf("%s%s size too large (overflow)", s, suff)
		}
		v = lo
	} else {
		return 0, fmt.Errorf("unknown size suffix %s", suff)
	}

	return v, nil
}

// Humanize turns a given uint64 into human readable string with a unit
// suffix of k, M, T, etc.
func HumanizeSize(sz uint64) string {
	for i := range orderedSizes {
		v := &orderedSizes[i]
		if sz < v.mult {
			continue
		}

		m := v.mult
		if b := sz % m; b > 0 {
			f := float64(sz) / float64(m)
			return fmt.Sprintf("%.02f %s", f, v.suff)
		}
		return fmt.Sprintf("%d %s", sz/m, v.suff)
	}

	return fmt.Sprintf("%d B", sz)
}
