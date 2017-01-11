// Tideland Go Cells - Behaviors - Pair
//
// Copyright (C) 2010-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors

//--------------------
// IMPORTS
//--------------------

import (
	"time"

	"github.com/tideland/gocells/cells"
)

//--------------------
// PAIR BEHAVIOR
//--------------------

// PairCriterion is used by the pair behavior and has to return true, if
// the passed event matches a criterion for rate measuring.
type PairCriterion func(event cells.Event, hitData interface{}) (interface{}, bool)

// pairBehavior checks if events occur in pairs.
type pairBehavior struct {
	cell     cells.Cell
	matches  PairCriterion
	duration time.Duration
	hit      *time.Time
	hitData  interface{}
	timeout  *time.Timer
}

// NewPairBehavior creates ...
func NewPairBehavior(matches PairCriterion, duration time.Duration) cells.Behavior {
	return &pairBehavior{
		cell:     nil,
		matches:  matches,
		duration: duration,
		hit:      nil,
		hitData:  nil,
		timeout:  nil,
	}
}

// Init the behavior.
func (b *pairBehavior) Init(c cells.Cell) error {
	b.cell = c
	return nil
}

// Terminate the behavior.
func (b *pairBehavior) Terminate() error {
	return nil
}

// ProcessEvent collects and re-emits events.
func (b *pairBehavior) ProcessEvent(event cells.Event) error {
	switch event.Topic() {
	case EventPairTimeoutTopic:
		if b.hit != nil && b.timeout != nil {
			// Received timeout event, check if the expected one.
			hit, ok := event.Payload.GetTine(EventPairFirstTimePayload)
			if ok && hit.Equal(*b.hit) {
				b.emitTimeout("timeout event")
				b.timeout = nil
			}
		}
	default:
		if hitData, ok := b.matches(event, b.hitData); ok {
			now := time.Now()
			if b.hit == nil {
				// First hit, store time and data and start timeout reminder.
				b.hit = &now
				b.hitData = hitData
				b.timeout = time.AfterFunc(b.duration, func() {
					b.cell.Environment().EmitNew(b.cell.ID(), EventPairTimeoutTopic, cells.PayloadValues{
						EventPairFirstTimePayload: *b.hit,
					})
				})
			} else {
				// Second hit earlier than timeout event.
				// Check if it is in time.
				b.timeout.Stop()
				b.timeout = nil
				if now.Sub(*b.hit) > b.duration {
					b.emitTimeout("hit but too late")
				} else {
					b.emitPair(now, hitData)
				}
			}
		}
	}
	return nil
}

// Recover from an error.
func (b *pairBehavior) Recover(err interface{}) error {
	return nil
}

// emitPair emits the event for a successful pair.
func (b *pairBehavior) emitPair(timestamp time.Time, data interface{}) {
	b.cell.EmitNew(EventPairTopic, cells.PayloadValues{
		EventPairFirstTimePayload:  *b.hit,
		EventPairFirstDataPayload:  b.hitData,
		EventPairSecondTimePayload: timestamp,
		EventPairSecondDataPayload: data,
	})
	b.hit = nil
}

// emitTimeout emits the event for a pairing timeout.
func (b *pairBehavior) emitTimeout(debug string) {
	b.cell.EmitNew(EventPairTimeoutTopic, cells.PayloadValues{
		EventPairFirstTimePayload: *b.hit,
		EventPairFirstDataPayload: b.hitData,
		EventPairTimeoutPayload:   time.Now(),
		"debug":                   debug,
	})
	b.hit = nil
}

// EOF
