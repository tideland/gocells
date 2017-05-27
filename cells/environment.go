// Tideland Go Cells - Environment
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
	"runtime"

	"github.com/tideland/golib/identifier"
	"github.com/tideland/golib/logger"
)

//--------------------
// ENVIRONMENT
//--------------------

// Environment implements the Environment interface.
type environment struct {
	id    string
	cells *registry
}

// NewEnvironment creates a new environment.
func NewEnvironment(idParts ...interface{}) Environment {
	var id string
	if len(idParts) == 0 {
		id = identifier.NewUUID().String()
	} else {
		id = identifier.Identifier(idParts...)
	}
	env := &environment{
		id:    id,
		cells: newRegistry(),
	}
	runtime.SetFinalizer(env, (*environment).Stop)
	logger.Infof("cells environment %q started", env.ID())
	return env
}

// ID implements the Environment interface.
func (env *environment) ID() string {
	return env.id
}

// StartCell implements the Environment interface.
func (env *environment) StartCell(id string, behavior Behavior) error {
	return env.cells.startCell(env, id, behavior)
}

// StopCell implements the Environment interface.
func (env *environment) StopCell(id string) error {
	return env.cells.stopCell(id)
}

// HasCell implements the Environment interface.
func (env *environment) HasCell(id string) bool {
	_, err := env.cells.cell(id)
	return err == nil
}

// Subscribe implements the Environment interface.
func (env *environment) Subscribe(emitterID string, subscriberIDs ...string) error {
	return env.cells.subscribe(emitterID, subscriberIDs...)
}

// Subscribers implements the Environment interface.
func (env *environment) Subscribers(id string) ([]string, error) {
	return env.cells.subscribers(id)
}

// Unsubscribe implements the Environment interface.
func (env *environment) Unsubscribe(emitterID string, subscriberIDs ...string) error {
	return env.cells.unsubscribe(emitterID, subscriberIDs...)
}

// Emit implements the Environment interface.
func (env *environment) Emit(id string, event Event) error {
	c, err := env.cells.cell(id)
	if err != nil {
		return err
	}
	return c.ProcessEvent(event)
}

// EmitNew implements the Environment interface.
func (env *environment) EmitNew(id, topic string, payload Payload) error {
	event, err := NewEvent(topic, payload)
	if err != nil {
		return err
	}
	return env.Emit(id, event)
}

// Stop implements the Environment interface.
func (env *environment) Stop() error {
	runtime.SetFinalizer(env, nil)
	if err := env.cells.stop(); err != nil {
		return err
	}
	logger.Infof("cells environment %q terminated", env.ID())
	return nil
}

// createQueue is a factory for the configured type of queues.
func (env *environment) createQueue() Queue {
	// Currently only one queue type is supported.
	return newInMemoryQueue()
}

// EOF
