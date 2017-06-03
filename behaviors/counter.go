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

// Counter analyzes the passed event and returns, which counters
// shall be incremented.
type Counter func(event cells.Event) []string

// counterBehavior counts events based on the counter function.
type counterBehavior struct {
	cell     cells.Cell
	count    Counter
	counters map[string]uint
}

// NewCounterBehavior creates a counter behavior based on the passed
// function. This function may increase, decrease, or set the counter
// values. Afterwards the counter values will be emitted. All values
// can be reset with the topic "reset!".
func NewCounterBehavior(counter Counter) cells.Behavior {
	return &counterBehavior{
		count:    counter,
		counters: map[string]uint{},
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
		b.counters = map[string]uint{}
	default:
		increments := b.count(event)
		for _, increment := range increments {
			b.counters[increment]++
		}
		values := cells.Values{}
		for counter, value := range b.counters {
			values[counter] = value
		}
		b.cell.EmitNew(cells.TopicCounted, values.Payload())
	}
	return nil
}

// Recover from an error.
func (b *counterBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
