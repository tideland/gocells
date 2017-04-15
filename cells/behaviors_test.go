// Tideland Go Cells - Unit Tests - Behaviors
//
// Copyright (C) 2010-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package cells_test

//--------------------
// IMPORTS
//--------------------

import (
	"time"

	"github.com/tideland/gocells/cells"
)

//--------------------
// TOPICS
//--------------------

const (
	// iterateTopic lets the test behavior iterate over its subscribers.
	iterateTopic = "iterate!"

	// ouchTopic is used in a request returning an error..
	ouchTopic = "ouch?"

	// panicTopic lets the test behavior panic to check recovering.
	panicTopic = "panic!"

	// subscribersTopic returns the current subscribers.
	subscribersTopic = "subscribers?"

	// emitTopic tells the cell to emit a test event.
	emitTopic = "emit!"

	// sleepTopic lets the cell sleep for a longer time so the queue gets full.
	sleepTopic = "sleep!"
)

//--------------------
// TEST BEHAVIORS
//--------------------

// nullBehavior does nothing.
type nullBehavior struct{}

var _ cells.Behavior = (*nullBehavior)(nil)

func (b *nullBehavior) Init(c cells.Cell) error { return nil }

func (b *nullBehavior) Terminate() error { return nil }

func (b *nullBehavior) ProcessEvent(event cells.Event) error { return nil }

func (b *nullBehavior) Recover(r interface{}) error { return nil }

// processingFunc defines the type of functions for the
// simpleBehavior.
type processingFunc func(cell cells.Cell, event cells.Event) (cells.Event, error)

// simpleBehavior allows to pass a processing function
// doing all the work.
type simpleBehavior struct {
	cell    cells.Cell
	process processingFunc
}

var _ cells.Behavior = (*simpleBehavior)(nil)

func newSimpleBehavior(pf processingFunc) *simpleBehavior {
	return &simpleBehavior{nil, pf}
}

func (b *simpleBehavior) Init(c cells.Cell) error {
	b.cell = c
	return nil
}

func (b *simpleBehavior) Terminate() error {
	return nil
}

func (b *simpleBehavior) ProcessEvent(event cells.Event) error {
	e, err := b.process(b.cell, event)
	if err != nil {
		return err
	}
	if e != nil {
		b.cell.Emit(e)
	}
	return nil
}

func (b *simpleBehavior) Recover(r interface{}) error {
	return nil
}

// collectBehavior collects and re-emits all events aand deletes
// all collected on the topic "reset".
type collectBehavior struct {
	cell cells.Cell
	sink cells.EventSink
}

var _ cells.Behavior = (*collectBehavior)(nil)

func newCollectBehavior(sink cells.EventSink) *collectBehavior {
	return &collectBehavior{nil, sink}
}

func (b *collectBehavior) Init(cell cells.Cell) error {
	b.cell = cell
	return nil
}

func (b *collectBehavior) Terminate() error {
	return nil
}

func (b *collectBehavior) ProcessEvent(event cells.Event) error {
	switch event.Topic() {
	case cells.TopicReset:
		b.sink.Clear()
	default:
		b.sink.Push(event)
		return b.cell.Emit(event)
	}
	return nil
}

func (b *collectBehavior) Recover(r interface{}) error {
	return nil
}

// recoveringFrequencyBehavior allows testing the setting
// of the recovering frequency.
type recoveringFrequencyBehavior struct {
	*collectBehavior

	number   int
	duration time.Duration
}

var _ cells.BehaviorRecoveringFrequency = (*recoveringFrequencyBehavior)(nil)

func newRecoveringFrequencyBehavior(number int, duration time.Duration, sink cells.EventSink) cells.Behavior {
	return &recoveringFrequencyBehavior{
		collectBehavior: newCollectBehavior(sink),
		number:          number,
		duration:        duration,
	}
}

func (b *recoveringFrequencyBehavior) RecoveringFrequency() (int, time.Duration) {
	return b.number, b.duration
}

// emitBehavior simply re-emits events.
type emitBehavior struct {
	cell cells.Cell
}

var _ cells.Behavior = (*emitBehavior)(nil)

func newEmitBehavior() *emitBehavior {
	return &emitBehavior{}
}

func (b *emitBehavior) Init(c cells.Cell) error {
	b.cell = c
	return nil
}

func (b *emitBehavior) Terminate() error {
	return nil
}

func (b *emitBehavior) ProcessEvent(event cells.Event) error {
	return b.cell.Emit(event)
}

func (b *emitBehavior) Recover(r interface{}) error {
	return nil
}

// EOF
