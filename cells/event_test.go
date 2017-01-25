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
	"context"
	"errors"
	"testing"
	"time"

	"github.com/tideland/golib/audit"

	"github.com/tideland/gocells/cells"
)

//--------------------
// TESTS
//--------------------

// TestPositiveWaitPayload waits for a payload.
func TestPositiveWaitPayload(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	waiter := cells.NewPayloadWaiter()

	go func() {
		time.Sleep(250 * time.Millisecond)
		waiter.Set(cells.NewPayload(4711))
		waiter.Set(cells.NewPayload(1174))
	}()

	ctx := context.Background()
	payload, err := waiter.Wait(ctx)
	assert.Nil(err)
	set := payload.GetInt(cells.DefaultPayload, 0)
	assert.Equal(set, 4711)
}

// TestWaitPayloadTimeout waits for a payload but
// timeout is faster.
func TestWaitPayloadTimeout(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	waiter := cells.NewPayloadWaiter()

	go func() {
		time.Sleep(500 * time.Millisecond)
		waiter.Set(cells.NewPayload(4711))
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	payload, err := waiter.Wait(ctx)
	cancel()
	assert.ErrorMatch(err, "context deadline exceeded")
	assert.Nil(payload)
}

// TestWaitPayloadCancel waits for a payload but
// it's canceled earlier.
func TestWaitPayloadCancel(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	waiter := cells.NewPayloadWaiter()

	go func() {
		time.Sleep(500 * time.Millisecond)
		waiter.Set(cells.NewPayload(4711))
	}()

	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(100*time.Millisecond, func() {
		cancel()
	})
	payload, err := waiter.Wait(ctx)
	assert.ErrorMatch(err, "context canceled")
	assert.Nil(payload)
}

// TestEventSink tests the simple event sink.
func TestEventSink(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	checkTopic := func(event cells.Event) {
		assert.Contents(event.Topic(), topics)
	}

	// Empty sink.
	sink := cells.NewEventSink(0)
	first, ok := sink.First()
	assert.Nil(first)
	assert.False(ok)
	last, ok := sink.Last()
	assert.Nil(last)
	assert.False(ok)
	at, ok := sink.At(-1)
	assert.Nil(at)
	assert.False(ok)
	at, ok = sink.At(4711)
	assert.Nil(at)
	assert.False(ok)

	// Limited number of events.
	sink = cells.NewEventSink(5)
	addEvents(assert, 10, sink)
	assert.Length(sink, 5)

	// Unlimited number of events.
	sink = cells.NewEventSink(0)
	addEvents(assert, 10, sink)
	assert.Length(sink, 10)

	first, ok = sink.First()
	assert.True(ok)
	checkTopic(first)
	last, ok = sink.Last()
	assert.True(ok)
	checkTopic(last)

	for i := 0; i < sink.Len(); i++ {
		at, ok = sink.At(i)
		assert.True(ok)
		checkTopic(at)
	}
}

// TestEventSinkIteration tests the event sink iteration.
func TestEventSinkIteration(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	sink := cells.NewEventSink(0)
	addEvents(assert, 10, sink)

	assert.Length(sink, 10)
	err := sink.Do(func(index int, event cells.Event) error {
		assert.Contents(event.Topic(), topics)
		payload := event.Payload().GetInt(cells.DefaultPayload, -1)
		assert.Range(payload, 1, 10)
		return nil
	})
	assert.Nil(err)
	ok, err := sink.Match(func(index int, event cells.Event) (bool, error) {
		topicOK := event.Topic() >= "a" && event.Topic() <= "j"
		payload := event.Payload().GetInt(cells.DefaultPayload, -1)
		payloadOK := payload >= 1 && payload <= 10
		return topicOK && payloadOK, nil
	})
	assert.Nil(err)
	assert.True(ok)
}

// TestEventSinkIterationError tests the event sink iteration error.
func TestEventSinkIterationError(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	sink := cells.NewEventSink(0)
	addEvents(assert, 10, sink)

	err := sink.Do(func(index int, event cells.Event) error {
		return errors.New("ouch")
	})
	assert.ErrorMatch(err, "ouch")
	ok, err := sink.Match(func(index int, event cells.Event) (bool, error) {
		// The bool true won't be passed to outside.
		return true, errors.New("ouch")
	})
	assert.False(ok)
	assert.ErrorMatch(err, "ouch")
}

// TestCheckedEventSink tests the notification of a waiter
// when a criterion in the sink matches.
func TestCheckedEventSink(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	checker := func(events cells.EventSinkAccessor) (bool, cells.Payload, error) {
		wanted := []string{"f", "c", "c"}
		if events.Len() < len(wanted) {
			return false, nil, nil
		}
		ok, err := events.Match(func(index int, event cells.Event) (bool, error) {
			return event.Topic() == wanted[index], nil
		})
		if err != nil {
			return false, nil, err
		}
		if ok {
			first, _ := events.First()
			last, _ := events.Last()
			payload := cells.NewPayload(cells.PayloadValues{
				"first": first.Timestamp(),
				"last":  last.Timestamp(),
			})
			return true, payload, nil
		}
		return false, nil, nil
	}
	sink, waiter := cells.NewCheckedEventSink(3, checker)

	go addEvents(assert, 100, sink)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	payload, err := waiter.Wait(ctx)
	assert.Nil(err)
	first := payload.GetTime("first", time.Time{})
	last := payload.GetTime("last", time.Time{})
	assert.Logf("First: %v", first)
	assert.Logf("Last : %v", last)
	assert.Logf("Duration: %v", last.Sub(first))
	assert.True(last.UnixNano() > first.UnixNano())
	cancel()
}

// TestCheckedEventSinkFailing tests the missing notification of a waiter
// when a criterion in the sink does not match.
func TestCheckedEventSinkFailing(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	checker := func(events cells.EventSinkAccessor) (bool, cells.Payload, error) {
		return false, nil, nil
	}
	sink, waiter := cells.NewCheckedEventSink(3, checker)

	go addEvents(assert, 100, sink)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	payload, err := waiter.Wait(ctx)
	assert.Nil(payload)
	assert.ErrorMatch(err, "context deadline exceeded")
	cancel()
}

//--------------------
// HELPER
//--------------------

// topics contains the test topics.
var topics = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}

// addEvents adds a number of events to a sink.
func addEvents(assert audit.Assertion, count int, sink cells.EventSink) {
	generator := audit.NewGenerator(audit.FixedRand())
	for i := 0; i < count; i++ {
		topic := generator.OneStringOf(topics...)
		payload := generator.Int(1, 10)
		event, err := cells.NewEvent(nil, topic, payload)
		assert.Nil(err)
		n, err := sink.Add(event)
		assert.Nil(err)
		assert.True(n > 0)
		time.Sleep(2 * time.Millisecond)
	}
}

// EOF
