package main

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/aws/aws-sdk-go/service/ecs"
	"strings"
)

const DEPLOY_POLLING_RATE = 1 * time.Second
const HEADER_LENGTH = 120
const ATTEMPTS_TTL = 160

type DeployOptions struct {
	Verbose        bool
	OverrideBranch string
}

type DeployCmd struct {
	Options DeployOptions
	s       []*DeployState
	c       *Config
	cmd     *Cmd
	Env     *Environment
	branch  string
	head    string
}

// Runs through all steps of the deploy command
func RunDeployCmd(c *Config, options DeployOptions) error {
	var err error

	d := &DeployCmd{
		Options: options,
		branch:  options.OverrideBranch,
		c:       c,
	}

	d.head, err = d.cmd.getCurrentHead()

	if err != nil {
		return err
	}

	if d.Options.OverrideBranch == EmptyValue {
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
	d.s = make([]*DeployState, 0)

	err = d.UploadDockerImages()

	if err != nil {
		return err
	}

	d.InitDeployments()

	if err != nil {
		return err
	}

	ticker := time.NewTicker(DEPLOY_POLLING_RATE)

	for range ticker.C {
		ClearScreen()
		d.PrintStatus()

		if d.AllDeploymentsComplete() {
			ticker.Stop()

			break
		}
	}

	ClearScreen()
	d.PrintStatus()

	return nil
}

// Builds and pushes new docker images that will be used for service deployments
// This should only be called once and before all the service goroutines are run
func (d *DeployCmd) UploadDockerImages() error {
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

	return nil
}

// Run through all the configured services and create a goroutine for their deployment
// DeployCmd keeps an array of DeployState pointers for each service which will be deployed.
func (d *DeployCmd) InitDeployments() {
	for _, service := range d.Env.Services {
		s := &DeployState{
			ServiceName: service,
			LastStatus:  "Starting",
			Done:        false,
		}

		d.s = append(d.s, s)

		go d.deploy(service, s)
	}
}

// Run through the array of DeployStates to determine if they've completed or errored.
func (d *DeployCmd) AllDeploymentsComplete() bool {
	for _, s := range d.s {
		if ! s.Done && s.Error == nil {
			return false
		}
	}

	return true
}

// Run through the array of DeployStates and print a status of each one
func (d *DeployCmd) PrintStatus() {
	for _, s := range d.s {
		nameLength := len(s.ServiceName)
		sideLength := (HEADER_LENGTH - nameLength) / 2

		fmt.Printf("%s%s%s\n", strings.Repeat("=", sideLength), s.ServiceName, strings.Repeat("=", sideLength))

		if s.cluster != nil {
			fmt.Printf("Cluster: %s\n", *s.cluster.ClusterName)
		}

		if s.service != nil {
			fmt.Printf("Service: %s\n", *s.service.ServiceName)
		}

		fmt.Printf("Status: %s\n", s.LastStatus)
		fmt.Printf("Is complete: %t\n", s.Done)

		if s.Error != nil {
			fmt.Errorf("Encountered an error: %v", s.Error)
		}

		fmt.Printf("%s\n", strings.Repeat("=", (sideLength * 2) + nameLength))
	}
}

// Deploy an individual service
func (d *DeployCmd) deploy(service string, s *DeployState) error {
	s.UpdateStatus(fmt.Sprintf("Preparing to deploy branch %s to service %s on cluster %s.\n", d.Env.Branch, service, d.Env.Cluster))

	s.cluster = d.cmd.loadCluster(d.Env.Cluster)
	s.service, s.oldT = d.cmd.loadService(s.cluster, service)

	s.UpdateStatus(fmt.Sprintf("Beginning deployment to service %s.\n", service))
	t, err := d.cmd.UFO.Deploy(s.cluster, s.service, d.head)

	if err != nil {
		s.Error = err
		return err
	}

	s.newT = t

	s.UpdateStatus("Waiting for new containers to start.")
	err = d.awaitCompletion(s)

	if err != nil {
		s.Error = err
		return err
	}

	s.UpdateStatus(fmt.Sprintf("Successfully deployed. Your new task definition is %s:%d.", *s.newT.Family, *s.newT.Revision))

	return nil
}

// Create a Docker image from the vcs head
func (d *DeployCmd) buildImage() error {
	c := fmt.Sprintf("%s:%s", d.c.ImageRepositoryURL, d.head)
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

// Push the newly created docker image to the image repo
func (d *DeployCmd) pushImage() error {
	c := fmt.Sprintf("%s:%s", d.c.ImageRepositoryURL, d.head)
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

// Poll service update completion
// A service is updated when a container with the new task defintion is in a "RUNNING" state
func (d *DeployCmd) awaitCompletion(s *DeployState) error {
	attempts := 0
	waitTime := 2 * time.Second

	for !d.cmd.isDeployed(s.cluster, s.service, s.newT) {
		if attempts > ATTEMPTS_TTL {
			return ErrDeployTimeout
		}

		attempts++

		s.UpdateStatus("Waiting for deployment to complete.")
		time.Sleep(waitTime)
	}

	s.Done = true

	return nil
}

type DeployState struct {
	cluster     *ecs.Cluster
	service     *ecs.Service
	oldT        *ecs.TaskDefinition
	newT        *ecs.TaskDefinition
	ServiceName string
	LastStatus  string
	Done        bool
	Error       error
}

// Update the status of the DeployState
func (s *DeployState) UpdateStatus(status string) {
	s.LastStatus = status
}
