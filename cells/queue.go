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

//--------------------
// IN-MEMORY QUEUE
//--------------------

// inMemoryQueue implements Queue based on a simple channel.
type inMemoryQueue struct {
	queuec chan Event
}

// newInMemoryQueue creates the in-memory queue.
func newInMemoryQueue() Queue {
	return &inMemoryQueue{
		queuec: make(chan Event, 2048),
	}
}

// Emit implements the Queue interface.
func (q *inMemoryQueue) Emit(event Event) error {
	q.queuec <- event
	return nil
}

// Events implements the Queue interface.
func (q *inMemoryQueue) Events() <-chan Event {
	return q.queuec
}

// Close implements the Queue interface.
func (q *inMemoryQueue) Close() error {
	close(q.queuec)
	return nil
}

// EOF
