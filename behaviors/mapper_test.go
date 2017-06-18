// Tideland Go Cells - Behaviors - Unit Tests - Mapper
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
	"strings"
	"sync"
	"testing"

	"github.com/tideland/golib/audit"

	"github.com/tideland/gocells/behaviors"
	"github.com/tideland/gocells/cells"
)

//--------------------
// TESTS
//--------------------

// TestMapperBehavior tests the mapping of events.
func TestMapperBehavior(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	env := cells.NewEnvironment("mapper-behavior")
	defer env.Stop()

	mapper := func(event cells.Event) (cells.Event, error) {
		text := event.Payload().String()
		return cells.NewEvent(event.Topic(), strings.ToUpper(text))
	}

	var wg sync.WaitGroup

	processor := func(cell cells.Cell, event cells.Event) error {
		wg.Done()
		text := event.Payload().String()
		switch event.Topic() {
		case "a":
			assert.Equal(text, "ABC")
		case "b":
			assert.Equal(text, "DEF")
		case "c":
			assert.Equal(text, "GHI")
		default:
			assert.Fail("mapper didn't work: %s = %s", event.Topic(), text)
		}
		return nil
	}

	env.StartCell("mapper", behaviors.NewMapperBehavior(mapper))
	env.StartCell("processor", behaviors.NewSimpleProcessorBehavior(processor))
	env.Subscribe("mapper", "processor")

	wg.Add(3)
	env.EmitNew("mapper", "a", "abc")
	env.EmitNew("mapper", "b", "def")
	env.EmitNew("mapper", "c", "ghi")
	wg.Wait()
}

// EOF
