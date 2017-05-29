// Tideland Go Cells - Unit Tests
//
// Copyright (C) 2010-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package cells_test

//--------------------
// IMPORTS
//--------------------

import (
	"testing"
	"time"

	"github.com/tideland/golib/audit"
	"github.com/tideland/golib/errors"
	"github.com/tideland/golib/monitoring"

	"github.com/tideland/gocells/cells"
)

//--------------------
// TESTS
//--------------------

// TestEnvironment tests general environment methods.
func TestEnvironment(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)

	envOne := cells.NewEnvironment("part", 1, "of", "env", "ONE")
	defer envOne.Stop()

	id := envOne.ID()
	assert.Equal(id, "part:1:of:env:one")

	envTwo := cells.NewEnvironment("environment TWO")
	defer envTwo.Stop()

	id = envTwo.ID()
	assert.Equal(id, "environment-two")
}

// TestEnvironmentStartStopCell tests starting, checking and
// stopping of cells.
func TestEnvironmentStartStopCell(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)

	env := cells.NewEnvironment("start-stop")
	defer env.Stop()

	sink := cells.NewEventSink(0)
	err := env.StartCell("foo", newCollectBehavior(sink))
	assert.Nil(err)

	hasFoo := env.HasCell("foo")
	assert.True(hasFoo)

	err = env.StopCell("foo")
	assert.Nil(err)
	hasFoo = env.HasCell("foo")
	assert.False(hasFoo)

	hasBar := env.HasCell("bar")
	assert.False(hasBar)
	err = env.StopCell("bar")
	assert.True(errors.IsError(err, cells.ErrInvalidID))
	hasBar = env.HasCell("bar")
	assert.False(hasBar)
}

// TestBehaviorRecoveringFrequency tests the setting of
// the recovering frequency.
func TestBehaviorRecoveringFrequency(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)

	env := cells.NewEnvironment("recovering-frequency")
	defer env.Stop()

	sink := cells.NewEventSink(0)
	err := env.StartCell("negative", newRecoveringFrequencyBehavior(-1, time.Second, sink))
	assert.Nil(err)
	ci := cells.InspectCell(env, "negative")
	assert.Equal(ci.RecoveringNumber(), cells.MinRecoveringNumber)
	assert.Equal(ci.RecoveringDuration(), cells.MinRecoveringDuration)

	err = env.StartCell("low", newRecoveringFrequencyBehavior(10, time.Millisecond, sink))
	assert.Nil(err)
	ci = cells.InspectCell(env, "low")
	assert.Equal(ci.RecoveringNumber(), cells.MinRecoveringNumber)
	assert.Equal(ci.RecoveringDuration(), cells.MinRecoveringDuration)

	err = env.StartCell("high", newRecoveringFrequencyBehavior(12, time.Minute, sink))
	assert.Nil(err)
	ci = cells.InspectCell(env, "high")
	assert.Equal(ci.RecoveringNumber(), 12)
	assert.Equal(ci.RecoveringDuration(), time.Minute)
}

// TestEnvironmentSubscribeStop subscribing and stopping
func TestEnvironmentSubscribeStop(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)

	env := cells.NewEnvironment("subscribe-unsubscribe-stop")
	defer env.Stop()

	sink := cells.NewEventSink(0)
	assert.Nil(env.StartCell("foo", newCollectBehavior(sink)))
	assert.Nil(env.StartCell("bar", newCollectBehavior(sink)))
	assert.Nil(env.StartCell("baz", newCollectBehavior(sink)))

	assert.Nil(env.Subscribe("foo", "bar", "baz"))
	assert.Nil(env.Subscribe("bar", "foo", "baz"))

	assert.Nil(env.StopCell("bar"))
	assert.Nil(env.StopCell("foo"))
	assert.Nil(env.StopCell("baz"))
}

// TestEnvironmentSubscribeUnsubscribe tests subscribing,
// checking and unsubscribing of cells.
func TestEnvironmentSubscribeUnsubscribe(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)

	env := cells.NewEnvironment("subscribe-unsubscribe")
	defer env.Stop()

	sink := cells.NewEventSink(0)
	err := env.StartCell("foo", newCollectBehavior(sink))
	assert.Nil(err)
	err = env.StartCell("bar", newCollectBehavior(sink))
	assert.Nil(err)
	err = env.StartCell("baz", newCollectBehavior(sink))
	assert.Nil(err)
	err = env.StartCell("yadda", newCollectBehavior(sink))
	assert.Nil(err)

	err = env.Subscribe("humpf", "foo")
	assert.True(errors.IsError(err, cells.ErrInvalidID))
	err = env.Subscribe("foo", "humpf")
	assert.True(errors.IsError(err, cells.ErrInvalidID))

	err = env.Subscribe("foo", "bar", "baz")
	assert.Nil(err)
	subs, err := env.Subscribers("foo")
	assert.Nil(err)
	assert.Contents("bar", subs)
	assert.Contents("baz", subs)

	err = env.Unsubscribe("foo", "bar")
	assert.Nil(err)
	subs, err = env.Subscribers("foo")
	assert.Nil(err)
	assert.Contents("baz", subs)

	err = env.Unsubscribe("foo", "baz")
	assert.Nil(err)
	subs, err = env.Subscribers("foo")
	assert.Nil(err)
	assert.Empty(subs)
}

