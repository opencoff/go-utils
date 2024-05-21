// q.go - Fixed size circular queue
//
// (c) 2024 Sudhi Herle <sw-at-herle.net>
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
	wr, rd uint64
	mask   uint64 // size-1 (qhen size is a power-of-2

	q []T
}

// Make a new Queue instance to hold (at least) 'n' slots. If 'n' is
// NOT a power-of-2, this function will pick the next closest
// power-of-2.
func NewQ[T any](n int) *Q[T] {
	q := &Q[T]{}
	q.init(n)
	return q
}

// NewQFrom makes a new queue with contents from the initial list
func NewQFrom[T any](v []T) *Q[T] {
	q := &Q[T]{}
	q.init2(v)
	return q
}

func (q *Q[T]) init(n int) {
	z := nextpow2(uint64(n))

	q.rd, q.wr = 0, 0
	q.mask = z - 1
	q.q = make([]T, z)
}

func (q *Q[T]) init2(v []T) {
	q.init(len(v))
	n := copy(q.q[1:], v)
	q.wr = uint64(n)
}

// Empty the queue
func (q *Q[T]) Flush() {
	q.wr = 0
	q.rd = 0
}

// Insert new element; return false if queue full
func (q *Q[T]) Enq(x T) bool {
	wr := (1 + q.wr) & q.mask
	if wr == q.rd {
		return false
	}

	q.q[wr] = x
	q.wr = wr
	return true
}

// Remove oldest element; return false if queue empty
func (q *Q[T]) Deq() (T, bool) {
	rd := q.rd
	if rd == q.wr {
		var z T
		return z, false
	}

	rd = (rd + 1) & q.mask
	q.rd = rd
	return q.q[rd], true
}

// Return true if queue is empty
func (q *Q[T]) IsEmpty() bool {
	return qempty(q.rd, q.wr, q.mask)
}

// Return true if queue is full
func (q *Q[T]) IsFull() bool {
	return qfull(q.rd, q.wr, q.mask)
}

// Return number of valid/usable elements
func (q *Q[T]) Len() int {
	return qlen(q.rd, q.wr, q.mask)
}

// Return total capacity of the queue
func (q *Q[T]) Size() int {
	// Due to the q-full and q-empty conditions, we will
	// always have one unused slot.
	return len(q.q) - 1
}

// Dump queue in human readable form
func (q *Q[T]) String() string {
	return q.repr("Q")
}

func (q *Q[T]) repr(nm string) string {
	suff := qrepr(q.rd, q.wr, q.mask)

	return fmt.Sprintf("<%s %T %s>", nm, q, suff)
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
	q := &SyncQ[T]{}
	q.init(n)
	return q
}

// NewSyncQFrom makes a new queue with contents from the initial list
func NewSyncQFrom[T any](v []T) *SyncQ[T] {
	q := &SyncQ[T]{}
	q.init2(v)
	return q
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

// vim: ft=go:sw=8:ts=8:noexpandtab:tw=98:
