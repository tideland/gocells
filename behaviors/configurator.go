// Tideland Go Cells - Behaviors - Configurator
//
// Copyright (C) 2015-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors

//--------------------
// IMPORTS
//--------------------

import (
	"github.com/tideland/golib/errors"
	"github.com/tideland/golib/etc"
	"github.com/tideland/golib/logger"

	"github.com/tideland/gocells/cells"
)

//--------------------
// CONSTANTS
//--------------------

const (
	// TopicConfiguration contains the topic for the configuration behavior.
	TopicConfiguration = "configuration"

	// TopicConfigurationRead tells the configurator behavior to
	// read a configuration file.
	TopicConfigurationRead = "read-configuration!"

	// PayloadConfiguration contains the emitted configuration.
	PayloadConfiguration = "configuration"

	// PayloadConfigurationFilename contains the name of the
	// configuration file to read.
	PayloadConfigurationFilename = "configuration:filename"
)

//--------------------
// CONFIGURATOR BEHAVIOR
//--------------------

// ConfigurationValidator defines a function for the validation of
// a new read configuration.
type ConfigurationValidator func(etc.Etc) error

// configuratorBehavior implements the configurator behavior.
type configuratorBehavior struct {
	cell     cells.Cell
	validate ConfigurationValidator
}

// NewConfiguratorBehavior creates the configurator behavior. It loads a
// configuration file and emits the it to its subscribers. If a validator
// is passed the read configuration will be validated using it. Errors
// will be logged.
func NewConfiguratorBehavior(validator ConfigurationValidator) cells.Behavior {
	// TODO 2017-05-30 Mue Think about concept of configurator. Maybe better
	// read it lazy, keep it, re-read on demand, and emit only wanted
	// values of it.
	return &configuratorBehavior{
		validate: validator,
	}
}

// Init implements the cells.Behavior interface.
func (b *configuratorBehavior) Init(c cells.Cell) error {
	b.cell = c
	return nil
}

// Terminate implements the cells.Behavior interface.
func (b *configuratorBehavior) Terminate() error {
	return nil
}

// ProcessEvent reads, validates and emits a configuration.
func (b *configuratorBehavior) ProcessEvent(event cells.Event) error {
	switch event.Topic() {
	case TopicConfigurationRead:
		// Read configuration
		filename, ok := event.Payload().Get(PayloadConfigurationFilename)
		if !ok {
			logger.Errorf("cannot read configuration without filename")
			return nil
		}
		logger.Infof("reading configuration from %q", filename)
		cfg, err := etc.ReadFile(filename)
		if err != nil {
			return errors.Annotate(err, ErrCannotReadConfiguration, errorMessages)
		}
		// If wanted then validate it.
		if b.validate != nil {
			err = b.validate(cfg)
			if err != nil {
				return errors.Annotate(err, ErrCannotValidateConfiguration, errorMessages)
			}
		}
		// All done, emit it.
		payloadValues := cells.PayloadValues{
			PayloadConfiguration: cfg,
		}
		b.cell.EmitNew(TopicConfiguration, cells.NewPayload(payloadValues))
	}
	return nil
}

// Recover implements the cells.Behavior interface.
func (b *configuratorBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
