package ufo

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/ecr"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
)

type logger struct{}

func (l *logger) Printf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

type mockECRClient struct {
	ecriface.ECRAPI
}

type mockedDescribeImages struct {
	ecriface.ECRAPI
	Resp ecr.DescribeImagesOutput
}

type mockECSClient struct {
	ecsiface.ECSAPI
}

type mockedDescribeClusters struct {
	ecsiface.ECSAPI
	Resp ecs.DescribeClustersOutput
}

type mockedDescribeServices struct {
	ecsiface.ECSAPI
	Resp ecs.DescribeServicesOutput
}

type mockedDescribeTaskDefinition struct {
	ecsiface.ECSAPI
	Resp ecs.DescribeTaskDefinitionOutput
}

type mockedDescribeTasks struct {
	ecsiface.ECSAPI
	Resp ecs.DescribeTasksOutput
}

type mockedListClusters struct {
	ecsiface.ECSAPI
	Resp ecs.ListClustersOutput
}

type mockedListServices struct {
	ecsiface.ECSAPI
	Resp ecs.ListServicesOutput
}

type mockedListTasks struct {
	ecsiface.ECSAPI
	Resp ecs.ListTasksOutput
}

type mockedRunTask struct {
	ecsiface.ECSAPI
	Resp ecs.RunTaskOutput
}

type mockedDeploy struct {
	ecsiface.ECSAPI
	DescribeTaskDefResp ecs.DescribeTaskDefinitionOutput
	RegisterTaskDefResp ecs.RegisterTaskDefinitionOutput
	UpdateServiceResp   ecs.UpdateServiceOutput
}

func (m mockedDescribeClusters) DescribeClusters(in *ecs.DescribeClustersInput) (*ecs.DescribeClustersOutput, error) {
	return &m.Resp, nil
}

func (m mockedDescribeImages) DescribeImages(in *ecr.DescribeImagesInput) (*ecr.DescribeImagesOutput, error) {
	return &m.Resp, nil
}

func (m mockedDescribeServices) DescribeServices(in *ecs.DescribeServicesInput) (*ecs.DescribeServicesOutput, error) {
	return &m.Resp, nil
}

func (m mockedDescribeTaskDefinition) DescribeTaskDefinition(in *ecs.DescribeTaskDefinitionInput) (*ecs.DescribeTaskDefinitionOutput, error) {
	return &m.Resp, nil
}

func (m mockedDescribeTasks) DescribeTasks(in *ecs.DescribeTasksInput) (*ecs.DescribeTasksOutput, error) {
	return &m.Resp, nil
}

func (m mockedListClusters) ListClusters(in *ecs.ListClustersInput) (*ecs.ListClustersOutput, error) {
	return &m.Resp, nil
}

func (m mockedListServices) ListServices(in *ecs.ListServicesInput) (*ecs.ListServicesOutput, error) {
	return &m.Resp, nil
}

func (m mockedListTasks) ListTasks(in *ecs.ListTasksInput) (*ecs.ListTasksOutput, error) {
	return &m.Resp, nil
}

func (m mockedRunTask) RunTask(in *ecs.RunTaskInput) (*ecs.RunTaskOutput, error) {
	return &m.Resp, nil
}

func (m mockedDeploy) DescribeTaskDefinition(in *ecs.DescribeTaskDefinitionInput) (*ecs.DescribeTaskDefinitionOutput, error) {
	return &m.DescribeTaskDefResp, nil
}

func (m mockedDeploy) RegisterTaskDefinition(in *ecs.RegisterTaskDefinitionInput) (*ecs.RegisterTaskDefinitionOutput, error) {
	return &m.RegisterTaskDefResp, nil
}

