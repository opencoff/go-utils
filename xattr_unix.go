// xattr.go - extended attribute support for Walk
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

//go:build !dragonfly && !solaris && !illumos && !openbsd && !windows

package utils

import (
	"errors"
	"fmt"
	"golang.org/x/sys/unix"
)

// get xattr for file 'p'
func getxattr(p string) (Xattr, error) {
	attrs, err := listxattr(p)
	if err != nil {
		return nil, err
	}

	x := make(map[string]string)
	b := make([]byte, 1024)
	for _, a := range attrs {
		sz, err := unix.Lgetxattr(p, a, b)
		if errors.Is(err, unix.ERANGE) {
			sz, err = unix.Lgetxattr(p, a, nil)
			if err != nil {
				return nil, fmt.Errorf("%s: getxattr %s: %w", p, a, err)
			}
			b = make([]byte, sz)
			sz, err = unix.Lgetxattr(p, a, b)
		}
		if err != nil {
			return nil, fmt.Errorf("%s: getxattr %s: %w", p, a, err)
		}

		x[a] = string(b[:sz])
	}
	return Xattr(x), nil
}

// Set xattrs in 'x' for file 'p'
// This does not delete other xattrs already present
func setxattr(p string, x Xattr) error {
	for a, v := range x {
		err := unix.Lsetxattr(p, a, []byte(v), 0)
		if err != nil {
			return fmt.Errorf("%s: setxattr %s: %w", p, a, err)
		}
	}
	return nil
}

// remove xattrs in 'x' for file 'p'
func delxattr(p string, x Xattr) error {
	for a := range x {
		err := unix.Lremovexattr(p, a)
		if err != nil {
			return fmt.Errorf("%s: delxattr %s: %w", p, a, err)
		}
	}
	return nil
}
