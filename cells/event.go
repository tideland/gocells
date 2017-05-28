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
func NewEvent(topic string, payload Payload) (Event, error) {
	if topic == "" {
		return nil, errors.New(ErrNoTopic, errorMessages)
	}
	return &event{
		timestamp: time.Now().UTC(),
		topic:     topic,
		payload:   payload,
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
	Do(doer func(index int, event Event) error) error

	// Match checks if all events match the passed criterion.
	Match(matcher func(index int, event Event) (bool, error)) (bool, error)
}

// EventSinkChecker can be used to check sinks for a criterion.
type EventSinkChecker func(events EventSinkAccessor) error

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

// eventSink implements the EventSink interface.
type eventSink struct {
	mutex   sync.RWMutex
	max     int
	events  []Event
	checker EventSinkChecker
}

// NewEventSink creates a sink for events.
func NewEventSink(max int) EventSink {
	return &eventSink{
		max: max,
	}
}

// NewCheckedEventSink creates a sink for events.
func NewCheckedEventSink(max int, checker EventSinkChecker) EventSink {
	return &eventSink{
		max:     max,
		checker: checker,
	}
}

// Push implements the EventSink interface.
func (s *eventSink) Push(event Event) (int, error) {
	s.mutex.Lock()
	s.events = append(s.events, event)
	if s.max > 0 && len(s.events) > s.max {
		s.events = s.events[1:]
	}
	s.mutex.Unlock()
	return len(s.events), s.check()
}

// PullFirst implements the EventSink interface.
func (s *eventSink) PullFirst() (Event, error) {
	var event Event
	s.mutex.Lock()
	if len(s.events) > 0 {
		event = s.events[0]
		s.events = s.events[1:]
	}
	s.mutex.Unlock()
	return event, s.check()
}

// PullLast implements the EventSink interface.
func (s *eventSink) PullLast() (Event, error) {
	var event Event
	s.mutex.Lock()
	if len(s.events) > 0 {
		event = s.events[len(s.events)-1]
		s.events = s.events[:len(s.events)-1]
	}
	s.mutex.Unlock()
	return event, s.check()
}

// Clear implements tne EventSink interface.
func (s *eventSink) Clear() error {
	s.mutex.Lock()
	s.events = nil
	s.mutex.Unlock()
	return s.check()
}

// Len implements the EventSinkAccessor interface.
func (s *eventSink) Len() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.events)
}

// PeekFirst implements the EventSinkAccessor interface.
func (s *eventSink) PeekFirst() (Event, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if len(s.events) < 1 {
		return nil, false
	}
	return s.events[0], true
}

// PeekLast implements the EventSinkAccessor interface.
func (s *eventSink) PeekLast() (Event, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if len(s.events) < 1 {
		return nil, false
	}
	return s.events[len(s.events)-1], true
}

// PeekAt implements the EventSinkAccessor interface.
func (s *eventSink) PeekAt(index int) (Event, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if index < 0 || index > len(s.events)-1 {
		return nil, false
	}
	return s.events[index], true
}

// Do implements the EventSinkAccessor interface.
func (s *eventSink) Do(doer func(index int, event Event) error) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	for index, event := range s.events {
		if err := doer(index, event); err != nil {
			return err
		}
	}
	return nil
}

// Match implements the EventSinkAccessor interface.
func (s *eventSink) Match(matcher func(index int, event Event) (bool, error)) (bool, error) {
	match := true
	doer := func(mindex int, mevent Event) error {
		ok, err := matcher(mindex, mevent)
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

// check calls the checker if configured.
func (s *eventSink) check() error {
	if s.checker != nil {
		if err := s.checker(s); err != nil {
			return err
		}
	}
	return nil
}

// EOF
