package outback

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/koala-labs/outback/pkg/docker"
)

type Deployment struct {
	DeployDetails []*DeployDetail
	BuildDetail   BuildDetail
	Err           error
}

type DeployDetail struct {
	Cluster                  *ecs.Cluster
	Service                  *ecs.Service
	TaskDefinition           *ecs.TaskDefinition
	TaskDefinitionFamilyName string
	RevisionNumber           int
	Done                     bool
}

type BuildDetail struct {
	Repo            string
	CommitHash      string
	Dockerfile      string
	buildArgs       []string
	configBuildArgs []string
	cacheFrom       []string
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

func (d *DeployDetail) SetTaskDefinitionFamilyName(TaskDefinitionFamilyName string) {
	d.TaskDefinitionFamilyName = TaskDefinitionFamilyName
}

func (d *DeployDetail) SetRevisionNumber(revisionNumber int) {
	d.RevisionNumber = revisionNumber
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

func (d *Deployment) SetBuildArgs(buildArgs []string) {
	d.BuildDetail.buildArgs = buildArgs
}

func (d *Deployment) SetConfigBuildArgs(configBuildArgs []string) {
	d.BuildDetail.configBuildArgs = configBuildArgs
}

func (d *Deployment) SetBuildCacheFrom(cacheFrom []string) {
	d.BuildDetail.cacheFrom = cacheFrom
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

func (u *Outback) NewDeployDetail() *DeployDetail {
	return &DeployDetail{
		Done: false,
	}
}

func (u *Outback) AwaitServicesRunning(deployment *Deployment) chan *DeployDetail {
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

func (u *Outback) RollbackAll(deploy *Deployment, deployDetail *DeployDetail) <-chan error {
	var wg sync.WaitGroup
	errCh := make(chan error)

	wg.Add(len(deploy.DeployDetails))
	for _, detail := range deploy.DeployDetails {
		go func(detail *DeployDetail) {
			taskDefName, err := u.RollbackTaskDefinition(detail.Cluster, detail.Service, detail.TaskDefinition, deployDetail.RevisionNumber)

			if err != nil {
				errCh <- err
				wg.Done()
			}
			detail.SetTaskDefinitionFamilyName(taskDefName)

			wg.Done()
		}(detail)
	}

	wg.Wait()
	close(errCh)
	return errCh
}

func (u *Outback) DeployAll(deploy *Deployment) <-chan error {
	var wg sync.WaitGroup
	errCh := make(chan error)

	wg.Add(len(deploy.DeployDetails))
	for _, detail := range deploy.DeployDetails {
		go func(detail *DeployDetail) {
			taskDef, err := u.UpdateServiceWithNewTaskDefinition(detail.Cluster, detail.Service, deploy.BuildDetail.Repo, deploy.BuildDetail.CommitHash)

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

func (u *Outback) LoginBuildPushImage(info BuildDetail) error {
	var err error

	err = u.ECRLogin()

	if err != nil {
		return err
	}

	err = docker.ImageBuild(info.Repo, info.CommitHash, info.Dockerfile, info.buildArgs, info.configBuildArgs, info.cacheFrom)

	if err != nil {
		return err
	}

	err = docker.ImagePush(info.Repo, info.CommitHash)

	if err != nil {
		return err
	}

	return nil
}

func (u *Outback) LoginPullImage(repo string, tag string) error {
	var err error

	err = u.ECRLogin()

	if err != nil {
		return err
	}

	err = docker.ImagePull(repo, tag)

	if err != nil {
		return err
	}

	return nil
}
