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
	sigc := audit.MakeSigChan()
	env := cells.NewEnvironment("fsm-behavior")
	defer env.Stop()

	processor := func(accessor cells.EventSinkAccessor) (cells.Payload, error) {
		eventInfos := []string{}
		accessor.Do(func(index int, event cells.Event) error {
			eventInfos = append(eventInfos, event.Topic())
			return nil
		})
		sigc <- eventInfos
		return nil, nil
	}

	lockA := lockMachine{}
	lockB := lockMachine{}

	env.StartCell("lock-a", behaviors.NewFSMBehavior(behaviors.FSMStatus{"locked", lockA.Locked, nil}))
	env.StartCell("lock-b", behaviors.NewFSMBehavior(behaviors.FSMStatus{"locked", lockB.Locked, nil}))
	env.StartCell("restorer", newRestorerBehavior())
	env.StartCell("collector-a", behaviors.NewCollectorBehavior(10, processor))
	env.StartCell("collector-b", behaviors.NewCollectorBehavior(10, processor))

	env.Subscribe("lock-a", "restorer", "collector-a")
	env.Subscribe("lock-b", "restorer", "collector-b")

	// 1st run: emit not enough and press button.
	env.EmitNew("lock-a", "coin", 20)
	env.EmitNew("lock-a", "coin", 20)
	env.EmitNew("lock-a", "coin", 20)
	env.EmitNew("lock-a", "info", nil)
	env.EmitNew("lock-a", "press-button", nil)
	env.EmitNew("lock-a", "check-cents", nil)
	env.EmitNew("restorer", "grab", nil)

	time.Sleep(time.Second)

	env.EmitNew("collector-a", cells.TopicProcess, true)
	assert.Wait(sigc, []string{"status", "coins-dropped", "cents-checked"}, time.Second)

	// 2nd run: unlock the lock and lock it again.
	env.EmitNew("lock-a", "coin", 50)
	env.EmitNew("lock-a", "coin", 20)
	env.EmitNew("lock-a", "coin", 50)
	env.EmitNew("lock-a", "info", nil)
	env.EmitNew("lock-a", "press-button", nil)

	time.Sleep(time.Second)

	env.EmitNew("collector-a", cells.TopicProcess, true)
	assert.Wait(sigc, []string{"unlocked", "status", "coins-dropped"}, time.Second)

	// 3rd run: put a screwdriwer in the lock.
	env.EmitNew("lock-a", "plastic-chip", true)

	time.Sleep(time.Second)

	env.EmitNew("collector-a", cells.TopicProcess, true)
	assert.Wait(sigc, []string{"dunno"}, time.Second)

	// 4th run: try an illegal action.
	env.EmitNew("lock-b", "screwdriver", nil)

	time.Sleep(time.Second)

	env.EmitNew("collector-b", cells.TopicProcess, true)
	assert.Wait(sigc, []string{"error"}, time.Second)
}

//--------------------
// HELPERS
//--------------------

// cents retrieves the cents out of the payload of an event.
func payloadCents(event cells.Event) int {
	var cents int
	event.Payload().Unmarshal(&cents)
	return cents
}

// lockMachine will be unlocked if enough money is inserted. After
// that it can be locked again.
type lockMachine struct {
	cents int
}

// Locked represents the locked state receiving coins.
func (m *lockMachine) Locked(cell cells.Cell, event cells.Event) behaviors.FSMStatus {
	switch event.Topic() {
	case "check-cents":
		cell.EmitNew("cents-checked", fmt.Sprintf("%s: %d", cell.ID(), m.cents))
	case "info":
		cell.EmitNew("status", fmt.Sprintf("%s: locked with %d cents", cell.ID(), m.cents))
	case "coin":
		cents := payloadCents(event)
		if cents < 1 {
			return behaviors.FSMStatus{"locked-error", nil, fmt.Errorf("do not insert buttons")}
		}
		m.cents += cents
		if m.cents > 100 {
			m.cents -= 100
			cell.EmitNew("unlocked", fmt.Sprintf("%s: unlocked", cell.ID()))
			return behaviors.FSMStatus{"unlocked", m.Unlocked, nil}
		}
	case "press-button":
		if m.cents > 0 {
			cell.EmitNew("coins-dropped", m.cents)
			m.cents = 0
		}
	case "screwdriver":
		cell.EmitNew("error", 0)
		return behaviors.FSMStatus{event.Topic(), nil, fmt.Errorf("don't try to break me")}
	default:
		cell.EmitNew("dunno", 0)
	}
	return behaviors.FSMStatus{"locked", m.Locked, nil}
}

// Unlocked represents the unlocked state receiving coins.
func (m *lockMachine) Unlocked(cell cells.Cell, event cells.Event) behaviors.FSMStatus {
	switch event.Topic() {
	case "check-cents":
		cell.EmitNew("cents-checked", fmt.Sprintf("%s: %d", cell.ID(), m.cents))
	case "info":
		cell.EmitNew("status", fmt.Sprintf("%s: unlocked with %d cents", cell.ID(), m.cents))
	case "coin":
		cents := payloadCents(event)
		cell.EmitNew("coins-returned", cents)
	case "press-button":
		if m.cents > 0 {
			cell.EmitNew("coins-dropped", m.cents)
			m.cents = 0
		}
		return behaviors.FSMStatus{"locked", m.Locked, nil}
	default:
		cell.EmitNew("dunno", 0)
	}
	return behaviors.FSMStatus{"unlocked", m.Unlocked, nil}
}

type restorerBehavior struct {
	cell  cells.Cell
	cents int
}

func newRestorerBehavior() cells.Behavior {
	return &restorerBehavior{
		cents: 0,
	}
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
	case "grab-coins":
		b.cell.EmitNew("cents", b.cents)
		b.cents = 0
	case "drop-coins":
		b.cents += payloadCents(event)
	}
	return nil
}

func (b *restorerBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
