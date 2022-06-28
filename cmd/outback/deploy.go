package cmd

import (
	"fmt"
	"time"

	"github.com/koala-labs/outback/pkg/git"
	Outback "github.com/koala-labs/outback/pkg/outback"
	"github.com/koala-labs/outback/pkg/term"
	"github.com/spf13/cobra"
)

var deployBuildArgs []string

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Create a deployment",
	Long: `A cluster must be specified via the --cluster flag.
	The --verbose flag can be input to enable verbose output.
	The --login flag can be input to login to AWS ECR.`,
	RunE: runDeploy,
}

func runDeploy(cmd *cobra.Command, args []string) error {
	return deploy(flagCluster, flagTimeout)
}

func deploy(clusterName string, timeout int) error {
	outback := Outback.New(awsConfig)

	commit, err := git.GetCommit()
	if err != nil {
		return err
	}

	cluster, err := cfg.getCluster(clusterName)
	if err != nil {
		return err
	}

	configBuildArgs := cfg.getBuildArgs(clusterName)

	deployment := &Outback.Deployment{}
	deployment.SetCommitHash(commit)
	deployment.SetRepo(cfg.Repo)
	deployment.SetDockerfile(cluster.Dockerfile)
	deployment.SetBuildArgs(deployBuildArgs)
	deployment.SetConfigBuildArgs(configBuildArgs)

	for _, service := range cluster.Services {
		detail := outback.NewDeployDetail()

		// Get the ECS Cluster
		ecsCluster, err := outback.GetCluster(cluster.Name)
		if err != nil {
			return err
		}

		// Set the Cluster in the deployment detail
		detail.SetCluster(ecsCluster)

		// Get the ECS Service
		ecsService, err := outback.GetService(detail.Cluster, service)
		if err != nil {
			return err
		}

		// Set the Service in the deployment detail
		detail.SetService(ecsService)

		// Get the Service's TaskDefinition
		ecsTaskDef, err := outback.GetTaskDefinition(detail.Cluster, detail.Service)
		if err != nil {
			return err
		}

		// Set the TaskDefinition in the deployment detail
		detail.SetTaskDefinition(ecsTaskDef)

		// Get the commit from the last TaskDefinition if it exists
		commit, err := outback.GetLastDeployedCommit(*ecsTaskDef.TaskDefinitionArn)
		if err == nil {
			deployment.SetBuildCacheFrom([]string{fmt.Sprintf("%s:%s", deployment.BuildDetail.Repo, commit)})
			fmt.Printf("Will attempt to restore Docker cache for %s from commit: %s\n", service, commit)
		}

		deployment.DeployDetails = append(deployment.DeployDetails, detail)
	}

	// Build Docker image and push to repo
	err = outback.LoginBuildPushImage(deployment.BuildDetail)
	if err != nil {
		return err
	}

	term.Clear()

	errCh := outback.DeployAll(deployment)

	for err := range errCh {
		fmt.Printf("Deployment failed: %s \n", err)
		return err
	}

	fmt.Printf("Waiting for deployment(s) to services [ %s]\n", deployment.Services())
	doneCh := outback.AwaitServicesRunning(deployment)

	for i := 0; i < len(deployment.DeployDetails); i++ {
		select {
		case detail := <-doneCh:
			fmt.Printf("Service %s (%s) is now running \n", *detail.Service.ServiceName, detail.TaskDefinitionFamily())
		case <-time.After(time.Minute * time.Duration(timeout)):
			return ErrDeployTimeout
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.Flags().StringSliceVarP(&deployBuildArgs, "build-arg", "b", []string{}, "Set build-time variables")
}
