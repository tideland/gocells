// Tideland Go Cells - Behaviors - Key/Value Processor
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
// KEY/VALUE BEHAVIOR
//--------------------

// keyValueBehavior collects and emits payloads by topic.
type keyValueBehavior struct {
	cell      cells.Cell
	maximum   int
	collected map[string][][]byte
}

// NewKeyValueBehavior creates a behavior collecting the payloads
// of events by their topics and emits them.
func NewKeyValueBehavior(maximum int) cells.Behavior {
	return &keyValueBehavior{
		cell:      nil,
		maximum:   maximum,
		collected: make(map[string][][]byte),
	}
}

// Init the behavior.
func (b *keyValueBehavior) Init(c cells.Cell) error {
	b.cell = c
	return nil
}

// Terminate the behavior.
func (b *keyValueBehavior) Terminate() error {
	return nil
}

// ProcessEvent calls the simple processor function.
func (b *keyValueBehavior) ProcessEvent(event cells.Event) error {
	payloads := b.collected[event.Topic()]
	payloads = append(payloads, event.Payload().Bytes())
	if len(payloads) > b.maximum {
		payloads = payloads[1:]
	}
	b.collected[event.Topic()] = payloads
	return b.cell.EmitNew(event.Topic(), payloads)
}

// Recover from an error.
func (b *keyValueBehavior) Recover(err interface{}) error {
	b.collected = make(map[string][][]byte)
	return nil
}

// EOF
