// spscq.go - Fixed size SPSC circular queue
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
	"sync/atomic"
)

// SPSCQ[T] is a generic & bounded single-producer/single-consumer
// queue. This queue always has a power-of-2 size. For a queue
// with capacity 'N', it will store N-1 elements.
type SPSCQ[T any] struct {
	rd atomic.Uint64
	_   [7]uint64 // cache-line pad

	wr atomic.Uint64
	_   [7]uint64 // cache-line pad

	rdc uint64    // read-index cached
	_   [7]uint64 // cache-line pad

	wrc uint64    // write-index cached
	_   [7]uint64 // cache-line pad

	mask uint64
	q    []T
}

// Make a new SPSC-Q to hold at-least 'n' elements. If 'n'
// is not a power-of-2, this function will pick the next
// closest power-of-2.
func NewSPSCQ[T any](n int) *SPSCQ[T] {
	return newSPSCQ[T](n)
}

// NewSPSCQFrom makes a new SPSC-Q with the contents from the initial
// list 'v'
func NewSPSCQFrom[T any](v []T) *SPSCQ[T] {
	q := newSPSCQ[T](len(v))

	n := copy(q.q[1:], v)
	q.wr.Store(uint64(n))
	return q
}

func newSPSCQ[T any](n int) *SPSCQ[T] {
	q := &SPSCQ[T]{}
	z := nextpow2(uint64(n))

	q.mask = z - 1
	q.rd.Store(0)
	q.wr.Store(0)
	q.q = make([]T, z)
	return q
}

// Flush and empty the queue
func (q *SPSCQ[T]) Flush() {
	q.rd.Store(0)
	q.wr.Store(0)
}

// Enq enqueues a new element. Returns true on success
// and false when Q is full.
func (q *SPSCQ[T]) Enq(x T) bool {
	wr := (1 + q.wr.Load()) & q.mask
	if wr == q.rdc {
		if q.rdc = q.rd.Load(); wr == q.rdc {
			return false
		}
	}
	q.q[wr] = x
	q.wr.Store(wr)
	return true
}

// Deq dequeues an element from the queue. Returns false
// if the queue is empty, true otherwise.
func (q *SPSCQ[T]) Deq() (T, bool) {
	rd := q.rd.Load()
	if rd == q.wrc {
		if q.wrc = q.wr.Load(); rd == q.wrc {
			var z T
			return z, false
		}
	}

	rd = (1 + rd) & q.mask
	z := q.q[rd]
	q.rd.Store(rd)
	return z, true
}

// IsEmpty returns true if the queue is empty
func (q *SPSCQ[T]) IsEmpty() bool {
	rd := q.rd.Load()
	wr := q.wr.Load()
	return qempty(rd, wr, q.mask)
}

// IsFull returns true if the queue is full
func (q *SPSCQ[T]) IsFull() bool {
	rd := q.rd.Load()
	wr := q.wr.Load()
	return qfull(rd, wr, q.mask)
}

// Len returns the number of elements in the queue
func (q *SPSCQ[T]) Len() int {
	rd := q.rd.Load()
	wr := q.wr.Load()
	return qlen(rd, wr, q.mask)
}

// Size returns the capacity of the queue
func (q *SPSCQ[T]) Size() int {
	// Due to the q-full and q-empty conditions, we will
	// always have one unused slot.
	return len(q.q) - 1
}

// String returns a human readable description of the queue
func (q *SPSCQ[T]) String() string {
	suff := qrepr(q.rd.Load(), q.wr.Load(), q.mask)

	return fmt.Sprintf("<SPSCQ %T %s>", q, suff)
}
