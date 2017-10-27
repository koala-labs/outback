package main

import (
	"errors"
	"fmt"

	"log"
	"os"
	"sort"
	"time"

	"github.com/abiosoft/ishell"
	"github.com/aws/aws-sdk-go/service/ecs"
	"gitlab.fuzzhq.com/Web-Ops/ufo/pkg/ufo"
)

type AppState struct {
	c       *ecs.Cluster
	s       *ecs.Service
	oldT    *ecs.TaskDefinition
	newT    *ecs.TaskDefinition
	version string
}

type AppFlags struct {
	noInteractive bool
	cluster       string
	service       string
	version       string
}

type App struct {
	AppState
	f       *AppFlags
	Profile *string
	Region  *string
	UFO     *ufo.UFO
	Shell   *ishell.Shell
}

func (a *App) Init() {
	if a.f.noInteractive {
		a.runNoInteractive()
		return
	}

	a.Shell.Println("Welcome to UFO!")
	a.Shell.Println("Use `help` or `start` to continue.")

	a.Shell.AddCmd(&ishell.Cmd{
		Name: "start",
		Help: "Start the deployment process.",
		Func: a.Executor,
	})

	a.Shell.Run()
}

func (a *App) runNoInteractive() {
	log.Println("Running in non-interactive mode.")

	a.AppState.c = a.loadCluster(a.f.cluster)
	a.AppState.s, a.AppState.oldT = a.loadService(a.f.service)
	a.AppState.version = a.f.version

	a.AppState.newT = a.deploy()

	attempts := 0
	waitTime := 2 * time.Second

	for !a.IsDeployed(a.AppState.c, a.AppState.s, a.AppState.newT) {
		if attempts > 60 {
			a.HandleError(errors.New("Timed out waiting for task to start."))

			return
		}

		attempts++

		log.Print("Waiting.")
		time.Sleep(waitTime)
	}

	log.Printf("Successfully deployed. Your new task definition is %s:%d.\n", *a.AppState.newT.Family, *a.AppState.newT.Revision)
}

// Run the user through the required deployment steps
func (a *App) Executor(c *ishell.Context) {
	a.ChooseCluster(c)
	a.ChooseService(c)
	a.ChooseVersion(c)
	a.ConfirmDeployment(c)
	a.PollForStatus(c)

	os.Exit(0)
}

// HandleError is intended to be called with which error return to simplify error handling
// Usage:
// foo, err := GetFoo()
// a.HandleError(err)
// DoSomethingBecauseNoError()
func (a *App) HandleError(err error) {
	if err == nil {
		return
	}

	a.Shell.Printf("\nEncountered an error: %s", err.Error())

	os.Exit(1)
}

func (a *App) ChooseCluster(c *ishell.Context) {
	clusters, err := a.UFO.Clusters()

	a.HandleError(err)

	choice := c.MultiChoice(clusters, "Select a cluster: ")

	a.AppState.c = a.loadCluster(clusters[choice])
}

func (a *App) loadCluster(name string) *ecs.Cluster {
	awsCluster, err := a.UFO.GetCluster(name)

	a.HandleError(err)

	return awsCluster
}

func (a *App) ChooseService(c *ishell.Context) {
	services, err := a.UFO.Services(a.AppState.c)

	a.HandleError(err)

	choice := c.MultiChoice(services, "Select a service: ")

	a.AppState.s, a.AppState.oldT = a.loadService(services[choice])
}

func (a *App) loadService(name string) (*ecs.Service, *ecs.TaskDefinition) {
	awsService, err := a.UFO.GetService(a.AppState.c, name)

	a.HandleError(err)

	t, err := a.UFO.GetTaskDefinition(a.AppState.c, awsService)

	a.HandleError(err)

	return awsService, t
}

func (a *App) ChooseVersion(c *ishell.Context) {
	images, err := a.UFO.GetImages(a.AppState.oldT)

	sort.Slice(images, func(i, j int) bool {
		return images[i].ImagePushedAt.Unix() > images[j].ImagePushedAt.Unix()
	})

	a.HandleError(err)

	choices := make([]string, 0)

	for _, image := range images {
		choices = append(choices, fmt.Sprintf("%s: %s", image.ImagePushedAt, *image.ImageTags[0]))
	}

	choice := c.MultiChoice(choices, "Select a version to deploy: ")

	a.AppState.version = *images[choice].ImageTags[0]
}

func (a *App) ConfirmDeployment(c *ishell.Context) {
	c.Printf("Chosen cluster: %s\n", *a.AppState.c.ClusterName)
	c.Printf("Chosen service: %s\n", *a.AppState.s.ServiceName)
	c.Printf("Chosen version: %s\n", a.AppState.version)

	c.Println("Ready to deploy? (yes/no)")
	choice := c.ReadLine()

	if choice != "yes" {
		c.Println("Not ready to deploy, exiting.")

		return
	}

	a.AppState.newT = a.deploy()
}

func (a *App) deploy() *ecs.TaskDefinition {
	t, err := a.UFO.Deploy(a.AppState.c, a.AppState.s, a.AppState.version)

	a.HandleError(err)

	return t
}

func (a *App) PollForStatus(c *ishell.Context) {
	c.Println("Waiting for new task to deploy.")

	attempts := 0
	waitTime := 2 * time.Second

	for !a.IsDeployed(a.AppState.c, a.AppState.s, a.AppState.newT) {
		if attempts > 60 {
			a.HandleError(errors.New("Timed out waiting for task to start."))

			return
		}

		attempts++

		c.Print(".")
		time.Sleep(waitTime)
	}

	c.Printf("Successfully deployed. Your new task definition is %s:%d.\n", *a.AppState.newT.Family, *a.AppState.newT.Revision)
}

func (a *App) IsDeployed(c *ecs.Cluster, s *ecs.Service, t *ecs.TaskDefinition) bool {
	if *s.DesiredCount <= 0 {
		return false
	}

	runningTasks, err := a.UFO.RunningTasks(c, s)

	a.HandleError(err)

	if len(runningTasks) <= 0 {
		return false
	}

	tasks, err := a.UFO.GetTasks(c, runningTasks)

	a.HandleError(err)

	for _, task := range tasks.Tasks {
		if *task.TaskDefinitionArn == *t.TaskDefinitionArn && *task.LastStatus == "RUNNING" {
			return true
		}
	}

	return false
}
