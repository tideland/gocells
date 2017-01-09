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
type PairCriterion func(event cells.Event) bool

// pairBehavior checks if events occur in pairs.
type pairBehavior struct {
	cell     cells.Cell
	matches  PairCriterion
	duration time.Duration
	hit      *time.Time
	timeout  *time.Timer
}

// NewPairBehavior creates ...
func NewPairBehavior(matches PairCriterion, duration time.Duration) cells.Behavior {
	return &pairBehavior{nil, matches, duration, nil, nil}
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
		if b.hit != nil {
			// Otherwise it has been reset already, just a queued event.
			b.cell.EmitNew(EventPairTimeoutTopic, cells.PayloadValues{
				EventPairFirstPayload:   b.hit,
				EventPairTimeoutPayload: time.Now(),
			})
			b.hit = nil
		}
	default:
		if b.matches(event) {
			now := time.Now()
			if b.hit == nil {
				// First hit, store time and start timeout reminder.
				b.hit = &now
				b.timeout = time.AfterFunc(b.duration, func() {
					b.cell.Environment().EmitNew(b.cell.ID(), EventPairTimeoutTopic, nil)
				})
			} else {
				// Second hit earlier than timeout, fine.
				b.timeout.Stop()
				b.cell.EmitNew(EventPairTopic, cells.PayloadValues{
					EventPairFirstPayload:  b.hit,
					EventPairSecondPayload: now,
				})
				b.hit = nil
			}
		}
	}
	return nil
}

// Recover from an error.
func (b *pairBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
