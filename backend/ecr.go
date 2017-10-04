package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
)

// ECRService ...
var ECRService *ecr.ECR

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

func describeImages(repoName string) []*ecr.ImageDetail {
	input := &ecr.DescribeImagesInput{
		RepositoryName: aws.String(repoName),
	}

	result, err := ECRService.DescribeImages(input)
	handleECRErr(err)

	images := make([]*ecr.ImageDetail, 0)
	for _, image := range result.ImageDetails {
		if image.ImageTags != nil {
			images = append(images, image)
		}
	}

	return images
}

func filterImages(images []*ecr.ImageDetail) []string {
	versions := make([]string, 0)
	for _, image := range images {
		versions = append(versions, *image.ImageTags[0])
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
