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
// RATE BEHAVIOR
//--------------------

// RateWindowCriterion is used by the rate window behavior and has to return
// true, if the passed event matches a criterion for rate window measuring.
type RateWindowCriterion func(event cells.Event) bool

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
		if b.matches(event) {
			current := time.Now()
			b.timestamps.Push(current)
			if b.timestamps.Len() == b.timestamps.Cap() {
				// Collected timestamps are full, check duration.
				firstRaw, _ := b.timestamps.Peek()
				first := firstRaw.(time.Time)
				difference := current.Sub(first)
				if difference <= b.duration {
					// We've got a burst!
					b.cell.EmitNew(event.Context(), TopicRateWindow, cells.PayloadValues{
						PayloadRateWindowCount:     b.count,
						PayloadRateWindowFirstTime: first,
						PayloadRateWindowLastTime:  current,
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
