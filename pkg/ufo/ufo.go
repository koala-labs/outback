package ufo

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	"github.com/pkg/errors"
)

type AwsConfig struct {
	Profile string
	Region  string
}

type UFO struct {
	Config *AwsConfig
	ECS    ecsiface.ECSAPI
	ECR    ecriface.ECRAPI
	CWL    cloudwatchlogsiface.CloudWatchLogsAPI
}

// New creates a UFO session and connects to AWS to create a session
func New(awsConfig *AwsConfig) *UFO {
	var sess *session.Session
	if os.Getenv("AWS_ACCESS_KEY_ID") != "" && os.Getenv("AWS_SECRET_ACCESS_KEY") != "" {
		sess = session.Must(session.NewSessionWithOptions(session.Options{
			Config: aws.Config{Region: aws.String(awsConfig.Region)}}))
	} else {
		sess = session.Must(session.NewSessionWithOptions(session.Options{
			Config:  aws.Config{Region: aws.String(awsConfig.Region)},
			Profile: awsConfig.Profile,
		}))
	}

	app := &UFO{
		Config: awsConfig,
		ECS:    ecs.New(sess),
		ECR:    ecr.New(sess),
		CWL:    cloudwatchlogs.New(sess),
	}

	return app
}

// Clusters returns all ECS clusters
func (u *UFO) Clusters() ([]string, error) {
	res, err := u.ECS.ListClusters(&ecs.ListClustersInput{})

	if err != nil {
		return nil, errors.Wrap(err, errFailedToListClusters)
	}

	r := regexp.MustCompile(`([^\/]+)$`)
	clusters := make([]string, len(res.ClusterArns))

	// Amazon return ARNs which we then keep just the cluster name from
	for i, cluster := range res.ClusterArns {
		clusters[i] = r.FindString(*cluster)
	}

	return clusters, nil
}

// Services returns all services in a cluster
func (u *UFO) Services(c *ecs.Cluster) ([]string, error) {
	res, err := u.ECS.ListServices(&ecs.ListServicesInput{
		Cluster: c.ClusterArn,
	})

	if err != nil {
		return nil, errors.Wrap(err, errFailedToListServices)
	}

	r := regexp.MustCompile(`([^\/]+)$`)
	services := make([]string, len(res.ServiceArns))

	for i, service := range res.ServiceArns {
		services[i] = r.FindString(*service)
	}

	return services, nil
}

// RunningTasks gets all running tasks in a cluster and service
func (u *UFO) RunningTasks(c *ecs.Cluster, s *ecs.Service) ([]*string, error) {
	result, err := u.ECS.ListTasks(&ecs.ListTasksInput{
		Cluster:       c.ClusterName,
		ServiceName:   s.ServiceName,
		DesiredStatus: aws.String("RUNNING"),
	})

	if err != nil {
		return nil, errors.Wrap(err, errFailedToListRunningTasks)
	}

	return result.TaskArns, nil
}

// GetCluster returns a clusters detail
func (u *UFO) GetCluster(name string) (*ecs.Cluster, error) {
	res, err := u.ECS.DescribeClusters(&ecs.DescribeClustersInput{
		Clusters: []*string{
			&name,
		},
	})

	if err != nil {
		return nil, errors.Wrap(err, errCouldNotRetrieveCluster)
	}

	if len(res.Clusters) < 1 {
		return nil, errors.Wrap(err, errClusterNotFound)
	}

	return res.Clusters[0], nil
}

// GetService returns service details within a cluster by service name or ARN
func (u *UFO) GetService(c *ecs.Cluster, service string) (*ecs.Service, error) {
	res, err := u.ECS.DescribeServices(&ecs.DescribeServicesInput{
		Cluster: c.ClusterArn,
		Services: []*string{
			&service,
		},
	})

	if err != nil {
		return nil, errors.Wrap(err, errCouldNotRetrieveService)
	}

	if len(res.Services) < 1 {
		return nil, errors.Wrap(err, errServiceNotFound)
	}

	return res.Services[0], nil
}

// GetTaskDefinition returns details of a task definition in
// a cluster and service by service's current task definition
func (u *UFO) GetTaskDefinition(c *ecs.Cluster, s *ecs.Service) (*ecs.TaskDefinition, error) {
	result, err := u.ECS.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{
		TaskDefinition: s.TaskDefinition,
	})

	if err != nil {
		return nil, errors.Wrap(err, errCouldNotRetrieveTaskDefinition)
	}

	return result.TaskDefinition, nil
}

