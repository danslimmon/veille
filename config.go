package main

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v1"
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
