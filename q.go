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

// Notes:
//   - read from 'rd', write to 'wr+1'.
//   - queue size always a power-of-2
//   - for a queue of capacity N, it will store N-1 usable elements
//   - queue-empty: rd   == wr
//   - queue-full:  wr+1 == rd

// Q[T] is a generic fixed-size queue. This queue always has a
// power-of-2 size.  For a queue with capacity 'N', it will store
// N-1 queue elements.
type Q[T any] struct {
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
	w := &Q[T]{}
	w.init(n)
	return w
}

func (w *Q[T]) init(n int) {
	z := nextpow2(uint(n))

	w.rd, w.wr = 0, 0
	w.mask = z - 1
	w.q = make([]T, z)
}

func (w *Q[T]) init2(v []T) {
	w.init(len(v))
	n := copy(w.q[1:], v)
	w.wr = uint(n)
}

// NewQFrom makes a new queue with contents from the initial list
func NewQFrom[T any](v []T) *Q[T] {
	w := &Q[T]{}
	w.init2(v)
	return w
}

// Empty the queue
func (w *Q[T]) Flush() {
	w.wr = 0
	w.rd = 0
}

// Insert new element; return false if queue full
func (w *Q[T]) Enq(x T) bool {
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
	return w.rd == w.wr
}

// Return true if queue is full
func (w *Q[T]) IsFull() bool {
	return w.rd == ((1 + w.wr) & w.mask)
}

// Return number of valid/usable elements
func (w *Q[T]) Len() int {
	return w.size()
}

// Return total capacity of the queue
func (w *Q[T]) Size() int {
	// Due to the q-full and q-empty conditions, we will
	// always have one unused slot.
	return len(w.q) - 1
}

// Dump queue in human readable form
func (w *Q[T]) String() string {
	return w.repr("Q")
}

func (w *Q[T]) repr(nm string) string {
	full := w.rd == ((1 + w.wr) & w.mask)
	empty := w.rd == w.wr

	var p string = ""
	if full {
		p = "[FULL] "
	} else if empty {
		p = "[EMPTY] "
	}

	return fmt.Sprintf("<%s %T %scap=%d, size=%d wr=%d rd=%d",
		nm, w, p, w.mask, w.size(), w.wr, w.rd)
}

// SyncQ[T] is a generic, thread-safe, fixed-size queue. This queue
// always has a power-of-2 size. For a queue with capacity 'N', it will
// store N-1 queue elements.
type SyncQ[T any] struct {
	Q[T]
	sync.Mutex
}

// Make a new thread-safe queue instance to hold (at least) 'n' slots.
// If 'n' is NOT a power-of-2, this function will pick the next closest
// power-of-2.
func NewSyncQ[T any](n int) *SyncQ[T] {
	w := &SyncQ[T]{}
	w.init(n)
	return w
}

// NewSyncQFrom makes a new queue with contents from the initial list
func NewSyncQFrom[T any](v []T) *SyncQ[T] {
	w := &SyncQ[T]{}
	w.init2(v)
	return w
}

// Flush empties the queue
func (q *SyncQ[T]) Flush() {
	q.Lock()
	q.Q.Flush()
	q.Unlock()
}

// Enq enqueues a new element to the queue, return false if the queue is full
// and true otherwise.
func (q *SyncQ[T]) Enq(x T) bool {
	q.Lock()
	r := q.Q.Enq(x)
	q.Unlock()
	return r
}

// Deq dequeues an element from the queue and returns it. The bool retval is false
// if the queue is empty and true otherwise.
func (q *SyncQ[T]) Deq() (T, bool) {
	q.Lock()
	a, b := q.Q.Deq()
	q.Unlock()
	return a, b
}

// IsEmpty returns true if the queue is empty and false otherwise
func (q *SyncQ[T]) IsEmpty() bool {
	q.Lock()
	r := q.Q.IsEmpty()
	q.Unlock()
	return r
}

// IsFull returns true if the queue is full and false otherwise
func (q *SyncQ[T]) IsFull() bool {
	q.Lock()
	r := q.Q.IsFull()
	q.Unlock()
	return r
}

// Len returns the number of elements in the queue
func (q *SyncQ[T]) Len() int {
	q.Lock()
	r := q.Q.Len()
	q.Unlock()
	return r
}

// Size returns the capacity of the queue
func (q *SyncQ[T]) Size() int {
	q.Lock()
	r := q.Q.Size()
	q.Unlock()
	return r
}

// String prints a string representation of the queue
func (q *SyncQ[T]) String() string {
	q.Lock()
	s := q.Q.repr("SyncQ")
	q.Unlock()
	return s
}

// internal func to return queue size
// caller must hold lock
func (w *Q[T]) size() int {
	if w.wr == w.rd {
		return 0
	} else if w.rd < w.wr {
		return int(w.wr - w.rd)
	}
	return int((w.mask + 1) - w.rd + w.wr)
}

// vim: ft=go:sw=8:ts=8:noexpandtab:tw=98:
