package ufo

import (
	"fmt"
	"testing"

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

func (m mockECSClient) DescribeClusters(*ecs.DescribeClustersInput) (*ecs.DescribeClustersOutput, error) {
	// Only need to return mocked response output
	var count int64 = 1
	name := "test"
	status := "ACTIVE"
	arn := "arn:aws:ecs:us-east-1:111222333344:cluster/test"
	return &ecs.DescribeClustersOutput{
		Clusters: []*ecs.Cluster{&ecs.Cluster{
			ActiveServicesCount:               &count,
			ClusterArn:                        &arn,
			ClusterName:                       &name,
			PendingTasksCount:                 &count,
			RegisteredContainerInstancesCount: &count,
			RunningTasksCount:                 &count,
			Statistics:                        []*ecs.KeyValuePair{},
			Status:                            &status,
		},
		},
		Failures: []*ecs.Failure{},
	}, nil
}

func (m mockECSClient) ListClusters(*ecs.ListClustersInput) (*ecs.ListClustersOutput, error) {
	// Only need to return mocked response output
	var test1, test2, nextToken string = "test1", "test2", ""
	return &ecs.ListClustersOutput{
		ClusterArns: []*string{&test1, &test2},
		NextToken:   &nextToken,
	}, nil
}

func (m mockECSClient) ListServices(*ecs.ListServicesInput) (*ecs.ListServicesOutput, error) {
	// Only need to return mocked response output
	var test1, test2, nextToken string = "test1", "test2", ""
	return &ecs.ListServicesOutput{
		ServiceArns: []*string{&test1, &test2},
		NextToken:   &nextToken,
	}, nil
}

func TestUFOClusters(t *testing.T) {
	cases := []struct {
		Expected []string
	}{
		{
			Expected: []string{
				"test1",
				"test2",
			},
		},
	}

	u := UFO{
		l:     &logger{},
		State: &State{},
		ECS:   mockECSClient{},
		ECR:   mockECRClient{},
	}

	clusters, err := u.Clusters()

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
				"test1",
				"test2",
			},
		},
	}

	u := UFO{
		l:     &logger{},
		State: &State{},
		ECS:   mockECSClient{},
		ECR:   mockECRClient{},
	}

	clusters, err := u.Clusters()
	cluster, err := u.GetCluster(clusters[0])
	services, err := u.Services(cluster)

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
