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
	"testing"
	"time"

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
	assertPayload := func(accessor cells.EventSinkAccessor, index int, value string) {
		event, ok := accessor.PeekAt(index)
		assert.True(ok)
		upperText := event.Payload().GetString("upper-text", "<none>")
		assert.Equal(upperText, value)
	}
	env := cells.NewEnvironment("mapper-behavior")
	defer env.Stop()

	mapper := func(id string, event cells.Event) (cells.Event, error) {
		text := event.Payload().GetString(cells.PayloadDefault, "")
		pv := cells.PayloadValues{
			"upper-text": strings.ToUpper(text),
		}
		payload := event.Payload().Apply(pv)
		return cells.NewEvent(event.Context(), event.Topic(), payload)
	}

	env.StartCell("mapper", behaviors.NewMapperBehavior(mapper))
	env.StartCell("collector", behaviors.NewCollectorBehavior(10))
	env.Subscribe("mapper", "collector")

	env.EmitNew(ctx, "mapper", "a", "abc")
	env.EmitNew(ctx, "mapper", "b", "def")
	env.EmitNew(ctx, "mapper", "c", "ghi")

	time.Sleep(100 * time.Millisecond)

	accessor, err := behaviors.RequestCollectedAccessor(env, "collector", cells.DefaultTimeout)
	assert.Nil(err)
	assert.Length(accessor, 3)
	assertPayload(accessor, 0, "ABC")
	assertPayload(accessor, 1, "DEF")
	assertPayload(accessor, 2, "GHI")
}

// EOF
