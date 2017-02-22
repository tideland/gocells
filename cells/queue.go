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
	"sync"
)

//--------------------
// CHANNEL QUEUE
//--------------------

// channelQueue implements Queue based on simple Go channels.
type channelQueue struct {
	mutex  sync.RWMutex
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
	q.mutex.RLock()
	defer q.mutex.RUnlock()
	return nil
}

// EOF
