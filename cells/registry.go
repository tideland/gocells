// Tideland Go Cells - Registry
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

	"github.com/tideland/golib/errors"
)

//--------------------
// CELL REGISTRY
//--------------------

// registry manages the mapping of identifiers to cells.
type registry struct {
	mutex sync.RWMutex
	cells map[string]*cell
}

// newRegistry creates a new cell registry.
func newRegistry() *registry {
	return &registry{
		cells: make(map[string]*cell),
	}
}

// stop stops the registry.
func (r *registry) stop() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for _, rc := range r.cells {
		if err := rc.stop(); err != nil {
			return err
		}
	}
	r.cells = make(map[string]*cell)
	return nil
}

// startCell starts and adds a new cell to the registry if the
// ID does not already exist.
func (r *registry) startCell(env *environment, id string, behavior Behavior) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	// Check if the ID already exists.
	if _, ok := r.cells[id]; ok {
		return errors.New(ErrDuplicateID, errorMessages, id)
	}
	// Create and add.
	rc, err := newCell(env, id, behavior)
	if err != nil {
		return err
	}
	r.cells[id] = rc
	return nil
}

// stopCell stops a cell.
func (r *registry) stopCell(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	rc, ok := r.cells[id]
	if !ok {
		return errors.New(ErrInvalidID, errorMessages, id)
	}
	// Stop the cell.
	if err := rc.stop(); err != nil {
		return err
	}
	// Remove the cell from the registry.
	delete(r.cells, id)
	return nil
}

// subscribe subscribes cells to an emitter.
func (r *registry) subscribe(emitterID string, subscriberIDs ...string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	ec, ok := r.cells[emitterID]
	if !ok {
		return errors.New(ErrInvalidID, errorMessages, emitterID)
	}
	for _, subscriberID := range subscriberIDs {
		if sc, ok := r.cells[subscriberID]; ok {
			ec.subscribers.add(sc)
			sc.emitters.add(ec)
		} else {
			return errors.New(ErrInvalidID, errorMessages, subscriberID)
		}
	}
	return nil
}

// unsubscribe usubscribes cells from an emitter.
func (r *registry) unsubscribe(emitterID string, subscriberIDs ...string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	ec, ok := r.cells[emitterID]
	if !ok {
		return errors.New(ErrInvalidID, errorMessages, emitterID)
	}
	for _, subscriberID := range subscriberIDs {
		if sc, ok := r.cells[subscriberID]; ok {
			ec.subscribers.remove(subscriberID)
			sc.emitters.remove(emitterID)
		} else {
			return errors.New(ErrInvalidID, errorMessages, subscriberID)
		}
	}
	return nil
}

// subscribers returns the IDs of the subscribers of one cell.
func (r *registry) subscribers(emitterID string) ([]string, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	ec, ok := r.cells[emitterID]
	if !ok {
		return nil, errors.New(ErrInvalidID, errorMessages, emitterID)
	}
	return ec.subscribers.ids(), nil
}

// cell returns the cell with the given id.
func (r *registry) cell(id string) (*cell, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	c, ok := r.cells[id]
	if !ok {
		return nil, errors.New(ErrInvalidID, errorMessages, id)
	}
	return c, nil
}

// EOF
