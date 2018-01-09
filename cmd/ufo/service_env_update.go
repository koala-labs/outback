package cmd

import (
	"github.com/spf13/cobra"
)

var serviceEnvUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update environment variables",
}

func init() {
	serviceEnvCmd.AddCommand(serviceEnvUpdateCmd)
}
