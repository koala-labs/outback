package cmd

import (
	"errors"
	"fmt"
	"os"
)

// Config Errors
var (
	ErrNoEnvironments  = errors.New("No environments configured")
	ErrClusterNotFound = errors.New("Selected cluster could not be found. Please check your config")
	ErrServiceNotFound = errors.New("Selected service could not be found. Please check your config")
	ErrCommandNotFound = errors.New("Selected command could not be found. Please check your config")
)

// Deploy Errors
var (
	ErrDeployTimeout = errors.New("Timed out waiting for task to start")
	ErrDockerBuild   = errors.New("Could not build docker image")
	ErrDockerPush    = errors.New("Could not push docker image. Are you logged in to ECR? http://docs.aws.amazon.com/AmazonECR/latest/userguide/Registries.html#registry_auth\nHint: `$(aws ecr get-login --no-include-email --region us-west-1)`\nDon't forget your --profile if you use one")
)

// Version errors
var (
	ErrCouldNotAssertVersion = errors.New("Could not assert that UFO is up to date.")
	ErrUFOOutOfDate          = errors.New("UFO is out of date, please update to continue.")
)

// Init errors
var (
	ErrNoGitIgnore             = errors.New("No .gitignore exists")
	ErrNoWorkingDirectory      = errors.New("Could not resolve current working directory")
	ErrCouldNotLoadConfig      = errors.New("Could not load config file, please make sure it is valid JSON")
	ErrCouldNotOpenGitIgnore   = errors.New("Could not open .gitignore. Is there a .gitignore in the current directory?")
	ErrCouldNotCreateConfig    = errors.New("Could not create config file")
	ErrConfigFileAlreadyExists = errors.New("Config file already exists at the chosen location")
)

// Service errors
var (
	ErrInvalidEnvInput = errors.New("Input must be in the form of key=value")
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
