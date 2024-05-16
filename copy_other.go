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

//go:build !linux

package utils

import (
	"fmt"
	"os"

	"github.com/opencoff/go-mmap"
)

// TODO macOS, FreeBSD optimizations
// Use mmap(2) to copy src to dst.
func copyFile(dst, src *os.File) error {
	_, err := mmap.Reader(src, func(b []byte) error {
		_, err := fullWrite(dst, b)
		return err
	})
	if err != nil {
		return fmt.Errorf("safefile: can't read %s: %w", src.Name(), err)
	}
	_, err = dst.Seek(0, os.SEEK_SET)
	if err != nil {
		return fmt.Errorf("safefile: can't seek to start of %s: %w", dst.Name(), err)
	}
	return nil
}
