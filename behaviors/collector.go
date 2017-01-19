// Tideland Go Cells - Behaviors - Collector
//
// Copyright (C) 2010-2016 Frank Mueller / Tideland / Oldenburg / Germany
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
// COLLECTOR BEHAVIOR
//--------------------

// collectorBehavior collects events for debugging.
type collectorBehavior struct {
	cell      cells.Cell
	max       int
	collected []EventData
}

// NewCollectorBehavior creates a collector behavior. It collects
// a configured maximum number events emitted directly or by subscription.
// The event is passed through. The collected events can be requested with
// the topic "collected?" and will be stored in the scene store named in
// the events payload. Additionally the collection can be resetted with
// "reset!".
func NewCollectorBehavior(max int) cells.Behavior {
	return &collectorBehavior{nil, max, []EventData{}}
}

// Init the behavior.
func (b *collectorBehavior) Init(c cells.Cell) error {
	b.cell = c
	return nil
}

// Terminate the behavior.
func (b *collectorBehavior) Terminate() error {
	return nil
}

// ProcessEvent collects and re-emits events.
func (b *collectorBehavior) ProcessEvent(event cells.Event) error {
	switch event.Topic() {
	case cells.CollectedTopic:
		response := make([]EventData, len(b.collected))
		copy(response, b.collected)
		if err := event.Respond(response); err != nil {
			return err
		}
	case cells.ResetTopic:
		b.collected = []EventData{}
	default:
		b.collected = append(b.collected, newEventData(event))
		if len(b.collected) > b.max {
			b.collected = b.collected[1:]
		}
		b.cell.Emit(event)
	}
	return nil
}

// Recover from an error.
func (b *collectorBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
