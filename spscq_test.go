// queue test
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

	// sync barrier
	ch chan bool
	wg sync.WaitGroup
}

func newQ(n int) *myQ {
	spq := NewSPSCQ[uint64](n)

	myq := &myQ{
		q:  spq,
		ch: make(chan bool),
	}
	return myq
}

// Barrier wait -- just wait for chan to be closed
func (m *myQ) BarrierWait() {
	for _ = range m.ch {
	}
}

// lift the barrier by closing the chan
func (m *myQ) BarrierOpen() {
	close(m.ch)
}

var qsizes = []int{128, 1024, 4096, 16384}

func TestConcurrency(t *testing.T) {
	enq := func(myq *myQ, n uint64) {
		myq.BarrierWait()

		var v uint64 = 1
		var full uint64
		q := myq.q
		start := time.Now()
		for n > 0 {
			if q.Enq(v) {
				v++
			} else {
				full++
			}
			n -= 1
		}

		myq.prod = time.Now().Sub(start)
		myq.full = full
		myq.wg.Done()
	}

	deq := func(myq *myQ, n uint64) {
		myq.BarrierWait()

		var v uint64 = 1
		var err, empty uint64
		q := myq.q
		start := time.Now()
		for n > 0 {
			z, ok := q.Deq()
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

		myq.cons = time.Now().Sub(start)
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

		myq.BarrierOpen()
		myq.wg.Wait()

		pc := float64(myq.prod) / float64(iters)
		cc := float64(myq.cons) / float64(iters)
		t.Logf("Q size %6d: %d items; P %4.2fns/item (full %d), C %4.2fns/item (empty %d, errs %d)\n",
			qsize, iters, pc, myq.full, cc, myq.empty, myq.errs)
	}
}
