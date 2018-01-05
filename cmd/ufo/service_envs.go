package cmd

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/spf13/cobra"
	UFO "gitlab.fuzzhq.com/Web-Ops/ufo/ufo"
)

var serviceEnvsCmd = &cobra.Command{
	Use:   "envs",
	Short: "List envs for services in your cluster",
	Long:  `envs`,
	Run:   envsRun,
}

func envsRun(cmd *cobra.Command, args []string) {
	c, err := cfg.getSelectedCluster(flagCluster)

	handleError(err)

	s, err := cfg.getSelectedService(c.Services, flagService)

	handleError(err)

	printEnvsForService(c.Name, *s)
}

func printEnvsForService(cluster string, service string) {
	ufo := UFO.New(ufoCfg)

	c, err := ufo.GetCluster(cluster)

	handleError(err)

	s, err := ufo.GetService(c, service)

	handleError(err)

	t, err := ufo.GetTaskDefinition(c, s)

	handleError(err)

	printEnvTable(t)
}

func printEnvTable(t *ecs.TaskDefinition) {
	for _, containerDefinition := range t.ContainerDefinitions {
		longestName, longestValue := longNameAndValue(containerDefinition.Environment)
		nameDashes := strings.Repeat("-", longestName+2) // Adding two because of the table padding
		valueDashes := strings.Repeat("-", longestValue+2)

		fmt.Printf("+%s+%s+\n", nameDashes, valueDashes)

		for _, value := range containerDefinition.Environment {
			name := *value.Name
			value := *value.Value
			nameSpaces := longestName - len(name)
			valueSpaces := longestValue - len(value)
			spacesForName := strings.Repeat(" ", nameSpaces)
			spacesForValue := strings.Repeat(" ", valueSpaces)
			fmt.Printf("| %s%s | %s%s |\n", name, spacesForName, value, spacesForValue)
			fmt.Printf("+%s+%s+\n", nameDashes, valueDashes)
		}
	}
}

// longNameAndValue returns the length of the longest Name and Value
func longNameAndValue(e []*ecs.KeyValuePair) (longName int, longVal int) {
	for _, value := range e {
		nameLength := len(*value.Name)
		valueLength := len(*value.Value)
		if nameLength > longName {
			longName = nameLength
		}
		if valueLength > longVal {
			longVal = valueLength
		}
	}
	return longName, longVal
}

func init() {
	serviceCmd.AddCommand(serviceEnvsCmd)
}
