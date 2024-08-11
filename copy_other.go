// copy_other.go - non-Linux file copy
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

//go:build !linux && !darwin

package utils

import (
	"os"
)

func copyFile(dst, src *os.File) error {
	return copyViaMmap(dst, src)
}
