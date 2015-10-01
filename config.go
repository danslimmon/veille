package main

import (
	"io/ioutil"
	"path/filepath"
	"time"

	yaml "gopkg.in/yaml.v2"
)

var configRead bool
var config Config

// CrunchOptions describes the metrics that will be created by CrunchMetrics.
type Config struct {
	// A list of windows (in days) for which to generate rolling downtime metrics.
	//
	// For example, []int{30,90,365} would result in rolling 30-, 90-, and 365-day
	// downtime metrics.
	Windows []int `yaml:"roll_windows"`
	// Interval (in minutes) between data points for the rolling downtime metrics.
	TSInterval int `yaml:"ts_interval"`
	// False positive patterns
	FoPos []*FoPoPattern `yaml:"false_positive_patterns"`
}

// FoPoPatterns provided in the config are matched against alerts to determine whether
// those alerts should be ignored as false positives.
//
// FoPoPattern implements the yaml.Unmarshaler interface.
type FoPoPattern struct {
	Start           time.Time
	End             time.Time
	ServicePatterns []string
}

func (fpp *FoPoPattern) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var err error
	var intermediary struct {
		StartStr        string   `yaml:"start"`
		EndStr          string   `yaml:"end"`
		ServicePatterns []string `yaml:"service_patterns"`
	}
	err = unmarshal(&intermediary)
	if err != nil {
		return err
	}

	fpp.Start, err = time.Parse("2006-01-02 15:04", intermediary.StartStr)
	if err != nil {
		return err
	}
	fpp.End, err = time.Parse("2006-01-02 15:04", intermediary.EndStr)
	if err != nil {
		return err
	}
	fpp.ServicePatterns = intermediary.ServicePatterns
	return nil
}

// Match determines whether the given service name matches any of the patterns
// provided in the FoPoPattern.
func (fpp *FoPoPattern) Match(serviceName string) bool {
	var pattern string
	var matched bool
	for _, pattern = range fpp.ServicePatterns {
		matched, _ = filepath.Match(pattern, serviceName)
		if matched {
			return matched
		}
	}
	return false
}

// getConfig parses the config file and returns the corresponding Config struct.
func getConfig() (Config, error) {
	var configPath string
	var err error

	if !configRead {
		configPath, err = getConfigPath()
		if err != nil {
			return Config{}, err
		}
		config, err = parseConfig(configPath)
		if err != nil {
			return Config{}, err
		}
		configRead = true
	}

	return config, nil
}

// getConfigPath() returns the path to the veille.yaml file
func getConfigPath() (string, error) {
	return "./veille.yaml", nil
}

// parseConfig() parses the YAML file with the given path, returning a Config struct
func parseConfig(path string) (Config, error) {
	var err error
	var yamlBlob []byte
	var conf Config

	yamlBlob, err = ioutil.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	err = yaml.Unmarshal(yamlBlob, &conf)
	if err != nil {
		return Config{}, err
	}
	return conf, nil
}
