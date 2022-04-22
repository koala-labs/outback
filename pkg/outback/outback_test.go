package outback

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"

	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
)

type mockedECRClient struct {
	ecriface.ECRAPI
}

type mockedECSClient struct {
	ecsiface.ECSAPI
}

type mockedCWLClient struct {
	cloudwatchlogsiface.CloudWatchLogsAPI
}

type mockedDescribeImages struct {
	ecriface.ECRAPI
	Resp  *ecr.DescribeImagesOutput
	Error error
}

type mockedDescribeClusters struct {
	ecsiface.ECSAPI
	Resp  *ecs.DescribeClustersOutput
	Error error
}

type mockedDescribeServices struct {
	ecsiface.ECSAPI
	Resp  *ecs.DescribeServicesOutput
	Error error
}

type mockedDescribeTaskDefinition struct {
	ecsiface.ECSAPI
	Resp  *ecs.DescribeTaskDefinitionOutput
	Error error
}

type mockedDescribeTasks struct {
	ecsiface.ECSAPI
	Resp  *ecs.DescribeTasksOutput
	Error error
}

type mockedListClusters struct {
	ecsiface.ECSAPI
	Resp  *ecs.ListClustersOutput
	Error error
}

type mockedListServices struct {
	ecsiface.ECSAPI
	Resp  *ecs.ListServicesOutput
	Error error
}

type mockedListTasks struct {
	ecsiface.ECSAPI
	Resp  *ecs.ListTasksOutput
	Error error
}

type mockedRunTask struct {
	ecsiface.ECSAPI
	Resp  *ecs.RunTaskOutput
	Error error
}

type mockedRegisterTaskDefinition struct {
	ecsiface.ECSAPI
	Resp  *ecs.RegisterTaskDefinitionOutput
	Error error
}

type mockedDeploy struct {
	ecsiface.ECSAPI
	DescribeTaskDefResp  *ecs.DescribeTaskDefinitionOutput
	RegisterTaskDefResp  *ecs.RegisterTaskDefinitionOutput
	UpdateServiceResp    *ecs.UpdateServiceOutput
	DescribeTaskDefError error
	RegisterTaskDefError error
	UpdateServiceError   error
}

type mockedIsServiceRunning struct {
	ecsiface.ECSAPI
	ListTasksResp      *ecs.ListTasksOutput
	DescribeTasksResp  *ecs.DescribeTasksOutput
	ListTasksError     error
	DescribeTasksError error
}

func (m mockedDescribeClusters) DescribeClusters(in *ecs.DescribeClustersInput) (*ecs.DescribeClustersOutput, error) {
	return m.Resp, m.Error
}

func (m mockedDescribeImages) DescribeImages(in *ecr.DescribeImagesInput) (*ecr.DescribeImagesOutput, error) {
	return m.Resp, m.Error
}

func (m mockedDescribeServices) DescribeServices(in *ecs.DescribeServicesInput) (*ecs.DescribeServicesOutput, error) {
	return m.Resp, m.Error
}

func (m mockedDescribeTaskDefinition) DescribeTaskDefinition(in *ecs.DescribeTaskDefinitionInput) (*ecs.DescribeTaskDefinitionOutput, error) {
	return m.Resp, m.Error
}

func (m mockedDescribeTasks) DescribeTasks(in *ecs.DescribeTasksInput) (*ecs.DescribeTasksOutput, error) {
	return m.Resp, m.Error
}

func (m mockedListClusters) ListClusters(in *ecs.ListClustersInput) (*ecs.ListClustersOutput, error) {
	return m.Resp, m.Error
}

func (m mockedListServices) ListServices(in *ecs.ListServicesInput) (*ecs.ListServicesOutput, error) {
	return m.Resp, m.Error
}

func (m mockedListTasks) ListTasks(in *ecs.ListTasksInput) (*ecs.ListTasksOutput, error) {
	return m.Resp, m.Error
}

func (m mockedRunTask) RunTask(in *ecs.RunTaskInput) (*ecs.RunTaskOutput, error) {
	return m.Resp, m.Error
}

func (m mockedRegisterTaskDefinition) RegisterTaskDefinition(in *ecs.RegisterTaskDefinitionInput) (*ecs.RegisterTaskDefinitionOutput, error) {
	return m.Resp, m.Error
}

func (m mockedDeploy) DescribeTaskDefinition(in *ecs.DescribeTaskDefinitionInput) (*ecs.DescribeTaskDefinitionOutput, error) {
	return m.DescribeTaskDefResp, m.DescribeTaskDefError
}

func (m mockedDeploy) RegisterTaskDefinition(in *ecs.RegisterTaskDefinitionInput) (*ecs.RegisterTaskDefinitionOutput, error) {
	return m.RegisterTaskDefResp, m.RegisterTaskDefError
}

func (m mockedDeploy) UpdateService(in *ecs.UpdateServiceInput) (*ecs.UpdateServiceOutput, error) {
	return m.UpdateServiceResp, m.UpdateServiceError
}

func (m mockedIsServiceRunning) ListTasks(in *ecs.ListTasksInput) (*ecs.ListTasksOutput, error) {
	return m.ListTasksResp, m.ListTasksError
}

func (m mockedIsServiceRunning) DescribeTasks(in *ecs.DescribeTasksInput) (*ecs.DescribeTasksOutput, error) {
	return m.DescribeTasksResp, m.DescribeTasksError
}

func TestOutbackNew(t *testing.T) {
	cases := []struct {
		Expected *Outback
	}{
		{
			Expected: &Outback{},
		},
	}

	outback := New(&AwsConfig{
		Profile: "profile",
		Region:  "region",
	})

	for i, c := range cases {
		if a, e := reflect.TypeOf(outback), reflect.TypeOf(c.Expected); a != e {
			t.Errorf("%d, expected %v state, got %v", i, e, a)
		}
	}
}

func TestOutbackClusters(t *testing.T) {
	var cluster1, cluster2, nextToken string = "cluster1", "cluster2", ""
	cases := []struct {
		Resp     *ecs.ListClustersOutput
		Expected []string
	}{
		{
			Resp: &ecs.ListClustersOutput{
				ClusterArns: []*string{&cluster1, &cluster2},
				NextToken:   &nextToken,
			},
			Expected: []string{
				"cluster1",
				"cluster2",
			},
		},
	}

	for i, c := range cases {
		outback := Outback{
			ECS: &mockedListClusters{Resp: c.Resp},
			ECR: mockedECRClient{},
		}

		clusters, err := outback.Clusters()

		if err != nil {
			t.Fatalf("%d, unexpected error", err)
		}

		if a, e := len(clusters), len(c.Expected); a != e {
			t.Fatalf("%d, expected %d clusters, got %d", i, e, a)
		}

		for j, cluster := range clusters {
			if a, e := cluster, c.Expected[j]; a != e {
				t.Errorf("%d, expected %v cluster, got %v", i, e, a)
			}
		}
	}
}

