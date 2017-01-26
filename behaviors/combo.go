// Tideland Go Cells - Behaviors - Combo
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
// SEQUENCE BEHAVIOR
//--------------------

// ComboCriterion is used by the combo behavior and has to return
// true as first, if the passed events match a criterion for combo
// measuring. If the first result is false the event will be dropped
// again. The second bool signals if a combo is full and an event
// shall be emitted.
type ComboCriterion func(accessor cells.EventSinkAccessor) (bool, bool)

// comboBehavior implements the combo behavior.
type comboBehavior struct {
	cell    cells.Cell
	matches ComboCriterion
	sink    cells.EventSink
}

// NewComboBehavior creates an event sequence behavior. It checks the
// event stream for a combination of events defined by the criterion. In
// this case an event containing the combination is emitted.
func NewComboBehavior(matches ComboCriterion) cells.Behavior {
	return &comboBehavior{
		matches: matches,
		sink:    cells.NewEventSink(0),
	}
}

// Init implements the cells.Behavior interface.
func (b *comboBehavior) Init(c cells.Cell) error {
	b.cell = c
	return nil
}

// Terminate implements the cells.Behavior interface.
func (b *comboBehavior) Terminate() error {
	b.sink.Clear()
	return nil
}

// ProcessEvent implements the cells.Behavior interface.
func (b *comboBehavior) ProcessEvent(event cells.Event) error {
	switch event.Topic() {
	case ResetTopic:
		b.sink.Clear()
	default:
		b.sink.Push(event)
		matches, done := b.matches(b.sink)
		if !matches {
			// No match, so pull last.
			b.sink.PullLast()
			return nil
		}
		if done {
			// All matches collected.
			b.cell.EmitNew(EventComboTopic, cells.PayloadValues{
				EventComboEventsPayload: b.sink,
			})
			b.sink = cells.NewEventSink(0)
		}
	}
	return nil
}

// Recover implements the cells.Behavior interface.
func (b *comboBehavior) Recover(err interface{}) error {
	b.sink.Clear()
	return nil
}

// EOF
