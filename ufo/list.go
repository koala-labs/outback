package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/service/ecs"
)

type ListOptions struct {
	ListEnvs bool
}

type ListCmd struct {
	c       *Config
	cmd     *Cmd
	Env     *Environment
	Options ListOptions
}

func ListCommand(c *Config, options ListOptions) {
	l := &ListCmd{
		c:       c,
		Options: options,
	}

	PrintConfigInfo(c)

	b, err := l.cmd.getCurrentBranch()

	if err != nil {
		fmt.Println(err)
	}

	e, err := c.GetEnvironmentByBranch(b)

	if err != nil {
		fmt.Println(err)
	}

	l.Env = e

	l.cmd = l.cmd.initUFO(l.c.Profile, l.Env.Region)

	l.PrintInfoForAllEnvironments(c)
}

func PrintConfigInfo(c *Config) {
	fmt.Printf("AWS Profile:            %s\n", c.Profile)
	fmt.Printf("ECR Repository:         %s\n\n", c.ImageRepositoryUrl)
	fmt.Printf("Number of environments: %d\n", len(c.Env))
}

func PrintInfoForEnvironment(e *Environment) {
	CWD, err := os.Getwd()

	if err != nil {
		HandleError(ErrNoWorkingDirectory)
	}

	fmt.Printf("Branch:     %s\n", e.Branch)
	fmt.Printf("Region:     %s\n", e.Region)
	fmt.Printf("Cluster:    %s\n", e.Cluster)
	fmt.Printf("Service:    %s\n", e.Service)
	fmt.Printf("Dockerfile: %s\n", CWD+"/"+e.Dockerfile)
}

func (l *ListCmd) PrintInfoForAllEnvironments(c *Config) {
	for _, env := range c.Env {
		fmt.Println("================================")
		PrintInfoForEnvironment(env)
		if l.Options.ListEnvs {
			l.PrintEnvsForEnvironment(env)
		}
		fmt.Println("================================")
	}
}

func (l *ListCmd) PrintEnvsForEnvironment(e *Environment) {
	cluster := l.cmd.loadCluster(e.Cluster)
	_, taskDef := l.cmd.loadService(cluster, e.Service)
	for _, containerDefinition := range taskDef.ContainerDefinitions {
		longestName, longestValue := longestNameAndValue(containerDefinition.Environment)
		nameDashes := strings.Repeat("-", longestName+2)
		valueDashes := strings.Repeat("-", longestValue+2)
		fmt.Printf("+%s+%s+\n", nameDashes, valueDashes)
		for _, value := range containerDefinition.Environment {
			name := *value.Name
			value := *value.Value
			nameSpaces := longestName - len(name)
			valueSpaces := longestValue - len(value)
			spacesForName := strings.Repeat(" ", nameSpaces)
			spacesForValue := strings.Repeat(" ", valueSpaces)
			fmt.Printf("| %s%s | %s%s |\n", name, spacesForName, value, spacesForValue)
			fmt.Printf("+%s+%s+\n", nameDashes, valueDashes)
		}
	}
}

func longestNameAndValue(e []*ecs.KeyValuePair) (longestName int, longestValue int) {
	for _, value := range e {
		name := *value.Name
		value := *value.Value
		nameLength := len(name)
		valueLength := len(value)
		if nameLength > longestName {
			longestName = nameLength
		}
		if valueLength > longestValue {
			longestValue = valueLength
		}
	}
	return longestName, longestValue
}

func (l *ListCmd) PrintEnvsForAllEnvironments(c *Config) {
	for _, env := range c.Env {
		l.PrintEnvsForEnvironment(env)
	}
}
