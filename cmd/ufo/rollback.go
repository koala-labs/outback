package cmd

import (
	"fmt"
	"time"

	"github.com/fuzz-productions/ufo/pkg/term"
	UFO "github.com/fuzz-productions/ufo/pkg/ufo"
	"github.com/spf13/cobra"
)

var revisionNumber int

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback a deployment",
	Long: `A cluster must be specified via the --cluster flag.
	The --verbose flag can be input to enable verbose output.
	The --login flag can be input to login to AWS ECR.`,
	RunE: runRollback,
}

func runRollback(cmd *cobra.Command, args []string) error {
	return rollback(flagCluster, flagTimeout)
}

func rollback(clusterName string, timeout int) error {
	ufo := UFO.New(awsConfig)

	cluster, err := cfg.getCluster(clusterName)
	if err != nil {
		return err
	}

	deployDetail := &UFO.DeployDetail{}
	deployDetail.SetRevisionNumber(revisionNumber)

	deployment := &UFO.Deployment{}

	for _, service := range cluster.Services {
		detail := ufo.NewDeployDetail()

		// Get the ECS Cluster
		ecsCluster, err := ufo.GetCluster(cluster.Name)
		if err != nil {
			return err
		}

		// Set the Cluster in the deployment detail
		detail.SetCluster(ecsCluster)

		// Get the ECS Service
		ecsService, err := ufo.GetService(detail.Cluster, service)
		if err != nil {
			return err
		}

		// Set the Service in the deployment detail
		detail.SetService(ecsService)

		// Get the Service's TaskDefinition
		// if revision set here
		ecsTaskDef, err := ufo.GetTaskDefinition(detail.Cluster, detail.Service)
		if err != nil {
			return err
		}

		// Set the TaskDefinition in the deployment detail
		detail.SetTaskDefinition(ecsTaskDef)

		deployment.DeployDetails = append(deployment.DeployDetails, detail)
	}

	term.Clear()

	errCh := ufo.RollbackAll(deployment, deployDetail)

	for err := range errCh {
		return err
	}

	fmt.Printf("Waiting for deployment(s) to services [ %s]\n", deployment.Services())
	doneCh := ufo.AwaitServicesRunning(deployment)

	for i := 0; i < len(deployment.DeployDetails); i++ {
		select {
		case detail := <-doneCh:
			fmt.Printf("Service %s (%s) is now running \n", *detail.Service.ServiceName, detail.TaskDefinitionFamilyName)
		case <-time.After(time.Minute * time.Duration(timeout)):
			return ErrDeployTimeout
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(rollbackCmd)
	rollbackCmd.Flags().IntVarP(&revisionNumber, "revision", "r", 0, "Set the task revision number")
}
