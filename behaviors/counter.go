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
	"context"
	"time"

	"github.com/tideland/golib/errors"
	"github.com/tideland/golib/logger"

	"github.com/tideland/gocells/cells"
)

//--------------------
// COUNTER BEHAVIOR
//--------------------

// Counters is a set of named counters and their values.
type Counters map[string]int64

// CounterFunc is the signature of a function which analyzis
// an event and returns, which counters shall be incremented.
type CounterFunc func(id string, event cells.Event) []string

// counterBehavior counts events based on the counter function.
type counterBehavior struct {
	cell        cells.Cell
	counterFunc CounterFunc
	counters    Counters
}

// NewCounterBehavior creates a counter behavior based on the passed
// function. It increments and emits those counters named by the result
// of the counter function. The counters can be retrieved with the
// event "counters?" and a payload waiter as payload. It can be reset
// with "reset!".
func NewCounterBehavior(cf CounterFunc) cells.Behavior {
	return &counterBehavior{nil, cf, make(Counters)}
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
	case cells.TopicCounters:
		payload, ok := cells.HasWaiterPayload(event)
		if !ok {
			logger.Warningf("retrieving counters from '%s' not possible without payload waiter", b.cell.ID())
		}
		response := b.copyCounters()
		payload.GetWaiter().Set(response)
	case cells.TopicReset:
		b.counters = make(map[string]int64)
	default:
		cids := b.counterFunc(b.cell.ID(), event)
		if cids != nil {
			for _, cid := range cids {
				v, ok := b.counters[cid]
				if ok {
					b.counters[cid] = v + 1
				} else {
					b.counters[cid] = 1
				}
				topic := "counter:" + cid
				b.cell.EmitNew(topic, b.counters[cid])
			}
		}
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

//--------------------
// CONVENIENCE
//--------------------

// RequestCounterResults retrieves the results to the
// behaviors counters.
func RequestCounterResults(ctx context.Context, env cells.Environment, id string, timeout time.Duration) (Counters, error) {
	payload, err := env.Request(ctx, id, cells.TopicCounters, timeout)
	if err != nil {
		return nil, err
	}
	counters, ok := payload.GetDefault(nil).(Counters)
	if !ok {
		return nil, errors.New(ErrInvalidPayload, errorMessages, cells.PayloadDefault)
	}
	return counters, nil
}

// EOF
