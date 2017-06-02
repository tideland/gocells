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

// List of criterion match signals.
const (
	CriterionDone CriterionMatch = iota + 1
	CriterionKeep
	CriterionDropFirst
	CriterionDropLast
	CriterionClear
)

// Standard topics.
const (
	TopicCollected = "collected"
	TopicCounted   = "counted"
	TopicProcess   = "process"
	TopicProcessed = "processed"
	TopicReset     = "reset"
	TopicStatus    = "status"
	TopicTick      = "tick"
)

// Standard playload keys.
const (
	PayloadClear      = "clear"
	PayloadDefault    = "default"
	PayloadDone       = "done"
	PayloadError      = "error"
	PayloadTickerID   = "ticker:id"
	PayloadTickerTime = "ticker:time"
)

// Default timeout for requests to cells.
const (
	DefaultTimeout = 5 * time.Second
)

const (
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
