// Tideland Go Cells - Behaviors - Unit Tests - Event Rate Window
//
// Copyright (C) 2010-2017 Frank Mueller / Oldenburg / Germany
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

// TestRateWindowBehavior tests the event rate window behavior.
func TestRateWindowBehavior(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	sigc := audit.MakeSigChan()
	generator := audit.NewGenerator(audit.FixedRand())
	env := cells.NewEnvironment("rate-window-behavior")
	defer env.Stop()

	topics := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "bang"}
	duration := 50 * time.Millisecond
	matcher := func(event cells.Event) (bool, error) {
		match := event.Topic() == "bang"
		return match, nil
	}
	processor := func(accessor cells.EventSinkAccessor) (cells.Payload, error) {
		first, _ := accessor.PeekFirst()
		last, _ := accessor.PeekLast()
		difference := last.Timestamp().Sub(first.Timestamp())
		sigc <- difference
		return cells.NewPayload(difference)
	}
	oncer := func(cell cells.Cell, event cells.Event) error {
		var difference time.Duration
		err := event.Payload().Unmarshal(&difference)
		assert.Nil(err)
		assert.True(difference < duration)
		assert.Equal(event.Topic(), behaviors.TopicRateWindow)
		return nil
	}

	env.StartCell("windower", behaviors.NewRateWindowBehavior(matcher, 5, duration, processor))
	env.StartCell("oncer", behaviors.NewOnceBehavior(oncer))
	env.Subscribe("windower", "oncer")

	for i := 0; i < 100; i++ {
		topic := generator.OneStringOf(topics...)
		env.EmitNew("windower", topic, nil)
		time.Sleep(time.Millisecond)
	}

	assert.WaitTested(sigc, func(v interface{}) error {
		difference := v.(time.Duration)
		if difference > 50*time.Millisecond {
			return fmt.Errorf("diff %v", difference)
		}
		return nil
	}, 5*time.Second)
}

// EOF
