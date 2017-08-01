// Tideland Go Cells - Behaviors - Router
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
// ROUTER BEHAVIOR
//--------------------

// Router is a function type determining which subscribed
// cells shall receive the event.
type Router func(emitterID, subscriberID string, event cells.Event) (bool, error)

// routerBehavior check for each received event which subscriber will
// get it based on the router function.
type routerBehavior struct {
	cell       cells.Cell
	shallRoute Router
}

// NewRouterBehavior creates a router behavior using the passed function
// to determine to which subscriber the received event will be emitted.
func NewRouterBehavior(router Router) cells.Behavior {
	return &routerBehavior{
		shallRoute: router,
	}
}

// Init the behavior.
func (b *routerBehavior) Init(c cells.Cell) error {
	b.cell = c
	return nil
}

// Terminate the behavior.
func (b *routerBehavior) Terminate() error {
	return nil
}

// ProcessEvent emits the event to those ids returned by the router
// function.
func (b *routerBehavior) ProcessEvent(event cells.Event) error {
	return b.cell.SubscribersDo(func(s cells.Subscriber) error {
		ok, err := b.shallRoute(b.cell.ID(), s.ID(), event)
		if err != nil {
			return err
		}
		if ok {
			return s.ProcessEvent(event)
		}
		return nil
	})
}

// Recover from an error.
func (b *routerBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
