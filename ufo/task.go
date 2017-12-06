package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ecs"
	"gitlab.fuzzhq.com/Web-Ops/ufo/pkg/ufo"
)

type TaskOptions struct {
	Command        string
	OverrideBranch string
}

type TaskState struct {
	cluster        *ecs.Cluster
	taskDefinition *ecs.TaskDefinition
	command        string
}

type TaskCmd struct {
	branch  string
	Options TaskOptions
	c       *Config
	Env     *Environment
	s       *TaskState
	UFO     *ufo.UFO
}

func RunTask(c *Config, options TaskOptions) error {
	var err error

	r := &TaskCmd{
		c:       c,
		Options: options,
		s:       &TaskState{},
		branch:  options.OverrideBranch,
	}

	e, err := c.GetEnvironmentByBranch(r.branch)

	if err != nil {
		return err
	}

	r.Env = e

	r.initUFO()

	return r.run()
}

func (r *TaskCmd) run() error {
	fmt.Printf("Preparing to run desired task on cluster %s.\n", r.Env.Cluster)

	r.s.cluster = r.loadCluster(r.Env.Cluster)

	r.s.taskDefinition = r.loadTaskDefinition(r.s.cluster, r.Env.Service)

	r.s.command = r.Options.Command

	t, err := r.UFO.RunTask(r.s.cluster, r.s.taskDefinition, r.s.command)

	if err != nil {
		return err
	}

	fmt.Printf("Running task... %s", t)

	return nil
}

func (r *TaskCmd) initUFO() {
	c := ufo.UFOConfig{
		Profile: &r.c.Profile,
		Region:  &r.Env.Region,
	}

	r.UFO = ufo.Fly(c, &logger{})
}

func (r *TaskCmd) loadCluster(name string) *ecs.Cluster {
	awsCluster, err := r.UFO.GetCluster(name)

	HandleError(err)

	return awsCluster
}

func (r *TaskCmd) loadTaskDefinition(c *ecs.Cluster, name string) *ecs.TaskDefinition {
	awsService, err := r.UFO.GetService(c, name)

	HandleError(err)

	t, err := r.UFO.GetTaskDefinition(c, awsService)

	HandleError(err)

	return t
}
