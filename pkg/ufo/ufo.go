package ufo

import (
	"fmt"
	"regexp"
	"runtime"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecs"
	log "github.com/sirupsen/logrus"
)

// @todo set up log file

type UFOConfig struct {
	Profile *string
	Region  *string
}

type UFOState struct {
	Cluster        *ecs.Cluster
	Service        *ecs.Service
	TaskDefinition *ecs.TaskDefinition
}

type UFO struct {
	State   *UFOState
	Session *session.Session
	ECS     *ecs.ECS
	ECR     *ecr.ECR
}

// Alias for CreateUFO
func Fly(appConfig UFOConfig) *UFO {
	return CreateUFO(appConfig)
}

// Create a UFO session and connect to AWS to create a session
func CreateUFO(appConfig UFOConfig) *UFO {
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: appConfig.Region},
		Profile: *appConfig.Profile,
	}))

	app := &UFO{
		Session: awsSession,
		ECS:     ecs.New(awsSession),
		ECR:     ecr.New(awsSession),
		State:   &UFOState{},
	}

	log.SetFormatter(&log.JSONFormatter{})

	return app
}

// Set a cluster choice in UFO state
// @todo this may be extraneous but if we decide to leave it in, we should have funcs optionally require
//		the cluster/service/taskDef and if not passed, can use the ones stored in state.
func (u *UFO) UseCluster(c *ecs.Cluster) {
	u.State.Cluster = c
}

// Set a service choice in UFO state
func (u *UFO) UseService(s *ecs.Service) {
	u.State.Service = s
}

// Set a task definition choice in UFO state
func (u *UFO) UseTaskDefinition(t *ecs.TaskDefinition) {
	u.State.TaskDefinition = t
}

// Handle errors, nil or otherwise
// This func is intended to always be called after an error is returned from an AWS method call in UFO
func (u *UFO) logError(err error) {
	parsed, ok := err.(awserr.Error)

	if !ok {
		log.Errorf("Unable to parse error: %v.\n", err)
	}

	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	log.WithFields(log.Fields{
		"code":  parsed.Code(),
		"error": parsed.Error(),
		"frame": fmt.Sprintf("%s,:%d %s\n", frame.File, frame.Line, frame.Function),
	}).Error("Received an error from AWS.")
}

// View all ECS clusters in the account
func (u *UFO) Clusters() ([]string, error) {
	res, err := u.ECS.ListClusters(&ecs.ListClustersInput{})

	if err != nil {
		u.logError(err)

		return nil, err
	}

	r := regexp.MustCompile(`([^\/]+)$`)
	clusters := make([]string, 0)

	for _, cluster := range res.ClusterArns {
		clusterValue := *cluster
		clusters = append(clusters, r.FindString(clusterValue))
	}

	return clusters, nil
}

// View all services in a cluster
func (u *UFO) Services(c *ecs.Cluster) ([]string, error) {
	res, err := u.ECS.ListServices(&ecs.ListServicesInput{
		Cluster: c.ClusterArn,
	})

	if err != nil {
		u.logError(err)

		return nil, err
	}

	r := regexp.MustCompile(`([^\/]+)$`)
	services := make([]string, 0)

	for _, service := range res.ServiceArns {
		serviceValue := *service
		services = append(services, r.FindString(serviceValue))
	}

	return services, nil
}

// View all running tasks in a cluster and service
func (u *UFO) RunningTasks(c *ecs.Cluster, s *ecs.Service) ([]*string, error) {
	result, err := u.ECS.ListTasks(&ecs.ListTasksInput{
		Cluster:       c.ClusterName,
		ServiceName:   s.ServiceName,
		DesiredStatus: aws.String("RUNNING"),
	})

	if err != nil {
		u.logError(err)

		return nil, err
	}

	return result.TaskArns, nil
}

// Get cluster detail with cluster name or ARN
func (u *UFO) GetCluster(clusterName string) (*ecs.Cluster, error) {
	res, err := u.ECS.DescribeClusters(&ecs.DescribeClustersInput{
		Clusters: []*string{
			&clusterName,
		},
	})

	if err != nil {
		u.logError(err)

		return nil, err
	}

	return res.Clusters[0], nil
}

// Get service details within a cluster by service name or ARN
func (u *UFO) GetService(c *ecs.Cluster, service string) (*ecs.Service, error) {
	res, err := u.ECS.DescribeServices(&ecs.DescribeServicesInput{
		Cluster: c.ClusterArn,
		Services: []*string{
			&service,
		},
	})

	if err != nil {
		u.logError(err)

		return nil, err
	}

	return res.Services[0], nil
}

