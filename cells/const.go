// Tideland Go Cells - Constants
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
	"time"
)

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
)

const (
	// Often used standard topics.
	TopicCollected = "collected?"
	TopicCounters  = "counters?"
	TopicProcessed = "processed?"
	TopicReset     = "reset!"
	TopicStatus    = "status?"
	TopicTick      = "tick!"

	// Standard payload keys.
	PayloadDefault    = "default"
	PayloadTickerID   = "ticker:id"
	PayloadTickerTime = "ticker:time"

	// Default timeout for requests to cells.
	DefaultTimeout = 5 * time.Second

	// minEventBufferSize is the minimum size of the
	// event buffer per cell.
	minEventBufferSize = 16

	// minRecoveringNumber and minRecoveringDuration
	// control the default recovering frequency.
	minRecoveringNumber   = 10
	minRecoveringDuration = time.Second

	// minEmitTimeout is the minimum allowed timeout
	// for event emitting (see below).
	minEmitTimeout = 5 * time.Second

	// maxEmitTimeout is the maximum time to emit an
	// event into a cells event buffer before a timeout
	// error is returned to the emitter.
	maxEmitTimeout = 30 * time.Second
)

// EOF
