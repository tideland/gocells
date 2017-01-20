// Tideland Go Cells - Behaviors - Sequence
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
// SEQUENCE BEHAVIOR
//--------------------

// SequenceCriterion is used by the sequence behavior and has to return
// true, if the passed event matches a criterion for sequence measuring.
// The collected events help the criterion to decide, if the new one
// is a matching one. The second bool signals if a sequence is full and
// and event shall be emitted.
type SequenceCriterion func(event cells.Event, collected *cells.EventDatas) (bool, bool)

// sequenceBehavior implements the sequence behavior.
type sequenceBehavior struct {
	cell    cells.Cell
	matches SequenceCriterion
	events  *cells.EventDatas
}

// NewSequenceBehavior creates an event sequence behavior. It ...
func NewSequenceBehavior(matches SequenceCriterion) cells.Behavior {
	return &sequenceBehavior{
		matches: matches,
		events:  cells.NewEventDatas(0),
	}
}

// Init implements the cells.Behavior interface.
func (b *sequenceBehavior) Init(c cells.Cell) error {
	b.cell = c
	return nil
}

// Terminate implements the cells.Behavior interface.
func (b *sequenceBehavior) Terminate() error {
	return nil
}

// ProcessEvent implements the cells.Behavior interface.
func (b *sequenceBehavior) ProcessEvent(event cells.Event) error {
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
func (b *sequenceBehavior) Recover(err interface{}) error {
	b.events.Clear()
	return nil
}

// EOF
