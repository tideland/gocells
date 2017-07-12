// Tideland Go Cells - Behaviors - Countdown
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
// COUNTDOWN BEHAVIOR
//--------------------

// Zeroer is called when the countdown reaches zero. The collected
// events are passed, the returned event will be emitted, and the
// returned number sets a new countdown.
type Zeroer func(accessor cells.EventSinkAccessor) (cells.Event, int, error)

// countdownBehavior counts events based on the counter function.
type countdownBehavior struct {
	cell   cells.Cell
	sink   cells.EventSink
	t      int
	zeroer Zeroer
}

// NewCountdownBehavior creates a countdown behavior based on the passed
// t value and zeroer function.
func NewCountdownBehavior(t int, zeroer Zeroer) cells.Behavior {
	return &countdownBehavior{
		sink:   cells.NewEventSink(t),
		t:      t,
		zeroer: zeroer,
	}
}

// Init the behavior.
func (b *countdownBehavior) Init(cell cells.Cell) error {
	b.cell = cell
	return nil
}

// Terminate the behavior.
func (b *countdownBehavior) Terminate() error {
	return nil
}

// ProcessEvent counts the event for the return value of the counter func
// and emits this value.
func (b *countdownBehavior) ProcessEvent(event cells.Event) error {
	if b.t <= 0 {
		return nil
	}
	sl, err := b.sink.Push(event)
	if err != nil {
		return err
	}
	if sl == b.t {
		// T-0, call the zeroer, set t, and emit event.
		e, t, err := b.zeroer(b.sink)
		if err != nil {
			return err
		}
		b.sink.Clear()
		b.t = t
		return b.cell.Emit(e)
	}
	return nil
}

// Recover from an error.
func (b *countdownBehavior) Recover(err interface{}) error {
	return b.sink.Clear()
}

// EOF
