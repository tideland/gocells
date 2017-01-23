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
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/tideland/golib/errors"
)

//--------------------
// PAYLOAD
//--------------------

// PayloadValues is intended to set and get the information
// of a payload as bulk.
type PayloadValues map[string]interface{}

// Payload is a write-once/read-multiple container for the
// transport of additional information with events. In case
// one item is a reference type it's in the responsibility
// of the users to avoid concurrent changes of their values.
type Payload interface {
	fmt.Stringer

	// Len returns the number of values.
	Len() int

	// Get returns one of the payload values.
	Get(key string) (interface{}, bool)

	// GetBool returns one of the payload values
	// as bool. If it's no bool false is returned.
	GetBool(key string) (bool, bool)

	// GetInt returns one of the payload values
	// as int. If it's no int false is returned.
	GetInt(key string) (int, bool)

	// GetFloat64 returns one of the payload values
	// as float64. If it's no float64 false is returned.
	GetFloat64(key string) (float64, bool)

	// GetString returns one of the payload values
	// as string. If it's no string false is returned.
	GetString(key string) (string, bool)

	// GetTime returns one of the payload values
	// as time.Time. If it's no time false is returned.
	GetTime(key string) (time.Time, bool)

	// GetDuration returns one of the payload values as
	// time.Duration. If it's no duration false is returned.
	GetDuration(key string) (time.Duration, bool)

	// GetWaiter returns a payload waiter to be used
	// for answering with a payload.
	GetWaiter(key string) (PayloadWaiter, bool)

	// Keys return all keys of the payload.
	Keys() []string

	// Do iterates a function over all keys and values.
	Do(f func(key string, value interface{}) error) error

	// Apply creates a new payload containing the values
	// of this one and the passed values. Allowed are
	// PayloadValues, map[string]interface{}, and any
	// other single value. The latter will be stored
	// with the cells.DefaultPayload key. Values of this
	// payload are overwritten by those which are passed
	// if they share the key.
	Apply(values interface{}) Payload
}

// payload implements the Payload interface.
type payload struct {
	values PayloadValues
}

// NewPayload creates a new payload containing the passed
// values. In case of a Payload this is used directly, in
// case of a PayloadValues or a map[string]interface{} their
// content is used, and when passing any other type the
// value is stored with the key cells.DefaultPayload.
func NewPayload(values interface{}) Payload {
	if p, ok := values.(Payload); ok {
		return p
	}
	p := &payload{
		values: PayloadValues{},
	}
	switch vs := values.(type) {
	case PayloadValues:
		for key, value := range vs {
			p.values[key] = value
		}
	case map[string]interface{}:
		for key, value := range vs {
			p.values[key] = value
		}
	default:
		p.values[DefaultPayload] = values
	}
	return p
}

// Len implementes the Payload interface.
func (p *payload) Len() int {
	return len(p.values)
}

// Get implementes the Payload interface.
func (p *payload) Get(key string) (interface{}, bool) {
	value, ok := p.values[key]
	return value, ok
}

// GetBool implementes the Payload interface.
func (p *payload) GetBool(key string) (bool, bool) {
	raw, ok := p.Get(key)
	if !ok {
		return false, ok
	}
	value, ok := raw.(bool)
	return value, ok
}

// GetInt implementes the Payload interface.
func (p *payload) GetInt(key string) (int, bool) {
	raw, ok := p.Get(key)
	if !ok {
		return 0, ok
	}
	value, ok := raw.(int)
	return value, ok
}

// GetFloat64 implementes the Payload interface.
func (p *payload) GetFloat64(key string) (float64, bool) {
	raw, ok := p.Get(key)
	if !ok {
		return 0.0, ok
	}
	value, ok := raw.(float64)
	return value, ok
}

// GetString implementes the Payload interface.
func (p *payload) GetString(key string) (string, bool) {
	raw, ok := p.Get(key)
	if !ok {
		return "", ok
	}
	value, ok := raw.(string)
	return value, ok
}

// GetTime implementes the Payload interface.
func (p *payload) GetTime(key string) (time.Time, bool) {
	raw, ok := p.Get(key)
	if !ok {
		return time.Time{}, ok
	}
	value, ok := raw.(time.Time)
	return value, ok
}

// GetDuration implementes the Payload interface.
func (p *payload) GetDuration(key string) (time.Duration, bool) {
	raw, ok := p.Get(key)
	if !ok {
		return time.Duration(0), ok
	}
	value, ok := raw.(time.Duration)
	return value, ok
}

// GetWaiter implements the Payload interface.
func (p *payload) GetWaiter(key string) (PayloadWaiter, bool) {
	raw, ok := p.Get(key)
	if !ok {
		return nil, ok
	}
	value, ok := raw.(PayloadWaiter)
	return value, ok
}

// Keys is specified on the Payload interface.
func (p *payload) Keys() []string {
	keys := []string{}
	for key := range p.values {
		keys = append(keys, key)
	}
	return keys
}

// Do implementes the Payload interface.
func (p *payload) Do(f func(key string, value interface{}) error) error {
	for key, value := range p.values {
		if err := f(key, value); err != nil {
			return err
		}
	}
	return nil
}

// Apply implementes the Payload interface.
func (p *payload) Apply(values interface{}) Payload {
	applied := &payload{
		values: PayloadValues{},
	}
	for key, value := range p.values {
		applied.values[key] = value
	}
	switch vs := values.(type) {
	case Payload:
		vs.Do(func(key string, value interface{}) error {
			applied.values[key] = value
			return nil
		})
	case PayloadValues:
		for key, value := range vs {
			applied.values[key] = value
		}
	case map[string]interface{}:
		for key, value := range vs {
			applied.values[key] = value
		}
	default:
		applied.values[DefaultPayload] = values
	}
	return applied
}