// GetTasks gets all tasks in a cluster
func (u *UFO) GetTasks(c *ecs.Cluster, tasks []*string) ([]*ecs.Task, error) {
	result, err := u.ECS.DescribeTasks(&ecs.DescribeTasksInput{
		Cluster: c.ClusterName,
		Tasks:   tasks,
	})

	if err != nil {
		return nil, errors.Wrap(err, errCouldNotRetrieveTasks)
	}

	return result.Tasks, nil
}

// GetImages gets images for a task definition
// @todo how does this handle multiple images in a task def?
func (u *UFO) GetImages(t *ecs.TaskDefinition) ([]*ecr.ImageDetail, error) {
	currentImage := t.ContainerDefinitions[0].Image

	// Parse the repo name out of an image tag
	repoName := u.GetRepoFromImage(currentImage)

	result, err := u.ECR.DescribeImages(&ecr.DescribeImagesInput{
		RepositoryName: aws.String(repoName),
	})

	if err != nil {
		return nil, errors.Wrap(err, errCouldNotRetrieveImages)
	}

	images := make([]*ecr.ImageDetail, 0)

	for _, image := range result.ImageDetails {
		if image.ImageTags != nil {
			images = append(images, image)
		}
	}

	return images, nil
}

// GetLastDeployedCommit finds the most recent committed image for a taskDefinition
func (u *UFO) GetLastDeployedCommit(taskDefinition string) (string, error) {
	result, err := u.ECS.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{
		TaskDefinition: &taskDefinition,
	})

	if err != nil {
		return "", errors.Wrap(err, errCouldNotRetrieveTaskDefinition)
	}

	if len(result.TaskDefinition.ContainerDefinitions) < 1 {
		return "", errors.Wrap(err, errInvalidTaskDefinition)
	}

	repo := result.TaskDefinition.ContainerDefinitions[0].Image

	r := regexp.MustCompile(`\:(\S+)`)

	return r.FindStringSubmatch(*repo)[1], nil
}

// RegisterTaskDefinitionWithImage creates a new task definition with the provided tag
// This copies an existing task definition and only changes the tag used for the image
func (u *UFO) RegisterTaskDefinitionWithImage(c *ecs.Cluster, s *ecs.Service, tag string) (*ecs.TaskDefinition, error) {
	t, err := u.GetTaskDefinition(c, s)

	if err != nil {
		return nil, err
	}

	newTaskDef := u.UpdateTaskDefinitionImage(*t, tag)

	result, err := u.ECS.RegisterTaskDefinition(&ecs.RegisterTaskDefinitionInput{
		// Update the task definition to use the new docker image via UpdateTaskDefinitionImage
		Cpu:                     newTaskDef.Cpu,
		Family:                  newTaskDef.Family,
		Memory:                  newTaskDef.Memory,
		Volumes:                 newTaskDef.Volumes,
		NetworkMode:             newTaskDef.NetworkMode,
		ExecutionRoleArn:        newTaskDef.ExecutionRoleArn,
		TaskRoleArn:             newTaskDef.TaskRoleArn,
		ContainerDefinitions:    newTaskDef.ContainerDefinitions,
		RequiresCompatibilities: newTaskDef.RequiresCompatibilities,
	})

	if err != nil {
		return nil, errors.Wrap(err, errCouldNotRegisterTaskDefinition)
	}

	return result.TaskDefinition, nil
}

// RegisterTaskDefinitionWithEnvVars takes a task definition as an argument and updates its
// ContainerDefinitions field which contains environment variables
func (u *UFO) RegisterTaskDefinitionWithEnvVars(t *ecs.TaskDefinition) (*ecs.TaskDefinition, error) {
	result, err := u.ECS.RegisterTaskDefinition(&ecs.RegisterTaskDefinitionInput{
		Cpu:                     t.Cpu,
		Family:                  t.Family,
		Memory:                  t.Memory,
		Volumes:                 t.Volumes,
		NetworkMode:             t.NetworkMode,
		ExecutionRoleArn:        t.ExecutionRoleArn,
		TaskRoleArn:             t.TaskRoleArn,
		ContainerDefinitions:    t.ContainerDefinitions,
		RequiresCompatibilities: t.RequiresCompatibilities,
	})

	if err != nil {
		return nil, errors.Wrap(err, errCouldNotRegisterTaskDefinition)
	}

	return result.TaskDefinition, nil
}

