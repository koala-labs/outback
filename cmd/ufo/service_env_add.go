package cmd

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/spf13/cobra"
	UFO "gitlab.fuzzhq.com/Web-Ops/ufo/pkg/ufo"
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
	u := UFO.New(ufoCfg)

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

	updatedDefinition, err := updateTaskDefinition(t, flagServiceAddEnvVars)

	if err != nil {
		return err
	}

	registeredDefinition, err := u.RegisterTaskDefinitionWithEnvVars(updatedDefinition)

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

func updateTaskDefinition(t *ecs.TaskDefinition, inputs []string) (*ecs.TaskDefinition, error) {
	current := t.ContainerDefinitions[0].Environment

	parsed, err := stringsToKeyValue(inputs)

	if err != nil {
		return nil, err
	}

	t.ContainerDefinitions[0].Environment = updateEnvVars(current, parsed)

	return t, nil
}

func updateEnvVars(current []*ecs.KeyValuePair, updates []*ecs.KeyValuePair) []*ecs.KeyValuePair {
	for _, u := range updates {
		if i, ok := contains(current, u); ok {
			current[*i].Value = u.Value
		} else {
			current = append(current, u)
		}
	}

	return current
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

// contains returns an index and bool if keyVal.Name is in the keyVals slice
func contains(keyVals []*ecs.KeyValuePair, keyVal *ecs.KeyValuePair) (*int, bool) {
	for i, kv := range keyVals {
		if kv.Name == keyVal.Name {
			return &i, true
		}
	}
	return nil, false
}

func init() {
	serviceEnvCmd.AddCommand(serviceAddEnvCmd)

	serviceAddEnvCmd.Flags().StringSliceVarP(&flagServiceAddEnvVars, "env", "e", []string{}, "Environment variables to add e.g. key=value")
}
