// barrier.go - synchronization barrier for go-routines
//
// (c) 2024 Sudhi Herle <sw-at-herle.net>
//
// Placed in the Public Domain
// This software does not come with any express or implied
// warranty; it is provided "as is". No claim  is made to its
// suitability for any purpose.

package utils

import (
	"sync/atomic"
)

type Barrier struct {
	ch     chan bool
	closed atomic.Bool
}

func NewBarrier() *Barrier {
	b := &Barrier{
		ch: make(chan bool),
	}
	return b
}

func (b *Barrier) Wait() {
	if b.closed.Load() {
		return
	}
	for range b.ch {
	}
}

func (b *Barrier) Broadcast() {
	if !b.closed.Swap(true) {
		close(b.ch)
	}
}
