package cmd

import (
	"errors"
)

// Config Errors
var (
	ErrNoEnvironments  = errors.New("No environments configured")
	ErrClusterNotFound = errors.New("Selected cluster could not be found")
	ErrServiceNotFound = errors.New("Selected service could not be found")
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

var (
	ErrEmptyCluster = errors.New("Please enter a cluster")
)
