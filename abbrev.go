// abbrev.go - generate abbreviations from a wordlist
//
// (c) 2014 Sudhi Herle <sw-at-herle.net>
//
// Placed in the Public Domain
// This software does not come with any express or implied
// warranty; it is provided "as is". No claim  is made to its
// suitability for any purpose.

package utils

// Given a wordlist in 'words', generate unique abbreviations of it
// and return as a map[abbrev]word.
// e.g.,
//  given a wordlist ["hello", "help", "sync"],
//  Abbrev() returns:
//    {
//       "hello": "hello",
//       "hell":  "hell"
//       "help":  "help",
//       "sync":  "sync",
//       "syn":   "sync",
//       "sy":    "sync",
//       "s":     "sync"
//    }
func Abbrev(words []string) map[string]string {
	seen := make(map[string]int)
	table := make(map[string]string)

	for _, w := range words {
		table[w] = w
		for n := len(w) - 1; n > 0; n -= 1 {
			ab := w[:n]
			if _, ok := seen[ab]; ok {
				seen[ab] += 1
			} else {
				seen[ab] = 0
			}

			z := seen[ab]
			if z == 0 {
				table[ab] = w
			} else if z == 1 {
				delete(table, ab)
			} else {
				break
			}
		}
	}

	return table
}
