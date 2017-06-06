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
// CONSTANTS
//--------------------

const (
	// TopicSequence signals a complete sequence based on the criterion.
	TopicSequence = "sequence"
)

//--------------------
// SEQUENCE BEHAVIOR
//--------------------

// SequenceCriterion is used by the sequence behavior. It has to return
// CriterionDone when a sequence is complete, CriterionKeep when it is
// so far okay but not complete, and CriterionClear when the sequence
// doesn't match and has to be cleared.
type SequenceCriterion func(accessor cells.EventSinkAccessor) cells.CriterionMatch

// sequenceBehavior implements the sequence behavior.
type sequenceBehavior struct {
	cell    cells.Cell
	matches SequenceCriterion
	analyze cells.EventSinkAnalyzer
	sink    cells.EventSink
}

// NewSequenceBehavior creates an event sequence behavior. It checks the
// event stream for a sequence defined by the criterion. In this case an
// event containing the sequence is emitted.
func NewSequenceBehavior(criterion SequenceCriterion, analyzer cells.EventSinkAnalyzer) cells.Behavior {
	return &sequenceBehavior{
		matches: criterion,
		analyze: analyzer,
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
	case cells.TopicReset:
		b.sink.Clear()
	default:
		b.sink.Push(event)
		matches := b.matches(b.sink)
		switch matches {
		case cells.CriterionDone:
			// All done, process and start over.
			payload, err := b.analyze(b.sink)
			if err != nil {
				return err
			}
			b.cell.EmitNew(TopicSequence, payload)
			b.sink = cells.NewEventSink(0)
		case cells.CriterionKeep:
			// So far ok.
		default:
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