// TestEnvironmentStopUnsubscribe tests the unsubscribe of a cell when
// it is stopped.
func TestEnvironmentStopUnsubscribe(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)

	env := cells.NewEnvironment("stop-unsubscribe")
	defer env.Stop()

	fooSink := cells.NewEventSink(0)
	barSink := cells.NewEventSink(0)
	bazSink := cells.NewEventSink(0)
	err := env.StartCell("foo", newCollectBehavior(fooSink))
	assert.Nil(err)
	err = env.StartCell("bar", newCollectBehavior(barSink))
	assert.Nil(err)
	err = env.StartCell("baz", newCollectBehavior(bazSink))
	assert.Nil(err)

	err = env.Subscribe("foo", "bar", "baz")
	assert.Nil(err)

	err = env.StopCell("bar")
	assert.Nil(err)

	// Expect only baz because bar is stopped.
	// response, err := env.Request(ctx, "foo", subscribersTopic, time.Second)
	// assert.Nil(err)
	// ids := response.GetDefault([]string{})
	// assert.Equal(ids, []string{"baz"})
}

// TestEnvironmentSubscribersDo tests the iteration over
// the subscribers.
func TestEnvironmentSubscribersDo(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	env := cells.NewEnvironment("subscribers-do")
	defer env.Stop()

	fooc := audit.MakeSigChan()
	foo := func(cell cells.Cell, event cells.Event) (cells.Event, error) {
		fooc <- cell.ID() + "/" + event.Topic()
		return nil, nil
	}
	barc := audit.MakeSigChan()
	bar := func(cell cells.Cell, event cells.Event) (cells.Event, error) {
		barc <- cell.ID() + "/" + event.Topic()
		return nil, nil
	}
	iterator := func(cell cells.Cell, event cells.Event) (cells.Event, error) {
		err := cell.SubscribersDo(func(sub cells.Subscriber) error {
			return sub.ProcessEvent(event)
		})
		return nil, err
	}
	env.StartCell("foo", newSimpleBehavior(foo))
	env.StartCell("bar", newSimpleBehavior(bar))
	env.StartCell("iterator", newSimpleBehavior(iterator))

	err := env.Subscribe("iterator", "foo", "bar")
	assert.Nil(err)
	err = env.EmitNew("iterator", "ping", nil)
	assert.Nil(err)

	assert.Wait(fooc, "foo/ping", 2*time.Second)
	assert.Wait(barc, "bar/ping", 2*time.Second)
}

// TestEnvironmentScenario tests creating and using the
// environment in a simple way.
func TestEnvironmentScenario(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	env := cells.NewEnvironment("scenario")
	defer env.Stop()

	bazSink := cells.NewEventSink(0)
	sigc := audit.MakeSigChan()
	count := 0
	counter := func(cell cells.Cell, event cells.Event) (cells.Event, error) {
		count++
		if count == 2 {
			sigc <- count
		}
		return nil, nil
	}
	env.StartCell("foo", newEmitBehavior())
	env.StartCell("bar", newEmitBehavior())
	env.StartCell("baz", newCollectBehavior(bazSink))
	env.StartCell("counter", newSimpleBehavior(counter))

	err := env.Subscribe("foo", "bar")
	assert.Nil(err)
	err = env.Subscribe("bar", "baz")
	assert.Nil(err)
	err = env.Subscribe("baz", "counter")
	assert.Nil(err)

	err = env.EmitNew("foo", "lorem", cells.NewDefaultPayload("4711"))
	assert.Nil(err)
	err = env.EmitNew("foo", "ipsum", cells.NewDefaultPayload("1234"))
	assert.Nil(err)

	assert.Wait(sigc, 2, 2*time.Second)
	assert.Length(bazSink, 2)

	ok, err := bazSink.Match(func(index int, event cells.Event) (bool, error) {
		switch event.Topic() {
		case "lorem", "ipsum":
			return true, nil
		default:
			return false, nil
		}
	})
	assert.Nil(err)
	assert.True(ok)
}

//--------------------
// BENCHMARKS
//--------------------

// BenchmarkSmpleEmitNullMonitoring is a simple emitting to one cell
// with the null monitor.
func BenchmarkSmpleEmitNullMonitoring(b *testing.B) {
	monitoring.SetBackend(monitoring.NewNullBackend())
	env := cells.NewEnvironment("simple-emit-null")
	defer env.Stop()

	env.StartCell("null", &nullBehavior{})

	event, _ := cells.NewEvent("foo", cells.NewDefaultPayload("bar"))

	for i := 0; i < b.N; i++ {
		env.Emit("null", event)
	}
}

// BenchmarkSmpleEmitStandardMonitoring is a simple emitting to one cell
// with the standard monitor.
func BenchmarkSmpleEmitStandardMonitoring(b *testing.B) {
	monitoring.SetBackend(monitoring.NewStandardBackend())
	env := cells.NewEnvironment("simple-emit-standard")
	defer env.Stop()

	env.StartCell("null", &nullBehavior{})

	event, _ := cells.NewEvent("foo", cells.NewDefaultPayload("bar"))

	for i := 0; i < b.N; i++ {
		env.Emit("null", event)
	}
}

// EOF
