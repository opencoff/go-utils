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
	rd, wr atomic.Uint64
	mask   uint64
	q      []T
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
	z := nextpow2(uint(n))

	q.mask = uint64(z - 1)
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
	old := q.wr.Load()
	wr := (1 + old) & q.mask
	if wr == q.rd.Load() {
		return false
	}
	q.q[wr] = x
	q.wr.Store(wr)
	return true
}

// Deq dequeues an element from the queue. Returns false
// if the queue is empty, true otherwise.
func (q *SPSCQ[T]) Deq() (T, bool) {
	rd := q.rd.Load()
	if rd == q.wr.Load() {
		var z T
		return z, false
	}

	rd = (1 + rd) & q.mask
	q.rd.Store(rd)
	z := q.q[rd]
	return z, true
}

// IsEmpty returns true if the queue is empty
func (q *SPSCQ[T]) IsEmpty() bool {
	rd := q.rd.Load()
	wr := q.wr.Load()
	return rd == wr
}

// IsFull returns true if the queue is full
func (q *SPSCQ[T]) IsFull() bool {
	rd := q.rd.Load()
	wr := q.wr.Load()
	return rd == ((1 + wr) & q.mask)
}

func qlen(rd, wr, mask uint64) int {
	if wr == rd {
		return 0
	} else if rd < wr {
		return int(wr - rd)
	}
	return int((mask + 1) - rd + wr)
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
	rd := q.rd.Load()
	wr := q.wr.Load()
	n := qlen(rd, wr, q.mask)

	full := rd == ((1 + wr) & q.mask)
	empty := rd == wr

	var p string = ""
	if full {
		p = "[FULL] "
	} else if empty {
		p = "[EMPTY] "
	}

	return fmt.Sprintf("<SPSCQ %T %scap=%d, size=%d wr=%d rd=%d",
		q, p, q.mask, n, wr, rd)
}
