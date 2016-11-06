// Tideland Go Cells - Behaviors - Rate
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
	"time"

	"github.com/tideland/gocells/cells"
)

//--------------------
// RATE BEHAVIOR
//--------------------

// RateCriterion is used by the rate behavior and has to return true, if
// the passed evend matches a criterion for rate measuring.
type RateCriterion func(event cells.Event) bool

// rateBehavior calculates the average rate of event matching a criterion.
type rateBehavior struct {
	cell       cells.Cell
	matches    RateCriterion
	count      int
	timestamps []time.Time
}

// NewRateBehavior creates an even rate measuiring behavior. Each time the
// criterion function returns true for a received event a timestamp is
// stored and a moving average of the times between these events is emitted.
func NewRateBehavior(matches RateCriterion, count int) cells.Behavior {
	return &rateBehavior{nil, matches, count, []time.Time{}}
}

// Init the behavior.
func (b *rateBehavior) Init(c cells.Cell) error {
	b.cell = c
	return nil
}

// Terminate the behavior.
func (b *rateBehavior) Terminate() error {
	return nil
}

// ProcessEvent collects and re-emits events.
func (b *rateBehavior) ProcessEvent(event cells.Event) error {
	if b.matches(event) {
		b.timestamps = append(b.timestamps, time.Now())
		if len(b.timestamps) > b.count {
			b.timestamps = b.timestamps[1:]
		}
		d := b.timestamps[len(b.timestamps)-1].Sub(b.timestamps[0])
		avg := d / time.Duration(len(b.timestamps))
		hi := 0 * time.Nanosecond
		lo := d
		for i := 1; i < len(b.timestamps); i++ {
			d = b.timestamps[i].Sub(b.timestamps[i-1])
			if d > hi {
				hi = d
			}
			if d < lo {
				lo = d
			}
		}
		return b.cell.EmitNew(EventRateTopic, cells.PayloadValues{
			EventRateAveragePayload: avg,
			EventRateHighPayload:    hi,
			EventRateLowPayload:     lo,
		})
	}
	return nil
}

// Recover from an error.
func (b *rateBehavior) Recover(err interface{}) error {
	b.timestamps = []time.Time{}
	return nil
}

// EOF
