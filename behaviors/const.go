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

// Topics and payload keys.
const (
	ResetTopic             = "reset!"
	ReadConfigurationTopic = "read-configuration!"
	ConfigurationTopic     = "configuration"
	TickerTopic            = "tick!"
	EventPairTopic         = "event-pair!"
	EventPairTimeoutTopic  = "event-pair-timeout!"
	EventRateTopic         = "event-rate!"

	ConfigurationFilenamePayload = "configuration:filename"
	ConfigurationPayload         = "configuration"
	TickerIDPayload              = "ticker:id"
	TickerTimePayload            = "ticker:time"
	EventPairFirstPayload        = "event-pair:first"
	EventPairSecondPayload       = "event-pair:second"
	EventPairTimeoutPayload      = "event-pair:timeout"
	EventRateTimePayload         = "event-rate:time"
	EventRateDurationPayload     = "event-rate:duration"
	EventRateAveragePayload      = "event-rate:average"
	EventRateHighPayload         = "event-rate:high"
	EventRateLowPayload          = "event-rate:low"
)

// EOF
