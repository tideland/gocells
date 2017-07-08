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

// EOF
