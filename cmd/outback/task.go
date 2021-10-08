package cmd

import (
	"github.com/spf13/cobra"
)

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Run a one off task",
	Long: `Tasks are one-time executions of your container.
	Instances of your task are run until you manually stop them either through AWS APIs,
	the AWS Management Console, or fargate task stop, or until they are interrupted for any reason.`,
}

func init() {
	rootCmd.AddCommand(taskCmd)

	taskCmd.MarkPersistentFlagFilename("cluster")
	taskCmd.MarkPersistentFlagFilename("service")
}