// RollbackTaskDefinition updates the task definition to the desired revision number
func (u *UFO) RollbackTaskDefinition(c *ecs.Cluster, s *ecs.Service, t *ecs.TaskDefinition, n int) (string, error) {

	var taskFamilyRevision string

	r := regexp.MustCompile(`([^\/]+)$`)
	x := regexp.MustCompile(`([^\/]+)`)

	currentTaskDefinitionArn := *t.TaskDefinitionArn
	currentTaskDefinitionFamilyRevision := r.FindString(currentTaskDefinitionArn)
	currentTaskDefinitionArnName := x.FindString(currentTaskDefinitionArn)
	split := strings.Split(currentTaskDefinitionFamilyRevision, ":")
	taskFamily, taskRevision := split[0], split[1]

	if n != 0 {
		taskFamilyRevision = strings.Join([]string{taskFamily, ":", strconv.Itoa(n)}, "")
	} else {
		i, _ := strconv.Atoi(taskRevision)
		i--
		taskFamilyRevision = strings.Join([]string{taskFamily, ":", strconv.Itoa(i)}, "")
	}

	taskFamilyRevisionArn := currentTaskDefinitionArnName + "/" + taskFamilyRevision
	*t.TaskDefinitionArn = taskFamilyRevisionArn
	_, err := u.RollbackService(c, s, taskFamilyRevision)

	return taskFamilyRevision, err
}

// UpdateTaskDefinitionImage copies a task definition and update its image tag
func (u *UFO) UpdateTaskDefinitionImage(t ecs.TaskDefinition, tag string) ecs.TaskDefinition {
	r := regexp.MustCompile(`(\S+):`)
	currentImage := *t.ContainerDefinitions[0].Image

	repo := r.FindStringSubmatch(currentImage)[1]
	newImage := fmt.Sprintf("%s:%s", repo, tag)

	*t.ContainerDefinitions[0].Image = newImage

	return t
}

// GetRepoFromImage parses an image URL:tag and reads its repo
func (u *UFO) GetRepoFromImage(image *string) string {
	// Parse the repo  out of an image tag
	r := regexp.MustCompile(`\/(\S+):`) // /repo:sha256:
	repo := r.FindStringSubmatch(*image)[1]

	return repo
}

// RollbackService updates the ECS service with the desired rollback revision
func (u *UFO) RollbackService(c *ecs.Cluster, s *ecs.Service, t string) (*ecs.UpdateServiceOutput, error) {
	result, err := u.ECS.UpdateService(&ecs.UpdateServiceInput{
		Cluster:        c.ClusterArn,
		Service:        s.ServiceArn,
		TaskDefinition: aws.String(t),
	})

	if err != nil {
		return nil, errors.Wrap(err, errCouldNotUpdateService)
	}

	return result, nil
}

// UpdateService updates a service in a cluster with a new task definition
func (u *UFO) UpdateService(c *ecs.Cluster, s *ecs.Service, t *ecs.TaskDefinition) (*ecs.UpdateServiceOutput, error) {
	result, err := u.ECS.UpdateService(&ecs.UpdateServiceInput{
		Cluster:        c.ClusterArn,
		Service:        s.ServiceArn,
		TaskDefinition: t.TaskDefinitionArn,
	})

	if err != nil {
		return nil, errors.Wrap(err, errCouldNotUpdateService)
	}

	return result, nil
}

// UpdateServiceWithNewTaskDefinition registers a task definition with a tag and updates a service
// with the newly registered task definition
func (u *UFO) UpdateServiceWithNewTaskDefinition(c *ecs.Cluster, s *ecs.Service, tag string) (*ecs.TaskDefinition, error) {
	t, err := u.RegisterTaskDefinitionWithImage(c, s, tag)

	if err != nil {
		return nil, err
	}

	_, err = u.UpdateService(c, s, t)

	if err != nil {
		return nil, err
	}

	return t, err
}

// RunTask runs a specified task in a cluster
func (u *UFO) RunTask(c *ecs.Cluster, t *ecs.TaskDefinition, cmd string) (*ecs.RunTaskOutput, error) {
	splitString := strings.Split(cmd, " ")

	result, err := u.ECS.RunTask(&ecs.RunTaskInput{
		Cluster:        c.ClusterName,
		TaskDefinition: t.TaskDefinitionArn,
		Overrides: &ecs.TaskOverride{
			ContainerOverrides: []*ecs.ContainerOverride{&ecs.ContainerOverride{
				Command: aws.StringSlice(splitString),
				Name:    t.ContainerDefinitions[0].Name,
			}},
		},
	})

	if err != nil {
		return nil, errors.Wrap(err, errCouldNotRunTask)
	}

	return result, nil
}

