// utimes_darwin.go -- set file times for macOS/Darwin
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

//go:build darwin || freebsd || netbsd

package utils

import (
	"fmt"
	"io/fs"
	"os"
	"syscall"
	"time"
)

func utimes(dest string, _ string, fi fs.FileInfo) error {
	if st, ok := fi.Sys().(*syscall.Stat_t); ok {
		at := time.Unix(st.Atimespec.Sec, st.Atimespec.Nsec)
		mt := time.Unix(st.Mtimespec.Sec, st.Mtimespec.Nsec)
		if err := os.Chtimes(dest, at, mt); err != nil {
			return fmt.Errorf("utimes: set: %w", err)
		}
	}
	return nil
}
