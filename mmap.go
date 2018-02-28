// mmap.go - Better interface to mmap(2) on go
//
// (c) 2016 Sudhi Herle <sudhi@herle.net>
//
// Licensing Terms: GPLv2
//
// If you need a commercial license for this work, please contact
// the author.
//
// This software does not come with any express or implied
// warranty; it is provided "as is". No claim  is made to its
// suitability for any purpose.

package util

import (
	"fmt"
	"io"
	"os"
	"syscall"
)

// A mmap'd file reader that processes an already open file in large chunks.
// The default chunk-size is 1GB (1024 x 1024 x 1024 bytes). Data is
// read and written to the provided io.Writer without any data copy.
// The writer is provided as an argument because we need to unmap
// the chunks after they are processed - otherwise, have to punt
// that function to the caller and it gets messy.
//
// Returns: number of bytes written (mapped in total)
//
// This function can be used to efficiently hash very large files:
//
//    h := sha256.New()
//    err := MmapReader(fd, 0, 0, h)
func MmapReader(fd *os.File, off, sz int64, wr io.Writer) (int64, error) {
	// Mmap'ing large files won't work. We need to do it in 1 or 2G
	// chunks.
	const chunk int64 = 1 * 1024 * 1024 * 1024
	var fsz int64

	st, err := fd.Stat()
	if err != nil {
		return 0, fmt.Errorf("mmap: can't stat: %s", err)
	}

	fsz = st.Size()

	if sz == 0 {
		sz = fsz
	}

	if off > fsz {
		return 0, fmt.Errorf("can't mmap offset %v outside filesize %v", off, fsz)
	}

	// Don't mmap outside the available size
	// This is a benign error?
	if (sz + off) > fsz {
		sz = fsz - off
	}

	var z int64

	for sz > 0 {
		var n = int(sz)

		if sz > chunk {
			n = int(chunk)
		}

		mem, err := syscall.Mmap(int(fd.Fd()), off, n, syscall.PROT_READ, syscall.MAP_SHARED)
		if err != nil {
			return 0, fmt.Errorf("can't mmap %v bytes at %v: %s", n, off, err)
		}

		wr.Write(mem)
		syscall.Munmap(mem)

		off += int64(n)
		sz -= int64(n)
		z += int64(n)
	}

	return z, nil
}

// EOF
// vim: ft=go:sw=8:ts=8:noexpandtab:tw=98:
