// Tideland Go Cells - Cell
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
	"sync"
	"time"

	"github.com/tideland/golib/errors"
	"github.com/tideland/golib/identifier"
	"github.com/tideland/golib/logger"
	"github.com/tideland/golib/loop"
	"github.com/tideland/golib/monitoring"
)

//--------------------
// CONNECTIONS
//--------------------

// connections manages the connections to connected
// cells.
type connections struct {
	mutex sync.RWMutex
	cells []*cell
}

// newConnections creates an instance of the
// connection manager.
func newConnections() *connections {
	return &connections{
		cells: []*cell{},
	}
}

// add adds a new cell with a given identifier to the
// connections.
func (cs *connections) add(c *cell) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	for _, csc := range cs.cells {
		if csc == c {
			return
		}
	}
	cs.cells = append(cs.cells, c)
}

// remove deletes the identified cell.
func (cs *connections) remove(id string) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	remaining := []*cell{}
	for _, csc := range cs.cells {
		if csc.id != id {
			remaining = append(remaining, csc)
		}
	}
	cs.cells = remaining
}

// ids returns the identifiers of the connected cells.
func (cs *connections) ids() []string {
	var ids []string
	for _, csc := range cs.cells {
		ids = append(ids, csc.id)
	}
	return ids
}

// do executes the passed function for all connected cells
// and collects potential errors.
func (cs *connections) do(f func(c *cell) error) error {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()
	var errs []error
	for _, csc := range cs.cells {
		if err := f(csc); err != nil {
			errs = append(errs, err)
		}
	}
	switch len(errs) {
	case 0:
		return nil
	case 1:
		return errs[0]
	default:
		return errors.Collect(errs...)
	}
}

//--------------------
// CELL
//--------------------

// cell for event processing.
type cell struct {
	env                *environment
	id                 string
	measuringID        string
	queue              Queue
	behavior           Behavior
	emitters           *connections
	subscribers        *connections
	recoveringNumber   int
	recoveringDuration time.Duration
	loop               loop.Loop
}

// newCell create a new cell around a behavior.
func newCell(env *environment, id string, behavior Behavior) (*cell, error) {
	logger.Infof("cell '%s' starts", id)
	// Init cell runtime.
	c := &cell{
		env:         env,
		id:          id,
		measuringID: identifier.Identifier("cells", env.id, "cell", id),
		queue:       env.createQueue(),
		behavior:    behavior,
		emitters:    newConnections(),
		subscribers: newConnections(),
	}
	// Set configuration.
	if brf, ok := behavior.(BehaviorRecoveringFrequency); ok {
		number, duration := brf.RecoveringFrequency()
		if duration.Seconds()/float64(number) < 0.1 {
			number = minRecoveringNumber
			duration = minRecoveringDuration
		}
		c.recoveringNumber = number
		c.recoveringDuration = duration
	} else {
		c.recoveringNumber = minRecoveringNumber
		c.recoveringDuration = minRecoveringDuration
	}
	// Init behavior.
	if err := behavior.Init(c); err != nil {
		return nil, errors.Annotate(err, ErrCellInit, errorMessages, id)
	}
	// Start backend.
	c.loop = loop.GoRecoverable(c.backendLoop, c.checkRecovering, id)
	return c, nil
}

// Environment implements the Cell interface.
func (c *cell) Environment() Environment {
	return c.env
}

// ID implements the Cell interface.
func (c *cell) ID() string {
	return c.id
}

// Emit implements the Cell interface.
func (c *cell) Emit(event Event) error {
	return c.SubscribersDo(func(cs Subscriber) error {
		return cs.ProcessEvent(event)
	})
}

// EmitNew implements the Cell interface.
func (c *cell) EmitNew(topic string, payload Payload) error {
	event, err := NewEvent(topic, payload)
	if err != nil {
		return err
	}
	return c.Emit(event)
}

// ProcessEvent implements the Subscriber interface.
func (c *cell) ProcessEvent(event Event) error {
	return c.queue.Emit(event)
}

// ProcessNewEvent implements the Subscriber interface.
func (c *cell) ProcessNewEvent(topic string, payload Payload) error {
	event, err := NewEvent(topic, payload)
	if err != nil {
		return err
	}
	return c.ProcessEvent(event)
}

// SubscribersDo implements the Subscriber interface.
func (c *cell) SubscribersDo(f func(s Subscriber) error) error {
	return c.subscribers.do(func(sc *cell) error { return f(Subscriber(sc)) })
}

// stop terminates the cell.
func (c *cell) stop() error {
	// Terminate connactions to emitters and subscribers.
	c.emitters.do(func(ec *cell) error {
		ec.subscribers.remove(c.id)
		return nil
	})
	c.subscribers.do(func(sc *cell) error {
		sc.emitters.remove(c.id)
		return nil
	})
	// Stop own backend.
	err := c.loop.Stop()
	if err != nil {
		logger.Errorf("cell '%s' stopped with error: %v", c.id, err)
	} else {
		logger.Infof("cell '%s' stopped", c.id)
	}
	return err
}

// backendLoop is the backend for the processing of messages.
func (c *cell) backendLoop(l loop.Loop) error {
	totalCellsID := identifier.Identifier("cells", c.env.ID(), "total-cells")
	monitoring.IncrVariable(totalCellsID)
	defer monitoring.DecrVariable(totalCellsID)

	for {
		select {
		case <-l.ShallStop():
			return c.behavior.Terminate()
		case event := <-c.queue.Events():
			if event == nil {
				panic("received illegal nil event!")
			}
			measuring := monitoring.BeginMeasuring(c.measuringID)
			err := c.behavior.ProcessEvent(event)
			measuring.EndMeasuring()
			if err != nil {
				logger.Errorf("cell %q processed event %q with error: %v", c.id, event.Topic(), err)
				return err
			}
		}
	}
}

// checkRecovering checks if the cell may recover after a panic. It will
// signal an error and let the cell stop working if there have been 12 recoverings
// during the last minute or the behaviors Recover() signals, that it cannot
// handle the error.
func (c *cell) checkRecovering(rs loop.Recoverings) (loop.Recoverings, error) {
	logger.Warningf("recovering cell %q after error: %v", c.id, rs.Last().Reason)
	// Check frequency.
	if rs.Frequency(c.recoveringNumber, c.recoveringDuration) {
		err := errors.New(ErrRecoveredTooOften, errorMessages, rs.Last().Reason)
		logger.Errorf("recovering frequency of cell %q too high", c.id)
		return nil, err
	}
	// Try to recover.
	if err := c.behavior.Recover(rs.Last().Reason); err != nil {
		err := errors.Annotate(err, ErrEventRecovering, errorMessages, rs.Last().Reason)
		logger.Errorf("recovering of cell %q failed: %v", c.id, err)
		return nil, err
	}
	logger.Infof("successfully recovered cell %q", c.id)
	return rs.Trim(c.recoveringNumber), nil
}

// EOF
