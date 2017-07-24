// Tideland Go Cells - Behaviors - Unit Tests - Key/Value
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
	"testing"
	"time"

	"github.com/tideland/golib/audit"

	"github.com/tideland/gocells/behaviors"
	"github.com/tideland/gocells/cells"
)

//--------------------
// TESTS
//--------------------

// TestKeyValueBehavior tests the key/value behavior.
func TestKeyValueBehavior(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	generator := audit.NewGenerator(audit.FixedRand())
	sigc := audit.MakeSigChan()
	env := cells.NewEnvironment("key-value-behavior")
	defer env.Stop()

	processor := func(accessor cells.EventSinkAccessor) (cells.Payload, error) {
		err := accessor.Do(func(index int, event cells.Event) error {
			var payloads cells.Payloads
			err := event.Payload().Unmarshal(&payloads)
			assert.Nil(err)
			assert.Range(len(payloads), 1, 5)
			// for _, payload := range payloads {
			// 	var value int
			// 	err = payload.Unmarshal(&value)
			// 	assert.Nil(err)
			// 	assert.Range(value, 1, 5)
			// }
			return nil
		})
		sigc <- err
		return nil, err
	}

	env.StartCell("keyvalue", behaviors.NewKeyValueBehavior(5))
	env.StartCell("collector", behaviors.NewCollectorBehavior(10, processor))
	env.Subscribe("keyvalue", "collector")

	topics := []string{"alpha", "beta", "gamma"}
	payloads := []int{1, 2, 3, 4, 5}

	for i := 0; i < 50; i++ {
		topic := generator.OneStringOf(topics...)
		payload := generator.OneIntOf(payloads...)
		env.EmitNew("keyvalue", topic, payload)
	}

	env.EmitNew("collector", cells.TopicProcess, nil)

	assert.Wait(sigc, true, 5*time.Second)
}

// EOF
