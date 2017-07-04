// Tideland Go Cells - Event
//
// Copyright (C) 2010-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package cells

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"sync"
	"time"

	"github.com/tideland/golib/errors"
)

//--------------------
// EVENT
//--------------------

// Event transports what to process.
type Event interface {
	fmt.Stringer

	// Timestamp returns the UTC time the event has been created.
	Timestamp() time.Time

	// Topic returns the topic of the event.
	Topic() string

	// Payload returns the payload of the event.
	Payload() Payload
}

// event implements the Event interface.
type event struct {
	timestamp time.Time
	topic     string
	payload   Payload
}

// NewEvent creates a new event with the given topic and payload.
func NewEvent(topic string, payload interface{}) (Event, error) {
	if topic == "" {
		return nil, errors.New(ErrNoTopic, errorMessages)
	}
	p, err := NewPayload(payload)
	if err != nil {
		return nil, err
	}
	return &event{
		timestamp: time.Now().UTC(),
		topic:     topic,
		payload:   p,
	}, nil
}

// Timestamp implements the Event interface.
func (e *event) Timestamp() time.Time {
	return e.timestamp
}

// Topic implements the Event interface.
func (e *event) Topic() string {
	return e.topic
}

// Payload implements the Event interface.
func (e *event) Payload() Payload {
	if e.payload == nil {
		// Fallback to empty one.
		return newEmptyPayload()
	}
	return e.payload
}

// String implements the Stringer interface.
func (e *event) String() string {
	timeStr := e.timestamp.Format(time.RFC3339Nano)
	payloadStr := "none"
	if e.payload != nil {
		payloadStr = fmt.Sprintf("%v", e.payload)
	}
	return fmt.Sprintf("<timestamp: %s / topic: '%s' / payload: %s>", timeStr, e.topic, payloadStr)
}

//--------------------
// EVENT SINK
//--------------------

// EventSinkDoer performs an operation on an event.
type EventSinkDoer func(index int, event Event) error

// EventSinkFilter checks if an event matches a criterium.
type EventSinkFilter func(index int, event Event) (bool, error)

// EventSinkFolder allows to reduce (fold) events.
type EventSinkFolder func(index int, acc interface{}, event Event) (interface{}, error)

// EventSinkAccessor can be used to read the events in a sink.
type EventSinkAccessor interface {
	// Len returns the number of stored events.
	Len() int

	// PeekFirst returns the first of the collected events.
	PeekFirst() (Event, bool)

	// PeekLast returns the last of the collected event datas.
	PeekLast() (Event, bool)

	// PeekAt returns an event at a given index and true if it
	// exists, otherwise nil and false.
	PeekAt(index int) (Event, bool)

	// Do iterates over all collected events.
	Do(doer EventSinkDoer) error

	// Filter creates a new accessor containing only the filtered
	// events.
	Filter(filter EventSinkFilter) (EventSinkAccessor, error)

	// Match checks if all events match the passed criterion.
	Match(matcher EventSinkFilter) (bool, error)

	// Fold reduces (folds) the events of the sink.
	Fold(initial interface{}, folder EventSinkFolder) (interface{}, error)
}

// EventSinkProcessor can be used as a checker function but also inside of
// behaviors to process the content of an event sink and return a new payload.
type EventSinkProcessor func(events EventSinkAccessor) (Payload, error)

// EventSink stores a number of events ordered by adding them at the end. To
// be used in behaviors for collecting sets of events and operate on them.
type EventSink interface {
	// Push adds a new event to the sink.
	Push(event Event) (int, error)

	// PullFirst returns and removed the first event of the sink.
	PullFirst() (Event, error)

	// PullLast returns and removed the last event of the sink.
	PullLast() (Event, error)

	// Clear removes all collected events.
	Clear() error

	EventSinkAccessor
}

// eventSink implements EventSink.
type eventSink struct {
	mutex  sync.RWMutex
	max    int
	events []Event
	check  EventSinkProcessor
}

// NewEventSink creates a sink for events.
func NewEventSink(max int) EventSink {
	return &eventSink{
		max: max,
	}
}

// NewCheckedEventSink creates a sink for events.
func NewCheckedEventSink(max int, checker EventSinkProcessor) EventSink {
	return &eventSink{
		max:   max,
		check: checker,
	}
}

// Push implements EventSink.
func (s *eventSink) Push(event Event) (int, error) {
	s.mutex.Lock()
	s.events = append(s.events, event)
	if s.max > 0 && len(s.events) > s.max {
		s.events = s.events[1:]
	}
	s.mutex.Unlock()
	return len(s.events), s.performCheck()
}

// PullFirst implements EventSink.
func (s *eventSink) PullFirst() (Event, error) {
	var event Event
	s.mutex.Lock()
	if len(s.events) > 0 {
		event = s.events[0]
		s.events = s.events[1:]
	}
	s.mutex.Unlock()
	return event, s.performCheck()
}

// PullLast implements EventSink.
func (s *eventSink) PullLast() (Event, error) {
	var event Event
	s.mutex.Lock()
	if len(s.events) > 0 {
		event = s.events[len(s.events)-1]
		s.events = s.events[:len(s.events)-1]
	}
	s.mutex.Unlock()
	return event, s.performCheck()
}

