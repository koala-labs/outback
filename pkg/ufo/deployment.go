package ufo

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/service/ecs"
	"gitlab.fuzzhq.com/Web-Ops/ufo/pkg/docker"
)

type Deployment struct {
	DeployDetails []*DeployDetail
	BuildDetail   BuildDetail
	Err           error
}

type DeployDetail struct {
	Cluster        *ecs.Cluster
	Service        *ecs.Service
	TaskDefinition *ecs.TaskDefinition
	Done           bool
}

type BuildDetail struct {
	Repo       string
	CommitHash string
	Dockerfile string
}

func (d *DeployDetail) SetCluster(cluster *ecs.Cluster) {
	d.Cluster = cluster
}

func (d *DeployDetail) SetService(service *ecs.Service) {
	d.Service = service
}

func (d *DeployDetail) SetTaskDefinition(taskDef *ecs.TaskDefinition) {
	d.TaskDefinition = taskDef
}

func (d *DeployDetail) SetDone(done bool) {
	d.Done = done
}

func (d *Deployment) SetRepo(repo string) {
	d.BuildDetail.Repo = repo
}

func (d *Deployment) SetCommitHash(commit string) {
	d.BuildDetail.CommitHash = commit
}

func (d *Deployment) SetDockerfile(dockerfile string) {
	d.BuildDetail.Dockerfile = dockerfile
}

func (d *Deployment) TaskDefinitions() string {
	var out strings.Builder
	for _, detail := range d.DeployDetails {
		out.WriteString(fmt.Sprintf("%s ", *detail.TaskDefinition.TaskDefinitionArn))
	}
	return out.String()
}

func (d *Deployment) Services() string {
	var out strings.Builder
	for _, detail := range d.DeployDetails {
		out.WriteString(fmt.Sprintf("%s ", *detail.Service.ServiceName))
	}
	return out.String()
}

func (d *DeployDetail) TaskDefinitionFamily() string {
	r := regexp.MustCompile(`([^\/]+)$`)
	return r.FindString(*d.TaskDefinition.TaskDefinitionArn)
}

func (u *UFO) NewDeployDetail() *DeployDetail {
	return &DeployDetail{
		Done: false,
	}
}

func (u *UFO) AwaitServicesRunning(deployment *Deployment) chan *DeployDetail {
	waitTime := time.Second * 2
	doneCh := make(chan *DeployDetail)
	for _, detail := range deployment.DeployDetails {
		go func(detail *DeployDetail) {
			for !detail.Done {
				done := u.IsServiceRunning(detail)
				detail.SetDone(done)
				time.Sleep(waitTime)
			}
			doneCh <- detail
		}(detail)
	}

	return doneCh
}

func (u *UFO) DeployAll(deploy *Deployment) <-chan error {
	var wg sync.WaitGroup
	errCh := make(chan error)

	wg.Add(len(deploy.DeployDetails))
	for _, detail := range deploy.DeployDetails {
		go func(detail *DeployDetail) {
			taskDef, err := u.UpdateServiceWithNewTaskDefinition(detail.Cluster, detail.Service, deploy.BuildDetail.CommitHash)

			if err != nil {
				errCh <- err
				wg.Done()
			}

			// Set TaskDefinition to the updated services new TaskDefinition
			detail.SetTaskDefinition(taskDef)
			wg.Done()
		}(detail)
	}

	wg.Wait()
	close(errCh)
	return errCh
}

func (u *UFO) LoginBuildPushImage(info BuildDetail) error {
	var err error

	err = u.ECRLogin()

	if err != nil {
		return err
	}

	err = docker.ImageBuild(info.Repo, info.CommitHash, info.Dockerfile)

	if err != nil {
		return err
	}

	err = docker.ImagePush(info.Repo, info.CommitHash)

	if err != nil {
		return err
	}

	return nil
}
