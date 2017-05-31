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
	stderr "errors"
	"testing"
	"time"

	"github.com/tideland/golib/audit"
	"github.com/tideland/golib/errors"

	"github.com/tideland/gocells/cells"
)

//--------------------
// TESTS
//--------------------

// TestEvent tests the event construction.
func TestEvent(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	now := time.Now().UTC()

	event, err := cells.NewEvent("foo", cells.NewDefaultPayload("bar"))
	assert.Nil(err)
	assert.True(event.Timestamp().After(now))
	assert.True(time.Now().UTC().After(event.Timestamp()))
	assert.Equal(event.Topic(), "foo")

	bar, ok := event.Payload().Get("default")
	assert.True(ok)
	assert.Equal(bar, "bar")

	_, err = cells.NewEvent("", nil)
	assert.True(errors.IsError(err, cells.ErrNoTopic))

	_, err = cells.NewEvent("yadda", nil)
	assert.Nil(err)
}

// TestPayload tests the payload creation and access.
func TestPayload(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	now := time.Now()
	dur := 30 * time.Second
	pl := cells.Values{
		"bool":     true,
		"int":      42,
		"float64":  47.11,
		"string":   "hello, world",
		"time":     now,
		"duration": dur,
	}.Payload()

	s, ok := pl.Get("string")
	assert.True(ok)
	assert.Equal(s, "hello, world")
	s, ok = pl.Get("no-string")
	assert.False(ok)
	assert.Equal(s, "")

	b, ok := pl.GetBool("bool")
	assert.True(ok)
	assert.True(b)
	b, ok = pl.GetBool("no-bool")
	assert.False(ok)
	assert.False(b)

	i, ok := pl.GetInt("int")
	assert.True(ok)
	assert.Equal(i, 42)
	i, ok = pl.GetInt("no-int")
	assert.False(ok)
	assert.Equal(i, 0)

	f, ok := pl.GetFloat64("float64")
	assert.True(ok)
	assert.Equal(f, 47.11)
	f, ok = pl.GetFloat64("no-float64")
	assert.False(ok)
	assert.Equal(f, 0.0)

	tt, ok := pl.GetTime("time")
	assert.True(ok)
	assert.Equal(tt, now)
	tt, ok = pl.GetTime("no-time")
	assert.False(ok)
	assert.Equal(tt, time.Time{})

	td, ok := pl.GetDuration("duration")
	assert.True(ok)
	assert.Equal(td, dur)
	td, ok = pl.GetDuration("no-duration")
	assert.False(ok)
	assert.Equal(td, 0*time.Second)

	pln := pl.Apply(cells.Values{
		cells.PayloadDefault: "foo",
	})
	s, ok = pln.Get(cells.PayloadDefault)
	assert.True(ok)
	assert.Equal(s, "foo")
	assert.Length(pln, 7)
}

// TestEventSink tests the simple event sink.
func TestEventSink(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	checkTopic := func(event cells.Event) {
		assert.Contents(event.Topic(), topics)
	}

	// Empty sink.
	sink := cells.NewEventSink(0)
	first, ok := sink.PeekFirst()
	assert.Nil(first)
	assert.False(ok)
	last, ok := sink.PeekLast()
	assert.Nil(last)
	assert.False(ok)
	at, ok := sink.PeekAt(-1)
	assert.Nil(at)
	assert.False(ok)
	at, ok = sink.PeekAt(4711)
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

	first, ok = sink.PeekFirst()
	assert.True(ok)
	checkTopic(first)
	last, ok = sink.PeekLast()
	assert.True(ok)
	checkTopic(last)

	for i := 0; i < sink.Len(); i++ {
		at, ok = sink.PeekAt(i)
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
		payload, ok := event.Payload().GetInt(cells.PayloadDefault)
		assert.True(ok)
		assert.Range(payload, 1, 10)
		return nil
	})
	assert.Nil(err)
	ok, err := sink.Match(func(index int, event cells.Event) (bool, error) {
		topicOK := event.Topic() >= "a" && event.Topic() <= "j"
		payload, ok := event.Payload().GetInt(cells.PayloadDefault)
		assert.True(ok)
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
		return stderr.New("ouch")
	})
	assert.ErrorMatch(err, "ouch")
	ok, err := sink.Match(func(index int, event cells.Event) (bool, error) {
		// The bool true won't be passed to outside.
		return true, stderr.New("ouch")
	})
	assert.False(ok)
	assert.ErrorMatch(err, "ouch")
}

// TestCheckedEventSink tests the notification of a waiter
// when a criterion in the sink matches.
func TestCheckedEventSink(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	payloadc := make(chan cells.Payload, 1)
	donec := make(chan struct{})
	count := 0
	wanted := []string{"f", "c", "c"}
	checker := func(events cells.EventSinkAccessor) error {
		count++
		defer func() {
			if count == 100 {
				donec <- struct{}{}
			}
		}()
		if events.Len() < len(wanted) {
			return nil
		}
		ok, err := events.Match(func(index int, event cells.Event) (bool, error) {
			return event.Topic() == wanted[index], nil
		})
		if err != nil {
			return err
		}
		if ok {
			first, _ := events.PeekFirst()
			last, _ := events.PeekLast()
			payload := cells.Values{
				"first": first.Timestamp(),
				"last":  last.Timestamp(),
			}.Payload()
			payloadc <- payload
		}
		return nil
	}
	sink := cells.NewCheckedEventSink(3, checker)

	go addEvents(assert, 100, sink)

	for {
		select {
		case payload := <-payloadc:
			first, ok := payload.GetTime("first")
			assert.True(ok)
			last, ok := payload.GetTime("last")
			assert.True(ok)
			assert.True(last.UnixNano() > first.UnixNano())
		case <-donec:
			return
		case <-time.After(5 * time.Second):
			assert.Fail()
		}
	}
}

//--------------------
// HELPER
//--------------------

// topics contains the test topics.
var topics = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}

// values contains the test values.
var values = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}

// addEvents adds a number of events to a sink.
func addEvents(assert audit.Assertion, count int, sink cells.EventSink) {
	generator := audit.NewGenerator(audit.FixedRand())
	for i := 0; i < count; i++ {
		topic := generator.OneStringOf(topics...)
		payload := cells.NewDefaultPayload(generator.OneStringOf(values...))
		event, err := cells.NewEvent(topic, payload)
		assert.Nil(err)
		n, err := sink.Push(event)
		assert.Nil(err)
		assert.True(n > 0)
		time.Sleep(2 * time.Millisecond)
	}
}

// EOF
