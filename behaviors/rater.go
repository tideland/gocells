// Tideland Go Cells - Behaviors - Rater
//
// Copyright (C) 2010-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors

//--------------------
// IMPORTS
//--------------------

import "github.com/tideland/gocells/cells"

//--------------------
// RATER BEHAVIOR
//--------------------

// Rater is a function returning a rate for each received event.
type Rater func(event cells.Event) (float64, error)

// raterBehavior implements the rater behavior.
type raterBehavior struct {
	cell    cells.Cell
	rater   Rater
	counter int
	minRate float64
	maxRate float64
	avgRate float64
}

// NewRaterBehavior creates a behavior rating received events based
// on the passed function. This function returns a rate, minimum, maximum,
// average, and number of events are emitted.
func NewRaterBehavior(rater Rater) cells.Behavior {
	return &raterBehavior{
		rater:   rater,
		counter: 0,
		minRate: 0.0,
		maxRate: 0.0,
		avgRate: 0.0,
	}
}

// Init the behavior.
func (b *raterBehavior) Init(c cells.Cell) error {
	b.cell = c
	return nil
}

// Terminate the behavior.
func (b *raterBehavior) Terminate() error {
	return nil
}

// ProcessEvent calls the simple processor function.
func (b *raterBehavior) ProcessEvent(event cells.Event) error {
	switch event.Topic() {
	case cells.ResetTopic:
		b.counter = 0
		b.minRate = 0.0
		b.maxRate = 0.0
		b.avgRate = 0.0
	default:
		rate, err := b.rater(event)
		if err != nil {
			return err
		}
		// Calculate values.
		if b.counter == 0 {
			b.counter = 1
			b.minRate = rate
			b.maxRate = rate
			b.avgRate = rate
		} else {
			b.avgRate = (b.avgRate*float64(b.counter) + rate) / float64(b.counter+1)
			b.counter = b.counter + 1
			if rate > b.maxRate {
				b.maxRate = rate
			}
			if rate < b.minRate {
				b.minRate = rate
			}
		}
		// Emit value.
	}
	return nil
}

// Recover from an error.
func (b *raterBehavior) Recover(err interface{}) error {
	b.counter = 0
	b.minRate = 0.0
	b.maxRate = 0.0
	b.avgRate = 0.0
	return nil
}

// EOF
