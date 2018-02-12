package cmd

import (
	"bufio"
	"fmt"
	"os/exec"
	"time"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/spf13/cobra"
	"gitlab.fuzzhq.com/Web-Ops/ufo/pkg/git"
	"gitlab.fuzzhq.com/Web-Ops/ufo/pkg/term"
	UFO "gitlab.fuzzhq.com/Web-Ops/ufo/pkg/ufo"
)

const (
	attemptsTTL       = 60
	headerLength      = 120
	deployPollingRate = 2 * time.Second
)

var (
	flagDeployVerbose bool
	flagECRLogin      bool
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Create a deployment",
	Long: `A cluster must be specified via the --cluster flag. 
	The --verbose flag can be input to enable verbose output. 
	The --login flag can be input to login to AWS ECR.`,
	RunE: deployRun,
}

type DeployOperation struct {
	s      []*DeployState
	branch string
	head   string
}

type DeployState struct {
	cluster     *ecs.Cluster
	service     *ecs.Service
	oldDef      *ecs.TaskDefinition
	newDef      *ecs.TaskDefinition
	ServiceName string
	LastStatus  string
	Done        bool
	Error       error
}

func deployRun(cmd *cobra.Command, args []string) error {
	var err error

	c, err := cfg.getCluster(flagCluster)

	if err != nil {
		fmt.Printf("Error: %v", err)
		return err
	}

	op := &DeployOperation{}

	op.head, err = git.GetHead()

	if err != nil {
		return err
	}

	op.branch, err = git.GetBranch()

	if err != nil {
		return err
	}

	op.s = make([]*DeployState, len(c.Services))

	if flagECRLogin {
		op.ecrLogin(cfg.Profile, cfg.Region)
	}

	err = op.buildPushImage(cfg.Repo, op.head, c.Dockerfile)

	if err != nil {
		return err
	}

	op.InitDeployments(flagCluster)

	ticker := time.NewTicker(deployPollingRate)

	for range ticker.C {
		term.Clear()

		op.PrintStatus()

		if op.AllDeploymentsComplete() {
			ticker.Stop()

			break
		}
	}

	return nil
}

// AllDeploymentsComplete runs through the array of DeployStates to determine
// if they have successfully deployed or errored
func (d *DeployOperation) AllDeploymentsComplete() bool {
	for _, s := range d.s {
		if !s.Done && s.Error == nil {
			return false
		}
	}

	return true
}

// PrintStatus runs through the array of DeployStates and print a status of each one
func (d *DeployOperation) PrintStatus() {
	for _, s := range d.s {
		if s.cluster != nil {
			fmt.Printf("Cluster: %s\n", *s.cluster.ClusterName)
		}

		if s.service != nil {
			fmt.Printf("Service: %s\n", *s.service.ServiceName)
		}

		fmt.Printf("Status: %s\n", s.LastStatus)

		if s.Error != nil {
			fmt.Printf("Encountered an error: %v", s.Error)
		}
	}
}

// InitDeployments runs through all the configured services and creates a goroutine for their deployment.
// DeployOperation keeps an array of DeployState pointaers for each service which will be deployed
func (d *DeployOperation) InitDeployments(cluster string) {
	c, err := cfg.getCluster(cluster)

	handleError(err)

	for i, service := range c.Services {
		s := &DeployState{
			ServiceName: service,
			LastStatus:  "Starting",
			Done:        false,
		}

		d.s[i] = s

		go d.deploy(cluster, service, s)
	}
}

// buildPushImage builds and pushes a new image
// This should only be called once and before all the service goroutines are run
func (d *DeployOperation) buildPushImage(repo string, tag string, dockerfile string) error {

	err := d.buildImage(repo, tag, dockerfile)

	if err != nil {
		return err
	}

	err = d.pushImage(repo, tag)

	if err != nil {
		return err
	}

	return nil
}