func TestOutbackClustersError(t *testing.T) {
	cases := []struct {
		Error    error
		Expected error
	}{
		{
			Error:    errors.New("test-error"),
			Expected: errors.Wrap(errors.New("test-error"), errFailedToListClusters),
		},
	}

	for i, c := range cases {
		outback := Outback{
			ECS: &mockedListClusters{Resp: nil, Error: c.Error},
			ECR: mockedECRClient{},
		}

		_, err := outback.Clusters()

		if a, e := err, c.Expected; a.Error() != e.Error() {
			t.Errorf("%d, expected %v, got %v", i, e, a)
		}
	}
}

func TestOutbackServices(t *testing.T) {
	var service1, service2, nextToken string = "service1", "service2", ""
	cases := []struct {
		Resp     *ecs.ListServicesOutput
		Expected []string
	}{
		{
			Resp: &ecs.ListServicesOutput{
				ServiceArns: []*string{&service1, &service2},
				NextToken:   &nextToken,
			},
			Expected: []string{
				"service1",
				"service2",
			},
		},
	}

	for i, c := range cases {
		outback := Outback{
			ECS: &mockedListServices{Resp: c.Resp},
			ECR: mockedECRClient{},
		}

		services, err := outback.Services(&ecs.Cluster{})

		if err != nil {
			t.Fatalf("%d, unexpected error", err)
		}

		if a, e := len(services), len(c.Expected); a != e {
			t.Fatalf("%d, expected %d services, got %d", i, e, a)
		}

		for j, service := range services {
			if a, e := service, c.Expected[j]; a != e {
				t.Errorf("%d, expected %v service, got %v", i, e, a)
			}
		}
	}
}

func TestOutbackServicesError(t *testing.T) {
	cases := []struct {
		Error    error
		Expected error
	}{
		{
			Error:    errors.New("test-error"),
			Expected: errors.Wrap(errors.New("test-error"), errFailedToListServices),
		},
	}

	for i, c := range cases {
		outback := Outback{
			ECS: &mockedListServices{Resp: nil, Error: c.Error},
			ECR: mockedECRClient{},
		}

		_, err := outback.Services(&ecs.Cluster{})

		if a, e := err, c.Expected; a.Error() != e.Error() {
			t.Errorf("%d, expected %v error, got %v", i, e, a)
		}
	}
}

func TestOutbackRunningTasks(t *testing.T) {
	var task1, task2, nextToken string = "task1", "task2", ""
	cases := []struct {
		Resp     *ecs.ListTasksOutput
		Expected []string
	}{
		{
			Resp: &ecs.ListTasksOutput{
				TaskArns:  []*string{&task1, &task2},
				NextToken: &nextToken,
			},
			Expected: []string{
				"task1",
				"task2",
			},
		},
	}

	for i, c := range cases {
		outback := Outback{
			ECS: &mockedListTasks{Resp: c.Resp},
			ECR: mockedECRClient{},
		}

		tasks, err := outback.RunningTasks(&ecs.Cluster{}, &ecs.Service{})

		if err != nil {
			t.Fatalf("%d, unexpected error", err)
		}

		if a, e := len(tasks), len(c.Expected); a != e {
			t.Fatalf("%d, expected %d tasks, got %d", i, e, a)
		}

		for j, task := range tasks {
			if a, e := *task, c.Expected[j]; a != e {
				t.Errorf("%d, expected %v task, got %v", i, e, a)
			}
		}
	}
}

func TestOutbackRunningTasksError(t *testing.T) {
	cases := []struct {
		Error    error
		Expected error
	}{
		{
			Error:    errors.New("test-error"),
			Expected: errors.Wrap(errors.New("test-error"), errFailedToListRunningTasks),
		},
	}

	for i, c := range cases {
		outback := Outback{
			ECS: &mockedListTasks{Resp: nil, Error: c.Error},
			ECR: mockedECRClient{},
		}

		_, err := outback.RunningTasks(&ecs.Cluster{}, &ecs.Service{})

		if a, e := err, c.Expected; a.Error() != e.Error() {
			t.Errorf("%d, expected %v, got %v", i, e, a)
		}
	}
}

func TestOutbackGetCluster(t *testing.T) {
	cases := []struct {
		Resp     *ecs.DescribeClustersOutput
		Expected *ecs.Cluster
	}{
		{
			Resp: &ecs.DescribeClustersOutput{
				Clusters: []*ecs.Cluster{{
					ClusterName: aws.String("test-cluster"),
					ClusterArn:  aws.String("test-clusterarn"),
				}},
				Failures: []*ecs.Failure{},
			},
			Expected: &ecs.Cluster{
				ClusterName: aws.String("test-cluster"),
				ClusterArn:  aws.String("test-clusterarn"),
			},
		},
	}

	for i, c := range cases {
		outback := Outback{
			ECS: &mockedDescribeClusters{Resp: c.Resp},
			ECR: mockedECRClient{},
		}

		cluster, err := outback.GetCluster("test-cluster")

		if err != nil {
			t.Fatalf("%d, unexpected error", err)
		}

		if a, e := *cluster.ClusterName, *c.Expected.ClusterName; a != e {
			t.Errorf("%d, expected %v cluster, got %v", i, e, a)
		}
	}
}

func TestOutbackGetClusterError(t *testing.T) {
	cases := []struct {
		Resp     *ecs.DescribeClustersOutput
		Expected error
	}{
		{
			Resp: &ecs.DescribeClustersOutput{
				Clusters: []*ecs.Cluster{},
			},
			Expected: errors.Wrap(nil, errClusterNotFound),
		},
	}

	for i, c := range cases {
		outback := Outback{
			ECS: &mockedDescribeClusters{Resp: c.Resp},
			ECR: mockedECRClient{},
		}

		_, err := outback.GetCluster("test-cluster")

		if a, e := err, c.Expected; a != e {
			t.Errorf("%d, expected %v, got %v", i, e, a)
		}
	}
}

func TestOutbackGetClusterError2(t *testing.T) {
	cases := []struct {
		Error    error
		Expected error
	}{
		{
			Error:    errors.New("test-error"),
			Expected: errors.Wrap(errors.New("test-error"), errCouldNotRetrieveCluster),
		},
	}

	for i, c := range cases {
		outback := Outback{
			ECS: &mockedDescribeClusters{Resp: nil, Error: c.Error},
			ECR: mockedECRClient{},
		}

		_, err := outback.GetCluster("test-cluster")

		if a, e := err, c.Expected; a.Error() != e.Error() {
			t.Errorf("%d, expected %v, got %v", i, e, a)
		}
	}
}

