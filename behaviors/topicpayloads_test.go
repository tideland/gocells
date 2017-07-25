// Tideland Go Cells - Behaviors - Unit Tests - Topic/Payloads
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

// TestTopicPayloadsBehavior tests the topic/payloads behavior.
func TestTopicPayloadsBehavior(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	generator := audit.NewGenerator(audit.FixedRand())
	sigc := audit.MakeSigChan()
	env := cells.NewEnvironment("topic-payloads-behavior")
	defer env.Stop()

	tpProcessor := func(topic string, payloads []cells.Payload) (cells.Payload, error) {
		total := 0
		for _, payload := range payloads {
			var value int
			if err := payload.Unmarshal(&value); err != nil {
				return nil, err
			}
			total += value
		}
		return cells.NewPayload(total)
	}
	cProcessor := func(accessor cells.EventSinkAccessor) (cells.Payload, error) {
		err := accessor.Do(func(index int, event cells.Event) error {
			var total int
			if err := event.Payload().Unmarshal(&total); err != nil {
				return err
			}
			assert.Range(total, 5, 25)
			return nil
		})
		sigc <- true
		return nil, err
	}

	env.StartCell("topic-payloads", behaviors.NewTopicPayloadsBehavior(5, tpProcessor))
	env.StartCell("collector", behaviors.NewCollectorBehavior(10, cProcessor))
	env.Subscribe("topic-payloads", "collector")

	topics := []string{"alpha", "beta", "gamma"}
	payloads := []int{1, 2, 3, 4, 5}

	for i := 0; i < 50; i++ {
		topic := generator.OneStringOf(topics...)
		payload := generator.OneIntOf(payloads...)
		env.EmitNew("topic-payloads", topic, payload)
	}

	env.EmitNew("collector", cells.TopicProcess, nil)

	assert.Wait(sigc, true, 5*time.Second)
}

// EOF
