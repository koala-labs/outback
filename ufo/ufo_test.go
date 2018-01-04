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

type mockECSClient struct {
	ecsiface.ECSAPI
}

type mockECRClient struct {
	ecriface.ECRAPI
}

type logger struct{}

func (l *logger) Printf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

func (m mockECRClient) DescribeImages(s *ecr.DescribeImagesInput) (*ecr.DescribeImagesOutput, error) {
	var tag1, tag2 string = "tag1", "tag2"

	return &ecr.DescribeImagesOutput{
		ImageDetails: []*ecr.ImageDetail{
			&ecr.ImageDetail{
				ImageTags: []*string{&tag1},
			},
			&ecr.ImageDetail{
				ImageTags: []*string{&tag2},
			}},
	}, nil
}

func (m mockECSClient) DescribeClusters(s *ecs.DescribeClustersInput) (*ecs.DescribeClustersOutput, error) {
	name := "cluster1"
	arn := "clusterarn"

	return &ecs.DescribeClustersOutput{
		Clusters: []*ecs.Cluster{&ecs.Cluster{
			ClusterName: &name,
			ClusterArn:  &arn,
		}},
		Failures: []*ecs.Failure{},
	}, nil
}

func (m mockECSClient) DescribeServices(s *ecs.DescribeServicesInput) (*ecs.DescribeServicesOutput, error) {
	name := "service1"
	arn := "servicearn"

	return &ecs.DescribeServicesOutput{
		Services: []*ecs.Service{&ecs.Service{
			ServiceName: &name,
			ServiceArn:  &arn,
		}},
		Failures: []*ecs.Failure{},
	}, nil
}

func (m mockECSClient) DescribeTaskDefinition(s *ecs.DescribeTaskDefinitionInput) (*ecs.DescribeTaskDefinitionOutput, error) {
	fam := "family1"
	command := "echo"
	commands := []*string{&command, &command}
	arn := "111222333444.dkr.ecr.us-west-1.amazonaws.com/task"
	image := "111222333444.dkr.ecr.us-west-1.amazonaws.com/image:ea13366"

	return &ecs.DescribeTaskDefinitionOutput{
		TaskDefinition: &ecs.TaskDefinition{
			ContainerDefinitions: []*ecs.ContainerDefinition{&ecs.ContainerDefinition{
				Command: commands,
				Image:   &image,
			}},
			Family:            &fam,
			TaskDefinitionArn: &arn,
		},
	}, nil
}

func (m mockECSClient) DescribeTasks(s *ecs.DescribeTasksInput) (*ecs.DescribeTasksOutput, error) {
	lastStatus := "PENDING"
	taskDefArn := "111222333444.dkr.ecr.us-west-1.amazonaws.com/task"

	return &ecs.DescribeTasksOutput{
		Tasks: []*ecs.Task{&ecs.Task{
			LastStatus:        &lastStatus,
			TaskDefinitionArn: &taskDefArn,
		}},
	}, nil
}

func (m mockECSClient) ListClusters(s *ecs.ListClustersInput) (*ecs.ListClustersOutput, error) {
	var cluster1, cluster2, nextToken string = "cluster1", "cluster2", ""

	return &ecs.ListClustersOutput{
		ClusterArns: []*string{&cluster1, &cluster2},
		NextToken:   &nextToken,
	}, nil
}

func (m mockECSClient) ListServices(s *ecs.ListServicesInput) (*ecs.ListServicesOutput, error) {
	var service1, service2, nextToken string = "service1", "service2", ""

	return &ecs.ListServicesOutput{
		ServiceArns: []*string{&service1, &service2},
		NextToken:   &nextToken,
	}, nil
}

func (m mockECSClient) ListTasks(s *ecs.ListTasksInput) (*ecs.ListTasksOutput, error) {
	var task1, task2, nextToken string = "task1", "task2", ""

	return &ecs.ListTasksOutput{
		TaskArns:  []*string{&task1, &task2},
		NextToken: &nextToken,
	}, nil
}

func (m mockECSClient) UpdateService(s *ecs.UpdateServiceInput) (*ecs.UpdateServiceOutput, error) {
	return &ecs.UpdateServiceOutput{
		Service: &ecs.Service{},
	}, nil
}

func (m mockECSClient) RegisterTaskDefinition(s *ecs.RegisterTaskDefinitionInput) (*ecs.RegisterTaskDefinitionOutput, error) {
	emptyValue := ""
	fam := "family1"
	command := "echo"
	commands := []*string{&command, &command}
	image := "111222333444.dkr.ecr.us-west-1.amazonaws.com/image:cbd0d9c"

	return &ecs.RegisterTaskDefinitionOutput{
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
	}, nil
}

