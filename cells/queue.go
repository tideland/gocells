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
	"time"

	"github.com/tideland/golib/errors"
)

//--------------------
// CHANNEL QUEUE
//--------------------

// channelQueue implements Queue based on simple Go channels.
type channelQueue struct {
	eventc chan Event
}

// newChannelQueue creates the channel based queue with a
// defined buffer size.
func newChannelQueue(size int) Queue {
	if size < minEventBufferSize {
		size = minEventBufferSize
	}
	q := &channelQueue{
		eventc: make(chan Event, size),
	}
	return q
}

// Emit implements the Queue interface.
func (q *channelQueue) Emit(event Event) error {
	d := 2 * time.Millisecond
	for i := 0; i < 5; i++ {
		select {
		case q.eventc <- event:
			return nil
		case <-time.After(d):
			d = 2 * d
		}
	}
	return errors.New(ErrCannotEmit, errorMessages)
}

// Events implements the Queue interface.
func (q *channelQueue) Events() <-chan Event {
	return q.eventc
}

// EOF