// String returns the payload represented as string.
func (p *payload) String() string {
	ps := []string{}
	for key, value := range p.values {
		ps = append(ps, fmt.Sprintf("<%q: %v>", key, value))
	}
	return strings.Join(ps, ", ")
}

// PayloadWaiter can be sent with an event as payload.
// Once a payload is set by a behavior Wait() continues
// and returns it.
type PayloadWaiter interface {
	// Set sets the payload somebody is waiting for.
	Set(p Payload)

	// Wait waits until the payload is set. A deadline
	// or timeout set by the context may cancel the
	// waiting.
	Wait(ctx context.Context) (Payload, error)
}

// payloadWaiter implements the PayloadWaiter interface.
type payloadWaiter struct {
	payloadc chan Payload
	once     sync.Once
}

// NewPayloadWaiter creates a new waiter for a payload
// returned by a behavior.
func NewPayloadWaiter() PayloadWaiter {
	return &payloadWaiter{
		payloadc: make(chan Payload, 1),
	}
}

// Set implements the PayloadWaiter interface.
func (w *payloadWaiter) Set(p Payload) {
	w.once.Do(func() {
		w.payloadc <- p
	})
}

// Wait implements the PayloadWaiter interface.
func (w *payloadWaiter) Wait(ctx context.Context) (Payload, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	for {
		select {
		case pl := <-w.payloadc:
			return pl, nil
		case <-ctx.Done():
			err := ctx.Err()
			return nil, err
		}
	}
}

//--------------------
// EVENT
//--------------------

// Event transports what to process.
type Event interface {
	fmt.Stringer

	// Context returns a Context that possibly has been
	// emitted with the event.
	Context() context.Context

	// Timestamp returns the UTC time the event has been created.
	Timestamp() time.Time

	// Topic returns the topic of the event.
	Topic() string

	// Payload returns the payload of the event.
	Payload() Payload
}

// event implements the Event interface.
type event struct {
	ctx       context.Context
	timestamp time.Time
	topic     string
	payload   Payload
}

// NewEvent creates a new event with the given topic and payload.
func NewEvent(ctx context.Context, topic string, payload interface{}) (Event, error) {
	if topic == "" {
		return nil, errors.New(ErrNoTopic, errorMessages)
	}
	p := NewPayload(payload)
	return &event{
		ctx:       ctx,
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
	return e.payload
}

// Context implements the Event interface.
func (e *event) Context() context.Context {
	return e.ctx
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

// EventSinkIterator can be used to check the events in a sink.
type EventSinkIterator interface {
	// Do iterates over all collected events.
	Do(doer func(index int, event Event) error) error

	// Match checks if all events match the passed criterion.
	Match(matcher func(index int, event Event) (bool, error)) (bool, error)
}

// EventSinkChecker can be used to check sinks for a criterion.
type EventSinkChecker func(events EventSinkIterator) (bool, error)

// EventSink stores a number of events ordered by adding. To be used
// in behaviors for collecting sets of events and operate on them.
type EventSink interface {
	// Add adds a new event data based on the passed event.
	Add(event Event) int

	// Len returns the number of stored events.
	Len() int

	// First returns the first of the collected events.
	First() (Event, bool)

	// Last returns the last of the collected event datas.
	Last() (Event, bool)

	// At returns an event at a given index and true if it
	// exists, otherwise nil and false.
	At(index int) (Event, bool)

	// Clear removes all collected events.
	Clear()

	EventSinkIterator
}

// eventSink implements the EventSink interface.
type eventSink struct {
	mutex   sync.RWMutex
	max     int
	events  []Event
	checker EventSinkChecker
	waiter  PayloadWaiter
}

// NewEventSink creates a sink for events.
func NewEventSink(max int) EventSink {
	return &eventSink{
		max: max,
	}
}

// NewCheckedEventSink creates a sink running a checker
// after each change.
func NewCheckedEventSink(max int, checker EventSinkChecker) (EventSink, PayloadWaiter) {
	waiter := NewPayloadWaiter()
	return &eventSink{
		max:     max,
		checker: checker,
		waiter:  waiter,
	}, waiter
}

// Add implements the EventSink interface.
func (s *eventSink) Add(event Event) int {
	s.mutex.Lock()
	s.events = append(s.events, event)
	if s.max > 0 && len(s.events) > s.max {
		s.events = s.events[1:]
	}
	s.mutex.Unlock()
	if s.checker != nil {
		ok, err := s.checker(s)
		if err != nil {
			// TODO
			return 0
		}
		if ok {
			s.waiter.Set(NewPayload(s))
		}
	}
	return len(s.events)
}

// Len implements the EventSink interface.
func (s *eventSink) Len() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.events)
}

// First implements the EventSink interface.
func (s *eventSink) First() (Event, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if len(s.events) < 1 {
		return nil, false
	}
	return s.events[0], true
}

// Last implements the EventSink interface.
func (s *eventSink) Last() (Event, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if len(s.events) < 1 {
		return nil, false
	}
	return s.events[len(s.events)-1], true
}

// At implements the EventSink interface.
func (s *eventSink) At(index int) (Event, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if index < 0 || index > len(s.events)-1 {
		return nil, false
	}
	return s.events[index], true
}

// Clear implements tne EventSink interface.
func (s *eventSink) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.events = nil
}

// Do implements the EventSinkIterator interface.
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

// Match implements the EventSinkIterator interface.
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

// EOF
