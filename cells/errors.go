// Tideland Go Cells - Errors
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
	"github.com/tideland/golib/errors"
)

//--------------------
// CONSTANTS
//--------------------

// Error codes of the cells package.
const (
	ErrCellInit = iota + 1
	ErrCannotRecover
	ErrCannotEmit
	ErrDuplicateID
	ErrInvalidID
	ErrExecuteID
	ErrEventRecovering
	ErrRecoveredTooOften
	ErrNoTopic
	ErrNoRequest
	ErrInactive
	ErrStopping
	ErrTimeout
	ErrMissingScene
)

// Error messages of the cells package.
var errorMessages = map[int]string{
	ErrCellInit:          "cell %q cannot initialize",
	ErrCannotRecover:     "cannot recover cell %q: %v",
	ErrCannotEmit:        "cannot emit event into queue",
	ErrDuplicateID:       "cell with ID %q is already registered",
	ErrInvalidID:         "cell with ID %q does not exist",
	ErrExecuteID:         "cannot %s with cell %q",
	ErrEventRecovering:   "cell cannot recover after error %v",
	ErrRecoveredTooOften: "cell needs too much recoverings, last error",
	ErrNoTopic:           "event has no topic",
	ErrNoRequest:         "cannot respond, event is no request",
	ErrInactive:          "cell %q is inactive",
	ErrStopping:          "%s is stopping",
	ErrTimeout:           "needed too long for %v",
	ErrMissingScene:      "missing scene for request",
}

//--------------------
// ERROR CHECKING
//--------------------

// NewCannotRecoverError returns an error showing that a cell cannot
// recover from errors.
func NewCannotRecoverError(id string, err interface{}) error {
	return errors.New(ErrCannotRecover, errorMessages, id, err)
}

// EOF
