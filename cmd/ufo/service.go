package cmd

import (
	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Run various tasks for a service",
	Long:  `service`,
}

func init() {
	rootCmd.AddCommand(serviceCmd)
}
