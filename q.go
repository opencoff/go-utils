// q.go - Fixed size circular queue
//
// (c) 2014 Sudhi Herle <sw-at-herle.net>
//
// Placed in the Public Domain
// This software does not come with any express or implied
// warranty; it is provided "as is". No claim  is made to its
// suitability for any purpose.

package utils

import (
	"fmt"
	"sync"
)

// Thread-safe, fixed-size circular queue.
// Stores interface{} in each queue slot.
//
// Notes:
//   - read from 'rd', write to 'wr+1'.
//   - queue size always a power-of-2
//   - for a queue of capacity N, it will store N-1 usable elements
//   - queue-empty: rd   == wr
//   - queue-full:  wr+1 == rd
type Q[T any] struct {
	sync.Mutex

	wr, rd uint
	mask   uint // size-1 (when size is a power-of-2

	q []T
}

// return strictly _next_ power of 2
func nextpow2(z uint) uint {
	if z == 0 {
		return 2
	}

	n := z - 1
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	if n += 1; n == z {
		n <<= 1
	}
	return n
}

// Make a new Queue instance to hold (at least) 'n' slots. If 'n' is
// NOT a power-of-2, this function will pick the next closest
// power-of-2.
func NewQ[T any](n int) *Q[T] {
	z := nextpow2(uint(n))
	w := &Q[T]{
		rd:   0,
		wr:   0,
		mask: z - 1,
		q:    make([]T, z),
	}
	return w
}

// NewQFrom makes a new queue with contents from the initial list
func NewQFrom[T any](v []T) *Q[T] {
	n := uint(len(v))
	z := nextpow2(n)
	w := &Q[T]{
		rd:   0,
		wr:   n,
		mask: z - 1,
		q:    make([]T, z),
	}

	for i, a := range v {
		w.q[i+1] = a
	}

	return w
}

// Empty the queue
func (w *Q[T]) Flush() {
	w.Lock()
	w.wr = 0
	w.rd = 0
	w.Unlock()
}

// Insert new element; return false if queue full
func (w *Q[T]) Enq(x T) bool {
	w.Lock()
	defer w.Unlock()

	wr := (1 + w.wr) & w.mask
	if wr == w.rd {
		return false
	}

	w.q[wr] = x
	w.wr = wr
	return true
}

// Remove oldest element; return false if queue empty
func (w *Q[T]) Deq() (T, bool) {
	w.Lock()
	defer w.Unlock()

	rd := w.rd
	if rd == w.wr {
		var z T
		return z, false
	}

	rd = (rd + 1) & w.mask
	w.rd = rd
	return w.q[rd], true
}

// Return true if queue is empty
func (w *Q[T]) IsEmpty() bool {
	w.Lock()
	defer w.Unlock()
	return w.rd == w.wr
}

// Return true if queue is full
func (w *Q[T]) IsFull() bool {
	w.Lock()
	defer w.Unlock()
	return w.rd == (1+w.wr)&w.mask
}

// Return number of valid/usable elements
func (w *Q[T]) Len() int {
	w.Lock()
	defer w.Unlock()

	return w.size()
}

// Dump queue in human readable form
func (w *Q[T]) String() string {
	w.Lock()
	defer w.Unlock()
	s := fmt.Sprintf("<Q-%T cap=%d, siz=%d wr=%d rd=%d>",
		w, w.mask+1, w.size(), w.wr, w.rd)

	return s
}

// internal func to return queue size
// caller must hold lock
func (w *Q[T]) size() int {
	if w.wr == w.rd {
		return 0
	} else if w.rd < w.wr {
		return int(w.wr - w.rd)
	} else {
		return int((w.mask + 1) - w.rd + w.wr)
	}
}

// vim: ft=go:sw=8:ts=8:noexpandtab:tw=98:
