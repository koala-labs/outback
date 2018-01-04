package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	UFO "gitlab.fuzzhq.com/Web-Ops/ufo/ufo"
)

var (
	flagTaskCluster         string
	flagTaskService         string
	flagTaskCommandOverride string
)

var taskRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a one off task",
	Long:  `Run a one off task based off the task definition of a service. Override the task definitions command with -o flag`,
	Run:   taskRun,
}

func taskRun(cmd *cobra.Command, args []string) {
	env, err := getSelectedEnv(flagTaskCluster)

	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	service, err := getSelectedService(env.Services, flagTaskService)

	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	run(env.Cluster, *service)
}

func run(cluster string, service string) error {
	ufo := UFO.New(ufoCfg)
	c, err := ufo.GetCluster(cluster)

	if err != nil {
		return err
	}

	s, err := ufo.GetService(c, service)

	if err != nil {
		return err
	}

	t, err := ufo.GetTaskDefinition(c, s)

	if err != nil {
		return err
	}

	_, err = ufo.RunTask(c, t, flagTaskCommandOverride)

	if err != nil {
		return err
	}

	fmt.Printf("Running task...")

	return nil
}

func init() {

	taskCmd.AddCommand(taskRunCmd)

	taskRunCmd.Flags().StringVarP(&flagTaskCluster, "cluster", "c", "", "cluster")
	taskRunCmd.Flags().StringVarP(&flagTaskService, "service", "s", "", "service to use base image from")
	taskRunCmd.Flags().StringVarP(&flagTaskCommandOverride, "override", "o", "", "command to override")
}
