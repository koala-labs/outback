package cmd

import (
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
	}

	clusterAnswer := struct {
		Cluster string `survey:"cluster"`
		Login   string `survey:"login"`
	}{}

	err = survey.Ask(clusterQuestion, &clusterAnswer)
	if err != nil {
		return err
	}

	var confirmQuestion = []*survey.Question{
		{
			Name: "confirm",
			Prompt: &survey.Select{
				Message: "Are you sure you want to deploy?",
				Options: []string{"no", "yes"},
			},
		},
	}

	confirmAnswer := struct {
		Confirm string `suvery:"confirm"`
	}{}

	err = survey.Ask(confirmQuestion, &confirmAnswer)
	if err != nil {
		return err
	}

	if toBool(confirmAnswer.Confirm) {
		return deploy(clusterAnswer.Cluster)
	}

	return nil
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
