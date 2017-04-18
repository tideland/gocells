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

	// PayloadAggregatorValue ppoints to the aggregated value.
	PayloadAggregatorValue = "aggregator:value"
)

//--------------------
// AGGREGATOR BEHAVIOR
//--------------------

// Aggregator is a function receiving the current aggregate value
// and event and returns the next aggregate value.
type Aggregator func(value interface{}, event cells.Event) (interface{}, error)

// aggregatorBehavior implements the aggregator behavior.
type aggregatorBehavior struct {
	cell      cells.Cell
	aggregate Aggregator
	value     interface{}
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
func (b *aggregatorBehavior) Init(c cells.Cell) error {
	b.cell = c
	return nil
}

// Terminate the behavior.
func (b *aggregatorBehavior) Terminate() error {
	return nil
}

// ProcessEvent aggregates the event.
func (b *aggregatorBehavior) ProcessEvent(event cells.Event) error {
	switch event.Topic() {
	case cells.TopicReset:
		b.value = nil
	default:
		value, err := b.aggregate(b.value, event)
		if err != nil {
			return err
		}
		b.value = value
		b.cell.EmitNew(TopicAggregator, cells.PayloadValues{
			PayloadAggregatorValue: b.value,
		})
	}
	return nil
}

// Recover from an error.
func (b *aggregatorBehavior) Recover(err interface{}) error {
	b.value = nil
	return nil
}

// EOF