// Get task definition details in a cluster and service by service's current task definition
func (u *UFO) GetTaskDefinition(c *ecs.Cluster, s *ecs.Service) (*ecs.TaskDefinition, error) {
	result, err := u.ECS.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{
		TaskDefinition: s.TaskDefinition,
	})

	if err != nil {
		u.logError(err)

		return nil, err
	}

	return result.TaskDefinition, nil
}

// Get all tasks in a cluster
// @todo return []*Task
func (u *UFO) GetTasks(c *ecs.Cluster, tasks []*string) (*ecs.DescribeTasksOutput, error) {
	result, err := u.ECS.DescribeTasks(&ecs.DescribeTasksInput{
		Cluster: c.ClusterName,
		Tasks:   tasks,
	})

	if err != nil {
		u.logError(err)

		return nil, err
	}

	return result, nil
}

// Get images for a task definition
// @todo how does this handle multiple images in a task def?
func (u *UFO) GetImages(t *ecs.TaskDefinition) ([]*ecr.ImageDetail, error) {
	currentImage := t.ContainerDefinitions[0].Image

	// Parse the repo name out of an image tag
	repoName := u.GetRepoNameFromImage(currentImage)

	result, err := u.ECR.DescribeImages(&ecr.DescribeImagesInput{
		RepositoryName: aws.String(repoName),
	})

	if err != nil {
		u.logError(err)

		return nil, err
	}

	images := make([]*ecr.ImageDetail, 0)

	for _, image := range result.ImageDetails {
		if image.ImageTags != nil {
			images = append(images, image)
		}
	}

	return images, nil
}

// Find the most recent committed image for a taskDefinition
func (u *UFO) GetLastDeployedCommit(taskDefinition string) (string, error) {
	result, err := u.ECS.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{
		TaskDefinition: &taskDefinition,
	})

	if err != nil {
		u.logError(err)

		return "", err
	}

	repo := result.TaskDefinition.ContainerDefinitions[0].Image

	r := regexp.MustCompile(`\:(\S+)`)

	return r.FindStringSubmatch(*repo)[1], nil
}

// Create a new task definition with an image at a specific version
// This copies an existing task definition and only updates the tag used for the image
func (u *UFO) RegisterNewTaskDefinition(c *ecs.Cluster, s *ecs.Service, version string) (*ecs.TaskDefinition, error) {
	t, err := u.GetTaskDefinition(c, s)

	if err != nil {
		u.logError(err)

		return nil, err // @todo simplify return
	}

	result, err := u.ECS.RegisterTaskDefinition(&ecs.RegisterTaskDefinitionInput{
		// Update the task definition to use the new docker image via updateTaskDefinition
		ContainerDefinitions: u.UpdateTaskDefinitionImage(*t, version),
		Family:               t.Family,
	})

	if err != nil {
		u.logError(err)

		return nil, err
	}

	return result.TaskDefinition, nil
}

// Copy a task definition and update its image tag
func (u *UFO) UpdateTaskDefinitionImage(t ecs.TaskDefinition, version string) []*ecs.ContainerDefinition {
	r := regexp.MustCompile(`(\S+):`)
	currentImage := *t.ContainerDefinitions[0].Image

	repo := r.FindStringSubmatch(currentImage)[1]
	newImage := fmt.Sprintf("%s:%s", repo, version)

	*t.ContainerDefinitions[0].Image = newImage

	return t.ContainerDefinitions
}

// Parse an image URL:tag and read its repo name
func (u *UFO) GetRepoNameFromImage(image *string) string {
	// Parse the repo name out of an image tag
	r := regexp.MustCompile(`\/(\S+):`) // /repoName:sha256:
	repoName := r.FindStringSubmatch(*image)[1]

	return repoName
}

// Update a service in a cluster with a new task definition
func (u *UFO) UpdateService(c *ecs.Cluster, s *ecs.Service, t *ecs.TaskDefinition) (*ecs.UpdateServiceOutput, error) {
	result, err := u.ECS.UpdateService(&ecs.UpdateServiceInput{
		Cluster:        c.ClusterArn,
		Service:        s.ServiceArn,
		TaskDefinition: t.TaskDefinitionArn,
	})

	if err != nil {
		u.logError(err)

		return nil, err
	}

	return result, nil
}

// Deploy a version to a service in a cluster
func (u *UFO) Deploy(c *ecs.Cluster, s *ecs.Service, version string) (*ecs.TaskDefinition, error) {
	t, err := u.RegisterNewTaskDefinition(c, s, version)

	if err != nil {
		u.logError(err)

		return nil, err
	}

	u.UpdateService(c, s, t)

	return t, nil
}
