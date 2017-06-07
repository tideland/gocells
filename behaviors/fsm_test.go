// Tideland Go Cells - Behaviors - Unit Tests - Finite State Machine
//
// Copyright (C) 2010-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors_test

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"testing"
	"time"

	"github.com/tideland/golib/audit"

	"github.com/tideland/gocells/behaviors"
	"github.com/tideland/gocells/cells"
)

//--------------------
// TESTS
//--------------------

// TestFSMBehavior tests the finite state machine behavior.
func TestFSMBehavior(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	env := cells.NewEnvironment("fsm-behavior")
	defer env.Stop()

	processor := func(accessor cells.EventSinkAccessor) error {
		return nil
	}

	lockA := lockMachine{}
	lockB := lockMachine{}

	env.StartCell("lock-a", behaviors.NewFSMBehavior(lockA.Locked))
	env.StartCell("lock-b", behaviors.NewFSMBehavior(lockB.Locked))
	env.StartCell("restorer", newRestorerBehavior())
	env.StartCell("collector", behaviors.NewCollectorBehavior(10, processor))

	env.Subscribe("lock-a", "restorer", "collector")
	env.Subscribe("lock-b", "restorer", "collector")

	// 1st run: emit not enough and press button.
	env.EmitNew("lock-a", "coin", 20)
	env.EmitNew("lock-a", "coin", 20)
	env.EmitNew("lock-a", "coin", 20)
	env.EmitNew("lock-a", "button-press", nil)
	env.EmitNew("lock-a", "check-cents", nil)
	env.EmitNew("restorer", "grab", nil)

	// TODO 2017-06-07 Mue Add asserts.

	// 2nd run: unlock the lock and lock it again.
	env.EmitNew("lock-a", "coin", 50)
	env.EmitNew("lock-a", "coin", 20)
	env.EmitNew("lock-a", "coin", 50)
	env.EmitNew("lock-a", "info", nil)
	env.EmitNew("lock-a", "button-press", nil)

	// TODO 2017-06-07 Mue Add asserts.

	// 3rd run: put a screwdriwer in the lock.
	env.EmitNew("lock-a", "screwdriver", nil)

	// TODO 2017-06-07 Mue Add asserts.

	// 4th run: try an illegal action.
	env.EmitNew("lock-b", "chewing-gum", nil)

	// TODO 2017-06-07 Mue Add asserts.
}

//--------------------
// HELPERS
//--------------------

// cents retrieves the cents out of the payload of an event.
func payloadCents(event cells.Event) int {
	return event.Payload().GetInt(cells.PayloadDefault, -1)
}

// lockMachine will be unlocked if enough money is inserted. After
// that it can be locked again.
type lockMachine struct {
	cents int
}

// Locked represents the locked state receiving coins.
func (m *lockMachine) Locked(cell cells.Cell, event cells.Event) (behaviors.FSMState, string, error) {
	switch event.Topic() {
	case "check-cents":
		cell.EmitNew("cents", fmt.Sprintf("check-cents %s: %d", cell.ID(), m.cents))
		return m.Locked, "locked", nil
	case "info?":
		info := fmt.Sprintf("state 'locked' with %d cents", m.cents)
		payload, ok := cells.HasWaiterPayload(event)
		if ok {
			payload.GetWaiter().Set(info)
		}
		return m.Locked, nil
	case "coin!":
		cents := payloadCents(event)
		if cents < 1 {
			return nil, fmt.Errorf("do not insert buttons")
		}
		m.cents += cents
		if m.cents > 100 {
			m.cents -= 100
			return m.Unlocked, nil
		}
		return m.Locked, nil
	case "button-press!":
		if m.cents > 0 {
			cell.Environment().EmitNew(event.Context(), "restorer", "drop!", m.cents)
			m.cents = 0
		}
		return m.Locked, nil
	case "screwdriver!":
		// Allow a screwdriver to bring the lock into an undefined state.
		return nil, nil
	}
	return m.Locked, fmt.Errorf("illegal topic in state 'locked': %s", event.Topic())
}

// Unlocked represents the unlocked state receiving coins.
func (m *lockMachine) Unlocked(cell cells.Cell, event cells.Event) (behaviors.FSMState, error) {
	switch event.Topic() {
	case "cents?":
		payload, ok := cells.HasWaiterPayload(event)
		if ok {
			payload.GetWaiter().Set(m.cents)
		}
		return m.Unlocked, nil
	case "info?":
		info := fmt.Sprintf("state 'unlocked' with %d cents", m.cents)
		payload, ok := cells.HasWaiterPayload(event)
		if ok {
			payload.GetWaiter().Set(info)
		}
		return m.Unlocked, nil
	case "coin!":
		cents := payloadCents(event)
		cell.EmitNew(event.Context(), "return", cents)
		return m.Unlocked, nil
	case "button-press!":
		cell.Environment().EmitNew(event.Context(), "restorer", "drop!", m.cents)
		m.cents = 0
		return m.Locked, nil
	}
	return m.Unlocked, fmt.Errorf("illegal topic in state 'unlocked': %s", event.Topic())
}

type restorerBehavior struct {
	cell  cells.Cell
	cents int
}

func newRestorerBehavior() cells.Behavior {
	return &restorerBehavior{nil, 0}
}

func (b *restorerBehavior) Init(c cells.Cell) error {
	b.cell = c
	return nil
}

func (b *restorerBehavior) Terminate() error {
	return nil
}

func (b *restorerBehavior) ProcessEvent(event cells.Event) error {
	switch event.Topic() {
	case "grab!":
		cents := b.cents
		b.cents = 0
		payload, ok := cells.HasWaiterPayload(event)
		if ok {
			payload.GetWaiter().Set(cents)
		}
	case "drop!":
		b.cents += payloadCents(event)
	}
	return nil
}

func (b *restorerBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
