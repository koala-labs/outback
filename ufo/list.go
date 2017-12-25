package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/service/ecs"
)

type listOptions struct {
	ListEnvs bool
}

type listCmd struct {
	c       *Config
	cmd     *Cmd
	Env     *Environment
	Options listOptions
}

func ListCommand(c *Config, options listOptions) {
	l := &listCmd{
		c:       c,
		Options: options,
	}

	printConfigInfo(c)

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

	l.printInfoForAllEnvironments(c)
}

func printConfigInfo(c *Config) {
	fmt.Printf("AWS Profile:            %s\n", c.Profile)
	fmt.Printf("ECR Repository:         %s\n\n", c.ImageRepositoryURL)
	fmt.Printf("Number of environments: %d\n", len(c.Env))
}

func printInfoForEnvironment(e *Environment) {
	CWD, err := os.Getwd()

	if err != nil {
		HandleError(ErrNoWorkingDirectory)
	}

	fmt.Printf("Branch:     %s\n", e.Branch)
	fmt.Printf("Region:     %s\n", e.Region)
	fmt.Printf("Cluster:    %s\n", e.Cluster)
	fmt.Printf("Services:    %s\n", strings.Join(e.Services, ","))
	fmt.Printf("Dockerfile: %s\n", CWD+"/"+e.Dockerfile)
}

func (l *listCmd) printInfoForAllEnvironments(c *Config) {
	for _, env := range c.Env {
		fmt.Println("================================")
		printInfoForEnvironment(env)
		if l.Options.ListEnvs {
			l.printEnvsForEnvironment(env)
		}
		fmt.Println("================================")
	}
}

func (l *listCmd) printEnvsForEnvironment(e *Environment) {
	cluster := l.cmd.loadCluster(e.Cluster)

	for _, service := range e.Services {
		fmt.Printf("Service: %s.", service)

		_, taskDef := l.cmd.loadService(cluster, service)

		for _, containerDefinition := range taskDef.ContainerDefinitions {
			longestName, longestValue := longestNameAndValue(containerDefinition.Environment)
			nameDashes := strings.Repeat("-", longestName+2) // Adding two because of the table padding
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

func (l *listCmd) printEnvsForAllEnvironments(c *Config) {
	for _, env := range c.Env {
		l.printEnvsForEnvironment(env)
	}
}
