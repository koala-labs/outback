package cmd

import (
	"fmt"
	"time"

	"github.com/fuzz-productions/ufo/pkg/git"
	"github.com/fuzz-productions/ufo/pkg/term"
	UFO "github.com/fuzz-productions/ufo/pkg/ufo"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Create a deployment",
	Long: `A cluster must be specified via the --cluster flag.
	The --verbose flag can be input to enable verbose output.
	The --login flag can be input to login to AWS ECR.`,
	RunE: runDeploy,
}

func runDeploy(cmd *cobra.Command, args []string) error {
	return deploy(flagCluster)
}

func deploy(clusterName string) error {
	ufo := UFO.New(awsConfig)

	commit, err := git.GetCommit()
	if err != nil {
		return err
	}

	cluster, err := cfg.getCluster(clusterName)
	if err != nil {
		return err
	}

	deployment := &UFO.Deployment{}
	deployment.SetCommitHash(commit)
	deployment.SetRepo(cfg.Repo)
	deployment.SetDockerfile(cluster.Dockerfile)

	// Build Docker image and push to repo
	err = ufo.LoginBuildPushImage(deployment.BuildDetail)
	if err != nil {
		return err
	}

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
		ecsTaskDef, err := ufo.GetTaskDefinition(detail.Cluster, detail.Service)
		if err != nil {
			return err
		}

		// Set the TaskDefinition in the deployment detail
		detail.SetTaskDefinition(ecsTaskDef)

		deployment.DeployDetails = append(deployment.DeployDetails, detail)
	}

	term.Clear()

	errCh := ufo.DeployAll(deployment)

	for err := range errCh {
		return err
	}

	fmt.Printf("Waiting for deployment(s) to services [ %s]\n", deployment.Services())
	doneCh := ufo.AwaitServicesRunning(deployment)

	for i := 0; i < len(deployment.DeployDetails); i++ {
		select {
		case detail := <-doneCh:
			fmt.Printf("Service %s (%s) is now running \n", *detail.Service.ServiceName, detail.TaskDefinitionFamily())
		case <-time.After(time.Minute * 5):
			return ErrDeployTimeout
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(deployCmd)
}
