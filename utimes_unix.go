// utimes_unix.go -- set file times for unixish platforms
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

//go:build unix && !darwin && !freebsd && !netbsd

package utils

import (
	"fmt"
	"io/fs"
	"syscall"

	"golang.org/x/sys/unix"
)

func utimes(dest string, _ string, fi fs.FileInfo) error {
	if st, ok := fi.Sys().(*syscall.Stat_t); ok {
		tv := []unix.Timeval{
			unix.NsecToTimeval(st.Atim.Nano()),
			unix.NsecToTimeval(st.Mtim.Nano()),
		}

		if err := unix.Lutimes(dest, tv); err != nil {
			return fmt.Errorf("utimes: set: %w", err)
		}
	}
	return nil
}
