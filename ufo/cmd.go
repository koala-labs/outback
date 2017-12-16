package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go/service/ecs"
	"gitlab.fuzzhq.com/Web-Ops/ufo/ufo-core"
)

type logger struct{}

func (l *logger) Printf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

type Cmd struct {
	UFO *ufo.UFO
}

func (cmd *Cmd) initUFO(profile string, region string) *Cmd {
	c := ufo.UFOConfig{
		Profile: &profile,
		Region:  &region,
	}

	command := &Cmd{
		UFO: ufo.Fly(c, &logger{}),
	}

	return command
}

func (cmd *Cmd) loadService(c *ecs.Cluster, name string) (*ecs.Service, *ecs.TaskDefinition) {
	awsService, err := cmd.UFO.GetService(c, name)

	HandleError(err)

	t, err := cmd.UFO.GetTaskDefinition(c, awsService)

	HandleError(err)

	return awsService, t
}

func (cmd *Cmd) loadCluster(name string) *ecs.Cluster {
	awsCluster, err := cmd.UFO.GetCluster(name)

	HandleError(err)

	return awsCluster
}

func (cmd *Cmd) loadTaskDefinition(c *ecs.Cluster, name string) *ecs.TaskDefinition {
	awsService, err := cmd.UFO.GetService(c, name)

	HandleError(err)

	t, err := cmd.UFO.GetTaskDefinition(c, awsService)

	HandleError(err)

	return t
}

func (cmd *Cmd) isDeployed(c *ecs.Cluster, s *ecs.Service, t *ecs.TaskDefinition) bool {
	if *s.DesiredCount <= 0 {
		return false
	}

	runningTasks, err := cmd.UFO.RunningTasks(c, s)

	HandleError(err)

	if len(runningTasks) <= 0 {
		return false
	}

	tasks, err := cmd.UFO.GetTasks(c, runningTasks)

	HandleError(err)

	for _, task := range tasks.Tasks {
		if *task.TaskDefinitionArn == *t.TaskDefinitionArn && *task.LastStatus == "RUNNING" {
			return true
		}
	}

	return false
}

func (cmd *Cmd) getCurrentHead() (string, error) {
	command := exec.Command("git", "rev-parse", "--short", "HEAD")

	r, err := command.Output()

	if err != nil {
		return "", ErrGitError
	}

	return strings.Trim(string(r), "\n"), nil
}

func (cmd *Cmd) getCurrentBranch() (string, error) {
	command := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")

	r, err := command.Output()

	if err != nil {
		return "", ErrGitError
	}

	return strings.Trim(string(r), "\n"), nil
}