// Clear implements EventSink.
func (s *eventSink) Clear() error {
	s.mutex.Lock()
	s.events = nil
	s.mutex.Unlock()
	return s.performCheck()
}

// Len implements EventSinkAccessor.
func (s *eventSink) Len() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.events)
}

// PeekFirst implements EventSinkAccessor.
func (s *eventSink) PeekFirst() (Event, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if len(s.events) < 1 {
		return nil, false
	}
	return s.events[0], true
}

// PeekLast implements EventSinkAccessor.
func (s *eventSink) PeekLast() (Event, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if len(s.events) < 1 {
		return nil, false
	}
	return s.events[len(s.events)-1], true
}

// PeekAt implements EventSinkAccessor.
func (s *eventSink) PeekAt(index int) (Event, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if index < 0 || index > len(s.events)-1 {
		return nil, false
	}
	return s.events[index], true
}

// Do implements EventSinkAccessor.
func (s *eventSink) Do(doer EventSinkDoer) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	for index, event := range s.events {
		if err := doer(index, event); err != nil {
			return err
		}
	}
	return nil
}

// Filter implements EventSinkAccessor.
func (s *eventSink) Filter(filter EventSinkFilter) (EventSinkAccessor, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	accessor := NewEventSink(s.Len())
	for index, event := range s.events {
		ok, err := filter(index, event)
		if err != nil {
			return nil, err
		}
		if ok {
			accessor.Push(event)
		}
	}
	return accessor, nil
}

// Match implements EventSinkAccessor.
func (s *eventSink) Match(matcher EventSinkFilter) (bool, error) {
	match := true
	doer := func(index int, event Event) error {
		ok, err := matcher(index, event)
		if err != nil {
			match = false
			return err
		}
		match = match && ok
		return nil
	}
	err := s.Do(doer)
	return match, err
}

// Fold implements EventSinkAccessor.
func (s *eventSink) Fold(inject interface{}, folder EventSinkFolder) (interface{}, error) {
	acc := inject
	doer := func(index int, event Event) error {
		facc, err := folder(index, acc, event)
		if err != nil {
			acc = nil
			return err
		}
		acc = facc
		return nil
	}
	err := s.Do(doer)
	return acc, err
}

// performCheck calls the checker if configured.
func (s *eventSink) performCheck() error {
	if s.check != nil {
		if _, err := s.check(s); err != nil {
			return err
		}
	}
	return nil
}

//--------------------
// EVENT SINK ANALYZER
//--------------------

// EventSinkAnalyzer describes a helpful type to analyze
// the events collected inside a sink. It's intended to
// make the life a behavior developer more simple.
type EventSinkAnalyzer interface {
	// TotalDuration returns the duration between the first
	// and the last event.
	TotalDuration() time.Duration

	// MinMaxDuration returns the minimum and maximum
	// durations between two individual events.
	MinMaxDuration() (time.Duration, time.Duration)

	// TopicQuantities returns a map of collected topics and
	// their quantity.
	TopicQuantities() map[string]int

	// TopicFolds reduces the events per topic.
	TopicFolds(folder EventSinkFolder) (map[string]interface{}, error)
}

// eventSinkAnalyzer implements EventSinkAnalyzer.
type eventSinkAnalyzer struct {
	accessor EventSinkAccessor
}

// NewEventSinkAnalyzer creates an analyzer for the
// given sink accessor.
func NewEventSinkAnalyzer(accessor EventSinkAccessor) EventSinkAnalyzer {
	return &eventSinkAnalyzer{
		accessor: accessor,
	}
}

// TotalDuration implements EventSinkAnalyzer.
func (esa *eventSinkAnalyzer) TotalDuration() time.Duration {
	first, firstOK := esa.accessor.PeekFirst()
	last, lastOK := esa.accessor.PeekLast()
	if !firstOK || !lastOK {
		return 0
	}
	return last.Timestamp().Sub(first.Timestamp())
}

// MinMaxDuration implements EventSinkAnalyzer.
func (esa *eventSinkAnalyzer) MinMaxDuration() (time.Duration, time.Duration) {
	minDuration := esa.TotalDuration()
	maxDuration := 0 * time.Nanosecond
	lastTimestamp := time.Time{}
	esa.accessor.Do(func(index int, event Event) error {
		if index > 0 {
			duration := event.Timestamp().Sub(lastTimestamp)
			if duration < minDuration {
				minDuration = duration
			}
			if duration > maxDuration {
				maxDuration = duration
			}
		}
		lastTimestamp = event.Timestamp()
		return nil
	})
	return minDuration, maxDuration
}

// TopicQuantities implements EventSinkAnalyzer.
func (esa *eventSinkAnalyzer) TopicQuantities() map[string]int {
	topics := map[string]int{}
	esa.accessor.Do(func(index int, event Event) error {
		topics[event.Topic()] = topics[event.Topic()] + 1
		return nil
	})
	return topics
}

// TopicFolds implements EventSinkAnalyzer.
func (esa *eventSinkAnalyzer) TopicFolds(folder EventSinkFolder) (map[string]interface{}, error) {
	folds := map[string]interface{}{}
	err := esa.accessor.Do(func(index int, event Event) error {
		facc, err := folder(index, folds[event.Topic()], event)
		if err != nil {
			folds = nil
			return err
		}
		folds[event.Topic()] = facc
		return nil
	})
	return folds, err
}

// EOF
