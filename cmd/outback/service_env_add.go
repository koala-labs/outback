package cmd

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	Outback "github.com/koala-labs/outback/pkg/outback"
	"github.com/spf13/cobra"
)

var (
	flagServiceAddEnvVars []string
)

var serviceAddEnvCmd = &cobra.Command{
	Use:   "add",
	Short: "Add/Update environment variables",
	Long: `At least one environment variable must be specified via the --env flag. Specify
	--env with a key=value parameter multiple times to add multiple variables.`,
	RunE: addEnvVar,
}

func addEnvVar(cmd *cobra.Command, args []string) error {
	u := Outback.New(awsConfig)

	c, err := u.GetCluster(flagCluster)

	if err != nil {
		return err
	}

	s, err := u.GetService(c, flagService)

	if err != nil {
		return err
	}

	t, err := u.GetTaskDefinition(c, s)

	if err != nil {
		return err
	}

	parsedEnvVars, err := stringsToKeyValue(flagServiceAddEnvVars)

	if err != nil {
		return err
	}
	updatedDefinition := u.UpdateContainerDefinitionEnvVars(*t, parsedEnvVars, cfg.Repo)

	registeredDefinition, err := u.RegisterTaskDefinitionWithEnvVars(&updatedDefinition)

	if err != nil {
		return err
	}

	_, err = u.UpdateService(c, s, registeredDefinition)

	if err != nil {
		return err
	}

	fmt.Println("Environment variable(s) " + strings.Join(flagServiceAddEnvVars, ", ") + " will be added")

	return nil
}

func stringsToKeyValue(inputs []string) ([]*ecs.KeyValuePair, error) {
	keyVals := make([]*ecs.KeyValuePair, 0)

	if len(inputs) == 0 {
		return keyVals, nil
	}

	for _, in := range inputs {
		splitIn := strings.SplitN(in, "=", 2)

		if len(splitIn) != 2 {
			return nil, ErrInvalidEnvInput
		}

		keyVal := &ecs.KeyValuePair{
			Name:  aws.String(strings.ToUpper(splitIn[0])),
			Value: aws.String(splitIn[1]),
		}

		keyVals = append(keyVals, keyVal)
	}

	return keyVals, nil
}

func init() {
	serviceEnvCmd.AddCommand(serviceAddEnvCmd)

	serviceAddEnvCmd.Flags().StringSliceVarP(&flagServiceAddEnvVars, "env", "e", []string{}, "Environment variables to add e.g. key=value")
}
