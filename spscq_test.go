// spsc queue test
//
// (c) 2024 Sudhi Herle <sw-at-herle.net>
//
// Placed in the Public Domain
// This software does not come with any express or implied
// warranty; it is provided "as is". No claim  is made to its
// suitability for any purpose.
package utils

import (
	"sync"
	"testing"
	"time"
)

func TestFunctionality(t *testing.T) {
	assert := newAsserter(t)

	var z int
	var ok bool

	q := NewSPSCQ[int](3)

	ok = q.Enq(100)
	assert(ok, "can't enq 100")
	ok = q.Enq(200)
	assert(ok, "can't enq 200")
	ok = q.Enq(300)
	assert(ok, "can't enq 300")

	ok = q.Enq(400)
	assert(!ok, "expected q full\n%s", q)

	z, ok = q.Deq()
	assert(ok, "can't deq 100")
	assert(z == 100, "exp 100, saw %d", z)
	z, ok = q.Deq()
	assert(ok, "can't deq 200")
	assert(z == 200, "exp 200, saw %d", z)
	z, ok = q.Deq()
	assert(ok, "can't deq 300")
	assert(z == 300, "exp 300, saw %d", z)

	z, ok = q.Deq()
	assert(!ok, "expected q empty\n%s", q)
}

type myQ struct {
	q *SPSCQ[uint64]

	// time taken by producer
	prod time.Duration

	// time taken by consumer
	cons time.Duration

	// number of times q was full/empty, had seq errs
	full, empty, errs uint64

	b  *Barrier
	wg sync.WaitGroup
}

func newQ(n int) *myQ {
	myq := &myQ{
		q: NewSPSCQ[uint64](n),
		b: NewBarrier(),
	}
	return myq
}

var qsizes = []int{128, 1024, 4096, 16384}

func TestConcurrency(t *testing.T) {
	enq := func(myq *myQ, n uint64) {
		myq.b.Wait()

		var v uint64 = 1
		var full uint64
		var tot time.Duration

		q := myq.q
		for n > 0 {
			start := time.Now()
			ok := q.Enq(v)
			tot += time.Now().Sub(start)
			if ok {
				v++
			} else {
				full++
			}
			n -= 1
		}

		myq.prod = tot
		myq.full = full
		myq.wg.Done()
	}

	deq := func(myq *myQ, n uint64) {
		myq.b.Wait()

		var v uint64 = 1
		var err, empty uint64
		var tot time.Duration

		q := myq.q
		for n > 0 {
			start := time.Now()
			z, ok := q.Deq()
			tot += time.Now().Sub(start)
			if ok {
				if v != z {
					err++
				}
				v++
			} else {
				empty++
			}
			n -= 1
		}

		myq.cons = tot
		myq.empty = empty
		myq.errs = err
		myq.wg.Done()
	}

	const iters uint64 = 10485760
	for _, qsize := range qsizes {
		myq := newQ(qsize)

		myq.wg.Add(2)
		go enq(myq, iters)
		go deq(myq, iters)

		myq.b.Broadcast()
		myq.wg.Wait()

		pc := float64(myq.prod) / float64(iters)
		cc := float64(myq.cons) / float64(iters)
		t.Logf("Q size %6d: %d items; P %4.2fns/item (full %d), C %4.2fns/item (empty %d, errs %d)\n",
			qsize, iters, pc, myq.full, cc, myq.empty, myq.errs)
	}
}
