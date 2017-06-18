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
// CONSTANTS
//--------------------

const (
	// TopicPair signals a detected pair of events.
	TopicPair = "pair"

	// TopicPairTimeout signals a timeout during waiting for a
	// pair of events.
	TopicPairTimeout = "pair:timeout"
)

//--------------------
// PAIR BEHAVIOR
//--------------------

// PairCriterion is used by the pair behavior and has to return true, if
// the passed event matches a criterion for rate measuring. The returned
// data in case of a first hit is stored and then passed as argument to
// each further call of the pair criterion. In case of a pair event both
// returned datas are part of the emitted event payload.
type PairCriterion func(event cells.Event, hit cells.Payload) (cells.Payload, bool)

// Pair contains event pair information.
type Pair struct {
	FirstTime     time.Time
	FirstPayload  cells.Payload
	SecondTime    time.Time
	SecondPayload cells.Payload
	Timeout       time.Time
}

// pairBehavior checks if events occur in pairs.
type pairBehavior struct {
	cell       cells.Cell
	matches    PairCriterion
	duration   time.Duration
	hitTime    *time.Time
	hitPayload cells.Payload
	timeout    *time.Timer
}

// NewPairBehavior creates a behavior checking if two events match a criterion
// defined by the PairCriterion function and the duration between them is not
// longer than the passed duration. In case of a positive pair match an according
// event containing both timestamps and both returned datas is emitted. In case
// of a timeout a timeout event is emitted. It's payload is the first timestamp,
// the first data, and the timestamp of the timeout.
func NewPairBehavior(criterion PairCriterion, duration time.Duration) cells.Behavior {
	return &pairBehavior{
		cell:       nil,
		matches:    criterion,
		duration:   duration,
		hitTime:    nil,
		hitPayload: nil,
		timeout:    nil,
	}
}

// Init implements the cells.Behavior interface.
func (b *pairBehavior) Init(c cells.Cell) error {
	b.cell = c
	return nil
}

// Terminate implements the cells.Behavior interface.
func (b *pairBehavior) Terminate() error {
	return nil
}

// ProcessEvent implements the cells.Behavior interface.
func (b *pairBehavior) ProcessEvent(event cells.Event) error {
	switch event.Topic() {
	case TopicPairTimeout:
		if b.hitTime != nil && b.timeout != nil {
			// Received timeout event, check if the expected one.
			var first time.Time
			err := event.Payload().Unmarshal(&first)
			if err != nil {
				return err
			}
			if first.Equal(*b.hitTime) {
				b.emitTimeout()
				b.timeout = nil
			}
		}
	default:
		if payload, ok := b.matches(event, b.hitPayload); ok {
			now := time.Now()
			if b.hitTime == nil {
				// First hit, store time and data and start timeout reminder.
				b.hitTime = &now
				b.hitPayload = payload
				b.timeout = time.AfterFunc(b.duration, func() {
					b.cell.Environment().EmitNew(b.cell.ID(), TopicPairTimeout, now)
				})
			} else {
				// Second hit earlier than timeout event.
				// Check if it is in time.
				b.timeout.Stop()
				b.timeout = nil
				if now.Sub(*b.hitTime) > b.duration {
					b.emitTimeout()
				} else {
					b.emitPair(now, payload)
				}
			}
		}
	}
	return nil
}

// Recover implements the cells.Behavior interface.
func (b *pairBehavior) Recover(err interface{}) error {
	return nil
}

// emitPair emits the event for a successful pair.
func (b *pairBehavior) emitPair(timestamp time.Time, payload cells.Payload) {
	b.cell.EmitNew(TopicPair, Pair{
		FirstTime:     *b.hitTime,
		FirstPayload:  b.hitPayload,
		SecondTime:    timestamp,
		SecondPayload: payload,
	})
	b.hitTime = nil
}

// emitTimeout emits the event for a pairing timeout.
func (b *pairBehavior) emitTimeout() {
	b.cell.EmitNew(TopicPairTimeout, Pair{
		FirstTime:    *b.hitTime,
		FirstPayload: b.hitPayload,
		Timeout:      time.Now(),
	})
	b.hitTime = nil
}

// EOF
