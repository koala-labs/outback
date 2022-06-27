package docker

import "errors"

var (
	ErrImageBuild = errors.New("Could not build docker image")
	ErrImagePush  = errors.New("Could not push docker image. Are you logged in to ECR? http://docs.aws.amazon.com/AmazonECR/latest/userguide/Registries.html#registry_auth\nHint: `$(aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin)`\nDon't forget your --profile if you use one")
	ErrImagePull  = errors.New("Could not push docker image. Are you logged in to ECR? http://docs.aws.amazon.com/AmazonECR/latest/userguide/Registries.html#registry_auth\nHint: `$(aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin)`\nDon't forget your --profile if you use one")
)
