package cmd

import (
	"fmt"
	Outback "github.com/koala-labs/outback/pkg/outback"
	"github.com/spf13/cobra"
)

var pullLatestCmd = &cobra.Command{
	Use:   "pull-latest",
	Short: "Pull the latest docker image from a deployment",
	Long: `A cluster must be specified via the --cluster flag, and a service must be specified via the --cluster flag.
	The --verbose flag can be input to enable verbose output.
	The --login flag can be input to login to AWS ECR.`,
	RunE: runPullLatest,
}

func runPullLatest(cmd *cobra.Command, args []string) error {
	return pullLatest(flagCluster, flagService)
}

func pullLatest(clusterName string, serviceName string) error {
	outback := Outback.New(awsConfig)

	cluster, err := cfg.getCluster(clusterName)
	if err != nil {
		return err
	}

	ecsCluster, err := outback.GetCluster(cluster.Name)
	if err != nil {
		return err
	}

	ecsService, err := outback.GetService(ecsCluster, serviceName)
	if err != nil {
		return err
	}

	// Get the Service's TaskDefinition
	ecsTaskDef, err := outback.GetTaskDefinition(ecsCluster, ecsService)
	if err != nil {
		return err
	}

	// Get the commit from the last TaskDefinition if it exists
	commit, err := outback.GetLastDeployedCommit(*ecsTaskDef.TaskDefinitionArn)
	if err != nil {
		return err
	}

	// Pull latest image
	err = outback.LoginPullImage(cfg.Repo, commit)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully pulled the latest image from :\n\t%s:%s\n", cfg.Repo, commit)

	return nil
}

func init() {
	rootCmd.AddCommand(pullLatestCmd)
}
