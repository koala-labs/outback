package cmd

import (
	"github.com/spf13/cobra"
)

var serviceEnvAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add environment variables",
}

func init() {
	serviceEnvCmd.AddCommand(serviceEnvAddCmd)
}