// IsServiceRunning is meant to be called after a service update. This function checks if the newly
// started task has the status "RUNNING"
func (u *UFO) IsServiceRunning(detail *DeployDetail) bool {
	if *detail.Service.DesiredCount <= 0 {
		return false
	}

	runningTasks, err := u.RunningTasks(detail.Cluster, detail.Service)

	if err != nil {
		return false
	}

	if len(runningTasks) <= 0 {
		return false
	}

	tasks, err := u.GetTasks(detail.Cluster, runningTasks)

	if err != nil {
		return false
	}

	for _, task := range tasks {
		if *task.TaskDefinitionArn == *detail.TaskDefinition.TaskDefinitionArn && *task.LastStatus == "RUNNING" {
			return true
		}
	}

	return false
}

func (u *UFO) IsTaskRunning(cluster *string, task *string) error {
	err := u.ECS.WaitUntilTasksStoppedWithContext(aws.BackgroundContext(), &ecs.DescribeTasksInput{
		Cluster: cluster,
		Tasks:   []*string{task},
	}, func(w *request.Waiter) {
		w.Delay = request.ConstantWaiterDelay(time.Second * 2)
	})

	return err
}

// ECRLogin uses an AWS region & profile to login to ECR
func (u *UFO) ECRLogin() error {
	input := &ecr.GetAuthorizationTokenInput{}

	resp, err := u.ECR.GetAuthorizationToken(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecr.ErrCodeServerException:
				fmt.Println(ecr.ErrCodeServerException, aerr.Error())
			case ecr.ErrCodeInvalidParameterException:
				fmt.Println(ecr.ErrCodeInvalidParameterException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(aerr.Error())
		}

	}

	auth := resp.AuthorizationData
	decode, err := base64.StdEncoding.DecodeString(*auth[0].AuthorizationToken)
	if err != nil {
		return err
	}

	token := strings.SplitN(string(decode), ":", 2)
	user := token[0]
	password := token[1]
	endpoint := *auth[0].ProxyEndpoint

	cmd := fmt.Sprintf("docker login -u %s -p %s %s", user, password, endpoint)
	login := exec.Command("bash", "-c", cmd)
	loginErr := login.Run()
	if loginErr != nil {
		return errors.Wrap(err, errECRLogin)
	}

	return nil
}

type GetLogsInput struct {
	Filter         string
	LogGroupName   string
	LogStreamNames []string
	EndTime        time.Time
	StartTime      time.Time
}

type LogLine struct {
	EventID       string
	LogStreamName string
	Message       string
	Timestamp     time.Time
}

func (u *UFO) GetLogs(i *GetLogsInput) ([]LogLine, error) {
	var logLines []LogLine

	input := &cloudwatchlogs.FilterLogEventsInput{
		LogGroupName: aws.String(i.LogGroupName),
		Interleaved:  aws.Bool(true),
	}

	if !i.StartTime.IsZero() {
		input.SetStartTime(i.StartTime.UTC().UnixNano() / int64(time.Millisecond))
	}

	if !i.EndTime.IsZero() {
		input.SetEndTime(i.EndTime.UTC().UnixNano() / int64(time.Millisecond))
	}

	if i.Filter != "" {
		input.SetFilterPattern(i.Filter)
	}

	if len(i.LogStreamNames) > 0 {
		input.SetLogStreamNames(aws.StringSlice(i.LogStreamNames))
	}

	err := u.CWL.FilterLogEventsPages(
		input,
		func(resp *cloudwatchlogs.FilterLogEventsOutput, lastPage bool) bool {
			for _, event := range resp.Events {
				logLines = append(logLines,
					LogLine{
						EventID:       aws.StringValue(event.EventId),
						Message:       aws.StringValue(event.Message),
						LogStreamName: aws.StringValue(event.LogStreamName),
						Timestamp:     time.Unix(0, aws.Int64Value(event.Timestamp)*int64(time.Millisecond)),
					},
				)
			}

			return true
		},
	)

	if err != nil {
		return nil, errors.Wrap(err, errCouldNotGetLogs)
	}

	return logLines, nil
}