func TestOutbackGetService(t *testing.T) {
	name := "test-service"
	arn := "test-servicearn"
	cases := []struct {
		Resp     *ecs.DescribeServicesOutput
		Expected *ecs.Service
	}{
		{
			Resp: &ecs.DescribeServicesOutput{
				Services: []*ecs.Service{{
					ServiceName: &name,
					ServiceArn:  &arn,
				}},
				Failures: []*ecs.Failure{},
			},
			Expected: &ecs.Service{
				ServiceName: &name,
			},
		},
	}

	for i, c := range cases {
		outback := Outback{
			ECS: &mockedDescribeServices{Resp: c.Resp},
			ECR: mockedECRClient{},
		}

		service, err := outback.GetService(&ecs.Cluster{}, "test-service")

		if err != nil {
			t.Fatalf("%d, unexpected error", err)
		}

		if a, e := *service.ServiceName, *c.Expected.ServiceName; a != e {
			t.Errorf("%d, expected %v service, got %v", i, e, a)
		}
	}
}

func TestOutbackGetServiceError(t *testing.T) {
	cases := []struct {
		Resp     *ecs.DescribeServicesOutput
		Expected error
	}{
		{
			Resp: &ecs.DescribeServicesOutput{
				Services: []*ecs.Service{},
			},
			Expected: fmt.Errorf("'%s' %s", "test-service", errServiceNotFound),
		},
	}

	for i, c := range cases {
		outback := Outback{
			ECS: &mockedDescribeServices{Resp: c.Resp},
			ECR: mockedECRClient{},
		}

		_, err := outback.GetService(&ecs.Cluster{}, "test-service")

		if a, e := err, c.Expected; a.Error() != e.Error() {
			t.Errorf("%d, expected %v, got %v", i, e, a)
		}
	}
}

func TestOutbackGetServiceError2(t *testing.T) {
	cases := []struct {
		Error    error
		Expected error
	}{
		{
			Error:    errors.New("test-error"),
			Expected: errors.Wrap(errors.New("test-error"), errCouldNotRetrieveService),
		},
	}

	for i, c := range cases {
		outback := Outback{
			ECS: &mockedDescribeServices{Resp: nil, Error: c.Error},
			ECR: mockedECRClient{},
		}

		_, err := outback.GetService(&ecs.Cluster{}, "test-cluster")

		if a, e := err, c.Expected; a.Error() != e.Error() {
			t.Errorf("%d, expected %v, got %v", i, e, a)
		}
	}
}

func TestOutbackGetTaskDefinition(t *testing.T) {
	fam := "test-family"
	subcommand := "echo"
	command := []*string{&subcommand, &subcommand}
	arn := "test-taskdefinitionarn"
	image := "111222333444.dkr.ecr.us-west-1.amazonaws.com/image:ea13366"
	cases := []struct {
		Resp     *ecs.DescribeTaskDefinitionOutput
		Expected *ecs.TaskDefinition
	}{
		{
			Resp: &ecs.DescribeTaskDefinitionOutput{
				TaskDefinition: &ecs.TaskDefinition{
					ContainerDefinitions: []*ecs.ContainerDefinition{{
						Command: command,
						Image:   &image,
					}},
					Family:            &fam,
					TaskDefinitionArn: &arn,
				},
			},
			Expected: &ecs.TaskDefinition{
				ContainerDefinitions: []*ecs.ContainerDefinition{{
					Command: command,
					Image:   &image,
				}},
				Family: &fam,
			},
		},
	}

	for i, c := range cases {
		outback := Outback{
			ECS: &mockedDescribeTaskDefinition{Resp: c.Resp},
			ECR: mockedECRClient{},
		}

		taskDef, err := outback.GetTaskDefinition(&ecs.Cluster{}, &ecs.Service{})

		if err != nil {
			t.Fatalf("%d, unexpected error", err)
		}

		if a, e := *taskDef.Family, *c.Expected.Family; a != e {
			t.Errorf("%d, expected %v family, got %v", i, e, a)
		}

		if a, e := *taskDef.ContainerDefinitions[0].Image, *c.Expected.ContainerDefinitions[0].Image; a != e {
			t.Errorf("%d, expected %v image, got %v", i, e, a)
		}

		actualContainerCommand := strings.Join(aws.StringValueSlice(taskDef.ContainerDefinitions[0].Command), " ")
		expectedContainerCommand := strings.Join(aws.StringValueSlice(c.Expected.ContainerDefinitions[0].Command), " ")

		if a, e := actualContainerCommand, expectedContainerCommand; a != e {
			t.Errorf("%d, expected %v command, got %v", i, e, a)
		}
	}
}

func TestOutbackGetTaskDefinitionError(t *testing.T) {
	cases := []struct {
		Error    error
		Expected error
	}{
		{
			Error:    errors.New("test-error"),
			Expected: errors.Wrap(errors.New("test-error"), errCouldNotRetrieveTaskDefinition),
		},
	}

	for i, c := range cases {
		outback := Outback{
			ECS: &mockedDescribeTaskDefinition{Resp: nil, Error: c.Error},
			ECR: mockedECRClient{},
		}

		_, err := outback.GetTaskDefinition(&ecs.Cluster{}, &ecs.Service{})

		if a, e := err, c.Expected; a.Error() != e.Error() {
			t.Errorf("%d, expected %v, got %v", i, e, a)
		}
	}
}

func TestOutbackGetTasks(t *testing.T) {
	var lastStatus, taskDefArn string = "PENDING", "111222333444.dkr.ecr.us-west-1.amazonaws.com/task"
	cases := []struct {
		Resp     *ecs.DescribeTasksOutput
		Expected *ecs.DescribeTasksOutput
	}{
		{
			Resp: &ecs.DescribeTasksOutput{
				Tasks: []*ecs.Task{{
					LastStatus:        &lastStatus,
					TaskDefinitionArn: &taskDefArn,
				}},
			},
			Expected: &ecs.DescribeTasksOutput{
				Tasks: []*ecs.Task{
					{
						LastStatus:        &lastStatus,
						TaskDefinitionArn: &taskDefArn,
					}},
			},
		},
	}

	for i, c := range cases {
		outback := Outback{
			ECS: &mockedDescribeTasks{Resp: c.Resp},
			ECR: mockedECRClient{},
		}

		tasks, err := outback.GetTasks(&ecs.Cluster{}, []*string{})

		if err != nil {
			t.Fatalf("%d, unexpected error", err)
		}

		if a, e := len(tasks), len(c.Expected.Tasks); a != e {
			t.Fatalf("%d, expected %d tasks, got %d", i, e, a)
		}

		for j, task := range tasks {
			if a, e := task, c.Expected.Tasks[j]; *a.LastStatus != *e.LastStatus {
				t.Errorf("%d, expected %v LastStatus, got %v", i, e, a)
			}

			if a, e := task, c.Expected.Tasks[j]; *a.TaskDefinitionArn != *e.TaskDefinitionArn {
				t.Errorf("%d, expected %v TaskDefinitionARN, got %v", i, e, a)
			}
		}
	}
}

