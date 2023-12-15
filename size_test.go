package utils

import (
	"testing"
)

type sizeTest struct {
	in  string
	out uint64
	err bool
}

var sizesTests = []sizeTest{
	{"", 0, false},
	{"10", 10, false},
	{"4k", 4096, false},
	{"10M", 10 * 1048576, false},
	{"80G", 80 * _GB, false},
	{"10T", 10 * _TB, false},

	{"1048576E", 0, true},
	{"4x", 0, true},
	{"boo", 0, true},
}

func TestSize(t *testing.T) {
	assert := newAsserter(t)

	for i, t := range sizesTests {
		v, err := ParseSize(t.in)
		if t.err {
			assert(err != nil, "%2d: %s: expected to fail", i, t.in)
		} else {
			assert(err == nil, "%2d: %s: unexpected err: %s", i, t.in, err)
			assert(t.out == v, "%2d: %s: exp %v, saw %v", i, t.in, t.out, v)
		}
	}
}

var humanTests = []struct {
	in  uint64
	out string
}{
	{1024, "1 kB"},
	{4659472483, "4.34 GB"},
	{3180811731, "2.96 GB"},
	{541321380, "516.24 MB"},
	{1724652, "1.64 MB"},
	{1426449, "1.36 MB"},
	{1048576, "1 MB"},
	{1024137, "1000.13 kB"},
	{1023497, "999.51 kB"},
	{1048, "1.02 kB"},
	{1000, "1000 B"},
}

func TestHumanize(t *testing.T) {
	assert := newAsserter(t)

	for i := range humanTests {
		h := &humanTests[i]
		s := HumanizeSize(h.in)
		assert(s == h.out, "%2d: %d: exp %s, saw %s", i, h.in, h.out, s)
	}
}
