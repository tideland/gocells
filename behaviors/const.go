// Tideland Go Cells - Behaviors - Constants
//
// Copyright (C) 2010-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors

//--------------------
// CONSTANTS
//--------------------

// CriterionMatch signals, how a criterion matches.
type CriterionMatch int

const (
	// Criterion matches.
	CriterionDone CriterionMatch = iota + 1
	CriterionKeep
	CriterionDropFirst
	CriterionDropLast
	CriterionClear

	// Topics.
	TopicSequence = "sequence"
	TopicTicker   = "tick!"

	// Payload keys.
	PayloadSequenceEvents = "sequence:events"
	PayloadTickerID       = "ticker:id"
	PayloadTickerTime     = "ticker:time"
)

// EOF
