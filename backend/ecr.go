package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
)

// ECRService ...
var ECRService *ecr.ECR

func getLastDeployedCommit(images []*ecr.ImageIdentifier, digest string) string {
	var lastCommit string
	for _, image := range images {
		if *image.ImageDigest == digest {
			lastCommit = *image.ImageTag
		}
	}
	return lastCommit
}

func getStageDigest(images []*ecr.ImageIdentifier) string {
	const STAGE = "latest"
	var latestDigest string
	for _, image := range images {
		if *image.ImageTag == STAGE {
			latestDigest = *image.ImageDigest
		}
	}
	return latestDigest
}

func listImages(repoName string) []*ecr.ImageIdentifier {
	input := &ecr.ListImagesInput{
		RepositoryName: aws.String(repoName),
	}

	result, err := ECRService.ListImages(input)
	handleECRErr(err)

	images := make([]*ecr.ImageIdentifier, 0)
	for _, image := range result.ImageIds {
		if image.ImageTag != nil {
			images = append(images, image)
		}
	}

	return images
}

func filterImages(images []*ecr.ImageIdentifier) []string {
	versions := make([]string, 0)
	for _, image := range images {
		versions = append(versions, *image.ImageTag)
	}
	return versions
}

func getRepoURI(repoName string) string {
	var repoURI string
	repos := listRepos()
	for _, repo := range repos {
		repoValue := *repo.RepositoryName
		if repoValue == repoName {
			repoURIValue := *repo.RepositoryUri
			repoURI = repoURIValue
		}
	}
	return repoURI
}

func listRepos() []*ecr.Repository {
	input := &ecr.DescribeRepositoriesInput{}

	result, err := ECRService.DescribeRepositories(input)
	handleECRErr(err)

	return result.Repositories
}
