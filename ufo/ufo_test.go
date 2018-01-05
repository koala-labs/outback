package ufo

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecr"

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

func TestCreateFlyUFO(t *testing.T) {
	cases := []struct {
		Expected *State
	}{
		{
			Expected: &State{},
		},
	}

	config := Config{
		Profile: aws.String("profile"),
		Region:  aws.String("region"),
	}

	ufo := Fly(config)

	for i, c := range cases {
		if a, e := *ufo.State, *c.Expected; a != e {
			t.Errorf("%d, expected %v state, got %v", i, e, a)
		}
	}
}

func TestUFOUseCluster(t *testing.T) {
	name := "cluster1"
	cases := []struct {
		Expected *ecs.Cluster
	}{
		{
			Expected: &ecs.Cluster{
				ClusterName: &name,
			},
		},
	}

	ufo := UFO{
		l:     &logger{},
		State: &State{},
		ECS:   mockedECSClient{},
		ECR:   mockedECRClient{},
	}

	for i, c := range cases {
		ufo.UseCluster(c.Expected)

		if a, e := *ufo.State.Cluster.ClusterName, *c.Expected.ClusterName; a != e {
			t.Errorf("%d, expected %v cluster, got %v", i, e, a)
		}
	}
}
func TestUFOUseService(t *testing.T) {
	name := "service1"
	cases := []struct {
		Expected *ecs.Service
	}{
		{
			Expected: &ecs.Service{
				ServiceName: &name,
			},
		},
	}

	ufo := UFO{
		l:     &logger{},
		State: &State{},
		ECS:   mockedECSClient{},
		ECR:   mockedECRClient{},
	}

	for i, c := range cases {
		ufo.UseService(c.Expected)

		if a, e := *ufo.State.Service.ServiceName, *c.Expected.ServiceName; a != e {
			t.Errorf("%d, expected %v service, got %v", i, e, a)
		}
	}
}
func TestUFOUseTaskDefinition(t *testing.T) {
	fam := "family1"
	command := "echo"
	commands := []*string{&command, &command}
	image := "111222333444.dkr.ecr.us-west-1.amazonaws.com/image:ea13366"
	cases := []struct {
		Expected *ecs.TaskDefinition
	}{
		{
			Expected: &ecs.TaskDefinition{
				ContainerDefinitions: []*ecs.ContainerDefinition{&ecs.ContainerDefinition{
					Command: commands,
					Image:   &image,
				}},
				Family: &fam,
			},
		},
	}

	ufo := UFO{
		l:     &logger{},
		State: &State{},
		ECS:   mockedECSClient{},
		ECR:   mockedECRClient{},
	}

	for i, c := range cases {
		ufo.UseTaskDefinition(c.Expected)
		actualContainerDefinition := ufo.State.TaskDefinition.ContainerDefinitions[0]
		expectedContainerDefinition := c.Expected.ContainerDefinitions[0]

		if a, e := *ufo.State.TaskDefinition.Family, *c.Expected.Family; a != e {
			t.Errorf("%d, expected %v task definition, got %v", i, e, a)
		}

		if a, e := *ufo.State.TaskDefinition.Family, *c.Expected.Family; a != e {
			t.Errorf("%d, expected %v family, got %v", i, e, a)
		}

		if a, e := *actualContainerDefinition.Image, *expectedContainerDefinition.Image; a != e {
			t.Errorf("%d, expected %v image, got %v", i, e, a)
		}
	}
}

