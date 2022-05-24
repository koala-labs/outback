package cmd

import (
	"fmt"

	"github.com/koala-labs/outback/pkg/git"
	Outback "github.com/koala-labs/outback/pkg/outback"
	"github.com/spf13/cobra"
)

var buildDockerArgs []string

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build and push a docker image a deployment",
	Long: `A cluster must be specified via the --cluster flag.
	The --verbose flag can be input to enable verbose output.
	The --login flag can be input to login to AWS ECR.`,
	RunE: runBuild,
}

func runBuild(cmd *cobra.Command, args []string) error {
	return build(flagCluster, flagTimeout)
}

func build(clusterName string, timeout int) error {
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
	deployment.SetBuildArgs(buildArgs)
	deployment.SetConfigBuildArgs(configBuildArgs)

	fmt.Println("Building image...")

	// Build Docker image and push to repo
	err = outback.LoginBuildPushImage(deployment.BuildDetail)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully built and pushed image to repository:\n\t%s:%s\n", cfg.Repo, commit)

	return nil
}

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().StringSliceVarP(&buildDockerArgs, "build-arg", "b", []string{}, "Set build-time variables")
}
