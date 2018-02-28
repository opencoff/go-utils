// uuid.go -- simple and fast uuid generation using crypto random
//
// (c) 2015, 2016 - Sudhi Herle <sudhi@herle.net>
//
// License: Public Domain
//

package util

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// UUID abstraction
type UUID struct {
	b [16]byte
}

// Make a new random UUID
func NewUUID() *UUID {
	u := &UUID{}

	b := u.b[:]

	rand.Read(b)

	b[6] = 0x4f & (b[6] | 0x40)
	b[8] = 0xBf & (b[8] | 0x80)

	return u
}

// Make UUID from a raw byte stream.
func MakeUUID(b []byte) *UUID {
	if len(b) != 16 {
		return nil
	}

	u := &UUID{}
	copy(u.b[:], b)
	return u
}

// Return a serializable ASCII string
func (u *UUID) Marshal() string {
	return hex.EncodeToString(u.b[:])
}

// Unmarshal a string
func UnmarshalUUID(s string) *UUID {
	b, err := hex.DecodeString(s)
	if err == nil && len(b) == 16 {
		u := &UUID{}
		copy(u.b[:], b)
		return u
	}

	return nil
}

func (u *UUID) Bytes() []byte {
	return u.b[:]
}

func (u UUID) String() string {
	b := u.b[:]
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// EOF
// vim: ft=go:sw=8:ts=8:noexpandtab:tw=98:
