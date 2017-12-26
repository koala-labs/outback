package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ecs"
)

type TaskOptions struct {
	Command        string
	OverrideBranch string
	CommandName    string
}

type TaskState struct {
	cluster        *ecs.Cluster
	taskDefinition *ecs.TaskDefinition
	command        string
}

type TaskCmd struct {
	branch  string
	Options TaskOptions
	cmd     *Cmd
	c       *Config
	Env     *Environment
	s       *TaskState
}

func RunTask(c *Config, options TaskOptions) error {
	var err error

	if options.CommandName != EmptyValue {
		runTaskConfig, err := c.GetCommandForName(options.CommandName)

		if err != nil {
			return err
		}

		options.Command = runTaskConfig.Command
	}

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

	r.cmd = r.cmd.initUFO(r.c.Profile, r.Env.Region)

	return r.run()
}

func (r *TaskCmd) run() error {
	fmt.Printf("Preparing to run desired task on cluster %s.\n", r.Env.Cluster)

	r.s.cluster = r.cmd.loadCluster(r.Env.Cluster)

	// @todo we should be able to override the service we copy the task def from?
	// Would services in the same env have the same setup?
	r.s.taskDefinition = r.cmd.loadTaskDefinition(r.s.cluster, r.Env.Services[0])

	r.s.command = r.Options.Command

	_, err := r.cmd.UFO.RunTask(r.s.cluster, r.s.taskDefinition, r.s.command)

	if err != nil {
		return err
	}

	fmt.Printf("Running task...")

	return nil
}
