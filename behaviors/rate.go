// Tideland Go Cells - Behaviors - Rate
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
	"time"

	"github.com/tideland/gocells/cells"
)

//--------------------
// RATE BEHAVIOR
//--------------------

// RateCriterion is used by the rate behavior and has to return true, if
// the passed event matches a criterion for rate measuring.
type RateCriterion func(event cells.Event) bool

// rateBehavior calculates the average rate of events matching a criterion.
type rateBehavior struct {
	cell      cells.Cell
	matches   RateCriterion
	count     int
	last      time.Time
	durations []time.Duration
}

// NewRateBehavior creates an even rate measuiring behavior. Each time the
// criterion function returns true for a received event the duration between
// this and the last one is calculated and emitted together with the timestamp.
// Additionally a moving average, lowest, and highest duration is calculated
// and emitted too. A "reset!" as topic resets the stored values.
func NewRateBehavior(matches RateCriterion, count int) cells.Behavior {
	return &rateBehavior{nil, matches, count, time.Now(), []time.Duration{}}
}

// Init implements the cells.Behavior interface.
func (b *rateBehavior) Init(c cells.Cell) error {
	b.cell = c
	return nil
}

// Terminate implements the cells.Behavior interface.
func (b *rateBehavior) Terminate() error {
	return nil
}

// ProcessEvent implements the cells.Behavior interface.
func (b *rateBehavior) ProcessEvent(event cells.Event) error {
	switch event.Topic() {
	case TopicReset:
		b.last = time.Now()
		b.durations = []time.Duration{}
	default:
		if b.matches(event) {
			current := time.Now()
			duration := current.Sub(b.last)
			b.last = current
			b.durations = append(b.durations, duration)
			if len(b.durations) > b.count {
				b.durations = b.durations[1:]
			}
			total := 0 * time.Nanosecond
			low := 0x7FFFFFFFFFFFFFFF * time.Nanosecond
			high := 0 * time.Nanosecond
			for _, d := range b.durations {
				total += d
				if d < low {
					low = d
				}
				if d > high {
					high = d
				}
			}
			avg := total / time.Duration(len(b.durations))
			return b.cell.EmitNew(event.Context(), TopicRate, cells.PayloadValues{
				PayloadRateTime:     current,
				PayloadRateDuration: duration,
				PayloadRateAverage:  avg,
				PayloadRateHigh:     high,
				PayloadRateLow:      low,
			})
		}
	}
	return nil
}

// Recover implements the cells.Behavior interface.
func (b *rateBehavior) Recover(err interface{}) error {
	b.last = time.Now()
	b.durations = []time.Duration{}
	return nil
}

// EOF