func TestOutbackGetTasksError(t *testing.T) {
	cases := []struct {
		Error    error
		Expected error
	}{
		{
			Error:    errors.New("test-error"),
			Expected: errors.Wrap(errors.New("test-error"), errCouldNotRetrieveTasks),
		},
	}

	for i, c := range cases {
		outback := Outback{
			ECS: &mockedDescribeTasks{Resp: nil, Error: c.Error},
			ECR: mockedECRClient{},
		}

		_, err := outback.GetTasks(&ecs.Cluster{}, aws.StringSlice([]string{"task1", "task2"}))

		if a, e := err, c.Expected; a.Error() != e.Error() {
			t.Errorf("%d, expected %v, got %v", i, e, a)
		}
	}
}

func TestOutbackGetImages(t *testing.T) {
	var tag1, tag2 string = "tag1", "tag2"
	cases := []struct {
		Resp     *ecr.DescribeImagesOutput
		Expected []*ecr.ImageDetail
	}{
		{
			Resp: &ecr.DescribeImagesOutput{
				ImageDetails: []*ecr.ImageDetail{
					{
						ImageTags: []*string{&tag1},
					},
					{
						ImageTags: []*string{&tag2},
					}},
			},
			Expected: []*ecr.ImageDetail{
				{
					ImageTags: []*string{&tag1},
				},
				{
					ImageTags: []*string{&tag2},
				}},
		},
		{
			Resp:     &ecr.DescribeImagesOutput{},
			Expected: []*ecr.ImageDetail{},
		},
	}

	for i, c := range cases {
		outback := Outback{
			ECS: mockedECSClient{},
			ECR: &mockedDescribeImages{Resp: c.Resp},
		}

		images, err := outback.GetImages(&ecs.TaskDefinition{
			ContainerDefinitions: []*ecs.ContainerDefinition{{
				Command: aws.StringSlice([]string{"echo", "this"}),
				Image:   aws.String("111222333444.dkr.ecr.us-west-1.amazonaws.com/image:ea13366"),
			}},
			Family:            aws.String("family"),
			TaskDefinitionArn: aws.String("taskdefarn"),
		})

		if err != nil {
			t.Fatalf("%d, unexpected error", err)
		}

		if a, e := len(images), len(c.Expected); a != e {
			t.Fatalf("%d, expected %d images, got %d", i, e, a)
		}

		for j, image := range images {
			actualImageTag := strings.Join(aws.StringValueSlice(image.ImageTags), " ")
			expectedImageTag := strings.Join(aws.StringValueSlice(c.Expected[j].ImageTags), " ")
			if a, e := actualImageTag, expectedImageTag; a != e {
				t.Errorf("%d, expected %v image tag, got %v", i, e, a)
			}
		}
	}
}

func TestOutbackGetImagesError(t *testing.T) {
	cases := []struct {
		Error    error
		Expected error
	}{
		{
			Error:    errors.New("test-error"),
			Expected: errors.Wrap(errors.New("test-error"), errCouldNotRetrieveImages),
		},
	}

	for i, c := range cases {
		outback := Outback{
			ECS: mockedECSClient{},
			ECR: &mockedDescribeImages{Resp: nil, Error: c.Error},
		}

		_, err := outback.GetImages(&ecs.TaskDefinition{
			ContainerDefinitions: []*ecs.ContainerDefinition{{
				Command: aws.StringSlice([]string{"error", "this"}),
				Image:   aws.String("111222333444.dkr.ecr.us-west-1.amazonaws.com/image:ea13366"),
			}},
			Family:            aws.String("family"),
			TaskDefinitionArn: aws.String("taskdefarn"),
		})

		if a, e := err, c.Expected; a.Error() != e.Error() {
			t.Errorf("%d, expected %v, got %v", i, e, a)
		}
	}
}

func TestOutbackGetLastDeployedCommit(t *testing.T) {
	fam := "test-family"
	subcommand := "echo"
	command := []*string{&subcommand, &subcommand}
	arn := "test-taskdefinitionarn"
	image := "111222333444.dkr.ecr.us-west-1.amazonaws.com/image:ea13366"
	cases := []struct {
		Resp     *ecs.DescribeTaskDefinitionOutput
		Expected string
	}{
		{
			Resp: &ecs.DescribeTaskDefinitionOutput{
				TaskDefinition: &ecs.TaskDefinition{
					ContainerDefinitions: []*ecs.ContainerDefinition{{
						Command: command,
						Image:   &image,
					}},
					Family:            &fam,
					TaskDefinitionArn: &arn,
				},
			},
			Expected: "ea13366",
		},
	}

	for i, c := range cases {
		outback := Outback{
			ECS: &mockedDescribeTaskDefinition{Resp: c.Resp},
			ECR: mockedECRClient{},
		}

		commit, err := outback.GetLastDeployedCommit("111222333444.dkr.ecr.us-west-1.amazonaws.com/image:ea13366")

		if err != nil {
			t.Fatalf("%d, unexpected error", err)
		}

		if a, e := commit, c.Expected; a != e {
			t.Errorf("%d, expected %v commit, got %v", i, e, a)
		}
	}
}

func TestOutbackGetLastDeployedCommitError(t *testing.T) {
	cases := []struct {
		Resp     *ecs.DescribeTaskDefinitionOutput
		Expected error
	}{
		{
			Resp: &ecs.DescribeTaskDefinitionOutput{
				TaskDefinition: &ecs.TaskDefinition{
					ContainerDefinitions: []*ecs.ContainerDefinition{},
				}},
			Expected: errors.Wrap(nil, errInvalidTaskDefinition),
		},
	}

	for i, c := range cases {
		outback := Outback{
			ECS: &mockedDescribeTaskDefinition{Resp: c.Resp},
			ECR: mockedECRClient{},
		}

		_, err := outback.GetLastDeployedCommit("error-taskdef")

		if a, e := err, c.Expected; a != e {
			t.Errorf("%d, expected %v error, got %v", i, e, a)
		}
	}
}

func TestOutbackGetLastDeployedCommitError2(t *testing.T) {
	cases := []struct {
		Error    error
		Expected error
	}{
		{
			Error:    errors.New("test-error"),
			Expected: errors.Wrap(errors.New("test-error"), errCouldNotRetrieveTaskDefinition),
		},
	}

	for i, c := range cases {
		outback := Outback{
			ECS: mockedDescribeTaskDefinition{Resp: nil, Error: c.Error},
			ECR: mockedECRClient{},
		}

		_, err := outback.GetLastDeployedCommit("error")

		if a, e := err, c.Expected; a.Error() != e.Error() {
			t.Errorf("%d, expected %v, got %v", i, e, a)
		}
	}
}

func TestRegisterTaskDefinitionWithEnvVars(t *testing.T) {
	cases := []struct {
		Resp     *ecs.RegisterTaskDefinitionOutput
		Expected *ecs.TaskDefinition
		Input    *ecs.ContainerDefinition
		Error    error
	}{
		{
			Resp: &ecs.RegisterTaskDefinitionOutput{
				TaskDefinition: &ecs.TaskDefinition{
					ContainerDefinitions: []*ecs.ContainerDefinition{{
						Environment: []*ecs.KeyValuePair{{
							Name:  aws.String("KEY1"),
							Value: aws.String("VALUE1"),
						}},
					}},
				},
			},
			Input: &ecs.ContainerDefinition{},
			Error: nil,
			Expected: &ecs.TaskDefinition{
				ContainerDefinitions: []*ecs.ContainerDefinition{{
					Environment: []*ecs.KeyValuePair{{
						Name:  aws.String("KEY1"),
						Value: aws.String("VALUE1"),
					}},
				}},
			},
		},
	}

	for i, c := range cases {
		outback := Outback{
			ECS: mockedRegisterTaskDefinition{Resp: c.Resp, Error: c.Error},
			ECR: mockedECRClient{},
		}

		def, err := outback.RegisterTaskDefinitionWithEnvVars(&ecs.TaskDefinition{
			Cpu:                     aws.String("256"),
			Family:                  aws.String("Family"),
			Memory:                  aws.String("512"),
			Volumes:                 []*ecs.Volume{},
			NetworkMode:             aws.String("Bridge"),
			ExecutionRoleArn:        aws.String("Role"),
			TaskRoleArn:             aws.String("Role"),
			ContainerDefinitions:    []*ecs.ContainerDefinition{c.Input},
			RequiresCompatibilities: aws.StringSlice([]string{}),
		})

		if err != nil {
			t.Fatalf("%d, unexpected error", err)
		}

		a := *def.ContainerDefinitions[0].Environment[0].Value
		e := *c.Expected.ContainerDefinitions[0].Environment[0].Value

		if a != e {
			t.Errorf("%d, expected %v value, got %v", i, e, a)
		}
	}
}

func TestRegisterTaskDefinitionWithEnvVarsError(t *testing.T) {
	cases := []struct {
		Resp     *ecs.RegisterTaskDefinitionOutput
		Expected error
		Input    *ecs.ContainerDefinition
		Error    error
	}{
		{
			Resp:     &ecs.RegisterTaskDefinitionOutput{},
			Input:    &ecs.ContainerDefinition{},
			Error:    errors.New("test-error"),
			Expected: errors.Wrap(errors.New("test-error"), errCouldNotRegisterTaskDefinition),
		},
	}

	for i, c := range cases {
		outback := Outback{
			ECS: mockedRegisterTaskDefinition{Resp: c.Resp, Error: c.Error},
			ECR: mockedECRClient{},
		}

		_, err := outback.RegisterTaskDefinitionWithEnvVars(&ecs.TaskDefinition{})

		if a, e := err, c.Expected; a.Error() != e.Error() {
			t.Errorf("%d, expected %v, got %v", i, e, a)
		}
	}
}

func TestOutbackRunTask(t *testing.T) {
	cases := []struct {
		Resp     *ecs.RunTaskOutput
		Expected *ecs.RunTaskOutput
	}{
		{
			Resp: &ecs.RunTaskOutput{
				Tasks: []*ecs.Task{{
					ClusterArn:        aws.String("clusterarn"),
					TaskDefinitionArn: aws.String("taskdefarn"),
					Overrides: &ecs.TaskOverride{
						ContainerOverrides: []*ecs.ContainerOverride{
							{
								Command: aws.StringSlice([]string{"echo", "this"}),
								Name:    aws.String("task"),
							},
						},
					},
				}},
			},
			Expected: &ecs.RunTaskOutput{
				Tasks: []*ecs.Task{{
					ClusterArn:        aws.String("clusterarn"),
					TaskDefinitionArn: aws.String("taskdefarn"),
					Overrides: &ecs.TaskOverride{
						ContainerOverrides: []*ecs.ContainerOverride{
							{
								Command: aws.StringSlice([]string{"echo", "this"}),
								Name:    aws.String("task"),
							},
						},
					},
				}},
			},
		},
	}

	for i, c := range cases {
		outback := Outback{
			ECS: mockedRunTask{Resp: c.Resp},
			ECR: mockedECRClient{},
		}

		ranTasks, err := outback.RunTask(
			&ecs.Cluster{ClusterName: aws.String("test-cluster")},
			&ecs.TaskDefinition{
				TaskDefinitionArn: aws.String("taskdefarn"),
				ContainerDefinitions: []*ecs.ContainerDefinition{{
					Name: aws.String("test-container"),
				}},
			},
			"echo this",
		)

		if err != nil {
			t.Fatalf("%d, unexpected error", err)
		}

		for j, task := range ranTasks.Tasks {
			if a, e := *task.ClusterArn, *c.Expected.Tasks[j].ClusterArn; a != e {
				t.Errorf("%d, expected %v cluster arn, got %v", i, e, a)
			}

			if a, e := *task.TaskDefinitionArn, *c.Expected.Tasks[j].TaskDefinitionArn; a != e {
				t.Errorf("%d, expected %v task definition arn, got %v", i, e, a)
			}

			actualOverride := task.Overrides.ContainerOverrides[0]
			expectedOverride := c.Expected.Tasks[j].Overrides.ContainerOverrides[0]

			if a, e := *actualOverride.Name, *expectedOverride.Name; a != e {
				t.Errorf("%d, expected %v name, got %v", i, e, a)
			}

			actualOverrideCommand := strings.Join(aws.StringValueSlice(actualOverride.Command), " ")
			expectedOverrideCommand := strings.Join(aws.StringValueSlice(expectedOverride.Command), " ")

			if a, e := actualOverrideCommand, expectedOverrideCommand; a != e {
				t.Errorf("%d, expected %v command override, got %v", i, e, a)
			}
		}
	}
}

