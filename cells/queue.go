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
// CHANNEL QUEUE
//--------------------

// channelQueue implements Queue based on simple Go channels.
type channelQueue struct {
	eventc      chan Event
	subscribers map[Subscriber]struct{}
}

// newChannelQueue creates the channel based queue with a
// defined buffer size.
func newChannelQueue(size int) Queue {
	if size < minEventBufferSize {
		size = minEventBufferSize
	}
	q := &channelQueue{
		eventc:      make(chan Event, size),
		subscribers: make(map[Subscriber]struct{}),
	}
	return q
}

// Emit implements the Queue interface.
func (q *channelQueue) Emit(event Event) error {
	return nil
}

// Subscribe implements the Queue interface.
func (q *channelQueue) Subscribe(subscribers ...Subscriber) error {
	for _, subscriber := range subscribers {
		q.subscribers[subscriber] = struct{}{}
	}
	return nil
}

// Unsubscribe implements the Queue interface.
func (q *channelQueue) Unsubscribe(subscribers ...Subscriber) error {
	for _, subscriber := range subscribers {
		delete(q.subscribers, subscriber)
	}
	return nil
}

// EOF
