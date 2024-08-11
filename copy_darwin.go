// copy_macos.go - macOS specific file copy
//
// (c) 2021 Sudhi Herle <sudhi@herle.net>
//
// Licensing Terms: GPLv2
//
// If you need a commercial license for this work, please contact
// the author.
//
// This software does not come with any express or implied
// warranty; it is provided "as is". No claim  is made to its
// suitability for any purpose.

//go:build darwin

package utils

import (
	"os"

	"golang.org/x/sys/unix"
)

// macOS 10.12+ have clonefile(2)
func copyFile(dst, src *os.File) error {
	d := dst.Name()
	s := src.Name()

	err := unix.Clonefile(s, d, unix.CLONE_NOFOLLOW)
	if err == nil {
		return nil
	}

	return copyViaMmap(dst, src)
}
