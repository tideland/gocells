// Tideland Go Cells - Behaviors - Rate Window
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
	"github.com/tideland/golib/collections"
)

//--------------------
// CONSTANTS
//--------------------

const (
	// TopicRateWindow signals a detected event rate window.
	TopicRateWindow = "rate-window"

	// PayloadRateWindowCount contains the number of matching events.
	PayloadRateWindowCount = "rate-window:count"

	// PayloadRateWindowFirstTime contains the first time a
	// matching event has been detected.
	PayloadRateWindowFirstTime = "rate-window:first:time"

	// PayloadRateWindowLastTime contains the last time a
	// matching event has been detected.
	PayloadRateWindowLastTime = "rate-window:last:time"
)

//--------------------
// RATE BEHAVIOR
//--------------------

// RateWindowCriterion is used by the rate window behavior and has to return
// true, if the passed event matches a criterion for rate window measuring.
type RateWindowCriterion func(event cells.Event) (bool, error)

// RateWindow describes the time window of events matching the defined criterion.
// It contains the number of events, and the times of the first and last ones.
type RateWindow struct {
	Count int
	First time.Time
	Last  time.Time
}

// rateWindowBehavior implements the rate window behavior.
type rateWindowBehavior struct {
	cell       cells.Cell
	matches    RateWindowCriterion
	count      int
	duration   time.Duration
	timestamps collections.RingBuffer
}

// NewRateWindowBehavior creates an event rate window behavior. It checks
// if an event matches the passed criterion. If count events match during
// duration an according event containing the first time, the last time,
// and the number of matches is emitted. A "reset!" as topic resets the
// collected matches.
func NewRateWindowBehavior(matches RateWindowCriterion, count int, duration time.Duration) cells.Behavior {
	return &rateWindowBehavior{
		matches:    matches,
		count:      count,
		duration:   duration,
		timestamps: collections.NewRingBuffer(count),
	}
}

// Init implements the cells.Behavior interface.
func (b *rateWindowBehavior) Init(c cells.Cell) error {
	b.cell = c
	return nil
}

// Terminate implements the cells.Behavior interface.
func (b *rateWindowBehavior) Terminate() error {
	return nil
}

// ProcessEvent implements the cells.Behavior interface.
func (b *rateWindowBehavior) ProcessEvent(event cells.Event) error {
	switch event.Topic() {
	case cells.TopicReset:
		b.timestamps = collections.NewRingBuffer(b.count)
	default:
		ok, err := b.matches(event)
		if err != nil {
			return err
		}
		if ok {
			current := event.Timestamp()
			b.timestamps.Push(current)
			if b.timestamps.Len() == b.timestamps.Cap() {
				// Collected timestamps are full, check duration.
				firstRaw, _ := b.timestamps.Peek()
				first := firstRaw.(time.Time)
				difference := current.Sub(first)
				if difference <= b.duration {
					// We've got a burst!
					b.cell.EmitNew(TopicRateWindow, RateWindow{
						Count: b.count,
						First: first,
						Last:  current,
					})
				}
			}
		}
	}
	return nil
}

// Recover implements the cells.Behavior interface.
func (b *rateWindowBehavior) Recover(err interface{}) error {
	b.timestamps = collections.NewRingBuffer(b.count)
	return nil
}

// EOF
