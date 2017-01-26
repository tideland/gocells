// Tideland Go Cells - Behaviors - Collector
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
	"context"

	"github.com/tideland/gocells/cells"
)

//--------------------
// COLLECTOR BEHAVIOR
//--------------------

// collectorBehavior collects events for debugging.
type collectorBehavior struct {
	cell cells.Cell
	sink cells.EventSink
}

// NewCollectorBehavior creates a collector behavior. It collects
// a number of events in the passed sink. The event is passed through.
// The collected events can be requested with the topic "collected?"
// and a payload waiter as default payload. A cells.EventSinkAccessor
// will set in the waiter. Additionally the collection can be resetted
// with "reset!".
func NewCollectorBehavior(sink cells.EventSink) cells.Behavior {
	return &collectorBehavior{
		sink: sink,
	}
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
		waiter, ok := event.Payload().GetWaiter(cells.DefaultPayload)
		if ok {
			accessor := cells.EventSinkAccessor(b.sink)
			waiter.Set(cells.NewPayload(accessor))
		}
	case cells.ResetTopic:
		b.sink.Clear()
	default:
		b.sink.Push(event)
		b.cell.Emit(event)
	}
	return nil
}

// Recover from an error.
func (b *collectorBehavior) Recover(err interface{}) error {
	b.sink.Clear()
	return nil
}

// RequestCollectedAccessor retrieves the accessor to the
// collected events.
func RequestCollectedAccessor(env cells.Environment, id string) (cells.EventSinkAccessor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cells.DefaultTimeout)
	defer cancel()
	waiter := cells.NewPayloadWaiter()
	err := env.EmitNewContext(ctx, id, cells.CollectedTopic, waiter)
	if err != nil {
		return nil, err
	}
	payload, err := waiter.Wait(ctx)
	if err != nil {
		return nil, err
	}
	accessor := payload.Default(nil).(cells.EventSinkAccessor)
	return accessor, nil
}

// EOF
