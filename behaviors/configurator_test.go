// Tideland Go Cells - Behaviors - Unit Tests - Configurator
//
// Copyright (C) 2015 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors_test

//--------------------
// IMPORTS
//--------------------

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/tideland/golib/audit"
	"github.com/tideland/golib/etc"

	"github.com/tideland/gocells/behaviors"
	"github.com/tideland/gocells/cells"
)

//--------------------
// TESTS
//--------------------

// TestConfigurationRead tests the successful reading of a configuration.
func TestConfigurationRead(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	env := cells.NewEnvironment("configuration-read")
	defer env.Stop()
	tempDir, filename := createConfigurationFile(assert, "{etc {foo 42}}")
	defer tempDir.Restore()

	sigc := audit.MakeSigChan()
	spf := func(c cells.Cell, event cells.Event) error {
		if event.Topic() == behaviors.ConfigurationTopic {
			cfg := behaviors.Configuration(event)
			v := cfg.ValueAsString("foo", "0815")
			assert.Equal(v, "42")

			sigc <- true
		}
		return nil
	}

	env.StartCell("configurator", behaviors.NewConfiguratorBehavior(nil))
	env.StartCell("simple", behaviors.NewSimpleProcessorBehavior(spf))
	env.Subscribe("configurator", "simple")

	pvs := cells.PayloadValues{
		behaviors.ConfigurationFilenamePayload: filename,
	}
	env.EmitNew("configurator", behaviors.ReadConfigurationTopic, pvs, nil)
	assert.Wait(sigc, true, 100*time.Millisecond)
}

// TestConfigurationValidation tests the validation of a configuration.
func TestConfigurationValidation(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	env := cells.NewEnvironment("configuration-validation")
	defer env.Stop()
	tempDir, filename := createConfigurationFile(assert, "{etc {foo 42}}")
	defer tempDir.Restore()

	sigc := audit.MakeSigChan()
	spf := func(c cells.Cell, event cells.Event) error {
		sigc <- true
		return nil
	}
	var key string
	cv := func(cfg etc.Etc) error {
		v := cfg.ValueAsString(key, "[-default-]")
		if v == "[-default-]" {
			sigc <- false
		}
		return nil
	}

	env.StartCell("configurator", behaviors.NewConfiguratorBehavior(cv))
	env.StartCell("simple", behaviors.NewSimpleProcessorBehavior(spf))
	env.Subscribe("configurator", "simple")

	// First run with success as key has the valid value "foo"
	pvs := cells.PayloadValues{
		behaviors.ConfigurationFilenamePayload: filename,
	}
	key = "foo"
	env.EmitNew("configurator", behaviors.ReadConfigurationTopic, pvs, nil)
	assert.Wait(sigc, true, 100*time.Millisecond)

	// Second run also will succeed, even with "bar" as invalid value.
	// See definition of validator cv above. But validationError is not
	// nil.
	key = "bar"
	env.EmitNew("configurator", behaviors.ReadConfigurationTopic, pvs, nil)
	assert.Wait(sigc, false, 100*time.Millisecond)
}

//--------------------
// HELPER
//--------------------

// createConfigurationFile creates a temporary configuration file.
func createConfigurationFile(assert audit.Assertion, content string) (*audit.TempDir, string) {
	tempDir := audit.NewTempDir(assert)
	configFile, err := ioutil.TempFile(tempDir.String(), "configuration")
	assert.Nil(err)
	configFileName := configFile.Name()
	_, err = configFile.WriteString(content)
	assert.Nil(err)
	configFile.Close()

	return tempDir, configFileName
}

// EOF