func TestUFOClusters(t *testing.T) {
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
		ufo := UFO{
			l:     &logger{},
			State: &State{},
			ECS:   &mockedListClusters{Resp: c.Resp},
			ECR:   mockedECRClient{},
		}

		clusters, err := ufo.Clusters()

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

func TestUFOClustersError(t *testing.T) {
	cases := []struct {
		Error    error
		Expected error
	}{
		{
			Error:    awserr.New("0", "ERROR", fmt.Errorf("error")),
			Expected: ErrFailedToListClusters,
		},
	}

	for i, c := range cases {
		ufo := UFO{
			l:     &logger{},
			State: &State{},
			ECS:   &mockedListClusters{Resp: nil, Error: c.Error},
			ECR:   mockedECRClient{},
		}

		_, err := ufo.Clusters()

		if a, e := err, c.Expected; a != e {
			t.Errorf("%d, expected %v error, got %v", i, e, a)
		}
	}
}

func TestUFOServices(t *testing.T) {
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
		ufo := UFO{
			l:     &logger{},
			State: &State{},
			ECS:   &mockedListServices{Resp: c.Resp},
			ECR:   mockedECRClient{},
		}

		services, err := ufo.Services(&ecs.Cluster{})

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

func TestUFOServicesError(t *testing.T) {
	cases := []struct {
		Error    error
		Expected error
	}{
		{
			Error:    awserr.New("0", "ERROR", fmt.Errorf("error")),
			Expected: ErrFailedToListServices,
		},
	}

	for i, c := range cases {
		ufo := UFO{
			l:     &logger{},
			State: &State{},
			ECS:   &mockedListServices{Resp: nil, Error: c.Error},
			ECR:   mockedECRClient{},
		}

		_, err := ufo.Services(&ecs.Cluster{})

		if a, e := err, c.Expected; a != e {
			t.Errorf("%d, expected %v error, got %v", i, e, a)
		}
	}
}

func TestUFORunningTasks(t *testing.T) {
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
		ufo := UFO{
			l:     &logger{},
			State: &State{},
			ECS:   &mockedListTasks{Resp: c.Resp},
			ECR:   mockedECRClient{},
		}

		tasks, err := ufo.RunningTasks(&ecs.Cluster{}, &ecs.Service{})

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

func TestUFORunningTasksError(t *testing.T) {
	cases := []struct {
		Error    error
		Expected error
	}{
		{
			Error:    awserr.New("0", "ERROR", fmt.Errorf("error")),
			Expected: ErrFailedToListRunningTasks,
		},
	}

	for i, c := range cases {
		ufo := UFO{
			l:     &logger{},
			State: &State{},
			ECS:   &mockedListTasks{Resp: nil, Error: c.Error},
			ECR:   mockedECRClient{},
		}

		_, err := ufo.RunningTasks(&ecs.Cluster{}, &ecs.Service{})

		if a, e := err, c.Expected; a != e {
			t.Errorf("%d, expected %v error, got %v", i, e, a)
		}
	}
}

func TestUFOGetCluster(t *testing.T) {
	cases := []struct {
		Resp     *ecs.DescribeClustersOutput
		Expected *ecs.Cluster
	}{
		{
			Resp: &ecs.DescribeClustersOutput{
				Clusters: []*ecs.Cluster{&ecs.Cluster{
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
		ufo := UFO{
			l:     &logger{},
			State: &State{},
			ECS:   &mockedDescribeClusters{Resp: c.Resp},
			ECR:   mockedECRClient{},
		}

		cluster, err := ufo.GetCluster("test-cluster")

		if err != nil {
			t.Fatalf("%d, unexpected error", err)
		}

		if a, e := *cluster.ClusterName, *c.Expected.ClusterName; a != e {
			t.Errorf("%d, expected %v cluster, got %v", i, e, a)
		}
	}
}

func TestUFOGetClusterError(t *testing.T) {
	cases := []struct {
		Resp     *ecs.DescribeClustersOutput
		Expected error
	}{
		{
			Resp: &ecs.DescribeClustersOutput{
				Clusters: []*ecs.Cluster{},
			},
			Expected: ErrClusterNotFound,
		},
	}

	for i, c := range cases {
		ufo := UFO{
			l:     &logger{},
			State: &State{},
			ECS:   &mockedDescribeClusters{Resp: c.Resp},
			ECR:   mockedECRClient{},
		}

		_, err := ufo.GetCluster("test-cluster")

		if a, e := err, c.Expected; a != e {
			t.Errorf("%d, expected %v error, got %v", i, e, a)
		}
	}
}

func TestUFOGetClusterError2(t *testing.T) {
	cases := []struct {
		Error    error
		Expected error
	}{
		{
			Error:    awserr.New("0", "ERROR", fmt.Errorf("error")),
			Expected: ErrCouldNotRetrieveCluster,
		},
	}

	for i, c := range cases {
		ufo := UFO{
			l:     &logger{},
			State: &State{},
			ECS:   &mockedDescribeClusters{Resp: nil, Error: c.Error},
			ECR:   mockedECRClient{},
		}

		_, err := ufo.GetCluster("test-cluster")

		if a, e := err, c.Expected; a != e {
			t.Errorf("%d, expected %v error, got %v", i, e, a)
		}
	}
}

func TestUFOGetService(t *testing.T) {
	name := "test-service"
	arn := "test-servicearn"
	cases := []struct {
		Resp     *ecs.DescribeServicesOutput
		Expected *ecs.Service
	}{
		{
			Resp: &ecs.DescribeServicesOutput{
				Services: []*ecs.Service{&ecs.Service{
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
		ufo := UFO{
			l:     &logger{},
			State: &State{},
			ECS:   &mockedDescribeServices{Resp: c.Resp},
			ECR:   mockedECRClient{},
		}

		service, err := ufo.GetService(&ecs.Cluster{}, "test-service")

		if err != nil {
			t.Fatalf("%d, unexpected error", err)
		}

		if a, e := *service.ServiceName, *c.Expected.ServiceName; a != e {
			t.Errorf("%d, expected %v service, got %v", i, e, a)
		}
	}
}

func TestUFOGetServiceError(t *testing.T) {
	cases := []struct {
		Resp     *ecs.DescribeServicesOutput
		Expected error
	}{
		{
			Resp: &ecs.DescribeServicesOutput{
				Services: []*ecs.Service{},
			},
			Expected: ErrServiceNotFound,
		},
	}

	for i, c := range cases {
		ufo := UFO{
			l:     &logger{},
			State: &State{},
			ECS:   &mockedDescribeServices{Resp: c.Resp},
			ECR:   mockedECRClient{},
		}

		_, err := ufo.GetService(&ecs.Cluster{}, "test-service")

		if a, e := err, c.Expected; a != e {
			t.Errorf("%d, expected %v error, got %v", i, e, a)
		}
	}
}

func TestUFOGetServiceError2(t *testing.T) {
	cases := []struct {
		Error    error
		Expected error
	}{
		{
			Error:    awserr.New("0", "ERROR", fmt.Errorf("error")),
			Expected: ErrCouldNotRetrieveService,
		},
	}

	for i, c := range cases {
		ufo := UFO{
			l:     &logger{},
			State: &State{},
			ECS:   &mockedDescribeServices{Resp: nil, Error: c.Error},
			ECR:   mockedECRClient{},
		}

		_, err := ufo.GetService(&ecs.Cluster{}, "test-cluster")

		if a, e := err, c.Expected; a != e {
			t.Errorf("%d, expected %v error, got %v", i, e, a)
		}
	}
}

func TestUFOGetTaskDefinition(t *testing.T) {
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
					ContainerDefinitions: []*ecs.ContainerDefinition{&ecs.ContainerDefinition{
						Command: command,
						Image:   &image,
					}},
					Family:            &fam,
					TaskDefinitionArn: &arn,
				},
			},
			Expected: &ecs.TaskDefinition{
				ContainerDefinitions: []*ecs.ContainerDefinition{&ecs.ContainerDefinition{
					Command: command,
					Image:   &image,
				}},
				Family: &fam,
			},
		},
	}

	for i, c := range cases {
		ufo := UFO{
			l:     &logger{},
			State: &State{},
			ECS:   &mockedDescribeTaskDefinition{Resp: c.Resp},
			ECR:   mockedECRClient{},
		}

		taskDef, err := ufo.GetTaskDefinition(&ecs.Cluster{}, &ecs.Service{})

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

func TestUFOGetTaskDefinitionError(t *testing.T) {
	cases := []struct {
		Error    error
		Expected error
	}{
		{
			Error:    awserr.New("0", "ERROR", fmt.Errorf("error")),
			Expected: ErrCouldNotRetrieveTaskDefinition,
		},
	}

	for i, c := range cases {
		ufo := UFO{
			l:     &logger{},
			State: &State{},
			ECS:   &mockedDescribeTaskDefinition{Resp: nil, Error: c.Error},
			ECR:   mockedECRClient{},
		}

		_, err := ufo.GetTaskDefinition(&ecs.Cluster{}, &ecs.Service{})

		if a, e := err, c.Expected; a != e {
			t.Errorf("%d, expected %v error, got %v", i, e, a)
		}
	}
}

func TestUFOGetTasks(t *testing.T) {
	var lastStatus, taskDefArn string = "PENDING", "111222333444.dkr.ecr.us-west-1.amazonaws.com/task"
	cases := []struct {
		Resp     *ecs.DescribeTasksOutput
		Expected *ecs.DescribeTasksOutput
	}{
		{
			Resp: &ecs.DescribeTasksOutput{
				Tasks: []*ecs.Task{&ecs.Task{
					LastStatus:        &lastStatus,
					TaskDefinitionArn: &taskDefArn,
				}},
			},
			Expected: &ecs.DescribeTasksOutput{
				Tasks: []*ecs.Task{
					&ecs.Task{
						LastStatus:        &lastStatus,
						TaskDefinitionArn: &taskDefArn,
					}},
			},
		},
	}

	for i, c := range cases {
		ufo := UFO{
			l:     &logger{},
			State: &State{},
			ECS:   &mockedDescribeTasks{Resp: c.Resp},
			ECR:   mockedECRClient{},
		}

		tasks, err := ufo.GetTasks(&ecs.Cluster{}, []*string{})

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

func TestUFOGetTasksError(t *testing.T) {
	cases := []struct {
		Error    error
		Expected error
	}{
		{
			Error:    awserr.New("0", "ERROR", fmt.Errorf("error")),
			Expected: ErrCouldNotRetrieveTasks,
		},
	}

	for i, c := range cases {
		ufo := UFO{
			l:     &logger{},
			State: &State{},
			ECS:   &mockedDescribeTasks{Resp: nil, Error: c.Error},
			ECR:   mockedECRClient{},
		}

		_, err := ufo.GetTasks(&ecs.Cluster{}, aws.StringSlice([]string{"task1", "task2"}))

		if a, e := err, c.Expected; a != e {
			t.Errorf("%d, expected %v error, got %v", i, e, a)
		}
	}
}

func TestUFOGetImages(t *testing.T) {
	var tag1, tag2 string = "tag1", "tag2"
	cases := []struct {
		Resp     *ecr.DescribeImagesOutput
		Expected []*ecr.ImageDetail
	}{
		{
			Resp: &ecr.DescribeImagesOutput{
				ImageDetails: []*ecr.ImageDetail{
					&ecr.ImageDetail{
						ImageTags: []*string{&tag1},
					},
					&ecr.ImageDetail{
						ImageTags: []*string{&tag2},
					}},
			},
			Expected: []*ecr.ImageDetail{
				&ecr.ImageDetail{
					ImageTags: []*string{&tag1},
				},
				&ecr.ImageDetail{
					ImageTags: []*string{&tag2},
				}},
		},
		{
			Resp:     &ecr.DescribeImagesOutput{},
			Expected: []*ecr.ImageDetail{},
		},
	}

	for i, c := range cases {
		ufo := UFO{
			l:     &logger{},
			State: &State{},
			ECS:   mockedECSClient{},
			ECR:   &mockedDescribeImages{Resp: c.Resp},
		}

		images, err := ufo.GetImages(&ecs.TaskDefinition{
			ContainerDefinitions: []*ecs.ContainerDefinition{&ecs.ContainerDefinition{
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

func TestUFOGetImagesError(t *testing.T) {
	cases := []struct {
		Error    error
		Expected error
	}{
		{
			Error:    awserr.New("0", "ERROR", fmt.Errorf("error")),
			Expected: ErrCouldNotRetrieveImages,
		},
	}

	for i, c := range cases {
		ufo := UFO{
			l:     &logger{},
			State: &State{},
			ECS:   mockedECSClient{},
			ECR:   &mockedDescribeImages{Resp: nil, Error: c.Error},
		}

		_, err := ufo.GetImages(&ecs.TaskDefinition{
			ContainerDefinitions: []*ecs.ContainerDefinition{&ecs.ContainerDefinition{
				Command: aws.StringSlice([]string{"error", "this"}),
				Image:   aws.String("111222333444.dkr.ecr.us-west-1.amazonaws.com/image:ea13366"),
			}},
			Family:            aws.String("family"),
			TaskDefinitionArn: aws.String("taskdefarn"),
		})

		if a, e := err, c.Expected; a != e {
			t.Errorf("%d, expected %v error, got %v", i, e, a)
		}
	}
}

func TestUFOGetLastDeployedCommit(t *testing.T) {
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
					ContainerDefinitions: []*ecs.ContainerDefinition{&ecs.ContainerDefinition{
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
		ufo := UFO{
			l:     &logger{},
			State: &State{},
			ECS:   &mockedDescribeTaskDefinition{Resp: c.Resp},
			ECR:   mockedECRClient{},
		}

		commit, err := ufo.GetLastDeployedCommit("111222333444.dkr.ecr.us-west-1.amazonaws.com/image:ea13366")

		if err != nil {
			t.Fatalf("%d, unexpected error", err)
		}

		if a, e := commit, c.Expected; a != e {
			t.Errorf("%d, expected %v commit, got %v", i, e, a)
		}
	}
}

func TestUFOGetLastDeployedCommitError(t *testing.T) {
	cases := []struct {
		Resp     *ecs.DescribeTaskDefinitionOutput
		Expected error
	}{
		{
			Resp: &ecs.DescribeTaskDefinitionOutput{
				TaskDefinition: &ecs.TaskDefinition{
					ContainerDefinitions: []*ecs.ContainerDefinition{},
				}},
			Expected: ErrInvalidTaskDefinition,
		},
	}

	for i, c := range cases {
		ufo := UFO{
			l:     &logger{},
			State: &State{},
			ECS:   &mockedDescribeTaskDefinition{Resp: c.Resp},
			ECR:   mockedECRClient{},
		}

		_, err := ufo.GetLastDeployedCommit("error-taskdef")

		if a, e := err, c.Expected; a != e {
			t.Errorf("%d, expected %v error, got %v", i, e, a)
		}
	}
}

func TestUFOGetLastDeployedCommitError2(t *testing.T) {
	cases := []struct {
		Error    error
		Expected error
	}{
		{
			Error:    awserr.New("0", "ERROR", fmt.Errorf("error")),
			Expected: ErrCouldNotRetrieveTaskDefinition,
		},
	}

	for i, c := range cases {
		ufo := UFO{
			l:     &logger{},
			State: &State{},
			ECS:   mockedDescribeTaskDefinition{Resp: nil, Error: c.Error},
			ECR:   mockedECRClient{},
		}

		_, err := ufo.GetLastDeployedCommit("error")

		if a, e := err, c.Expected; a != e {
			t.Errorf("%d, expected %v error, got %v", i, e, a)
		}
	}
}

func TestUFORunTask(t *testing.T) {
	cases := []struct {
		Resp     *ecs.RunTaskOutput
		Expected *ecs.RunTaskOutput
	}{
		{
			Resp: &ecs.RunTaskOutput{
				Tasks: []*ecs.Task{&ecs.Task{
					ClusterArn:        aws.String("clusterarn"),
					TaskDefinitionArn: aws.String("taskdefarn"),
					Overrides: &ecs.TaskOverride{
						ContainerOverrides: []*ecs.ContainerOverride{
							&ecs.ContainerOverride{
								Command: aws.StringSlice([]string{"echo", "this"}),
								Name:    aws.String("task"),
							},
						},
					},
				}},
			},
			Expected: &ecs.RunTaskOutput{
				Tasks: []*ecs.Task{&ecs.Task{
					ClusterArn:        aws.String("clusterarn"),
					TaskDefinitionArn: aws.String("taskdefarn"),
					Overrides: &ecs.TaskOverride{
						ContainerOverrides: []*ecs.ContainerOverride{
							&ecs.ContainerOverride{
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
		ufo := UFO{
			l:     &logger{},
			State: &State{},
			ECS:   mockedRunTask{Resp: c.Resp},
			ECR:   mockedECRClient{},
		}

		ranTasks, err := ufo.RunTask(
			&ecs.Cluster{ClusterName: aws.String("test-cluster")},
			&ecs.TaskDefinition{
				TaskDefinitionArn: aws.String("taskdefarn"),
				ContainerDefinitions: []*ecs.ContainerDefinition{&ecs.ContainerDefinition{
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

func TestUFORunTaskError(t *testing.T) {
	cases := []struct {
		Error    error
		Expected error
	}{
		{
			Error:    awserr.New("0", "ERROR", fmt.Errorf("error")),
			Expected: ErrCouldNotRunTask,
		},
	}

	for i, c := range cases {
		ufo := UFO{
			l:     &logger{},
			State: &State{},
			ECS:   mockedRunTask{Resp: nil, Error: c.Error},
			ECR:   mockedECRClient{},
		}

		_, err := ufo.RunTask(
			&ecs.Cluster{ClusterName: aws.String("test-cluster")},
			&ecs.TaskDefinition{
				TaskDefinitionArn: aws.String("taskdefarn"),
				ContainerDefinitions: []*ecs.ContainerDefinition{&ecs.ContainerDefinition{
					Name: aws.String("test-container"),
				}},
			},
			"error",
		)

		if a, e := err, c.Expected; a != e {
			t.Errorf("%d, expected %v error, got %v", i, e, a)
		}
	}
}

func TestUFODeploy(t *testing.T) {
	emptyValue := ""
	fam := "family1"
	command := "echo"
	commands := []*string{&command, &command}
	image := "111222333444.dkr.ecr.us-west-1.amazonaws.com/image:cbd0d9c"
	commit := "8c018c8"
	newImage := fmt.Sprintf("111222333444.dkr.ecr.us-west-1.amazonaws.com/image:%s", commit)

	cases := []struct {
		DescribeTaskDefResp *ecs.DescribeTaskDefinitionOutput
		RegisterTaskDefResp *ecs.RegisterTaskDefinitionOutput
		UpdateServiceResp   *ecs.UpdateServiceOutput
		Expected            *ecs.TaskDefinition
	}{
		{
			DescribeTaskDefResp: &ecs.DescribeTaskDefinitionOutput{
				TaskDefinition: &ecs.TaskDefinition{
					ContainerDefinitions: []*ecs.ContainerDefinition{&ecs.ContainerDefinition{
						Command: commands,
						Image:   &image,
					}},
					Family:            &fam,
					TaskDefinitionArn: aws.String("task-definitionarn"),
				},
			},
			RegisterTaskDefResp: &ecs.RegisterTaskDefinitionOutput{
				TaskDefinition: &ecs.TaskDefinition{
					Cpu:              &emptyValue,
					Family:           &fam,
					Memory:           &emptyValue,
					Volumes:          []*ecs.Volume{},
					NetworkMode:      &emptyValue,
					ExecutionRoleArn: &emptyValue,
					TaskRoleArn:      &emptyValue,
					ContainerDefinitions: []*ecs.ContainerDefinition{&ecs.ContainerDefinition{
						Command: commands,
						Image:   &image,
					}},
					RequiresCompatibilities: []*string{&emptyValue},
				},
			},
			UpdateServiceResp: &ecs.UpdateServiceOutput{
				Service: &ecs.Service{},
			},
			Expected: &ecs.TaskDefinition{
				ContainerDefinitions: []*ecs.ContainerDefinition{&ecs.ContainerDefinition{
					Command: commands,
					Image:   &newImage,
				}},
				Family: &fam,
			},
		},
	}

	for i, c := range cases {
		ufo := UFO{
			l:     &logger{},
			State: &State{},
			ECS: mockedDeploy{
				DescribeTaskDefResp: c.DescribeTaskDefResp,
				RegisterTaskDefResp: c.RegisterTaskDefResp,
				UpdateServiceResp:   c.UpdateServiceResp,
			},
			ECR: mockedECRClient{},
		}

		newTaskDef, err := ufo.Deploy(&ecs.Cluster{}, &ecs.Service{}, commit)

		if err != nil {
			t.Fatalf("%d, unexpected error", err)
		}

		if a, e := *newTaskDef.Family, *c.Expected.Family; a != e {
			t.Errorf("%d, expected %v family, got %v", i, e, a)
		}

		if a, e := *newTaskDef.ContainerDefinitions[0].Image, *c.Expected.ContainerDefinitions[0].Image; a != e {
			t.Errorf("%d, expected %v image, got %v", i, e, a)
		}

		actualContainerCommand := strings.Join(aws.StringValueSlice(newTaskDef.ContainerDefinitions[0].Command), " ")
		expectedContainerCommand := strings.Join(aws.StringValueSlice(c.Expected.ContainerDefinitions[0].Command), " ")

		if a, e := actualContainerCommand, expectedContainerCommand; a != e {
			t.Errorf("%d, expected %v command, got %v", i, e, a)
		}
	}
}

func TestUFODeployError(t *testing.T) {
	fam := "family1"
	command := "echo"
	commands := []*string{&command, &command}
	image := "111222333444.dkr.ecr.us-west-1.amazonaws.com/image:cbd0d9c"
	commit := "8c018c8"

	cases := []struct {
		DescribeTaskDefResp  *ecs.DescribeTaskDefinitionOutput
		UpdateServiceResp    *ecs.UpdateServiceOutput
		RegisterTaskDefError error
		UpdateServiceError   error
		Expected             error
	}{
		{
			DescribeTaskDefResp: &ecs.DescribeTaskDefinitionOutput{
				TaskDefinition: &ecs.TaskDefinition{
					ContainerDefinitions: []*ecs.ContainerDefinition{&ecs.ContainerDefinition{
						Command: commands,
						Image:   &image,
					}},
					Family:            &fam,
					TaskDefinitionArn: aws.String("task-definitionarn"),
				},
			},
			RegisterTaskDefError: awserr.New("0", "ERROR", fmt.Errorf("error")),
			Expected:             ErrCouldNotRegisterTaskDefinition,
		},
	}

	for i, c := range cases {
		ufo := UFO{
			l:     &logger{},
			State: &State{},
			ECS: mockedDeploy{
				DescribeTaskDefResp:  c.DescribeTaskDefResp,
				RegisterTaskDefResp:  nil,
				RegisterTaskDefError: c.RegisterTaskDefError,
			},
			ECR: mockedECRClient{},
		}

		_, err := ufo.Deploy(&ecs.Cluster{}, &ecs.Service{}, commit)

		if a, e := err, c.Expected; a != e {
			t.Errorf("%d, expected %v error, got %v", i, e, a)
		}
	}
}

func TestUFODeployError2(t *testing.T) {
	fam := "family1"
	command := "echo"
	emptyValue := ""
	commands := []*string{&command, &command}
	image := "111222333444.dkr.ecr.us-west-1.amazonaws.com/image:cbd0d9c"
	commit := "8c018c8"

	cases := []struct {
		DescribeTaskDefResp *ecs.DescribeTaskDefinitionOutput
		RegisterTaskDefResp *ecs.RegisterTaskDefinitionOutput
		UpdateServiceError  error
		Expected            error
	}{
		{
			DescribeTaskDefResp: &ecs.DescribeTaskDefinitionOutput{
				TaskDefinition: &ecs.TaskDefinition{
					ContainerDefinitions: []*ecs.ContainerDefinition{&ecs.ContainerDefinition{
						Command: commands,
						Image:   &image,
					}},
					Family:            &fam,
					TaskDefinitionArn: aws.String("task-definitionarn"),
				},
			},
			RegisterTaskDefResp: &ecs.RegisterTaskDefinitionOutput{
				TaskDefinition: &ecs.TaskDefinition{
					Cpu:              &emptyValue,
					Family:           &fam,
					Memory:           &emptyValue,
					Volumes:          []*ecs.Volume{},
					NetworkMode:      &emptyValue,
					ExecutionRoleArn: &emptyValue,
					TaskRoleArn:      &emptyValue,
					ContainerDefinitions: []*ecs.ContainerDefinition{&ecs.ContainerDefinition{
						Command: commands,
						Image:   &image,
					}},
					RequiresCompatibilities: []*string{&emptyValue},
				},
			},
			UpdateServiceError: awserr.New("0", "ERROR", fmt.Errorf("error")),
			Expected:           ErrCouldNotUpdateService,
		},
	}

	for i, c := range cases {
		ufo := UFO{
			l:     &logger{},
			State: &State{},
			ECS: mockedDeploy{
				DescribeTaskDefResp: c.DescribeTaskDefResp,
				RegisterTaskDefResp: c.RegisterTaskDefResp,
				UpdateServiceError:  c.UpdateServiceError,
			},
			ECR: mockedECRClient{},
		}

		_, err := ufo.Deploy(&ecs.Cluster{}, &ecs.Service{}, commit)

		if a, e := err, c.Expected; a != e {
			t.Errorf("%d, expected %v error, got %v", i, e, a)
		}
	}
}

func TestUFODeployError3(t *testing.T) {
	commit := "8c018c8"

	cases := []struct {
		DescribeTaskDefError error
		Expected             error
	}{
		{
			DescribeTaskDefError: awserr.New("0", "ERROR", fmt.Errorf("error")),
			Expected:             ErrCouldNotRetrieveTaskDefinition,
		},
	}

	for i, c := range cases {
		ufo := UFO{
			l:     &logger{},
			State: &State{},
			ECS: mockedDeploy{
				DescribeTaskDefError: awserr.New("0", "ERROR", fmt.Errorf("error")),
			},
			ECR: mockedECRClient{},
		}

		_, err := ufo.Deploy(&ecs.Cluster{}, &ecs.Service{}, commit)

		if a, e := err, c.Expected; a != e {
			t.Errorf("%d, expected %v error, got %v", i, e, a)
		}
	}
}

func TestUFOIsServiceRunning(t *testing.T) {
	cases := []struct {
		ListTasksResp      *ecs.ListTasksOutput
		ListTasksError     error
		DescribeTasksResp  *ecs.DescribeTasksOutput
		DescribeTasksError error
		DesiredCount       *int64
		Expected           bool
	}{
		{
			ListTasksResp: &ecs.ListTasksOutput{
				TaskArns: aws.StringSlice([]string{"task1", "task2"}),
			},
			DescribeTasksResp: &ecs.DescribeTasksOutput{
				Tasks: []*ecs.Task{&ecs.Task{
					LastStatus:        aws.String("RUNNING"),
					TaskDefinitionArn: aws.String("taskdefarn"),
				}},
			},
			DesiredCount: aws.Int64(2),
			Expected:     true,
		},
		{
			ListTasksResp: &ecs.ListTasksOutput{
				TaskArns: aws.StringSlice([]string{}),
			},
			DesiredCount: aws.Int64(2),
			Expected:     false,
		},
		{
			ListTasksResp: &ecs.ListTasksOutput{
				TaskArns: aws.StringSlice([]string{"task1", "task2"}),
			},
			DescribeTasksResp: &ecs.DescribeTasksOutput{
				Tasks: []*ecs.Task{&ecs.Task{
					LastStatus:        aws.String("PENDING"),
					TaskDefinitionArn: aws.String("taskdefarn"),
				}},
			},
			DesiredCount: aws.Int64(2),
			Expected:     false,
		},
		{
			ListTasksResp: &ecs.ListTasksOutput{
				TaskArns: aws.StringSlice([]string{}),
			},
			DescribeTasksResp: &ecs.DescribeTasksOutput{
				Tasks: []*ecs.Task{&ecs.Task{}},
			},
			DesiredCount: aws.Int64(0),
			Expected:     false,
		},
	}

	for i, c := range cases {
		ufo := UFO{
			l:     &logger{},
			State: &State{},
			ECS: mockedIsServiceRunning{
				ListTasksResp:      c.ListTasksResp,
				DescribeTasksResp:  c.DescribeTasksResp,
				ListTasksError:     c.ListTasksError,
				DescribeTasksError: c.DescribeTasksError,
			},
			ECR: mockedECRClient{},
		}

		running, _ := ufo.IsServiceRunning(
			&ecs.Cluster{ClusterName: aws.String("test-cluster")},
			&ecs.Service{ServiceName: aws.String("test-service"), DesiredCount: c.DesiredCount},
			&ecs.TaskDefinition{TaskDefinitionArn: aws.String("taskdefarn")})

		if a, e := running, c.Expected; a != e {
			t.Errorf("%d, expected %v, got %v", i, e, a)
		}
	}
}

func TestUFOIsServiceRunningError(t *testing.T) {
	cases := []struct {
		ListTasksResp      *ecs.ListTasksOutput
		ListTasksError     error
		DescribeTasksResp  *ecs.DescribeTasksOutput
		DescribeTasksError error
		DesiredCount       *int64
		Expected           error
	}{
		{
			ListTasksError: awserr.New("0", "ERROR", fmt.Errorf("error")),
			DescribeTasksResp: &ecs.DescribeTasksOutput{
				Tasks: []*ecs.Task{&ecs.Task{
					LastStatus:        aws.String("RUNNING"),
					TaskDefinitionArn: aws.String("taskdefarn"),
				}},
			},
			DesiredCount: aws.Int64(1),
			Expected:     ErrFailedToListRunningTasks,
		},
		{
			ListTasksResp: &ecs.ListTasksOutput{
				TaskArns: aws.StringSlice([]string{"task1", "task2"}),
			},
			DescribeTasksError: awserr.New("0", "ERROR", fmt.Errorf("error")),
			DesiredCount:       aws.Int64(1),
			Expected:           ErrCouldNotRetrieveTasks,
		},
	}

	for i, c := range cases {
		ufo := UFO{
			l:     &logger{},
			State: &State{},
			ECS: mockedIsServiceRunning{
				ListTasksResp:      c.ListTasksResp,
				DescribeTasksResp:  c.DescribeTasksResp,
				ListTasksError:     c.ListTasksError,
				DescribeTasksError: c.DescribeTasksError,
			},
			ECR: mockedECRClient{},
		}

		_, err := ufo.IsServiceRunning(
			&ecs.Cluster{ClusterName: aws.String("test-cluster")},
			&ecs.Service{ServiceName: aws.String("test-service"), DesiredCount: c.DesiredCount},
			&ecs.TaskDefinition{TaskDefinitionArn: aws.String("taskdefarn")})

		if a, e := err, c.Expected; a != e {
			t.Errorf("%d, expected %v error, got %v", i, e, a)
		}
	}
}
