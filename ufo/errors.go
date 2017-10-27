package main

import (
	"errors"
	"fmt"
)

// Config Errors
var (
	ErrNoEnvironments = errors.New("No environments configured.")
)

var (
	ErrCouldNotLoadConfig = errors.New(fmt.Sprintf("Could not load config file. Please make sure it is located in %s.", UFO_CONFIG))
	ErrEnvironmentForBranchDoesNotExist = errors.New("Could not find environment for current branch.")
	ErrGitError = errors.New("Could not read git information. Please make sure you have git installed and are in a git repository.")
	ErrDeployTimeout = errors.New("Timed out waiting for task to start.")
	ErrDockerBuild = errors.New("Could not build docker image.")
	ErrDockerPush = errors.New("Could not push docker image. Are you logged in the ECR? http://docs.aws.amazon.com/AmazonECR/latest/userguide/Registries.html#registry_auth\nHint: `$(aws ecr get-login --no-include-email --region us-west-1)`\nDon't forget your --profile if you use one!")
)