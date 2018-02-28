// bufpool.go -- Buffer pool (blocking) abstraction
//
// (c) 2015, 2016 -- Sudhi Herle <sudhi@herle.net>
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

// A fixed-size buffer-pool backed by a channel and hence, callers
// are blocked if there are no more buffers available.
// Callers are expected to free the buffer back to its
// originating pool.
type Bufpool struct {
	Size int
	q    chan interface{}
}

// Default pool size
const Poolsize = 64

// NewBufpool creates a new Bufpool. The caller supplies a
// constructor for creating new buffers and filling the
// queue with initial elements.
func NewBufpool(sz int, ctor func() interface{}) *Bufpool {
	if sz <= 0 {
		sz = Poolsize
	}

	b := &Bufpool{Size: sz}
	b.q = make(chan interface{}, sz)

	for i := 0; i < sz; i++ {
		b.q <- ctor()
	}

	return b
}

// Put an item into the bufpool. This should not ever block; it
// indicates pool integrity failure (duplicates or erroneous Puts).
func (b *Bufpool) Put(o interface{}) {
	select {
	case b.q <- o:
		break
	default:
		panic("Bufpool put blocked. Queue corrupt?")
	}
}

// Get the next available item from the pool; block the caller if
// none are available.
func (b *Bufpool) Get() interface{} {
	o := <-b.q
	return o
}

// EOF
// vim: ft=go:sw=8:ts=8:noexpandtab:tw=98:
