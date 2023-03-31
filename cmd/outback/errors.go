package cmd

import (
	"errors"
	"fmt"
	"os"
)

// Config Errors
var (
	ErrClusterNotFound = errors.New("Selected cluster could not be found. Please check your config")
	ErrServiceNotFound = errors.New("Selected service could not be found. Please check your config")
	ErrCommandNotFound = errors.New("Selected command could not be found. Please check your config")
)

// Deploy Errors
var (
	ErrDeployTimeout = errors.New("Timed out waiting for task to start")
)

// Init errors
var (
	ErrCouldNotCreateConfig    = errors.New("Could not create config file")
	ErrConfigFileAlreadyExists = errors.New("Config file already exists at the chosen location")
)

// Service errors
var (
	ErrInvalidEnvInput       = errors.New("Input must be in the form of key=value")
	ErrKeyNotPresent         = errors.New("The key entered was not present in the environment variables for this service")
	ErrCouldNotParseTime     = errors.New("Could not parse the given time")
	ErrCantFollowWithEndTime = errors.New("Could not follow logs because an end time was given")
)

var (
	ErrorMissingClusterInput = errors.New("A cluster name must be specified using the --cluster flag")
	ErrorMissingServiceInput = errors.New("A service name must be specified using the --service flag")
)

// handleError is intended to be called with an error return to simplify error handling
// Usage:
// foo, err := GetFoo()
// HandleError(err)
// DoSomethingBecauseNoError()
func handleError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\nEncountered an error: %s\n", err.Error())

	os.Exit(1)
}
