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
)

//--------------------
// CONSTANTS
//--------------------

const (
	// TopicRateWindow signals a detected event rate window.
	TopicRateWindow = "rate-window"
)

//--------------------
// RATE BEHAVIOR
//--------------------

// RateWindowCriterion is used by the rate window behavior and has to return
// true, if the passed event matches a criterion for rate window measuring.
type RateWindowCriterion func(event cells.Event) (bool, error)

// rateWindowBehavior implements the rate window behavior.
type rateWindowBehavior struct {
	cell     cells.Cell
	sink     cells.EventSink
	matches  RateWindowCriterion
	count    int
	duration time.Duration
	process  cells.EventSinkProcessor
}

// NewRateWindowBehavior creates an event rate window behavior. It checks
// if an event matches the passed criterion. If count events match during
// duration the process function is called. Its returned payload is
// emitted as new event with topic "rate-window". A received "reset" as
// topic resets the collected matches.
func NewRateWindowBehavior(
	matches RateWindowCriterion,
	count int,
	duration time.Duration,
	process cells.EventSinkProcessor) cells.Behavior {
	return &rateWindowBehavior{
		sink:     cells.NewEventSink(count),
		matches:  matches,
		count:    count,
		duration: duration,
		process:  process,
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
		b.sink.Clear()
	default:
		ok, err := b.matches(event)
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
		b.sink.Push(event)
		if b.sink.Len() == b.count {
			// Got enough matches, check duration.
			first, _ := b.sink.PeekFirst()
			last, _ := b.sink.PeekLast()
			difference := last.Timestamp().Sub(first.Timestamp())
			if difference <= b.duration {
				// We've got a burst!
				payload, err := b.process(b.sink)
				if err != nil {
					return err
				}
				b.cell.EmitNew(TopicRateWindow, payload)
			}
			b.sink.PullFirst()
		}
	}
	return nil
}

// Recover implements the cells.Behavior interface.
func (b *rateWindowBehavior) Recover(err interface{}) error {
	b.sink = cells.NewEventSink(b.count)
	return nil
}

// EOF
