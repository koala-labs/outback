package main

import (
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
)

// ECSService ...
var ECSService *ecs.ECS

func updateService(cluster string, service string, taskDefinitionARN string) *ecs.UpdateServiceOutput {
	input := &ecs.UpdateServiceInput{
		Cluster:        aws.String(cluster),
		Service:        aws.String(service),
		TaskDefinition: aws.String(taskDefinitionARN),
	}

	result, err := ECSService.UpdateService(input)
	handleECSErr(err)

	return result
}

func describeService(cluster string, service string) *ecs.Service {
	input := &ecs.DescribeServicesInput{
		Cluster: aws.String(cluster),
		Services: []*string{
			aws.String(service),
		},
	}

	result, err := ECSService.DescribeServices(input)
	handleECSErr(err)

	return result.Services[0]
}

func registerNewTaskDefinition(service string, version string) (string, string) {
	latestDefinitions := describeLatestTaskDefinition(service)
	input := &ecs.RegisterTaskDefinitionInput{
		// Update the task definition to use the new docker image via updateTaskDefinition
		ContainerDefinitions: updateTaskDefinition(latestDefinitions.ContainerDefinitions, service, version),
		Family:               latestDefinitions.Family,
	}

	result, err := ECSService.RegisterTaskDefinition(input)
	handleECSErr(err)

	taskDefinitionValue := *result.TaskDefinition.TaskDefinitionArn

	return service, taskDefinitionValue
}

func updateTaskDefinition(definitions []*ecs.ContainerDefinition, image string, version string) []*ecs.ContainerDefinition {
	r := regexp.MustCompile(`\/(\S+):`)
	for _, containerDefinition := range definitions {
		containerDefinitionImage := *containerDefinition.Image
		if r.FindStringSubmatch(containerDefinitionImage)[1] == image {
			repo := getRepoURI(image)
			newImage := fmt.Sprintf("%s:%s", repo, version)
			*containerDefinition.Image = newImage
		}
	}
	return definitions
}

func getLastDeployedCommit(taskDefinition string) string {
	input := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(taskDefinition),
	}

	result, err := ECSService.DescribeTaskDefinition(input)
	handleECSErr(err)

	repo := result.TaskDefinition.ContainerDefinitions[0].Image
	r := regexp.MustCompile(`\:(\S+)`)
	return r.FindStringSubmatch(*repo)[1]
}

func describeLatestTaskDefinition(service string) *ecs.TaskDefinition {
	latestDefintionARN := getLatestDefinitionARN(service)
	latestDefintionARNValue := *latestDefintionARN
	input := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(latestDefintionARNValue),
	}

	result, err := ECSService.DescribeTaskDefinition(input)
	handleECSErr(err)

	return result.TaskDefinition
}

func getLatestDefinitionARN(service string) *string {
	input := &ecs.ListTaskDefinitionsInput{
		FamilyPrefix: aws.String(service),
	}

	result, err := ECSService.ListTaskDefinitions(input)
	handleECSErr(err)

	definitions := result.TaskDefinitionArns
	return definitions[len(definitions)-1]
}

func listECSClusters() []string {
	input := &ecs.ListClustersInput{}

	result, err := ECSService.ListClusters(input)
	handleECSErr(err)

	r := regexp.MustCompile(`([^\/]+)$`)

	clusters := make([]string, 0)
	for _, cluster := range result.ClusterArns {
		clusterValue := *cluster
		clusters = append(clusters, r.FindString(clusterValue))
	}
	return clusters
}

func listECSServices(clusterName string) []string {
	input := &ecs.ListServicesInput{
		Cluster: aws.String(clusterName),
	}

	result, err := ECSService.ListServices(input)
	handleECSErr(err)

	r := regexp.MustCompile(`([^\/]+)$`)

	services := make([]string, 0)
	for _, service := range result.ServiceArns {
		serviceValue := *service
		services = append(services, r.FindString(serviceValue))
	}
	return services
}
