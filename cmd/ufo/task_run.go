package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	UFO "gitlab.fuzzhq.com/Web-Ops/ufo/pkg/ufo"
)

var (
	flagTaskCommand string
)

var taskRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a one off tasks",
	Long: `You must specify a cluster, service, and command to run. The command will use the image described in the task definition for the service that is specified. When specifying a command, the task definitions current command will be overriden with the one specified. 
	There is also an option of creating command aliases in .ufo/config.json. Once a command alias is in the ufo config, specifying that alias via the --command flag will run the configured command.
	If the awslogs driver is configured for the service in which you base your task. Logs for that task will be sent to cloudwatch under the same log group and prefix as described in the task definition.`,
	Run: taskRun,
}

func taskRun(cmd *cobra.Command, args []string) {
	cfgCluster, err := cfg.getCluster(flagCluster)

	handleError(err)

	cfgService, err := cfg.getService(cfgCluster.Services, flagService)

	handleError(err)

	// Check if the command is available in the config as a shortcut
	command, err := cfg.getCommand(flagTaskCommand)

	// If the shortcut is not in the config, pass the command directly
	if err != nil {
		err = run(cfgCluster.Name, *cfgService, flagTaskCommand)
	} else {
		err = run(cfgCluster.Name, *cfgService, *command)
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