func (m mockECSClient) RunTask(s *ecs.RunTaskInput) (*ecs.RunTaskOutput, error) {
	name := "task"
	command := []string{"echo", "this"}
	commands := []*string{&command[0], &command[1]}
	lastStatus := "PENDING"
	taskDefArn := "111222333444.dkr.ecr.us-west-1.amazonaws.com/task"
	clusterArn := "clusterarn"

	return &ecs.RunTaskOutput{
		Tasks: []*ecs.Task{&ecs.Task{
			ClusterArn:        &clusterArn,
			TaskDefinitionArn: &taskDefArn,
			LastStatus:        &lastStatus,
			Overrides: &ecs.TaskOverride{
				ContainerOverrides: []*ecs.ContainerOverride{&ecs.ContainerOverride{
					Command: commands,
					Name:    &name,
				}},
			},
		}},
	}, nil
}

func mockUFO() UFO {
	u := UFO{
		l:     &logger{},
		State: &State{},
		ECS:   mockECSClient{},
		ECR:   mockECRClient{},
	}

	return u
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

	ufo := mockUFO()

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

	ufo := mockUFO()

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

	ufo := mockUFO()

	for i, c := range cases {
		ufo.UseTaskDefinition(c.Expected)
		if a, e := *ufo.State.TaskDefinition.Family, *c.Expected.Family; a != e {
			t.Errorf("%d, expected %v task definition, got %v", i, e, a)
		}
		if a, e := *ufo.State.TaskDefinition.Family, *c.Expected.Family; a != e {
			t.Errorf("%d, expected %v family, got %v", i, e, a)
		}
		actualContainerDefinition := ufo.State.TaskDefinition.ContainerDefinitions[0]
		expectedContainerDefinition := c.Expected.ContainerDefinitions[0]
		if a, e := *actualContainerDefinition.Image, *expectedContainerDefinition.Image; a != e {
			t.Errorf("%d, expected %v image, got %v", i, e, a)
		}
	}
}

