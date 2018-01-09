package cmd

import (
	"github.com/spf13/cobra"
)

var serviceEnvRmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove environment variables",
}

func init() {
	serviceEnvCmd.AddCommand(serviceEnvRmCmd)
}
