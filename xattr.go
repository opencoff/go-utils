// xattr.go - extended attribute support for Walk
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

package utils

import (
	"fmt"
	"strings"
)

type Xattr map[string]string

func (x Xattr) String() string {
	var s strings.Builder
	for k, v := range x {
		s.WriteString(fmt.Sprintf("%s=%s\n", k, v))
	}
	return s.String()
}
