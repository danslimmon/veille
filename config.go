package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

// All-caps version will be used as prefix for environment variables
const appName = "veille"

var Config *ConfStruct

type ConfStruct struct {
	Listen             string
	AwsAccessKeyId     string `envconfig:"aws_access_key_id"`
	AwsSecretAccessKey string `envconfig:"aws_secret_access_key"`
}

// LoadConfig loads the configuration from environment variables.
func LoadConfig() error {
	if Config != nil {
		return nil
	}

	Config = new(ConfStruct)
	if err := envconfig.Process(appName, Config); err != nil {
		return err
	}

	if Config.Listen == "" {
		listenDefault := ":8080"
		log.Printf("Environment variable `VEILLE_LISTEN` undefined; using default of '%s'", listenDefault)
		Config.Listen = listenDefault
	}
	if Config.AwsAccessKeyId == "" {
		log.Fatal("Environment variable `AWS_ACCESS_KEY_ID` is required (but may be set to 'dummy')")
	}
	if Config.AwsSecretAccessKey == "" {
		log.Fatal("Environment variable `AWS_SECRET_ACCESS_KEY` is required (but may be set to 'dummy')")
	}

	return nil
}
