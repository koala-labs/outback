package cmd

import (
	"github.com/spf13/cobra"
)

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Run a one off task",
	Long:  `task`,
}

func init() {
	rootCmd.AddCommand(taskCmd)
}
