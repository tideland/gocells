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

import (
	"github.com/tideland/gocells/cells"
)

//--------------------
// SEQUENCE BEHAVIOR
//--------------------

// SequenceCriterion is used by the sequence behavior and has to return
// true, if the passed event datas matches partly or totally the wanted
// sequence.
type SequenceCriterion func(accessor cells.EventSinkAccessor) CriterionMatch

// sequenceBehavior implements the sequence behavior.
type sequenceBehavior struct {
	cell    cells.Cell
	matches SequenceCriterion
	sink    cells.EventSink
}

// NewSequenceBehavior creates an event sequence behavior. It checks the
// event stream for a sequence defined by the criterion. In this case an
// event containing the sequence is emitted.
func NewSequenceBehavior(matches SequenceCriterion) cells.Behavior {
	return &sequenceBehavior{
		matches: matches,
		sink:    cells.NewEventSink(0),
	}
}

// Init implements the cells.Behavior interface.
func (b *sequenceBehavior) Init(c cells.Cell) error {
	b.cell = c
	return nil
}

// Terminate implements the cells.Behavior interface.
func (b *sequenceBehavior) Terminate() error {
	b.sink.Clear()
	return nil
}

// ProcessEvent implements the cells.Behavior interface.
func (b *sequenceBehavior) ProcessEvent(event cells.Event) error {
	switch event.Topic() {
	case ResetTopic:
		b.sink.Clear()
	default:
		b.sink.Push(event)
		matches := b.matches(b.sink)
		switch matches {
		case CriterionDone:
			// All done, emit and start over.
			b.cell.EmitNew(EventSequenceTopic, cells.PayloadValues{
				EventSequenceEventsPayload: b.sink,
			})
			b.sink = cells.NewEventSink(0)
		case CriterionKeep:
			// So far ok.
		case CriterionMove:
			// Event window has to be moved.
			b.sink.PullFirst()
		case CriterionClear:
			// Have to start from beginning.
			b.sink.Clear()
		}
	}
	return nil
}

// Recover implements the cells.Behavior interface.
func (b *sequenceBehavior) Recover(err interface{}) error {
	b.sink.Clear()
	return nil
}

// EOF
