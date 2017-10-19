package ufo

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecr"
	log "github.com/sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"regexp"
	"fmt"
	"runtime"
)

// @todo set up log file
// @todo stuff should accept AWS types instead of strings

type UFOConfig struct {
	Profile *string
	Region *string
}

type UFOState struct {
	Cluster *ecs.Cluster
	Service *ecs.Service
	TaskDefinition *ecs.TaskDefinition
}

type UFO struct {
	State *UFOState
	Session *session.Session
	ECS *ecs.ECS
	ECR *ecr.ECR
}

func Fly(appConfig UFOConfig) *UFO {
	return CreateUFO(appConfig)
}

func CreateUFO(appConfig UFOConfig) *UFO {
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: appConfig.Region},
		Profile: *appConfig.Profile,
	}))

	app := &UFO{
		Session: awsSession,
		ECS: ecs.New(awsSession),
		ECR: ecr.New(awsSession),
		State: &UFOState{},
	}

	log.SetFormatter(&log.JSONFormatter{})

	return app
}

func (a *UFO) UseCluster(c *ecs.Cluster) {
	a.State.Cluster = c
}

func (a *UFO) UseService(s *ecs.Service) {
	a.State.Service = s
}

func (a *UFO) UseTaskDefinition(t *ecs.TaskDefinition) {
	a.State.TaskDefinition = t
}

func (a *UFO) logError(err error) {
	parsed, ok := err.(awserr.Error)

	if ! ok {
		log.Errorf("Unable to parse error: %v.\n", err)
	}

	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	log.WithFields(log.Fields{
		"code": parsed.Code(),
		"error": parsed.Error(),
		"frame": fmt.Sprintf("%s,:%d %s\n", frame.File, frame.Line, frame.Function),
	}).Error("Received an error from AWS.")
}

func (a *UFO) ListECSClusters() ([]string, error) {
	res, err := a.ECS.ListClusters(&ecs.ListClustersInput{})

	if err != nil {
		a.logError(err)

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

func (a *UFO) ListECSServices(c *ecs.Cluster) ([]string, error) {
	res, err := a.ECS.ListServices(&ecs.ListServicesInput{
		Cluster: c.ClusterArn,
	})

	if err != nil {
		a.logError(err)

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

func (a *UFO) DescribeCluster(clusterName string) (*ecs.Cluster, error) {
	res, err := a.ECS.DescribeClusters(&ecs.DescribeClustersInput{
		Clusters: []*string{
			&clusterName,
		},
	})

	if err != nil {
		a.logError(err)

		return nil, err
	}

	return res.Clusters[0], nil
}

func (a *UFO) DescribeService(c *ecs.Cluster, service string) (*ecs.Service, error) {
	res, err := a.ECS.DescribeServices(&ecs.DescribeServicesInput{
		Cluster: c.ClusterArn,
		Services: []*string{
			&service,
		},
	})

	if err != nil {
		a.logError(err)

		return nil, err
	}

	return res.Services[0], nil
}

func (a *UFO) DescribeTaskDefinition(c *ecs.Cluster, s *ecs.Service) (*ecs.TaskDefinition, error) {
	result, err := a.ECS.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{
		TaskDefinition: s.TaskDefinition,
	})

	if err != nil {
		a.logError(err)

		return nil, err
	}

	return result.TaskDefinition, nil
}

func (a *UFO) DescribeTasks(c *ecs.Cluster, tasks []*string) (*ecs.DescribeTasksOutput, error) {
	result, err := a.ECS.DescribeTasks(&ecs.DescribeTasksInput{
		Cluster: c.ClusterName,
		Tasks:   tasks,
	})

	if err != nil {
		a.logError(err)

		return nil, err
	}

	return result, nil
}

func (a *UFO) DescribeImages(t *ecs.TaskDefinition) ([]*ecr.ImageDetail, error) {
	currentImage := t.ContainerDefinitions[0].Image

	// Parse the repo name out of an image tag
	repoName := a.GetRepoNameFromImage(currentImage)

	result, err := a.ECR.DescribeImages(&ecr.DescribeImagesInput{
		RepositoryName: aws.String(repoName),
	})

	if err != nil {
		a.logError(err)

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

// @todo require the abstracted types instead of string names?
func (a *UFO) GetLastDeployedCommit(taskDefinition string) (string, error) {
	result, err := a.ECS.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{
		TaskDefinition: &taskDefinition,
	})

	if err != nil {
		a.logError(err)

		return "", err
	}

	repo := result.TaskDefinition.ContainerDefinitions[0].Image

	r := regexp.MustCompile(`\:(\S+)`)

	// @todo handle cases where regex does not return required matches
	return r.FindStringSubmatch(*repo)[1], nil
}

func (a *UFO) RegisterNewTaskDefinition(c *ecs.Cluster, s *ecs.Service, version string) (*ecs.TaskDefinition, error) {
	// @todo rename funcs to be more descriptive?
	t, err := a.DescribeTaskDefinition(c, s)

	if err != nil {
		a.logError(err)

		return nil, err // @todo simplify return
	}

	result, err := a.ECS.RegisterTaskDefinition(&ecs.RegisterTaskDefinitionInput{
		// Update the task definition to use the new docker image via updateTaskDefinition
		ContainerDefinitions: a.UpdateTaskDefinition(t, version),
		Family:               t.Family,
	})

	if err != nil {
		a.logError(err)

		return nil, err
	}

	return result.TaskDefinition, nil
}

func (a *UFO) ListRunningTasks(c *ecs.Cluster, s *ecs.Service) ([]*string, error) {
	result, err := a.ECS.ListTasks(&ecs.ListTasksInput{
		Cluster:       c.ClusterName,
		ServiceName:   s.ServiceName,
		DesiredStatus: aws.String("RUNNING"),
	})

	if err != nil {
		a.logError(err)

		return nil, err
	}

	return result.TaskArns, nil
}

func (a *UFO) GetRepoNameFromImage(image *string) string {
	// Parse the repo name out of an image tag
	r := regexp.MustCompile(`\/(\S+):`) // /repoName:sha256:
	repoName := r.FindStringSubmatch(*image)[1]

	return repoName
}

func (a *UFO) UpdateTaskDefinition(t *ecs.TaskDefinition, version string) []*ecs.ContainerDefinition {
	r := regexp.MustCompile(`(\S+):`)
	currentImage := *t.ContainerDefinitions[0].Image

	repo := r.FindStringSubmatch(currentImage)[1]
	newImage := fmt.Sprintf("%s:%s", repo, version)

	*t.ContainerDefinitions[0].Image = newImage

	return t.ContainerDefinitions
}

func (a *UFO) UpdateService(c *ecs.Cluster, s *ecs.Service, t *ecs.TaskDefinition) (*ecs.UpdateServiceOutput, error) {
	result, err := a.ECS.UpdateService(&ecs.UpdateServiceInput{
		Cluster:        c.ClusterArn,
		Service:        s.ServiceArn,
		TaskDefinition: t.TaskDefinitionArn,
	})

	if err != nil {
		a.logError(err)

		return nil, err
	}

	return result, nil
}

func (a *UFO) Deploy(c *ecs.Cluster, s *ecs.Service, version string) (*ecs.TaskDefinition, error) {
	t, err := a.RegisterNewTaskDefinition(c, s, version)

	if err != nil {
		a.logError(err)

		return nil, err
	}

	a.UpdateService(c, s, t)

	return t, nil
}