func TestOutbackRunTaskError(t *testing.T) {
	cases := []struct {
		Error    error
		Expected error
	}{
		{
			Error:    errors.New("test-error"),
			Expected: errors.Wrap(errors.New("test-error"), errCouldNotRunTask),
		},
	}

	for i, c := range cases {
		outback := Outback{
			ECS: mockedRunTask{Resp: nil, Error: c.Error},
			ECR: mockedECRClient{},
		}

		_, err := outback.RunTask(
			&ecs.Cluster{ClusterName: aws.String("test-cluster")},
			&ecs.TaskDefinition{
				TaskDefinitionArn: aws.String("taskdefarn"),
				ContainerDefinitions: []*ecs.ContainerDefinition{{
					Name: aws.String("test-container"),
				}},
			},
			"error",
		)

		if a, e := err, c.Expected; a.Error() != e.Error() {
			t.Errorf("%d, expected %v, got %v", i, e, a)
		}
	}
}

func TestOutbackUpdateTaskDefinitionImage(t *testing.T) {
	outback := Outback{
		ECS: mockedRunTask{},
		ECR: mockedECRClient{},
	}

	result := outback.UpdateTaskDefinitionImage(ecs.TaskDefinition{
		TaskDefinitionArn: aws.String("taskdefarn"),
		ContainerDefinitions: []*ecs.ContainerDefinition{{
			Name: aws.String("test-container"),
			Image: aws.String("test-container:100"),
		}},
	},"test-container", "123")

	updatedTaskDefinitionImage := result.ContainerDefinitions[0].Image
	expectedTaskDefinitionImage := "test-container:123"

	if *updatedTaskDefinitionImage != expectedTaskDefinitionImage {
		t.Errorf("expected %v value, got %v", expectedTaskDefinitionImage, *updatedTaskDefinitionImage)
	}
}

func TestOutbackUpdateTaskDefinitionImageWithMultipleContainers(t *testing.T) {
	outback := Outback{
		ECS: mockedRunTask{},
		ECR: mockedECRClient{},
	}

	result := outback.UpdateTaskDefinitionImage(ecs.TaskDefinition{
		TaskDefinitionArn: aws.String("taskdefarn"),
		ContainerDefinitions: []*ecs.ContainerDefinition{{
			Name: aws.String("other-container"),
			Image: aws.String("other-container:999"),
		}, {
			Name: aws.String("target-container"),
			Image: aws.String("target-container:100"),
		}},
	}, "target-container", "123")

	cases := []string {"other-container:999", "target-container:123"}

	for i, expectedTaskDefinitionImage := range cases {
		if *result.ContainerDefinitions[i].Image != expectedTaskDefinitionImage {
			t.Errorf("expected %v value, got %v", expectedTaskDefinitionImage, *result.ContainerDefinitions[i].Image)
		}
	}
}

// func TestOutbackDeploy(t *testing.T) {
// 	emptyValue := ""
// 	fam := "family1"
// 	command := "echo"
// 	commands := []*string{&command, &command}
// 	image := "111222333444.dkr.ecr.us-west-1.amazonaws.com/image:cbd0d9c"
// 	commit := "8c018c8"
// 	newImage := fmt.Sprintf("111222333444.dkr.ecr.us-west-1.amazonaws.com/image:%s", commit)

// 	cases := []struct {
// 		DescribeTaskDefResp *ecs.DescribeTaskDefinitionOutput
// 		RegisterTaskDefResp *ecs.RegisterTaskDefinitionOutput
// 		UpdateServiceResp   *ecs.UpdateServiceOutput
// 		Expected            *ecs.TaskDefinition
// 	}{
// 		{
// 			DescribeTaskDefResp: &ecs.DescribeTaskDefinitionOutput{
// 				TaskDefinition: &ecs.TaskDefinition{
// 					ContainerDefinitions: []*ecs.ContainerDefinition{&ecs.ContainerDefinition{
// 						Command: commands,
// 						Image:   &image,
// 					}},
// 					Family:            &fam,
// 					TaskDefinitionArn: aws.String("task-definitionarn"),
// 				},
// 			},
// 			RegisterTaskDefResp: &ecs.RegisterTaskDefinitionOutput{
// 				TaskDefinition: &ecs.TaskDefinition{
// 					Cpu:              &emptyValue,
// 					Family:           &fam,
// 					Memory:           &emptyValue,
// 					Volumes:          []*ecs.Volume{},
// 					NetworkMode:      &emptyValue,
// 					ExecutionRoleArn: &emptyValue,
// 					TaskRoleArn:      &emptyValue,
// 					ContainerDefinitions: []*ecs.ContainerDefinition{&ecs.ContainerDefinition{
// 						Command: commands,
// 						Image:   &image,
// 					}},
// 					RequiresCompatibilities: []*string{&emptyValue},
// 				},
// 			},
// 			UpdateServiceResp: &ecs.UpdateServiceOutput{
// 				Service: &ecs.Service{},
// 			},
// 			Expected: &ecs.TaskDefinition{
// 				ContainerDefinitions: []*ecs.ContainerDefinition{&ecs.ContainerDefinition{
// 					Command: commands,
// 					Image:   &newImage,
// 				}},
// 				Family: &fam,
// 			},
// 		},
// 	}

// 	for i, c := range cases {
// 		outback := Outback{
// 			ECS: mockedDeploy{
// 				DescribeTaskDefResp: c.DescribeTaskDefResp,
// 				RegisterTaskDefResp: c.RegisterTaskDefResp,
// 				UpdateServiceResp:   c.UpdateServiceResp,
// 			},
// 			ECR: mockedECRClient{},
// 		}

// 		newTaskDef, err := outback.Deploy(&ecs.Cluster{}, &ecs.Service{}, commit)

// 		if err != nil {
// 			t.Fatalf("%d, unexpected error", err)
// 		}

// 		if a, e := *newTaskDef.Family, *c.Expected.Family; a != e {
// 			t.Errorf("%d, expected %v family, got %v", i, e, a)
// 		}

// 		if a, e := *newTaskDef.ContainerDefinitions[0].Image, *c.Expected.ContainerDefinitions[0].Image; a != e {
// 			t.Errorf("%d, expected %v image, got %v", i, e, a)
// 		}

// 		actualContainerCommand := strings.Join(aws.StringValueSlice(newTaskDef.ContainerDefinitions[0].Command), " ")
// 		expectedContainerCommand := strings.Join(aws.StringValueSlice(c.Expected.ContainerDefinitions[0].Command), " ")

