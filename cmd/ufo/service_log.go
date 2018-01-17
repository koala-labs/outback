package cmd

import (
	"github.com/spf13/cobra"
)

var serviceLogCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage logs",
}

func init() {
	serviceCmd.AddCommand(serviceLogCmd)
}
