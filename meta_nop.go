// meta_nop.go -- metadata updates for unsupported systems
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

//go:build !unix

package utils

import (
	"fmt"
	"io/fs"
)

func utimes(dest string, _ string, fi fs.FileInfo) error {
	return fmt.Errorf("utimes: not supported")
}

func chown(dest string, _ string, fi fs.FileInfo) error {
	return fmt.Errorf("chown: not supported")
}

func chmod(dest string, _ string, fi fs.FileInfo) error {
	return fmt.Errorf("chmod: not supported")
}

func mknod(dest string, src string, fi fs.FileInfo) error {
	return fmt.Errorf("mknod: not supported")
}

// clone a symlink - ie we make the target point to the same one as src
func clonelink(dest string, src string, fi fs.FileInfo) error {
	return fmt.Errorf("clonelink: not supported")
}

func clonexattr(dest, src string, _ fs.FileInfo) error {
	return fmt.Errorf("clonexattr: not supported")
}