func TestUFOClusters(t *testing.T) {
	cases := []struct {
		Expected []string
	}{
		{
			Expected: []string{
				"cluster1",
				"cluster2",
			},
		},
	}

	ufo := mockUFO()

	clusters, err := ufo.Clusters()

	for i, c := range cases {
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
	cases := []struct {
		Expected []string
	}{
		{
			Expected: []string{
				"service1",
				"service2",
			},
		},
	}

	ufo := mockUFO()

	clusters, err := ufo.Clusters()
	cluster, err := ufo.GetCluster(clusters[0])
	services, err := ufo.Services(cluster)

	for i, c := range cases {
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
	cases := []struct {
		Expected []string
	}{
		{
			Expected: []string{
				"task1",
				"task2",
			},
		},
	}

	ufo := mockUFO()

	clusters, err := ufo.Clusters()
	cluster, err := ufo.GetCluster(clusters[0])
	services, err := ufo.Services(cluster)
	service, err := ufo.GetService(cluster, services[0])
	tasks, err := ufo.RunningTasks(cluster, service)

	for i, c := range cases {
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

	ufo := mockUFO()

	clusters, err := ufo.Clusters()
	cluster, err := ufo.GetCluster(clusters[0])

	for i, c := range cases {
		if err != nil {
			t.Fatalf("%d, unexpected error", err)
		}
		if a, e := *cluster.ClusterName, *c.Expected.ClusterName; a != e {
			t.Errorf("%d, expected %v cluster, got %v", i, e, a)
		}
	}
}

func TestUFOGetService(t *testing.T) {
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

	ufo := mockUFO()

	clusters, err := ufo.Clusters()
	cluster, err := ufo.GetCluster(clusters[0])
	services, err := ufo.Services(cluster)
	service, err := ufo.GetService(cluster, services[0])

	for i, c := range cases {
		if err != nil {
			t.Fatalf("%d, unexpected error", err)
		}
		if a, e := *service.ServiceName, *c.Expected.ServiceName; a != e {
			t.Errorf("%d, expected %v service, got %v", i, e, a)
		}
	}
}

func TestUFOGetTaskDefinition(t *testing.T) {
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

	ufo := mockUFO()

	clusters, err := ufo.Clusters()
	cluster, err := ufo.GetCluster(clusters[0])
	services, err := ufo.Services(cluster)
	service, err := ufo.GetService(cluster, services[0])
	taskDef, err := ufo.GetTaskDefinition(cluster, service)

	for i, c := range cases {
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
		Expected *ecs.DescribeTasksOutput
	}{
		{
			Expected: &ecs.DescribeTasksOutput{
				Tasks: []*ecs.Task{
					&ecs.Task{
						LastStatus:        &lastStatus,
						TaskDefinitionArn: &taskDefArn,
					}},
			},
		},
	}

	ufo := mockUFO()

	clusters, err := ufo.Clusters()
	cluster, err := ufo.GetCluster(clusters[0])
	services, err := ufo.Services(cluster)
	service, err := ufo.GetService(cluster, services[0])
	runningTasks, err := ufo.RunningTasks(cluster, service)
	tasks, err := ufo.GetTasks(cluster, runningTasks)

	for i, c := range cases {
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
		Expected []*ecr.ImageDetail
	}{
		{
			Expected: []*ecr.ImageDetail{
				&ecr.ImageDetail{
					ImageTags: []*string{&tag1},
				},
				&ecr.ImageDetail{
					ImageTags: []*string{&tag2},
				}},
		},
	}

	ufo := mockUFO()

	clusters, err := ufo.Clusters()
	cluster, err := ufo.GetCluster(clusters[0])
	services, err := ufo.Services(cluster)
	service, err := ufo.GetService(cluster, services[0])
	taskDef, err := ufo.GetTaskDefinition(cluster, service)
	images, err := ufo.GetImages(taskDef)

	for i, c := range cases {
		if err != nil {
			t.Fatalf("%d, unexpected error", err)
		}
		if a, e := len(images), len(c.Expected); a != e {
			t.Fatalf("%d, expected %d tasks, got %d", i, e, a)
		}
		for j, image := range images {
			actualImageTag := strings.Join(aws.StringValueSlice(image.ImageTags), " ")
			expectedImageTag := strings.Join(aws.StringValueSlice(c.Expected[j].ImageTags), " ")
			if a, e := actualImageTag, expectedImageTag; a != e {
				t.Errorf("%d, expected %v TaskDefinitionARN, got %v", i, e, a)
			}
		}
	}
}

func TestUFOGetLastDeployedCommit(t *testing.T) {
	cases := []struct {
		Expected string
	}{
		{
			Expected: "ea13366",
		},
	}

	ufo := mockUFO()

	clusters, err := ufo.Clusters()
	cluster, err := ufo.GetCluster(clusters[0])
	services, err := ufo.Services(cluster)
	service, err := ufo.GetService(cluster, services[0])
	taskDef, err := ufo.GetTaskDefinition(cluster, service)
	commit, err := ufo.GetLastDeployedCommit(*taskDef.ContainerDefinitions[0].Image)

	for i, c := range cases {
		if err != nil {
			t.Fatalf("%d, unexpected error", err)
		}
		if a, e := commit, c.Expected; a != e {
			t.Errorf("%d, expected %v commit, got %v", i, e, a)
		}
	}
}

func TestUFODeploy(t *testing.T) {
	fam := "family1"
	command := "echo"
	commands := []*string{&command, &command}
	commit := "cbd0d9c"
	newImage := fmt.Sprintf("111222333444.dkr.ecr.us-west-1.amazonaws.com/image:%s", commit)
	cases := []struct {
		Expected *ecs.TaskDefinition
	}{
		{
			Expected: &ecs.TaskDefinition{
				ContainerDefinitions: []*ecs.ContainerDefinition{&ecs.ContainerDefinition{
					Command: commands,
					Image:   &newImage,
				}},
				Family: &fam,
			},
		},
	}

	ufo := mockUFO()

	clusters, err := ufo.Clusters()
	cluster, err := ufo.GetCluster(clusters[0])
	services, err := ufo.Services(cluster)
	service, err := ufo.GetService(cluster, services[0])
	newTaskDef, err := ufo.Deploy(cluster, service, commit)

	for i, c := range cases {
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

func TestUFORunTask(t *testing.T) {
	name := "task"
	cmd := "echo this"
	command := []string{"echo", "this"}
	commands := []*string{&command[0], &command[1]}

	ufo := mockUFO()

	clusters, err := ufo.Clusters()
	cluster, err := ufo.GetCluster(clusters[0])
	services, err := ufo.Services(cluster)
	service, err := ufo.GetService(cluster, services[0])
	taskDef, err := ufo.GetTaskDefinition(cluster, service)

	cases := []struct {
		Expected *ecs.RunTaskOutput
	}{
		{
			Expected: &ecs.RunTaskOutput{
				Tasks: []*ecs.Task{&ecs.Task{
					ClusterArn:        cluster.ClusterArn,
					TaskDefinitionArn: taskDef.TaskDefinitionArn,
					Overrides: &ecs.TaskOverride{
						ContainerOverrides: []*ecs.ContainerOverride{&ecs.ContainerOverride{
							Command: commands,
							Name:    &name,
						}},
					},
				}},
			},
		},
	}

	ranTasks, err := ufo.RunTask(cluster, taskDef, cmd)

	for i, c := range cases {
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
