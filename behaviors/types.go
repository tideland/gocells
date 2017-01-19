// Tideland Go Cells - Behaviors - Types
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
// TYPES
//--------------------

// EventData represents the pure collected event data.
type EventData struct {
	Timestamp time.Time
	Topic     string
	Payload   cells.Payload
}

// newEventData returns the passed event as event data to collect.
func newEventData(event cells.Event) EventData {
	data := EventData{
		Timestamp: time.Now(),
		Topic:     event.Topic(),
		Payload:   event.Payload(),
	}
	return data
}

// EOF
