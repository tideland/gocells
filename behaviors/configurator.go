// Tideland Go Cells - Behaviors - Configurator
//
// Copyright (C) 2015-2016 Frank Mueller / Tideland / Oldenburg / Germany
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
// CONVENIENCE
//--------------------

// Configuration returns the configuration payload
// of the passed event or an empty configuration.
func Configuration(event cells.Event) etc.Etc {
	payload, ok := event.Payload().Get(ConfigurationPayload)
	if !ok {
		logger.Warningf("event does not contain configuration payload")
		cfg, _ := etc.ReadString("{etc}")
		return cfg
	}
	cfg, ok := payload.(etc.Etc)
	if !ok {
		logger.Warningf("configuration payload has illegal type")
		cfg, _ := etc.ReadString("{etc}")
		return cfg
	}
	return cfg
}

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
	case ReadConfigurationTopic:
		// Read configuration
		filename, ok := event.Payload().GetString(ConfigurationFilenamePayload)
		if !ok {
			logger.Errorf("cannot read configuration without filename payload")
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
		pvs := cells.PayloadValues{
			ConfigurationPayload: cfg,
		}
		b.cell.EmitNewContext(ConfigurationTopic, pvs, event.Context())
	}
	return nil
}

// Recover implements the cells.Behavior interface.
func (b *configuratorBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
