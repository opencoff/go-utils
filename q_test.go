// queue test
package utils

import (
	"testing"
)

// Basic sanity tests
func TestBasic(t *testing.T) {
	assert := newAsserter(t)

	var v bool

	q := NewQ[int](3)

	assert(q.IsEmpty(), "expected q to be empty")
	assert(!q.IsFull(), "expected q to not be full")

	v = q.Enq(10)
	assert(v, "enq-10 failed")

	v = q.Enq(20)
	assert(v, "enq-20 failed")

	v = q.Enq(30)
	assert(v, "enq-30 failed")

	assert(q.IsFull(), "expected q to be full")
	assert(!q.IsEmpty(), "expected q to not be empty")
	assert(q.Len() == 3, "qsize exp 3, saw %d", q.Len())

	// Now q will be full
	v = q.Enq(40)
	assert(!v, "enq-4 should have failed")

	// Pull items off the queue

	z, v := q.Deq()
	assert(v, "deq-0 failed")
	assert(z == 10, "deq-0 value mismatch, exp 10, saw %d", z)

	z, v = q.Deq()
	assert(v, "deq-1 failed")
	assert(z == 20, "deq-1 value mismatch, exp 20, saw %d", z)

	z, v = q.Deq()
	assert(v, "deq-2 failed")
	assert(z == 30, "deq-2 value mismatch, exp 30, saw %d", z)

	assert(q.IsEmpty(), "expected q to be empty")

	z, v = q.Deq()
	assert(!v, "expected deq to fail")
}

// Test wrap around
func TestWrapAround(t *testing.T) {
	assert := newAsserter(t)

	var v bool

	q := NewQ[int](3)

	v = q.Enq(10)
	assert(v, "enq-10 failed")

	v = q.Enq(20)
	assert(v, "enq-20 failed")

	v = q.Enq(30)
	assert(v, "enq-30 failed")

	assert(q.IsFull(), "expected q to be full")
	assert(!q.IsEmpty(), "expected q to not be empty")
	assert(q.Len() == 3, "qsize exp 3, saw %d", q.Len())

	z, v := q.Deq()
	assert(v, "deq-0 failed")
	assert(z == 10, "deq-0 value mismatch, exp 10, saw %d", z)

	z, v = q.Deq()
	assert(v, "deq-1 failed")
	assert(z == 20, "deq-1 value mismatch, exp 20, saw %d", z)

	// This will wrap around
	v = q.Enq(40)
	assert(v, "enq-40 failed")
	v = q.Enq(50)
	assert(v, "enq-50 failed")

	assert(q.Len() == 3, "q size mismatch, exp 3, saw %d", q.Len())
}

// Test wrap around
func TestWithInitialData(t *testing.T) {
	assert := newAsserter(t)

	var z = [...]int{
		33,
		44,
		55,
	}

	q := NewQFrom[int](z[:])

	assert(q.Len() == len(z), "initial q-size: exp %d, saw %d", len(z), q.Len())
	assert(q.IsFull(), "expected q to be full")

	for i, v := range z {
		a, ok := q.Deq()
		assert(ok, "deq-%d failed; exp success", i)
		assert(a == v, "deq-%d: exp %d, saw %d", i, v, a)
	}
	assert(q.IsEmpty(), "expected q to be empty")
}
