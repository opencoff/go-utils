// q.go - Fixed size circular queue
//
// (c) 2014 Sudhi Herle <sw-at-herle.net>
//
// Placed in the Public Domain
// This software does not come with any express or implied
// warranty; it is provided "as is". No claim  is made to its
// suitability for any purpose.

package util

import (
	"fmt"
	"sync"
)

// Thread-safe, fixed-size circular queue.
// Stores interface{} in each queue slot.
//
// Notes:
//  - read from 'rd', write to 'wr+1'.
//  - queue size always a power-of-2
//  - for a queue of capacity N, it will store N-1 usable elements
//  - queue-empty: rd   == wr
//  - queue-full:  wr+1 == rd
type Q struct {
	q      []interface{}
	wr, rd uint
	mask   uint // size-1 (when size is a power-of-2

	l sync.Mutex
}

// return next power of 2
// return n if already a power of 2
func nextpow2(n uint) uint {
	n -= 1
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	return n + 1
}

// Make a new Queue instance to hold (at least) 'n' slots. If 'n' is
// NOT a power-of-2, this function will pick the next closest
// power-of-2.
func NewQ(n int) *Q {
	w := &Q{}
	w.mask = nextpow2(uint(n)) - 1
	w.q = make([]interface{}, w.mask+1)
	w.wr = 0
	w.rd = 0

	return w
}

// Empty the queue
func (w *Q) Flush() {
	w.wr = 0
	w.rd = 0
}

// Insert new element; return false if queue full
func (w *Q) Enq(x interface{}) bool {
	w.l.Lock()
	defer w.l.Unlock()

	wr := (1 + w.wr) & w.mask
	if wr == w.rd {
		return false
	}

	w.q[wr] = x
	w.wr = wr
	return true
}

// Remove oldest element; return false if queue empty
func (w *Q) Deq() (interface{}, bool) {
	w.l.Lock()
	defer w.l.Unlock()
	rd := w.rd
	if rd == w.wr {
		return nil, false
	}

	rd = (rd + 1) & w.mask
	w.rd = rd
	item := w.q[rd]
	w.q[rd] = nil // needed to ensure GC picks up items
	return item, true
}

// Return true if queue is empty
func (w *Q) IsEmpty() bool {
	w.l.Lock()
	defer w.l.Unlock()
	return w.rd == w.wr
}

// Return true if queue is full
func (w *Q) IsFull() bool {
	w.l.Lock()
	defer w.l.Unlock()
	return w.rd == (1+w.wr)&w.mask
}

// Return number of valid/usable elements
func (w *Q) Size() int {
	w.l.Lock()
	defer w.l.Unlock()

	return w.size()
}

// Dump queue in human readable form
func (w *Q) String() string {
	w.l.Lock()
	defer w.l.Unlock()
	return w.repr()
}

// internal func to print string repr of queue
// caller must hold lock
func (w *Q) repr() string {
	s := fmt.Sprintf("<Q cap=%d, siz=%d wr=%d rd=%d>",
		w.mask+1, w.size(), w.wr, w.rd)

	return s
}

// internal func to return queue size
// caller must hold lock
func (w *Q) size() int {
	if w.wr == w.rd {
		return 0
	} else if w.rd < w.wr {
		return int(w.wr - w.rd)
	} else {
		return int((w.mask + 1) - w.rd + w.wr)
	}
}

// EOF
// vim: ft=go:sw=8:ts=8:noexpandtab:tw=98:
