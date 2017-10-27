package main

import (
	"os/exec"
	"github.com/aws/aws-sdk-go/service/ecs"
	"gitlab.fuzzhq.com/Web-Ops/ufo/pkg/ufo"
	"time"
	"log"
	"fmt"
	"strings"
)

type DeployState struct {
	cluster *ecs.Cluster
	service *ecs.Service
	oldT    *ecs.TaskDefinition
	newT    *ecs.TaskDefinition
}

type DeployCmd struct {
	verbose bool
	s       *DeployState
	c       *Config
	UFO     *ufo.UFO
	Env     *Environment
	branch  string
	head    string
}

func getCurrentHead() string {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")

	r, err := cmd.Output()

	if err != nil {
		HandleError(ErrGitError)
	}

	return strings.Trim(string(r), "\n")
}

func getCurrentBranch() string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")

	r, err := cmd.Output()

	if err != nil {
		HandleError(ErrGitError)
	}

	return strings.Trim(string(r), "\n")
}

func RunDeployCmd(c *Config, verbose bool) {
	d := &DeployCmd{
		verbose: verbose,
		branch: getCurrentBranch(),
		head: getCurrentHead(),
		c:    c,
		s:    &DeployState{},
	}

	e, err := c.GetEnvironmentByBranch(d.branch)
	HandleError(err)

	d.Env = e

	d.initUFO()
	d.deploy()
}

func (d *DeployCmd) initUFO() {
	c := ufo.UFOConfig{
		Profile: &d.c.Profile,
		Region:  &d.Env.Region,
	}

	d.UFO = ufo.Fly(c)
}

func (d *DeployCmd) deploy() {
	log.Printf("Preparing to deploy branch %s to service %s on cluster %s.\n", d.Env.Branch, d.Env.Service, d.Env.Cluster)

	// Push an image to docker repo
	log.Println("Building docker image.")
	d.buildImage()

	log.Println("Pushing docker image.")
	d.pushImage()

	d.s.cluster = d.loadCluster(d.Env.Cluster)
	d.s.service, d.s.oldT = d.loadService(d.s.cluster, d.Env.Service)

	log.Printf("Beginning deployment to service %s.\n", d.Env.Service)
	t, err := d.UFO.Deploy(d.s.cluster, d.s.service, d.head)
	HandleError(err)

	d.s.newT = t

	err = d.awaitCompletion()

	if err != nil {
		HandleError(ErrDeployTimeout)
	}

	log.Printf("Successfully deployed. Your new task definition is %s:%d.\n", *d.s.newT.Family, *d.s.newT.Revision)
}

func (d *DeployCmd) buildImage() {
	c := fmt.Sprintf("%s:%s", d.c.ImageRepositoryUrl, d.head)
	cmd := exec.Command("docker", "build", "-f", d.Env.Dockerfile, "--tag", c, ".")

	out, err := cmd.Output()

	if err != nil {
		log.Printf("%v", err)
		log.Printf("%v", string(out))
		HandleError(ErrDockerBuild)
	}

	if d.verbose {
		log.Printf("%s", string(out))
	}
}

func (d *DeployCmd) pushImage() {
	c := fmt.Sprintf("%s:%s", d.c.ImageRepositoryUrl, d.head)
	cmd := exec.Command("docker", "push", c)

	out, err := cmd.Output()

	if err != nil {
		log.Printf("%v", err)
		log.Printf("%v", string(out))
		HandleError(ErrDockerPush)
	}

	if d.verbose {
		log.Printf("%s", string(out))
	}
}

func (d *DeployCmd) awaitCompletion() error {
	attempts := 0
	waitTime := 2 * time.Second

	for ! d.IsDeployed(d.s.cluster, d.s.service, d.s.newT) {
		if attempts > 60 {
			return ErrDeployTimeout
		}

		attempts++

		log.Print("Waiting.")
		time.Sleep(waitTime)
	}

	return nil
}

func (d *DeployCmd) loadCluster(name string) *ecs.Cluster {
	awsCluster, err := d.UFO.GetCluster(name)

	HandleError(err)

	return awsCluster
}

func (d *DeployCmd) loadService(c *ecs.Cluster, name string) (*ecs.Service, *ecs.TaskDefinition) {
	awsService, err := d.UFO.GetService(c, name)

	HandleError(err)

	t, err := d.UFO.GetTaskDefinition(c, awsService)

	HandleError(err)

	return awsService, t
}

func (d *DeployCmd) IsDeployed(c *ecs.Cluster, s *ecs.Service, t *ecs.TaskDefinition) bool {
	if *s.DesiredCount <= 0 {
		return false
	}

	runningTasks, err := d.UFO.RunningTasks(c, s)

	HandleError(err)

	if len(runningTasks) <= 0 {
		return false
	}

	tasks, err := d.UFO.GetTasks(c, runningTasks)

	HandleError(err)

	for _, task := range tasks.Tasks {
		if *task.TaskDefinitionArn == *t.TaskDefinitionArn && *task.LastStatus == "RUNNING" {
			return true
		}
	}

	return false
}
