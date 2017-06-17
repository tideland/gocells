// Tideland Go Cells - Behaviors - Aggregator
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
// CONSTANTS
//--------------------

const (
	// TopicAggregator is used for events emitted by the aggregator behavior.
	TopicAggregator = "aggregator"
)

//--------------------
// AGGREGATOR BEHAVIOR
//--------------------

// AggregatorFunc is a function receiving the current aggregated payload
// and event and returns the next aggregated payload.
type Aggregator func(payload cells.Payload, event cells.Event) (cells.Payload, error)

// aggregatorBehavior implements the aggregator behavior.
type aggregatorBehavior struct {
	cell      cells.Cell
	aggregate Aggregator
	payload   cells.Payload
}

// NewAggregatorBehavior creates a behavior aggregating the received events
// and emits events with the new aggregate. A "reset!" topic resets the
// aggregate to nil again.
func NewAggregatorBehavior(aggregator Aggregator) cells.Behavior {
	return &aggregatorBehavior{
		aggregate: aggregator,
	}
}

// Init the behavior.
func (b *aggregatorBehavior) Init(cell cells.Cell) error {
	b.cell = cell
	return nil
}

// Terminate the behavior.
func (b *aggregatorBehavior) Terminate() error {
	return nil
}

// ProcessEvent aggregates the event.
func (b *aggregatorBehavior) ProcessEvent(event cells.Event) error {
	switch event.Topic() {
	case cells.TopicStatus:
		statusCell := event.Payload().String()
		b.cell.Environment().EmitNew(statusCell, b.cell.ID(), b.payload)
	case cells.TopicReset:
		b.payload = nil
	default:
		payload, err := b.aggregate(b.payload, event)
		if err != nil {
			return err
		}
		b.payload = payload
		b.cell.EmitNew(TopicAggregator, payload)
	}
	return nil
}

// Recover from an error.
func (b *aggregatorBehavior) Recover(err interface{}) error {
	b.payload = nil
	return nil
}

// EOF
