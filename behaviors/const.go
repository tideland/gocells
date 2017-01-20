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
	EventRateWindowTopic   = "event-rate-window!"
	EventSequenceTopic     = "event-sequence!"

	ConfigurationFilenamePayload    = "configuration:filename"
	ConfigurationPayload            = "configuration"
	TickerIDPayload                 = "ticker:id"
	TickerTimePayload               = "ticker:time"
	EventPairFirstTimePayload       = "event-pair:first:time"
	EventPairSecondTimePayload      = "event-pair:second:time"
	EventPairFirstDataPayload       = "event-pair:first:data"
	EventPairSecondDataPayload      = "event-pair:second:data"
	EventPairTimeoutPayload         = "event-pair:timeout"
	EventRateTimePayload            = "event-rate:time"
	EventRateDurationPayload        = "event-rate:duration"
	EventRateAveragePayload         = "event-rate:average"
	EventRateHighPayload            = "event-rate:high"
	EventRateLowPayload             = "event-rate:low"
	EventRateWindowCountPayload     = "event-rate-window:count"
	EventRateWindowFirstTimePayload = "event-rate-window:first:time"
	EventRateWindowLastTimePayload  = "event-rate-window:last:time"
	EventSequenceEventsPayload      = "event-sequence:events"
)

// EOF
