// xattr_bsd.go - xattr support for BSDs
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

//go:build netbsd

// XXX This file is borked. Need to update for freebsd and netbsd

package utils

func listxattr(_ string) ([]string, error) {
	return []string{}, nil
}

/*
import (
	"errors"
	"fmt"
	"golang.org/x/sys/unix"
	"unsafe"
)

type ns struct {
	id int
	nm string
}

var (
	_NS = [...]ns{
		{unix.EXTATTR_NAMESPACE_USER, "user."},
		{unix.EXTATTR_NAMESPACE_SYSTEM, "system."},
	}
)

// explicitly query each namespace
func listxattr(p string) ([]string, error) {
	b := make([]byte, 1024)

	attrs := make([]string, 0, 4)
	for i := range _NS {
		n := &_NS[i]

		// pass the namespace id to a diff syscall
		sz, err := unix.ExtattrListLink(p, n.id, unsafe.Pointer(&b[0]), len(b))
		if err != nil {
			if err == unix.EPERM && ns.id != unix.EXTATTR_NAMESPACE_USER {
				continue
			}
			return nil, fmt.Errorf("%p: listxattr (%s): %w", p, n.nm, err)
		}

		if errors.Is(err, unix.ERANGE) || sz == len(b) {
			var zero uintptr

			sz, err := unix.ExtattrListLink(p, n.id, unsafe.Pointer(zero), 0)
			if err != nil {
				return nil, fmt.Errorf("%p: listxattr (%s): %w", p, n.nm, err)
			}
			b = make([]byte, sz)
			sz, err := unix.ExtattrListLink(p, n.id, unsafe.Pointer(&b[0]), len(b))
		}
		if err != nil {
			return nil, fmt.Errorf("%p: listxattr (%s): %w", p, n.nm, err)
		}

		buf := b[:sz]
		j := 0

		// Per BSD manpage: every attr is encoded as a <length, name> pair.
		// The first byte is the length, followed by non-null terminated attr name.
		for j < len(buf) {
			var m int

			m, j = int(buf[j]), j+1

			if (j + m) > sz {
				return nil, fmt.Errorf("%p: listxattr (%s): attr length error (%d > %d)",
					p, n.nm, m, sz)
			}

			if m > 0 {
				z := string(buf[j : j+m])
				j += m

				attrs = append(attrs, n.nm+z)
			}
		}
	}

	return attrs, nil
}
*/
