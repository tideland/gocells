// Tideland Go Cells - Behaviors - Topic/Payloads Behavior
//
// Copyright (C) 2010-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors

//--------------------
// IMPORTS
//--------------------

import "github.com/tideland/gocells/cells"

//--------------------
// TOPIC/PAYLOADS BEHAVIOR
//--------------------

// TopicPayloadsProcessor processes the collected payloads of a topic and
// returns a new payload to emit.
type TopicPayloadsProcessor func(topic string, payloads []cells.Payload) (cells.Payload, error)

// topicPayloadsBehavior collects and processes payloads by topic.
type topicPayloadsBehavior struct {
	cell      cells.Cell
	maximum   int
	collected map[string][]cells.Payload
	process   TopicPayloadsProcessor
}

// NewTopicPayloadsBehavior creates a behavior collecting the payloads
// of events by their topics, processes them, and emit the processing
// result.
func NewTopicPayloadsBehavior(maximum int, processor TopicPayloadsProcessor) cells.Behavior {
	return &topicPayloadsBehavior{
		cell:      nil,
		maximum:   maximum,
		collected: make(map[string][]cells.Payload),
		process:   processor,
	}
}

// Init the behavior.
func (b *topicPayloadsBehavior) Init(c cells.Cell) error {
	b.cell = c
	return nil
}

// Terminate the behavior.
func (b *topicPayloadsBehavior) Terminate() error {
	return nil
}

// ProcessEvent calls the simple processor function.
func (b *topicPayloadsBehavior) ProcessEvent(event cells.Event) error {
	topic := event.Topic()
	payloads := b.collected[topic]
	payloads = append(payloads, event.Payload())
	if len(payloads) > b.maximum {
		payloads = payloads[1:]
	}
	b.collected[topic] = payloads
	payload, err := b.process(topic, b.collected[topic])
	if err != nil {
		return err
	}
	return b.cell.EmitNew(topic, payload)
}

// Recover from an error.
func (b *topicPayloadsBehavior) Recover(err interface{}) error {
	b.collected = make(map[string][]cells.Payload)
	return nil
}

// EOF