// 		if a, e := actualContainerCommand, expectedContainerCommand; a != e {
// 			t.Errorf("%d, expected %v command, got %v", i, e, a)
// 		}
// 	}
// }

// func TestOutbackDeployError(t *testing.T) {
// 	fam := "family1"
// 	command := "echo"
// 	commands := []*string{&command, &command}
// 	image := "111222333444.dkr.ecr.us-west-1.amazonaws.com/image:cbd0d9c"
// 	commit := "8c018c8"

// 	cases := []struct {
// 		DescribeTaskDefResp  *ecs.DescribeTaskDefinitionOutput
// 		UpdateServiceResp    *ecs.UpdateServiceOutput
// 		RegisterTaskDefError error
// 		UpdateServiceError   error
// 		Expected             error
// 	}{
// 		{
// 			DescribeTaskDefResp: &ecs.DescribeTaskDefinitionOutput{
// 				TaskDefinition: &ecs.TaskDefinition{
// 					ContainerDefinitions: []*ecs.ContainerDefinition{&ecs.ContainerDefinition{
// 						Command: commands,
// 						Image:   &image,
// 					}},
// 					Family:            &fam,
// 					TaskDefinitionArn: aws.String("task-definitionarn"),
// 				},
// 			},
// 			RegisterTaskDefError: errors.New("test-error"),
// 			Expected:             errors.Wrap(errors.New("test-error"), errCouldNotRegisterTaskDefinition),
// 		},
// 	}

// 	for i, c := range cases {
// 		outback := Outback{
// 			ECS: mockedDeploy{
// 				DescribeTaskDefResp:  c.DescribeTaskDefResp,
// 				RegisterTaskDefResp:  nil,
// 				RegisterTaskDefError: c.RegisterTaskDefError,
// 			},
// 			ECR: mockedECRClient{},
// 		}

// 		_, err := outback.Deploy(&ecs.Cluster{}, &ecs.Service{}, commit)

// 		if a, e := err, c.Expected; a.Error() != e.Error() {
// 			t.Errorf("%d, expected %v, got %v", i, e, a)
// 		}
// 	}
// }

// func TestOutbackDeployError2(t *testing.T) {
// 	fam := "family1"
// 	command := "echo"
// 	emptyValue := ""
// 	commands := []*string{&command, &command}
// 	image := "111222333444.dkr.ecr.us-west-1.amazonaws.com/image:cbd0d9c"
// 	commit := "8c018c8"

// 	cases := []struct {
// 		DescribeTaskDefResp *ecs.DescribeTaskDefinitionOutput
// 		RegisterTaskDefResp *ecs.RegisterTaskDefinitionOutput
// 		UpdateServiceError  error
// 		Expected            error
// 	}{
// 		{
// 			DescribeTaskDefResp: &ecs.DescribeTaskDefinitionOutput{
// 				TaskDefinition: &ecs.TaskDefinition{
// 					ContainerDefinitions: []*ecs.ContainerDefinition{&ecs.ContainerDefinition{
// 						Command: commands,
// 						Image:   &image,
// 					}},
// 					Family:            &fam,
// 					TaskDefinitionArn: aws.String("task-definitionarn"),
// 				},
// 			},
// 			RegisterTaskDefResp: &ecs.RegisterTaskDefinitionOutput{
// 				TaskDefinition: &ecs.TaskDefinition{
// 					Cpu:              &emptyValue,
// 					Family:           &fam,
// 					Memory:           &emptyValue,
// 					Volumes:          []*ecs.Volume{},
// 					NetworkMode:      &emptyValue,
// 					ExecutionRoleArn: &emptyValue,
// 					TaskRoleArn:      &emptyValue,
// 					ContainerDefinitions: []*ecs.ContainerDefinition{&ecs.ContainerDefinition{
// 						Command: commands,
// 						Image:   &image,
// 					}},
// 					RequiresCompatibilities: []*string{&emptyValue},
// 				},
// 			},
// 			UpdateServiceError: errors.New("test-error"),
// 			Expected:           errors.Wrap(errors.New("test-error"), errCouldNotUpdateService),
// 		},
// 	}

// 	for i, c := range cases {
// 		outback := Outback{
// 			ECS: mockedDeploy{
// 				DescribeTaskDefResp: c.DescribeTaskDefResp,
// 				RegisterTaskDefResp: c.RegisterTaskDefResp,
// 				UpdateServiceError:  c.UpdateServiceError,
// 			},
// 			ECR: mockedECRClient{},
// 		}

// 		_, err := outback.Deploy(&ecs.Cluster{}, &ecs.Service{}, commit)

// 		if a, e := err, c.Expected; a.Error() != e.Error() {
// 			t.Errorf("%d, expected %v, got %v", i, e, a)
// 		}
// 	}
// }

// func TestOutbackDeployError3(t *testing.T) {
// 	commit := "8c018c8"

// 	cases := []struct {
// 		DescribeTaskDefError error
// 		Expected             error
// 	}{
// 		{
// 			DescribeTaskDefError: errors.New("test-error"),
// 			Expected:             errors.Wrap(errors.New("test-error"), errCouldNotRetrieveTaskDefinition),
// 		},
// 	}

// 	for i, c := range cases {
// 		outback := Outback{
// 			ECS: mockedDeploy{
// 				DescribeTaskDefError: errors.New("test-error"),
// 			},
// 			ECR: mockedECRClient{},
// 		}

// 		_, err := outback.Deploy(&ecs.Cluster{}, &ecs.Service{}, commit)

// 		if a, e := err, c.Expected; a.Error() != e.Error() {
// 			t.Errorf("%d, expected %v, got %v", i, e, a)
// 		}
// 	}
// }

