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
// true, if the passed event matches a criterion for combo measuring.
// The collected events help the criterion to decide, if the new one
// is a matching one. The second bool signals if a combo is full and
// and an event shall be emitted.
type ComboCriterion func(event cells.Event, collected *cells.EventDatas) (bool, bool)

// comboBehavior implements the combo behavior.
type comboBehavior struct {
	cell    cells.Cell
	matches ComboCriterion
	events  *cells.EventDatas
}

// NewComboBehavior creates an event sequence behavior. It checks the
// event stream for a combination of events defined by the criterion. In
// this case an event containing the combination is emitted.
func NewComboBehavior(matches ComboCriterion) cells.Behavior {
	return &comboBehavior{
		matches: matches,
		events:  cells.NewEventDatas(0),
	}
}

// Init implements the cells.Behavior interface.
func (b *comboBehavior) Init(c cells.Cell) error {
	b.cell = c
	return nil
}

// Terminate implements the cells.Behavior interface.
func (b *comboBehavior) Terminate() error {
	return nil
}

// ProcessEvent implements the cells.Behavior interface.
func (b *comboBehavior) ProcessEvent(event cells.Event) error {
	switch event.Topic() {
	case ResetTopic:
		b.events.Clear()
	default:
		matches, done := b.matches(event, b.events)
		if !matches {
			// No match, so reset.
			b.events.Clear()
			return nil
		}
		b.events.Add(event)
		if done {
			// All matches collected.
			b.cell.EmitNew(EventSequenceTopic, cells.PayloadValues{
				EventSequenceEventsPayload: b.events,
			})
			b.events = cells.NewEventDatas(0)
		}
	}
	return nil
}

// Recover implements the cells.Behavior interface.
func (b *comboBehavior) Recover(err interface{}) error {
	b.events.Clear()
	return nil
}

// EOF
