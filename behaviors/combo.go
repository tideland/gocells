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
// CONSTANTS
//--------------------

const (
	// TopicComboComplete is used for events emitted by the combo in
	// case of a complete combination.
	TopicComboComplete = "combo-complete"
)

//--------------------
// SEQUENCE BEHAVIOR
//--------------------

// ComboCriterion is used by the combo behavior. It has to return
// CriterionDone when a combination is complete, CriterionKeep when it
// is so far okay but not complete, CriterionDropFirst when the first
// event shall be dropped, CriterionDropLast when the last event shall
// be dropped, and CriterionClear when the collected events have
// to be cleared for starting over.
type ComboCriterion func(accessor cells.EventSinkAccessor) cells.CriterionMatch

// ComboPayload is called in case of a complete combenation. It has to
// return the payload emitted with the complete event.
type ComboPayload func(accessor cells.EventSinkAccessor) cells.Payload

// comboBehavior implements the combo behavior.
type comboBehavior struct {
	cell    cells.Cell
	matches ComboCriterion
	payload ComboPayload
	sink    cells.EventSink
}

// NewComboBehavior creates an event sequence behavior. It checks the
// event stream for a combination of events defined by the criterion. In
// this case an event containing the combination is emitted.
func NewComboBehavior(matches ComboCriterion, payload ComboPayload) cells.Behavior {
	return &comboBehavior{
		matches: matches,
		payload: payload,
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
	case cells.TopicReset:
		b.sink.Clear()
	default:
		b.sink.Push(event)
		matches := b.matches(b.sink)
		switch matches {
		case cells.CriterionDone:
			// All done, emit and start over.
			b.cell.EmitNew(TopicComboComplete, b.payload(b.sink))
			b.sink = cells.NewEventSink(0)
		case cells.CriterionKeep:
			// So far ok.
		case cells.CriterionDropFirst:
			// First event doesn't match.
			b.sink.PullFirst()
		case cells.CriterionDropLast:
			// First event doesn't match.
			b.sink.PullLast()
		default:
			// Have to start from beginning.
			b.sink.Clear()
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
