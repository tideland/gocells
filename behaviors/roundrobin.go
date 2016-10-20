// Tideland Go Cells - Behaviors - Round-Robin
//
// Copyright (C) 2010-2016 Frank Mueller / Tideland / Oldenburg / Germany
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
// ROUND ROBIN BEHAVIOR
//--------------------

// roundRobinBehavior emit the received events round robin to its
// subscribers in a very simple way.
type roundRobinBehavior struct {
	cell    cells.Cell
	current int
}

// NewRoundRobinBehavior creates a behavior emitting the received events to
// its subscribers in a very simple way. Subscriptions or unsubscriptions
// during runtime may influence the order.
func NewRoundRobinBehavior() cells.Behavior {
	return &roundRobinBehavior{nil, 0}
}

// Init the behavior.
func (b *roundRobinBehavior) Init(c cells.Cell) error {
	b.cell = c
	return nil
}

// Terminate the behavior.
func (b *roundRobinBehavior) Terminate() error {
	return nil
}

// ProcessEvent emits the event round robin to the subscribers.
func (b *roundRobinBehavior) ProcessEvent(event cells.Event) error {
	loopCurrent := 0
	err := b.cell.SubscribersDo(func(s cells.Subscriber) error {
		if loopCurrent == b.current {
			if err := s.ProcessEvent(event); err != nil {
				return err
			}
		}
		loopCurrent++
		return nil
	})
	if err != nil {
		return err
	}
	if b.current < loopCurrent-1 {
		b.current++
	} else {
		b.current = 0
	}
	return nil
}

// Recover from an error.
func (b *roundRobinBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
