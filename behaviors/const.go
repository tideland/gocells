// Tideland Go Cells - Behaviors - Constants
//
// Copyright (C) 2010-2016 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors

//--------------------
// CONSTANTS
//--------------------

// Topics and payload keys.
const (
	ReadConfigurationTopic = "read-configuration!"
	ConfigurationTopic     = "configuration"
	TickerTopic            = "tick!"
	EventRateTopic         = "event-rate!"

	ConfigurationFilenamePayload = "configuration:filename"
	ConfigurationPayload         = "configuration"
	TickerIDPayload              = "ticker:id"
	TickerTimePayload            = "ticker:time"
	EventRateAveragePayload      = "event-rate:average"
	EventRateHighPayload         = "event-rate:high"
	EventRateLowPayload          = "event-rate:low"
)

// EOF
