// copy_linux.go - Linux specific file copy
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

//go:build linux

package utils

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

// Do copies in chunks of 1GB
const _ioChunkSize int = 1024 * 1048576

// try to use reflinks for copying where possible.
// Fallback to copy_file_range(2) which is available on all linuxes.
func copyFile(dst, src *os.File) error {
	d := int(dst.Fd())
	s := int(src.Fd())

	// First try to reflink.
	err := unix.IoctlFileClone(int(d), int(s))
	if err == nil {
		return nil
	}

	st, err := src.Stat()
	if err != nil {
		return fmt.Errorf("safefile: %w", err)
	}

	// Fallback to copy_file_range(2)
	var roff, woff int64
	sz := st.Size()
	for sz > 0 {
		var n int = _ioChunkSize
		if int(sz) < n {
			n = int(sz)
		}
		m, err := unix.CopyFileRange(s, &roff, d, &woff, n, 0)
		if err != nil {
			return fmt.Errorf("safefile: copy: %w", err)
		}
		if m == 0 {
			return fmt.Errorf("safefile: zero sized copy of %s", src.Name())
		}
		sz -= int64(m)
		roff += int64(m)
		woff += int64(m)
	}
	_, err = dst.Seek(0, os.SEEK_SET)
	if err != nil {
		return fmt.Errorf("safefile: seek: %w", err)
	}
	return nil
}
