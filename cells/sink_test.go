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
	"sort"
	"testing"
	"time"

	"github.com/tideland/golib/audit"

	"github.com/tideland/gocells/cells"
)

//--------------------
// TESTS
//--------------------

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
		var payload int
		err := event.Payload().Unmarshal(&payload)
		assert.Nil(err)
		assert.Range(payload, 1, 9)
		return nil
	})
	assert.Nil(err)
	ok, err := cells.NewEventSinkAnalyzer(sink).Match(func(index int, event cells.Event) (bool, error) {
		topicOK := event.Topic() >= "a" && event.Topic() <= "j"
		var payload int
		err := event.Payload().Unmarshal(&payload)
		assert.Nil(err)
		payloadOK := payload >= 1 && payload <= 9
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
	ok, err := cells.NewEventSinkAnalyzer(sink).Match(func(index int, event cells.Event) (bool, error) {
		// The bool true won't be passed to outside.
		return true, stderr.New("ouch")
	})
	assert.False(ok)
	assert.ErrorMatch(err, "ouch")
}

// TestCheckedEventSink tests the checking of new events.
func TestCheckedEventSink(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	payloadc := audit.MakeSigChan()
	donec := audit.MakeSigChan()
	count := 0
	wanted := []string{"f", "c", "c"}
	checker := func(events cells.EventSinkAccessor) (cells.Payload, error) {
		count++
		defer func() {
			if count == 100 {
				donec <- struct{}{}
			}
		}()
		if events.Len() < len(wanted) {
			return nil, nil
		}
		ok, err := cells.NewEventSinkAnalyzer(events).Match(func(index int, event cells.Event) (bool, error) {
			return event.Topic() == wanted[index], nil
		})
		if err != nil {
			return nil, err
		}
		if ok {
			first, _ := events.PeekFirst()
			last, _ := events.PeekLast()
			payload := last.Timestamp().Sub(first.Timestamp())
			payloadc <- payload
		}
		return nil, nil
	}
	sink := cells.NewCheckedEventSink(3, checker)

	go addEvents(assert, 100, sink)

	for {
		select {
		case payload := <-payloadc:
			d, ok := payload.(time.Duration)
			assert.True(ok)
			assert.True(d > 0)
		case <-donec:
			return
		case <-time.After(5 * time.Second):
			assert.Fail()
		}
	}
}

// TestEventSinkAnalyzer tests analyzing an event sink.
func TestEventSinkAnalyzer(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	sink := cells.NewEventSink(0)
	analyzer := cells.NewEventSinkAnalyzer(sink)

	addEvents(assert, 100, sink)

	// Check filtering.
	fchecker := func(index int, event cells.Event) (bool, error) {
		if event.Topic() == "f" {
			return true, nil
		}
		return false, nil
	}
	ferrchecker := func(index int, event cells.Event) (bool, error) {
		return false, stderr.New("ouch")
	}
	fs, err := analyzer.Filter(fchecker)
	assert.Nil(err)
	assert.True(fs.Len() < sink.Len(), "less events with topic f than total number")
	fs, err = analyzer.Filter(ferrchecker)
	assert.ErrorMatch(err, "ouch", "error is returned correctly")

	// Check matching.
	fs, err = analyzer.Filter(fchecker)
	assert.Nil(err)
	ok, err := cells.NewEventSinkAnalyzer(fs).Match(fchecker)
	assert.Nil(err)
	assert.True(ok, "all events in fs do have topic f")
	ok, err = cells.NewEventSinkAnalyzer(fs).Match(ferrchecker)
	assert.ErrorMatch(err, "ouch", "error is returned correctly")

	// Check folding.
	count := 0
	ffolder := func(index int, acc interface{}, event cells.Event) (interface{}, error) {
		if event.Topic() == "f" {
			count++
			fs, ok := acc.(int)
			if !ok {
				return nil, stderr.New("ouch")
			}
			return fs + 1, nil
		}
		return acc, nil
	}
	ferrfolder := func(index int, acc interface{}, event cells.Event) (interface{}, error) {
		return nil, stderr.New("ouch")
	}
	fcount, err := analyzer.Fold(0, ffolder)
	assert.Nil(err)
	assert.Equal(fcount, count, "accumulator has been updated correctly")
	fcount, err = analyzer.Fold(0, ferrfolder)
	assert.ErrorMatch(err, "ouch", "error is returned correctly")

	count = 0
	fpfolder := func(index int, acc cells.Payload, event cells.Event) (cells.Payload, error) {
		if event.Topic() == "f" {
			count++
			payload, err := cells.NewPayload(acc.String() + event.Topic())
			if err != nil {
				return nil, err
			}
			return payload, nil
		}
		return acc, nil
	}
	fperrfolder := func(index int, acc cells.Payload, event cells.Event) (cells.Payload, error) {
		return nil, stderr.New("ouch")
	}
	initial, err := cells.NewPayload("")
	assert.Nil(err)
	fpcount, err := analyzer.FoldPayload(initial, fpfolder)
	assert.Nil(err)
	assert.Length(fpcount, count, "payload accumulator has been updated correctly")
	initial, err = cells.NewPayload("")
	assert.Nil(err)
	fpcount, err = analyzer.FoldPayload(initial, fperrfolder)
	assert.ErrorMatch(err, "ouch", "error is returned correctly")

	// Check total duration.
	dsink := cells.NewEventSink(0)
	danalyzer := cells.NewEventSinkAnalyzer(dsink)
	duration := danalyzer.TotalDuration()
	assert.Equal(duration, 0*time.Nanosecond, "empty sink has no duration")

	addEvents(assert, 1, dsink)
	duration = danalyzer.TotalDuration()
	assert.Equal(duration, 0*time.Nanosecond, "sink containing one event has no duration")

	addEvents(assert, 1, dsink)
	first, ok := dsink.PeekFirst()
	assert.True(ok)
	last, ok := dsink.PeekLast()
	assert.True(ok)
	duration = danalyzer.TotalDuration()
	assert.Equal(duration, last.Timestamp().Sub(first.Timestamp()), "total duration calculated correctly")

	// Check minimum/maximum duration.
	durations := []time.Duration{}
	timestamp := time.Time{}
	sink.Do(func(index int, event cells.Event) error {
		if index == 0 {
			timestamp = event.Timestamp()
			return nil
		}
		duration = event.Timestamp().Sub(timestamp)
		durations = append(durations, duration)
		timestamp = event.Timestamp()
		return nil
	})
	sort.Slice(durations, func(i, j int) bool {
		return durations[i] < durations[j]
	})
	dmin := durations[0]
	dmax := durations[len(durations)-1]
	min, max := analyzer.MinMaxDuration()
	assert.Equal(min, dmin, "minimum duration is correct")
	assert.Equal(max, dmax, "maximum duration is correct")

	// Check topic quantities.
	quantities := analyzer.TopicQuantities()
	assert.Length(quantities, len(topics))
	for topic, quantity := range quantities {
		assert.Contents(topic, topics, "topic is one of the topics")
		assert.Range(quantity, 1, 100, "quantity is in range")
	}

	tfolder := func(index int, acc interface{}, event cells.Event) (interface{}, error) {
		if acc == nil {
			return 0, nil
		}
		count, ok := acc.(int)
		if !ok {
			return nil, stderr.New("ouch")
		}
		return count + 1, nil
	}
	terrfolder := func(index int, acc interface{}, event cells.Event) (interface{}, error) {
		return nil, stderr.New("ouch")
	}
	counts, err := analyzer.TopicFolds(tfolder)
	assert.Nil(err)
	assert.Length(quantities, len(topics))
	for topic, count := range counts {
		assert.Contents(topic, topics, "topic is one of the topics")
		assert.Range(count, 1, 100, "quantity is in range")
	}
	counts, err = analyzer.TopicFolds(terrfolder)
	assert.ErrorMatch(err, "ouch", "error is returned correctly")
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
		payload := generator.OneIntOf(1, 2, 3, 4, 5, 6, 7, 8, 9)
		sleep := generator.Duration(2*time.Millisecond, 4*time.Millisecond)
		event, err := cells.NewEvent(topic, payload)
		assert.Nil(err)
		n, err := sink.Push(event)
		assert.Nil(err)
		assert.True(n > 0)
		time.Sleep(sleep)
	}
}

// EOF
