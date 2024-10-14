// xattr_openbsd.go - xattr support for openbsd
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

//go:build openbsd

package utils

func listxattr(_ string) ([]string, error) {
	return []string{}, nil
}

func getxattr(_ string) (Xattr, error) {
	return Xattr{}, nil
}

func setxattr(_ string, _ Xattr) error {
	return nil
}
