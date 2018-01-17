package cmd

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/spf13/cobra"
	UFO "gitlab.fuzzhq.com/Web-Ops/ufo/pkg/ufo"
)

var (
	flagServiceEnvRmEnvVars []string
)

var serviceEnvRmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove environment variables and trigger a service update",
	RunE:  envRm,
}

func envRm(cmd *cobra.Command, args []string) error {
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

	newDefinition, err := removeEnvVarsFromTaskDefinition(t, flagServiceEnvRmEnvVars)

	if err != nil {
		return err
	}

	registeredDefinition, err := u.RegisterTaskDefinitionWithEnvVars(newDefinition)

	if err != nil {
		return err
	}

	_, err = u.UpdateService(c, s, registeredDefinition)

	if err != nil {
		return err
	}

	fmt.Println("The key(s) " + strings.Join(flagServiceEnvRmEnvVars, ", ") + " will be removed from your task definition")

	return nil
}

func removeEnvVarsFromTaskDefinition(t *ecs.TaskDefinition, removals []string) (*ecs.TaskDefinition, error) {
	encountered := map[string]bool{}

	for _, r := range removals {
		encountered[r] = true
	}

	current := t.ContainerDefinitions[0].Environment

	newSet := []*ecs.KeyValuePair{}

	for _, c := range current {
		if _, ok := encountered[*c.Name]; !ok {
			newSet = append(newSet, c)
		}
	}

	if len(current) == len(newSet) {
		return nil, ErrKeyNotPresent
	}

	t.ContainerDefinitions[0].Environment = newSet

	return t, nil
}

func init() {
	serviceEnvCmd.AddCommand(serviceEnvRmCmd)

	serviceEnvRmCmd.Flags().StringSliceVarP(&flagServiceEnvRmEnvVars, "key", "k", []string{}, "Environment variables to remove e.g. APP_ENV")
}
