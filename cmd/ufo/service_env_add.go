package cmd

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/spf13/cobra"
	UFO "gitlab.fuzzhq.com/Web-Ops/ufo/ufo"
)

var (
	flagServiceEnvAddEnvVars []string
)

var serviceEnvAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add/Update environment variables",
	RunE:  envAdd,
}

func envAdd(cmd *cobra.Command, args []string) error {
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

	updatedDef, err := updateTaskDefinition(t, flagServiceEnvAddEnvVars)

	if err != nil {
		return err
	}

	registeredDef, err := u.RegisterTaskDefinitionWithEnvVars(updatedDef)

	if err != nil {
		return err
	}

	_, err = u.UpdateService(c, s, registeredDef)

	if err != nil {
		return err
	}

	fmt.Println("Environment variables added")
	return nil
}

func updateTaskDefinition(t *ecs.TaskDefinition, inputs []string) (*ecs.TaskDefinition, error) {
	current := getEnvVars(t)

	parsed, err := parseEnvVars(inputs)

	if err != nil {
		return nil, err
	}

	t.ContainerDefinitions[0].Environment = updateEnvVars(current, parsed)

	return t, nil
}

func updateEnvVars(current []*ecs.KeyValuePair, updates []*ecs.KeyValuePair) []*ecs.KeyValuePair {
	// Loop through currently set EnvVars
	for _, u := range updates {
		if index, result := contains(current, u); result {
			current[*index].Value = u.Value
		} else {
			current = append(current, u)
		}
	}

	return current
}

// contains returns an index and bool if the value is in the slice
func contains(kvs []*ecs.KeyValuePair, kv *ecs.KeyValuePair) (*int, bool) {
	for i, v := range kvs {
		if v.Name == kv.Name {
			return &i, true
		}
	}
	return nil, false
}

func getEnvVars(t *ecs.TaskDefinition) []*ecs.KeyValuePair {
	return t.ContainerDefinitions[0].Environment
}

func parseEnvVars(inputs []string) ([]*ecs.KeyValuePair, error) {
	envVars := make([]*ecs.KeyValuePair, 0)

	if len(inputs) == 0 {
		return envVars, nil
	}

	for _, in := range inputs {
		splitIn := strings.SplitN(in, "=", 2)

		if len(splitIn) != 2 {
			return nil, ErrInvalidEnvInput
		}

		envVar := &ecs.KeyValuePair{
			Name:  aws.String(strings.ToUpper(splitIn[0])),
			Value: aws.String(splitIn[1]),
		}

		envVars = append(envVars, envVar)
	}

	return envVars, nil
}

func init() {
	serviceEnvCmd.AddCommand(serviceEnvAddCmd)

	serviceEnvAddCmd.Flags().StringSliceVarP(&flagServiceEnvAddEnvVars, "env", "e", []string{}, "Environment variables to add e.g. key=value")
}