// func TestOutbackIsServiceRunning(t *testing.T) {
// 	cases := []struct {
// 		ListTasksResp      *ecs.ListTasksOutput
// 		ListTasksError     error
// 		DescribeTasksResp  *ecs.DescribeTasksOutput
// 		DescribeTasksError error
// 		DesiredCount       *int64
// 		Expected           *ServiceStatus
// 	}{
// 		{
// 			ListTasksResp: &ecs.ListTasksOutput{
// 				TaskArns: aws.StringSlice([]string{"task1", "task2"}),
// 			},
// 			DescribeTasksResp: &ecs.DescribeTasksOutput{
// 				Tasks: []*ecs.Task{&ecs.Task{
// 					LastStatus:        aws.String("RUNNING"),
// 					TaskDefinitionArn: aws.String("taskdefarn"),
// 				}},
// 			},
// 			DesiredCount: aws.Int64(2),
// 			Expected: &ServiceStatus{
// 				arn:     "taskdefarn",
// 				running: true,
// 			},
// 		},
// 		{
// 			ListTasksResp: &ecs.ListTasksOutput{
// 				TaskArns: aws.StringSlice([]string{}),
// 			},
// 			DesiredCount: aws.Int64(2),
// 			Expected: &ServiceStatus{
// 				arn:     "",
// 				running: false,
// 			},
// 		},
// 		{
// 			ListTasksResp: &ecs.ListTasksOutput{
// 				TaskArns: aws.StringSlice([]string{"task1", "task2"}),
// 			},
// 			DescribeTasksResp: &ecs.DescribeTasksOutput{
// 				Tasks: []*ecs.Task{&ecs.Task{
// 					LastStatus:        aws.String("PENDING"),
// 					TaskDefinitionArn: aws.String("taskdefarn"),
// 				}},
// 			},
// 			DesiredCount: aws.Int64(2),
// 			Expected: &ServiceStatus{
// 				arn:     "taskdefarn",
// 				running: false,
// 			},
// 		},
// 		{
// 			ListTasksResp: &ecs.ListTasksOutput{
// 				TaskArns: aws.StringSlice([]string{}),
// 			},
// 			DescribeTasksResp: &ecs.DescribeTasksOutput{
// 				Tasks: []*ecs.Task{&ecs.Task{}},
// 			},
// 			DesiredCount: aws.Int64(0),
// 			Expected: &ServiceStatus{
// 				arn:     "",
// 				running: false,
// 			},
// 		},
// 	}

// 	for i, c := range cases {
// 		outback := Outback{
// 			ECS: mockedIsServiceRunning{
// 				ListTasksResp:      c.ListTasksResp,
// 				DescribeTasksResp:  c.DescribeTasksResp,
// 				ListTasksError:     c.ListTasksError,
// 				DescribeTasksError: c.DescribeTasksError,
// 			},
// 			ECR: mockedECRClient{},
// 		}

// 		status, _ := outback.IsServiceRunning(
// 			&ecs.Cluster{ClusterName: aws.String("test-cluster")},
// 			&ecs.Service{ServiceName: aws.String("test-service"), DesiredCount: c.DesiredCount},
// 			&ecs.TaskDefinition{TaskDefinitionArn: aws.String("taskdefarn")})

// 		if a, e := status, c.Expected; a.running != e.running {
// 			t.Errorf("%d, expected %v, got %v", i, e, a)
// 		}
// 	}
// }

// func TestOutbackIsServiceRunningError(t *testing.T) {
// 	cases := []struct {
// 		ListTasksResp      *ecs.ListTasksOutput
// 		ListTasksError     error
// 		DescribeTasksResp  *ecs.DescribeTasksOutput
// 		DescribeTasksError error
// 		DesiredCount       *int64
// 		Expected           error
// 	}{
// 		{
// 			ListTasksError: errors.New("test-error"),
// 			DescribeTasksResp: &ecs.DescribeTasksOutput{
// 				Tasks: []*ecs.Task{&ecs.Task{
// 					LastStatus:        aws.String("RUNNING"),
// 					TaskDefinitionArn: aws.String("taskdefarn"),
// 				}},
// 			},
// 			DesiredCount: aws.Int64(1),
// 			Expected:     errors.Wrap(errors.New("test-error"), errFailedToListRunningTasks),
// 		},
// 		{
// 			ListTasksResp: &ecs.ListTasksOutput{
// 				TaskArns: aws.StringSlice([]string{"task1", "task2"}),
// 			},
// 			DescribeTasksError: errors.New("test-error"),
// 			DesiredCount:       aws.Int64(1),
// 			Expected:           errors.Wrap(errors.New("test-error"), errCouldNotRetrieveTasks),
// 		},
// 	}

// 	for i, c := range cases {
// 		outback := Outback{
// 			ECS: mockedIsServiceRunning{
// 				ListTasksResp:      c.ListTasksResp,
// 				DescribeTasksResp:  c.DescribeTasksResp,
// 				ListTasksError:     c.ListTasksError,
// 				DescribeTasksError: c.DescribeTasksError,
// 			},
// 			ECR: mockedECRClient{},
// 		}

// 		_, err := outback.IsServiceRunning(
// 			&ecs.Cluster{ClusterName: aws.String("test-cluster")},
// 			&ecs.Service{ServiceName: aws.String("test-service"), DesiredCount: c.DesiredCount},
// 			&ecs.TaskDefinition{TaskDefinitionArn: aws.String("taskdefarn")})

// 		if a, e := err, c.Expected; a.Error() != e.Error() {
// 			t.Errorf("%d, expected %v, got %v", i, e, a)
// 		}
// 	}
// }

// func TestOutbackAwaitDeploymentCompletion(t *testing.T) {
// 	cases := []struct {
// 		ExpectedDoneCh     chan bool
// 		ExpectedErrCh      chan error
// 		ListTasksResp      *ecs.ListTasksOutput
// 		ListTasksError     error
// 		DescribeTasksResp  *ecs.DescribeTasksOutput
// 		DescribeTasksError error
// 		DesiredCount       *int64
// 	}{
// 		{
// 			ExpectedDoneCh: make(chan bool),
// 			ExpectedErrCh:  make(chan error),
// 			ListTasksError: errors.New("test-error"),
// 			DescribeTasksResp: &ecs.DescribeTasksOutput{
// 				Tasks: []*ecs.Task{&ecs.Task{
// 					LastStatus:        aws.String("RUNNING"),
// 					TaskDefinitionArn: aws.String("taskdefarn"),
// 				}},
// 			},
// 			DesiredCount: aws.Int64(1),
// 		},
// 	}

// 	for i, c := range cases {
// 		outback := Outback{
// 			,
// 			ECS: mockedIsServiceRunning{
// 				ListTasksResp:      c.ListTasksResp,
// 				DescribeTasksResp:  c.DescribeTasksResp,
// 				ListTasksError:     c.ListTasksError,
// 				DescribeTasksError: c.DescribeTasksError,
// 			},
// 			ECR: mockedECRClient{},
// 		}

// 		doneCh, errCh := outback.AwaitDeploymentCompletion(
// 			&ecs.Cluster{ClusterName: aws.String("test-cluster")},
// 			&ecs.Service{ServiceName: aws.String("test-service"), DesiredCount: c.DesiredCount},
// 			&ecs.TaskDefinition{TaskDefinitionArn: aws.String("taskdefarn")})

// 		// if doneCh, errCh := running, c.Expected; a != e {
// 		// 	t.Errorf("%d, expected %v, got %v", i, e, a)
// 		// }

// 		t.Errorf("%d, expected %v, got %v", i, <-doneCh, errCh)
// 	}
// }
