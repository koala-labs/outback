package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Interactively create a deployment by selecting from available configs and clusters",
	RunE:  interactiveDeploy,
}

func interactiveDeploy(cmd *cobra.Command, args []string) error {
	var configQuestion = []*survey.Question{
		{
			Name: "config",
			Prompt: &survey.Select{
				Message: "Choose a config:",
				Options: cfg.getConfigs(),
			},
		},
	}

	configAnswer := struct {
		Config string `survey:"config"`
	}{}

	err := survey.Ask(configQuestion, &configAnswer)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	viper.SetConfigName(configAnswer.Config)
	viper.ReadInConfig()
	viper.Unmarshal(&cfg)

	var clusterQuestion = []*survey.Question{
		{
			Name: "cluster",
			Prompt: &survey.Select{
				Message: "Choose a cluster:",
				Options: cfg.getClusters(),
			},
		},
		{
			Name: "login",
			Prompt: &survey.Select{
				Message: "Login to ECR?",
				Options: []string{"no", "yes"},
			},
		},
	}

	clusterAnswer := struct {
		Cluster string `survey:"cluster"`
		Login   string `survey:"login"`
	}{}

	err = survey.Ask(clusterQuestion, &clusterAnswer)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	var confirmQuestion = []*survey.Question{
		{
			Name: "confirm",
			Prompt: &survey.Select{
				Message: fmt.Sprintf("Are you sure you want to deploy services %v?", cfg.getServices(clusterAnswer.Cluster)),
				Options: []string{"no", "yes"},
			},
		},
	}

	confirmAnswer := struct {
		Confirm string `suvery:"confirm"`
	}{}

	err = survey.Ask(confirmQuestion, &confirmAnswer)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return fullDeploy(clusterAnswer.Cluster, toBool(clusterAnswer.Login))
}

func toBool(answer string) bool {
	if answer == "yes" {
		return true
	}

	return false
}

func init() {
	rootCmd.AddCommand(interactiveCmd)
}
