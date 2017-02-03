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
	TopicConfiguration     = "configuration"
	TopicConfigurationRead = "read-configuration!"
	TopicEvaluation        = "evaluation"
	TopicPair              = "pair"
	TopicPairTimeout       = "pair-timeout"
	TopicRate              = "rate"
	TopicRateWindow        = "rate-window"
	TopicReset             = "reset!"
	TopicSequence          = "sequence"
	TopicTicker            = "tick!"

	// Payload keys.
	PayloadConfiguration         = "configuration"
	PayloadConfigurationFilename = "configuration:filename"
	PayloadEvaluationCount       = "evaluation:count"
	PayloadEvaluationAvg         = "evaluation:avg"
	PayloadEvaluationMax         = "evaluation:max"
	PayloadEvaluationMin         = "evaluation:min"
	PayloadPairFirstData         = "pair:first:data"
	PayloadPairFirstTime         = "pair:first:time"
	PayloadPairSecondData        = "pair:second:data"
	PayloadPairSecondTime        = "pair:second:time"
	PayloadPairTimeout           = "pair:timeout"
	PayloadRateDuration          = "rate:duration"
	PayloadRateTime              = "rate:time"
	PayloadRateAverage           = "rate:average"
	PayloadRateHigh              = "rate:high"
	PayloadRateLow               = "rate:low"
	PayloadRateWindowCount       = "rate-window:count"
	PayloadRateWindowFirstTime   = "rate-window:first:time"
	PayloadRateWindowLastTime    = "rate-window:last:time"
	PayloadSequenceEvents        = "sequence:events"
	PayloadTickerID              = "ticker:id"
	PayloadTickerTime            = "ticker:time"
)

// EOF
