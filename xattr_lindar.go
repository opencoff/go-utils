// xattr_unix.go - xattr support for unix like systems
//
// (c) 2023- Sudhi Herle <sudhi@herle.net>
//
// Licensing Terms: GPLv2
//
// If you need a commercial license for this work, please contact
// the author.
//
// This software does not come with any express or implied
// warranty; it is provided "as is". No claim  is made to its
// suitability for any purpose.

//go:build linux || darwin

package utils

import (
	"errors"
	"fmt"
	"golang.org/x/sys/unix"
	"strings"
)

func listxattr(p string) ([]string, error) {
	b := make([]byte, 1024)

	sz, err := unix.Llistxattr(p, b)

	// darwin doesn't return ERANGE - so we will take an extra syscall in
	// case the attrbuf is exactly sized
	if errors.Is(err, unix.ERANGE) || sz == len(b) {
		sz, err = unix.Llistxattr(p, nil)
		if err != nil {
			return nil, fmt.Errorf("%s: listxattr: %w", p, err)
		}
		b = make([]byte, sz)
		sz, err = unix.Llistxattr(p, b)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: listxattr: %w", p, err)
	}

	// the xattr are a simple, unordered list of nul terminated strings.
	s := string(b[:sz])
	v := strings.Split(s, "\x00")
	return clean(v), nil
}

// remove empty strings in the list
func clean(v []string) []string {
	i := 0
	for _, s := range v {
		if s != "" {
			v[i] = s
			i++
		}
	}
	return v[:i]
}
