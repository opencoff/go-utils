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
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestSPSCFunctionality(t *testing.T) {
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

	// number of times q had seq errs
	errs uint64

	b  *Barrier
	wg sync.WaitGroup
}

// TestSPSCWrapAround ensures the ring buffer indices cycle correctly
// without crashing or losing data when they exceed the array size.
func TestSPSCWrapAround(t *testing.T) {
	assert := newAsserter(t)

	// Small queue to force frequent wrapping
	q := NewSPSCQ[int](2)
	cycles := 100

	for i := 0; i < cycles; i++ {
		// Enqueue 1
		ok := q.Enq(i)
		assert(ok, "failed to enq at cycle %d", i)

		// Dequeue 1
		val, ok := q.Deq()
		assert(ok, "failed to deq at cycle %d", i)
		assert(val == i, "cycle %d: exp %d, got %d", i, i, val)
	}
}

// TestSPSCZeroValue ensures that the queue handles the zero value of the type (0)
// correctly and distinguishes it from 'empty'.
func TestSPSCZeroValue(t *testing.T) {
	assert := newAsserter(t)
	q := NewSPSCQ[int](16)

	// Case 1: Enqueue 0 explicitly
	ok := q.Enq(0)
	assert(ok, "failed to enq 0")

	// Case 2: Dequeue 0
	val, ok := q.Deq()
	assert(ok, "should have received value 0, got empty")
	assert(val == 0, "exp 0, got %d", val)

	// Case 3: Queue should now be empty
	_, ok = q.Deq()
	assert(!ok, "queue should be empty after deq 0")
}

// TestSPSCTorture introduces random jitter (sleeps) to simulate real-world
// unpredictable thread scheduling (GC pauses, context switches).
func TestSPSCTorture(t *testing.T) {
	// A smaller number of items, but with sleeps
	const iters = 10_000
	q := NewSPSCQ[int](128)

	var wg sync.WaitGroup
	wg.Add(2)

	// Producer with Jitter
	go func() {
		defer wg.Done()
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for i := 0; i < iters; i++ {
			for !q.Enq(i) {
				// Busy wait, occasionally sleep
				if r.Float32() < 0.01 { // 1% chance to sleep
					time.Sleep(time.Microsecond)
				}
			}
			// Occasional sleep after success too
			if r.Float32() < 0.005 {
				time.Sleep(time.Microsecond)
			}
		}
	}()

	// Consumer with Jitter
	go func() {
		defer wg.Done()
		r := rand.New(rand.NewSource(time.Now().UnixNano() + 1))
		for i := 0; i < iters; {
			val, ok := q.Deq()
			if !ok {
				if r.Float32() < 0.01 {
					time.Sleep(time.Microsecond)
				}
				continue
			}

			// Verification
			if val != i {
				t.Errorf("Torture fail: exp %d, got %d", i, val)
				return
			}
			i++
		}
	}()

	wg.Wait()
}

func newQ(n int) *myQ {
	myq := &myQ{
		q: NewSPSCQ[uint64](n),
		b: NewBarrier(),
	}
	return myq
}

var qsizes = []int{128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536}

func TestSPSCConcurrency(t *testing.T) {
	enq := func(myq *myQ, n uint64) {
		myq.b.Wait()

		var v uint64 = 0

		q := myq.q
		start := time.Now()
		for v < n {
			if q.Enq(v) {
				v++
			}
		}

		myq.prod = time.Since(start)
		myq.wg.Done()
	}

	deq := func(myq *myQ, n uint64) {
		myq.b.Wait()

		var v uint64 = 0
		var err uint64
		q := myq.q
		start := time.Now()
		for v < n {
			if z, ok := q.Deq(); ok {
				if v != z {
					err++
				}
				v++
			}
		}

		myq.cons = time.Since(start)
		myq.errs = err
		myq.wg.Done()
	}

	const iters uint64 = 200 * 1048576
	for _, qsize := range qsizes {
		myq := newQ(qsize)

		myq.wg.Add(2)
		go enq(myq, iters)
		go deq(myq, iters)

		myq.b.Broadcast()
		myq.wg.Wait()

		pc := float64(myq.prod) / float64(iters)
		cc := float64(myq.cons) / float64(iters)
		t.Logf("Q size %6d: %d items; P %4.2f ns/item, C %4.2f ns/item (errs %d)\n",
			qsize, iters, pc, cc, myq.errs)
	}
}
