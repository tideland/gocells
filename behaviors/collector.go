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
	"github.com/tideland/gocells/cells"
)

//--------------------
// COLLECTOR BEHAVIOR
//--------------------

// collectorBehavior collects events for debugging.
type collectorBehavior struct {
	cell    cells.Cell
	sink    cells.EventSink
	process cells.EventSinkProcessor
}

// NewCollectorBehavior creates a collector behavior. It collects
// a maximum number of events, each event is passed through. If the
// maximum number is 0 it collects until the topic "reset". After
// receiving the topic "process" the processor will be called and
// the collected events will be reset afterwards.
func NewCollectorBehavior(max int, processor cells.EventSinkProcessor) cells.Behavior {
	return &collectorBehavior{
		sink:    cells.NewEventSink(max),
		process: processor,
	}
}

// Init the behavior.
func (b *collectorBehavior) Init(c cells.Cell) error {
	b.cell = c
	return nil
}

// Terminate the behavior.
func (b *collectorBehavior) Terminate() error {
	b.sink.Clear()
	return nil
}

// ProcessEvent collects and re-emits events.
func (b *collectorBehavior) ProcessEvent(event cells.Event) error {
	switch event.Topic() {
	case cells.TopicProcess:
		err := b.process(b.sink)
		if err != nil {
			return err
		}
		var clear bool
		if err := event.Payload().Unmarshal(&clear); err != nil {
			return err
		}
		if clear {
			b.sink.Clear()
		}
	case cells.TopicReset:
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

// EOF
