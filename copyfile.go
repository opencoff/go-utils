// copyfile.go - copy a file efficiently using platform specific
// primitives and fallback to simple mmap'd copy.
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

package utils

import (
	"fmt"
	"io/fs"
	"os"
)

// CopyFile copies files 'src' to 'dst' using the most efficient OS primitive
// available on the runtime platform. CopyFile will use copy-on-write
// facilities if the underlying file-system implements it. It will
// fallback to copying via memory mapping 'src' and writing the blocks
// to 'dst'.
func CopyFile(src, dst string, perm fs.FileMode) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}

	defer s.Close()

	// never overwrite an existing file.
	_, err = os.Stat(dst)
	if err == nil {
		return fmt.Errorf("copyfile: destination %s already exists", dst)
	}

	d, err := NewSafeFile(dst, OPT_COW, os.O_CREATE|os.O_RDWR|os.O_EXCL, perm)
	if err != nil {
		return fmt.Errorf("copyfile: %w", err)
	}

	defer d.Abort()
	if err = copyFile(d.File, s); err != nil {
		return fmt.Errorf("copyfile: %w", err)
	}

	return d.Close()
}

type op func(dest, src string, fi fs.FileInfo) error

// order of applying these is important; we can't update
// certain attributes if we're not the owner. So, we have
// to do it in the end.
var _Mdupdaters = []op{
	clonexattr,
	utimes,
	chmod,
	chown,
}

// update all the metadata
func updateMeta(dest, src string, fi fs.FileInfo) error {
	for _, fp := range _Mdupdaters {
		if err := fp(dest, src, fi); err != nil {
			return fmt.Errorf("clonefile: %w", err)
		}
	}
	return nil
}

// CloneFile copies src to dst - including all copyable file attributes
// and xattr. CloneFile will use the best available CoW facilities provided
// by the OS and Filesystem. It will fall back to using copy via mmap(2) on
// systems that don't have CoW semantics.
func CloneFile(src, dst string) error {
	// never overwrite an existing file.
	_, err := os.Stat(dst)
	if err == nil {
		return fmt.Errorf("clonefile: destination %s already exists", dst)
	}

	s, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("clonefile: %w", err)
	}

	defer s.Close()

	fi, err := s.Stat()
	if err != nil {
		return fmt.Errorf("clonefile: %w", err)
	}

	mode := fi.Mode()
	if mode.IsRegular() {
		return copyRegular(s, dst, fi)
	}

	switch mode.Type() {
	case fs.ModeDir:
		// We only update the metadata. Caller is responsible for cloning
		// the contents (ie deep copy)
		// now set mtime/ctime, mode etc.
		err = updateMeta(dst, src, fi)

	case fs.ModeSymlink:
		err = clonelink(dst, src, fi)

	case fs.ModeDevice, fs.ModeNamedPipe:
		err = mknod(dst, src, fi)

	//case ModeSocket: XXX Add named socket support

	default:
		err = fmt.Errorf("clonefile: %s: unsupported type %#x", src, mode)
	}

	return err
}

// copy a regular file to another regular file
func copyRegular(s *os.File, dst string, fi fs.FileInfo) error {

	// We create the file so that we can write to it; we'll update the perm bits
	// later on
	d, err := NewSafeFile(dst, OPT_COW, os.O_CREATE|os.O_RDWR|os.O_EXCL, 0600)
	if err != nil {
		return fmt.Errorf("clonefile: %w", err)
	}

	defer d.Abort()
	if err = copyFile(d.File, s); err != nil {
		return fmt.Errorf("clonefile: %w", err)
	}

	// now set mtime/ctime, mode etc.
	if err = updateMeta(d.Name(), s.Name(), fi); err != nil {
		return err
	}

	return d.Close()
}

// CopyFd copies open files 'src' to 'dst' using the most efficient OS
// primitive available on the runtime platform. CopyFile will use
// copy-on-write facailities if the underlying file-system implements it.
// It will fallback to copying via memory mapping 'src' and writing the
// blocks to 'dst'.
func CopyFd(src, dst *os.File) error {
	err := copyFile(dst, src)
	if err == nil {
		err = dst.Sync()
	}
	return err
}
