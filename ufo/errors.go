package main

import (
	"errors"
)

// Config Errors
var (
	ErrCouldNotCreateConfig    = errors.New("Could not create config file.")
	ErrConfigFileAlreadyExists = errors.New("Config file already exists at the chosen location.")
	ErrNoEnvironments          = errors.New("No environments configured.")
	ErrNoGitIgnore             = errors.New("No .gitignore exists.")
)

var (
	ErrNoWorkingDirectory = errors.New("Could not resolve current working directory.")
	ErrCouldNotLoadConfig = errors.New("Could not load config file, please make sure it is valid JSON.")
	ErrGitError           = errors.New("Could not read git information. Please make sure you have git installed and are in a git repository.")
	ErrDeployTimeout      = errors.New("Timed out waiting for task to start.")
	ErrDockerBuild        = errors.New("Could not build docker image.")
	ErrDockerPush         = errors.New("Could not push docker image. Are you logged in to ECR? http://docs.aws.amazon.com/AmazonECR/latest/userguide/Registries.html#registry_auth\nHint: `$(aws ecr get-login --no-include-email --region us-west-1)`\nDon't forget your --profile if you use one!")
)
