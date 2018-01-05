package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	UFO "gitlab.fuzzhq.com/Web-Ops/ufo/ufo"
)

var (
	flagTaskCommand string
)

var taskRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a one off task",
	Long:  `Run a one off task based off the task definition of a service. Override the task definitions command with -o flag`,
	Run:   taskRun,
}

func taskRun(cmd *cobra.Command, args []string) {
	c, err := cfg.getSelectedCluster(flagCluster)

	handleError(err)

	service, err := cfg.getSelectedService(c.Services, flagService)

	handleError(err)

	// Check if the command is available in the config as a shortcut
	command, err := cfg.getCommand(flagTaskCommand)

	// If the shortcut is not in the config, pass the command directly
	if err != nil {
		err = run(c.Name, *service, flagTaskCommand)
	} else {
		err = run(c.Name, *service, *command)
	}

	handleError(err)
}

func run(cluster string, service string, command string) error {
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

	_, err = ufo.RunTask(c, t, command)

	if err != nil {
		return err
	}

	fmt.Printf("Running task on cluster %s with command %s", cluster, command)

	return nil
}

func init() {

	taskCmd.AddCommand(taskRunCmd)

	taskRunCmd.Flags().StringVarP(&flagTaskCommand, "command", "n", "", "name of the command to run from your config or the command itself")
}
