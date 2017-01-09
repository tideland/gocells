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
	cell      cells.Cell
	matches   PairCriterion
	duration  time.Duration
	firstHit  *time.Time
	secondHit *time.Time
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
	return nil
}

// Recover from an error.
func (b *pairBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
