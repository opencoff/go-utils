package utils

import (
	"testing"
)

func TestSimple(t *testing.T) {
	assert := newAsserter(t)

	words := []string{"hello", "help", "sync", "uint", "uint16", "uint64"}
	ret := map[string]string{
		"hello": "hello",
		"hell":  "hello",
		"help":  "help",
		"sync":  "sync",
		"syn":   "sync",
		"sy":    "sync",
		"s":     "sync",
		"uint":  "uint",
		"uint16":  "uint16",
		"uint1":  "uint16",
		"uint64": "uint64",
		"uint6": "uint64",
	}

	ab := Abbrev(words)

	for k, v := range ab {
		x, ok := ret[k]
		assert(ok, "unexpected abbrev %s", k)
		assert(x == v, "abbrev %s: exp %s, saw %s", k, x, v)
	}

	for k, v := range ret {
		x, ok := ab[k]
		assert(ok, "unknown abbrev %s", k)
		assert(x == v, "abbrev %s: exp %s, saw %s", k, x, v)
	}

}
