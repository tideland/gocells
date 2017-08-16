// Tideland Go Cells - Queue
//
// Copyright (C) 2010-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package cells

//--------------------
// IMPORTS
//--------------------

import (
	"github.com/tideland/golib/loop"
)

//--------------------
// CONSTANTS
//--------------------

// TODO(mue) maxPending will later limit the queue size.
const maxPending = 65536

//--------------------
// IN-MEMORY QUEUE
//--------------------

// inMemoryQueue implements Queue based on a simple channel.
type inMemoryQueue struct {
	inc  chan Event
	outc chan Event
	loop loop.Loop
}

// newInMemoryQueue creates the in-memory queue.
func newInMemoryQueue() Queue {
	q := &inMemoryQueue{
		inc:  make(chan Event),
		outc: make(chan Event),
	}
	q.loop = loop.Go(q.backendLoop)
	return q
}

// Emit implements the Queue interface.
func (q *inMemoryQueue) Emit(event Event) error {
	q.inc <- event
	return nil
}

// Events implements the Queue interface.
func (q *inMemoryQueue) Events() <-chan Event {
	return q.outc
}

// Close implements the Queue interface.
func (q *inMemoryQueue) Close() error {
	return q.loop.Stop()
}

// backendLoop runs the queue goroutine.
func (q *inMemoryQueue) backendLoop(l loop.Loop) error {
	defer close(q.outc)
	defer close(q.inc)

	var pending []Event

	for {
		var first Event
		var outc chan Event

		if len(pending) > 0 {
			first = pending[0]
			outc = q.outc
		}

		select {
		case <-l.ShallStop():
			return nil
		case event := <-q.inc:
			// TODO(mue) Limit queue size, have to think about strategy.
			if len(pending) < maxPending {
				pending = append(pending, event)
			}
		case outc <- first:
			pending = pending[1:]
		}
	}
}

// EOF
