package cmd

import (
	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage an ECS service",
	Long: `Services manage long-lived instances of your containers that are run on AWS
	ECS. If your container exits for any reason, the service scheduler will
	restart your containers and ensure your service has the desired number of
	tasks running. Services can be used in concert with a load balancer to
	distribute traffic amongst the tasks in your service.`,
}

func init() {
	rootCmd.AddCommand(serviceCmd)
}
