package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
)

// All-caps version will be used as prefix for environment variables
const appName = "veille"

var Config *ConfStruct

// ConfStruct represents the configuration for Veille, including any
// input-mode-specific subconfigurations.
type ConfStruct struct {
	// LogLevel determines how much log output to produce. Currently
	// 'debug', 'info' (the default), and 'fatal'.
	LogLevel string `envconfig:"LOG_LEVEL"`
	// InputMode indicates the mode that will be used to read Nagios
	// log files.
	InputMode string `envconfig:"INPUT_MODE"`
	// TailPath contains the path of the file to tail (if InputMode
	// is "tail".
	TailPath string `envconfig:"TAIL_PATH"`
}

// LoadConfig loads the configuration from environment variables.
func LoadConfig() error {
	var err error
	if Config != nil {
		return nil
	}

	Config = new(ConfStruct)
	if err = envconfig.Process(appName, Config); err != nil {
		return err
	}

	switch Config.LogLevel {
	case "info", "":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	default:
		panic(fmt.Sprintf("unknown log level '%s'", Config.LogLevel))
	}

	switch Config.InputMode {
	default:
		panic(fmt.Sprintf("Unknown input mode '%s'", Config.InputMode))
	case "tail":
		err = TailValidateConfig(Config)
		log.WithFields(log.Fields{
			"input_mode": Config.InputMode,
		}).Debug("Loaded configuration from environment")
	}

	return nil
}
