package utils

import (
	"testing"
)

func TestSimple(t *testing.T) {
	assert := newAsserter(t)

	words := []string{"hello", "help", "sync"}
	ret := map[string]string{
		"hello": "hello",
		"hell":  "hello",
		"help":  "help",
		"sync":  "sync",
		"syn":   "sync",
		"sy":    "sync",
		"s":     "sync",
	}

	ab := Abbrev(words)

	for k, v := range ab {
		x, ok := ret[k]
		assert(ok, "unexpected abbrev %s", k)
		assert(x == v, "abbrev %s: %s != %s", k, x, v)
	}

	for k, v := range ret {
		x, ok := ab[k]
		assert(ok, "unknown abbrev %s", k)
		assert(x == v, "abbrev %s: %s != %s", k, x, v)
	}

}
