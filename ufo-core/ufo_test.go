package ufo

import (
	"fmt"
	"strings"
	"testing"

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

func (m mockECSClient) DescribeClusters(s *ecs.DescribeClustersInput) (*ecs.DescribeClustersOutput, error) {
	name := "cluster1"
	return &ecs.DescribeClustersOutput{
		Clusters: []*ecs.Cluster{&ecs.Cluster{
			ClusterName: &name,
		}},
		Failures: []*ecs.Failure{},
	}, nil
}

func (m mockECSClient) DescribeServices(s *ecs.DescribeServicesInput) (*ecs.DescribeServicesOutput, error) {
	name := "service1"
	return &ecs.DescribeServicesOutput{
		Services: []*ecs.Service{&ecs.Service{
			ServiceName: &name,
		}},
		Failures: []*ecs.Failure{},
	}, nil
}

func (m mockECSClient) DescribeTaskDefinition(s *ecs.DescribeTaskDefinitionInput) (*ecs.DescribeTaskDefinitionOutput, error) {
	fam := "family1"
	command := "echo"
	commands := []*string{&command, &command}
	image := "111222333444.dkr.ecr.us-west-1.amazonaws.com/test"
	return &ecs.DescribeTaskDefinitionOutput{
		TaskDefinition: &ecs.TaskDefinition{
			ContainerDefinitions: []*ecs.ContainerDefinition{&ecs.ContainerDefinition{
				Command: commands,
				Image:   &image,
			}},
			Family: &fam,
		},
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

func mockUFO() UFO {
	u := UFO{
		l:     &logger{},
		State: &State{},
		ECS:   mockECSClient{},
		ECR:   mockECRClient{},
	}

	return u
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

func TestUfoGetTaskDefinition(t *testing.T) {
	fam := "family1"
	command := "echo"
	commands := []*string{&command, &command}
	image := "111222333444.dkr.ecr.us-west-1.amazonaws.com/test"
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
		actualContainerCommand := strings.Join(aws.StringValueSlice(taskDef.ContainerDefinitions[0].Command), " ")
		expectedContainerCommand := strings.Join(aws.StringValueSlice(c.Expected.ContainerDefinitions[0].Command), " ")
		if a, e := actualContainerCommand, expectedContainerCommand; a != e {
			t.Errorf("%d, expected %v command, got %v", i, e, a)
		}
	}
}
