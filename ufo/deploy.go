package main

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/aws/aws-sdk-go/service/ecs"
)

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
	cmd     *Cmd
	Env     *Environment
	branch  string
	head    string
}

func RunDeployCmd(c *Config, options DeployOptions) error {
	var err error

	d := &DeployCmd{
		Options: options,
		branch:  options.OverrideBranch,
		c:       c,
		s:       &DeployState{},
	}

	d.head, err = d.cmd.getCurrentHead()

	if err != nil {
		return err
	}

	if d.Options.OverrideBranch == EMPTY_VALUE {
		d.branch, err = d.cmd.getCurrentBranch()

		if err != nil {
			return err
		}
	}

	e, err := c.GetEnvironmentByBranch(d.branch)

	if err != nil {
		return err
	}

	d.Env = e

	d.cmd = d.cmd.initUFO(d.c.Profile, d.Env.Region)

	return d.deploy()
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

	d.s.cluster = d.cmd.loadCluster(d.Env.Cluster)
	d.s.service, d.s.oldT = d.cmd.loadService(d.s.cluster, d.Env.Service)

	fmt.Printf("Beginning deployment to service %s.\n", d.Env.Service)
	t, err := d.cmd.UFO.Deploy(d.s.cluster, d.s.service, d.head)

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

	for !d.cmd.isDeployed(d.s.cluster, d.s.service, d.s.newT) {
		if attempts > 60 {
			return ErrDeployTimeout
		}

		attempts++

		fmt.Println("Waiting.")
		time.Sleep(waitTime)
	}

	return nil
}
