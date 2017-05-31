// Tideland Go Cells - Behaviors - Counter
//
// Copyright (C) 2010-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors

//--------------------
// IMPORTS
//--------------------

import (
	"github.com/tideland/gocells/cells"
)

//--------------------
// COUNTER BEHAVIOR
//--------------------

// Counters is a set of named counters and their values.
type Counters map[string]int64

// Counter changes the counter values based on the received
// event. Those values will be emitted afterwards.
type Counter func(event cells.Event, counters Counters) Counters

// counterBehavior counts events based on the counter function.
type counterBehavior struct {
	cell        cells.Cell
	count Counter
	counters    Counters
}

// NewCounterBehavior creates a counter behavior based on the passed
// function. This function may increase, decrease, or set the counter
// values. Afterwards the counter values will be emitted. All values
// can be reset with the topic "reset!".
func NewCounterBehavior(counter Counter) cells.Behavior {
	return &counterBehavior{
		count: counter,
		counters: make(Counters),
	}
}

// Init the behavior.
func (b *counterBehavior) Init(c cells.Cell) error {
	b.cell = c
	return nil
}

// Terminate the behavior.
func (b *counterBehavior) Terminate() error {
	return nil
}

// ProcessEvent counts the event for the return value of the counter func
// and emits this value.
func (b *counterBehavior) ProcessEvent(event cells.Event) error {
	switch event.Topic() {
	case cells.TopicReset:
		b.counters = make(Counters)
	default:
		b.counters = b.count(event, b.counters)
		payloadValues := cells.PayloadValues{}
		for counter, value := range b.counters {
			payloadValues[counter] = value
		}
		b.cell.EmitNew(cells.TopicCounted, cells.NewPayload(payloadValues))
	}
	return nil
}

// Recover from an error.
func (b *counterBehavior) Recover(err interface{}) error {
	return nil
}

// copyCounters copies the counters for a request.
func (b *counterBehavior) copyCounters() Counters {
	copiedCounters := make(Counters)
	for key, value := range b.counters {
		copiedCounters[key] = value
	}
	return copiedCounters
}

// EOF