func (m mockedDeploy) UpdateService(in *ecs.UpdateServiceInput) (*ecs.UpdateServiceOutput, error) {
	return &m.UpdateServiceResp, nil
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
		ECS:   mockECSClient{},
		ECR:   mockECRClient{},
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
		ECS:   mockECSClient{},
		ECR:   mockECRClient{},
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
		ECS:   mockECSClient{},
		ECR:   mockECRClient{},
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
		Resp     ecs.ListClustersOutput
		Expected []string
	}{
		{
			Resp: ecs.ListClustersOutput{
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
			ECS:   mockedListClusters{Resp: c.Resp},
			ECR:   mockECRClient{},
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

func TestUFOServices(t *testing.T) {
	var service1, service2, nextToken string = "service1", "service2", ""
	cases := []struct {
		Resp     ecs.ListServicesOutput
		Expected []string
	}{
		{
			Resp: ecs.ListServicesOutput{
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
			ECS:   mockedListServices{Resp: c.Resp},
			ECR:   mockECRClient{},
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

func TestUFORunningTasks(t *testing.T) {
	var task1, task2, nextToken string = "task1", "task2", ""
	cases := []struct {
		Resp     ecs.ListTasksOutput
		Expected []string
	}{
		{
			Resp: ecs.ListTasksOutput{
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
			ECS:   mockedListTasks{Resp: c.Resp},
			ECR:   mockECRClient{},
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

func TestUFOGetCluster(t *testing.T) {
	cases := []struct {
		Resp     ecs.DescribeClustersOutput
		Expected *ecs.Cluster
	}{
		{
			Resp: ecs.DescribeClustersOutput{
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

	//@TODO Handle Error
	for i, c := range cases {
		ufo := UFO{
			l:     &logger{},
			State: &State{},
			ECS:   mockedDescribeClusters{Resp: c.Resp},
			ECR:   mockECRClient{},
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

func TestUFOGetService(t *testing.T) {
	name := "test-service"
	arn := "test-servicearn"
	cases := []struct {
		Resp     ecs.DescribeServicesOutput
		Expected *ecs.Service
	}{
		{
			Resp: ecs.DescribeServicesOutput{
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
			ECS:   mockedDescribeServices{Resp: c.Resp},
			ECR:   mockECRClient{},
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

func TestUFOGetTaskDefinition(t *testing.T) {
	fam := "test-family"
	subcommand := "echo"
	command := []*string{&subcommand, &subcommand}
	arn := "test-taskdefinitionarn"
	image := "111222333444.dkr.ecr.us-west-1.amazonaws.com/image:ea13366"
	cases := []struct {
		Resp     ecs.DescribeTaskDefinitionOutput
		Expected *ecs.TaskDefinition
	}{
		{
			Resp: ecs.DescribeTaskDefinitionOutput{
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
			ECS:   mockedDescribeTaskDefinition{Resp: c.Resp},
			ECR:   mockECRClient{},
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

func TestUFOGetTasks(t *testing.T) {
	var lastStatus, taskDefArn string = "PENDING", "111222333444.dkr.ecr.us-west-1.amazonaws.com/task"
	cases := []struct {
		Resp     ecs.DescribeTasksOutput
		Expected *ecs.DescribeTasksOutput
	}{
		{
			Resp: ecs.DescribeTasksOutput{
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
			ECS:   mockedDescribeTasks{Resp: c.Resp},
			ECR:   mockECRClient{},
		}

		tasks, err := ufo.GetTasks(&ecs.Cluster{}, []*string{})

		if err != nil {
			t.Fatalf("%d, unexpected error", err)
		}

		if a, e := len(tasks.Tasks), len(c.Expected.Tasks); a != e {
			t.Fatalf("%d, expected %d tasks, got %d", i, e, a)
		}

		for j, task := range tasks.Tasks {
			if a, e := task, c.Expected.Tasks[j]; *a.LastStatus != *e.LastStatus {
				t.Errorf("%d, expected %v LastStatus, got %v", i, e, a)
			}

			if a, e := task, c.Expected.Tasks[j]; *a.TaskDefinitionArn != *e.TaskDefinitionArn {
				t.Errorf("%d, expected %v TaskDefinitionARN, got %v", i, e, a)
			}
		}
	}
}

func TestUFOGetImages(t *testing.T) {
	var tag1, tag2 string = "tag1", "tag2"
	cases := []struct {
		Resp     ecr.DescribeImagesOutput
		Expected []*ecr.ImageDetail
	}{
		{
			Resp: ecr.DescribeImagesOutput{
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
			Resp:     ecr.DescribeImagesOutput{},
			Expected: []*ecr.ImageDetail{},
		},
	}

	for i, c := range cases {
		ufo := UFO{
			l:     &logger{},
			State: &State{},
			ECS:   mockECSClient{},
			ECR:   mockedDescribeImages{Resp: c.Resp},
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

func TestUFOGetLastDeployedCommit(t *testing.T) {
	fam := "test-family"
	subcommand := "echo"
	command := []*string{&subcommand, &subcommand}
	arn := "test-taskdefinitionarn"
	image := "111222333444.dkr.ecr.us-west-1.amazonaws.com/image:ea13366"
	cases := []struct {
		Resp     ecs.DescribeTaskDefinitionOutput
		Expected string
	}{
		{
			Resp: ecs.DescribeTaskDefinitionOutput{
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
			ECS:   mockedDescribeTaskDefinition{Resp: c.Resp},
			ECR:   mockECRClient{},
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

func TestUFORunTask(t *testing.T) {
	cases := []struct {
		Resp     ecs.RunTaskOutput
		Expected *ecs.RunTaskOutput
	}{
		{
			Resp: ecs.RunTaskOutput{
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
			ECR:   mockECRClient{},
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

func TestUFODeploy(t *testing.T) {
	emptyValue := ""
	fam := "family1"
	command := "echo"
	commands := []*string{&command, &command}
	image := "111222333444.dkr.ecr.us-west-1.amazonaws.com/image:cbd0d9c"
	commit := "cbd0d9c"
	newImage := fmt.Sprintf("111222333444.dkr.ecr.us-west-1.amazonaws.com/image:%s", commit)

	cases := []struct {
		DescribeTaskDefResp ecs.DescribeTaskDefinitionOutput
		RegisterTaskDefResp ecs.RegisterTaskDefinitionOutput
		UpdateServiceResp   ecs.UpdateServiceOutput
		Expected            *ecs.TaskDefinition
	}{
		{
			DescribeTaskDefResp: ecs.DescribeTaskDefinitionOutput{
				TaskDefinition: &ecs.TaskDefinition{
					ContainerDefinitions: []*ecs.ContainerDefinition{&ecs.ContainerDefinition{
						Command: commands,
						Image:   &image,
					}},
					Family:            &fam,
					TaskDefinitionArn: aws.String("task-definitionarn"),
				},
			},
			RegisterTaskDefResp: ecs.RegisterTaskDefinitionOutput{
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
			UpdateServiceResp: ecs.UpdateServiceOutput{
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
			ECR: mockECRClient{},
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
