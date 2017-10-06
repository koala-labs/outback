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

func getServiceTaskDefinition(service *ecs.Service) *string {
	return service.TaskDefinition
}

func registerNewTaskDefinition(cluster string, service string, version string) (string, string) {
	latestDefinitions := describeTaskDefinition(cluster, service)
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

func describeTaskDefinition(cluster string, service string) *ecs.TaskDefinition {
	currentService := describeService(cluster, service)
	latestDefinition := *getServiceTaskDefinition(currentService)
	input := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(latestDefinition),
	}

	result, err := ECSService.DescribeTaskDefinition(input)
	handleECSErr(err)

	return result.TaskDefinition
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
