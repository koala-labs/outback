package main

import (
	"os/exec"
	"github.com/aws/aws-sdk-go/service/ecs"
	"gitlab.fuzzhq.com/Web-Ops/ufo/pkg/ufo"
	"time"
	"fmt"
	"strings"
)

type logger struct {}

func (l *logger) Printf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

type DeployOptions struct {
	Verbose        bool
	OverrideBranch string
}

type DeployState struct {
	cluster *ecs.Cluster
	service *ecs.Service
	oldT    *ecs.TaskDefinition
	newT    *ecs.TaskDefinition
}

type DeployCmd struct {
	Options DeployOptions
	s       *DeployState
	c       *Config
	UFO     *ufo.UFO
	Env     *Environment
	branch  string
	head    string
}

func getCurrentHead() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")

	r, err := cmd.Output()

	if err != nil {
		return "", ErrGitError
	}

	return strings.Trim(string(r), "\n"), nil
}

func getCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")

	r, err := cmd.Output()

	if err != nil {
		return "", ErrGitError
	}

	return strings.Trim(string(r), "\n"), nil
}

func RunDeployCmd(c *Config, options DeployOptions) error {
	var err error

	d := &DeployCmd{
		Options: options,
		branch:  options.OverrideBranch,
		c:       c,
		s:       &DeployState{},
	}

	d.head, err = getCurrentHead()

	if err != nil {
		return err
	}

	if d.Options.OverrideBranch == EMPTY_VALUE {
		d.branch, err = getCurrentBranch()

		if err != nil {
			return err
		}
	}

	e, err := c.GetEnvironmentByBranch(d.branch)

	if err != nil {
		return err
	}

	d.Env = e

	d.initUFO()

	return d.deploy()
}

func (d *DeployCmd) initUFO() {
	c := ufo.UFOConfig{
		Profile: &d.c.Profile,
		Region:  &d.Env.Region,
	}

	d.UFO = ufo.Fly(c, &logger{})
}

func (d *DeployCmd) deploy() error {
	fmt.Printf("Preparing to deploy branch %s to service %s on cluster %s.\n", d.Env.Branch, d.Env.Service, d.Env.Cluster)

	// Push an image to docker repo
	fmt.Println("Building docker image.")
	err := d.buildImage()

	if err != nil {
		return err
	}

	fmt.Println("Pushing docker image.")
	err = d.pushImage()

	if err != nil {
		return err
	}

	d.s.cluster = d.loadCluster(d.Env.Cluster)
	d.s.service, d.s.oldT = d.loadService(d.s.cluster, d.Env.Service)

	fmt.Printf("Beginning deployment to service %s.\n", d.Env.Service)
	t, err := d.UFO.Deploy(d.s.cluster, d.s.service, d.head)

	if err != nil {
		return err
	}

	d.s.newT = t

	err = d.awaitCompletion()

	if err != nil {
		return ErrDeployTimeout
	}

	fmt.Printf("Successfully deployed. Your new task definition is %s:%d.\n", *d.s.newT.Family, *d.s.newT.Revision)

	return nil
}

func (d *DeployCmd) buildImage() error {
	c := fmt.Sprintf("%s:%s", d.c.ImageRepositoryUrl, d.head)
	cmd := exec.Command("docker", "build", "-f", d.Env.Dockerfile, "--tag", c, ".")

	out, err := cmd.Output()

	if err != nil {
		fmt.Printf("%v", err)
		fmt.Printf("%v", string(out))
		return ErrDockerBuild
	}

	if d.Options.Verbose {
		fmt.Printf("%s", string(out))
	}

	return nil
}

func (d *DeployCmd) pushImage() error {
	c := fmt.Sprintf("%s:%s", d.c.ImageRepositoryUrl, d.head)
	cmd := exec.Command("docker", "push", c)

	out, err := cmd.Output()

	if err != nil {
		fmt.Printf("%v", err)
		fmt.Printf("%v", string(out))
		return ErrDockerPush
	}

	if d.Options.Verbose {
		fmt.Printf("%s", string(out))
	}

	return nil
}

func (d *DeployCmd) awaitCompletion() error {
	attempts := 0
	waitTime := 2 * time.Second

	for ! d.IsDeployed(d.s.cluster, d.s.service, d.s.newT) {
		if attempts > 60 {
			return ErrDeployTimeout
		}

		attempts++

		fmt.Println("Waiting.")
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
