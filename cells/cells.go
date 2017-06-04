// Tideland Go Cells
//
// Copyright (C) 2010-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package cells

//--------------------
// IMPORTS
//--------------------

import (
	"time"
)

//--------------------
// ENVIRONMENT
//--------------------

// Environment is a set of networked cells.
type Environment interface {
	// ID returns the ID of the environment. When creating the environment
	// the ID can by set manually or is generated automatically.
	ID() string

	// StartCell starts a new cell with a given ID and its behavior.
	StartCell(id string, behavior Behavior) error

	// StopCell stops and removes the cell with the given ID.
	StopCell(id string) error

	// HasCell returns true if the cell with the given ID exists.
	HasCell(id string) bool

	// Subscribe assigns cells as receivers of the emitted
	// events of the first cell.
	Subscribe(emitterID string, subscriberIDs ...string) error

	// Subscribers returns the subscribers of the passed ID.
	Subscribers(id string) ([]string, error)

	// Unsubscribe removes the assignment of emitting und subscribed cells.
	Unsubscribe(emitterID string, unsubscriberIDs ...string) error

	// Emit emits an event to the cell with a given ID.
	Emit(id string, event Event) error

	// EmitNew creates an event and emits it to the cell
	// with a given ID.
	EmitNew(id, topic string, payload interface{}) error

	// Stop manages the proper finalization of an environment.
	Stop() error
}

//--------------------
// SUBSCRIBER
//--------------------

// Subscriber describes a subscriber cell for an emitting cell.
type Subscriber interface {
	// ID returns the ID of the subscriber.
	ID() string

	// ProcessEvent tells the subscriber to process an event.
	ProcessEvent(event Event) error

	// ProcessNewEvent creates an event and tells the subscriber to process it.
	ProcessNewEvent(topic string, payload Payload) error
}

//--------------------
// QUEUE
//--------------------

// Queue transports events to its cell.
type Queue interface {
	// Emit emits an event into the queue.
	Emit(event Event) error

	// Events delivers events via a channel.
	Events() <-chan Event

	// Close closes the queue.
	Close() error
}

//--------------------
// CELL
//--------------------

// Cell gives a behavior access to the cell which is running the behavior.
type Cell interface {
	// Environment returns the environment the cell is running in.
	Environment() Environment

	// ID returns the ID used during the start of the cell. The same behavior
	// can be started multiple times but has to use different IDs.
	ID() string

	// Emit emits an event to all subscribers of a cell.
	Emit(event Event) error

	// EmitNew creates an event and emits it to all subscribers of a cell.
	EmitNew(topic string, payload interface{}) error

	// SubscribersDo calls the passed function for each subscriber.
	SubscribersDo(f func(s Subscriber) error) error
}

//--------------------
// BEHAVIOR
//--------------------

// Behavior is the interface that has to be implemented
// for the usage inside of cells.
type Behavior interface {
	// Init is called to initialize the behavior inside the environment.
	// The passed context allows the behavior to interact with this
	// environment and to emit events to subscribers during ProcessEvent().
	// So if this is needed the context should be stored inside the behavior.
	Init(c Cell) error

	// Terminate is called when a cell is stopped.
	Terminate() error

	// ProcessEvent is called to process the passed event. If during this
	// processing one or more events shall be emitted to the subscribers
	// the context passed during Init() is needed.
	ProcessEvent(event Event) error

	// Recover is called in case of an error or panic during the processing
	// of an event. Here the behavior can check if it can recover and establish
	// a valid state. If it's not possible the implementation has to return
	// an error documenting the reason.
	Recover(r interface{}) error
}

// BehaviorRecoveringFrequency is an additional optional interface for a behavior to
// set the allowed frequency for recoverings by returning the according number and
// duration (will never below once per second).
type BehaviorRecoveringFrequency interface {
	RecoveringFrequency() (int, time.Duration)
}

// EOF