// deploy deploys an individual service and awaits for the newly created task to be "RUNNING"
func (d *DeployOperation) deploy(cluster string, service string, s *DeployState) error {
	var err error

	ufo := UFO.New(ufoCfg)

	s.UpdateStatus(fmt.Sprintf("Preparing to deploy branch %s to service %s on cluster %s\n",
		d.branch, service, cluster))

	s.cluster, err = ufo.GetCluster(cluster)

	if err != nil {
		s.Error = err
		return err
	}

	s.service, err = ufo.GetService(s.cluster, service)

	if err != nil {
		s.Error = err
		return err
	}

	s.UpdateStatus(fmt.Sprintf("Beginning deployment to service %s\n", service))

	s.newDef, err = ufo.Deploy(s.cluster, s.service, d.head)

	if err != nil {
		s.Error = err
		return err
	}

	s.UpdateStatus("Waiting for deployment to complete")

	err = d.awaitCompletion(s)

	if err != nil {
		s.Error = err
		return err
	}

	s.UpdateStatus(fmt.Sprintf("Successfully deployed. \nYour new task definition is %s:%d\n", *s.newDef.Family, *s.newDef.Revision))

	return nil
}

func (d *DeployOperation) ecrLogin(profile string, region string) error {
	fmt.Println("Logging into ECR")
	cmd := fmt.Sprintf("$(aws ecr get-login --no-include-email --region %s --profile %s)", region, profile)
	getLogin := exec.Command("bash", "-c", cmd)

	if err := printStdout(getLogin); err != nil {
		return ErrECRLogin
	}

	return nil
}

// buildImage builds a docker image based on the configured dockerfile for
// the cluster you are deploying to and tags the image with the vcs head
func (d *DeployOperation) buildImage(repo string, tag string, dockerfile string) error {
	fmt.Println("Building docker image")

	image := fmt.Sprintf("%s:%s", repo, tag)

	cmd := exec.Command("docker", "build", "-f", dockerfile, "-t", image, ".")

	if err := printStdout(cmd); err != nil {
		return ErrDockerBuild
	}

	return nil
}

// pushImage pushes the image built from buildImage to the configured repository
func (d *DeployOperation) pushImage(repo string, tag string) error {
	fmt.Println("Pushing docker image")

	image := fmt.Sprintf("%s:%s", repo, tag)

	cmd := exec.Command("docker", "push", image)

	if err := printStdout(cmd); err != nil {
		return ErrDockerPush
	}

	return nil
}

// awaitCompletion polls a services tasks for its status until its status is "RUNNING"
func (d *DeployOperation) awaitCompletion(s *DeployState) error {
	var err error

	attempts := 0
	waitTime := 5 * time.Second

	ufo := UFO.New(ufoCfg)

	ok := false

	for !ok {
		ok, err = ufo.IsServiceRunning(s.cluster, s.service, s.newDef)

		if err != nil {
			return err
		}

		if attempts > attemptsTTL {
			return ErrDeployTimeout
		}

		attempts++

		s.UpdateStatus(fmt.Sprintf("Waiting for deployment of %s:%d to complete", *s.newDef.Family, *s.newDef.Revision))

		time.Sleep(waitTime)
	}

	s.Done = true

	return nil
}

// UpdateStatus updates the status of the DeployState
func (s *DeployState) UpdateStatus(status string) {
	s.LastStatus = status
}

func printStdout(command *exec.Cmd) error {
	stdout, err := command.StdoutPipe()

	if err != nil {
		fmt.Printf("%v", err)
		return fmt.Errorf("Error creating a stdout pipe")
	}

	if err := command.Start(); err != nil {
		fmt.Printf("%v", err)
		return fmt.Errorf("Error starting the command")
	}

	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			out := scanner.Text()
			fmt.Println(out)
		}
	}()

	if err := command.Wait(); err != nil {
		fmt.Printf("%v", err)
		return fmt.Errorf("Error waiting on releases of executed command")
	}

	return nil
}

func init() {
	rootCmd.AddCommand(deployCmd)

	deployCmd.Flags().BoolVarP(&flagDeployVerbose, "verbose", "v", false, "Shows output of the deployment process")
	deployCmd.Flags().BoolVarP(&flagECRLogin, "login", "l", false, "Log in to ECR during the deployment process")
}
