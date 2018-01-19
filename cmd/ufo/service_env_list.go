package cmd

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/spf13/cobra"
	UFO "gitlab.fuzzhq.com/Web-Ops/ufo/pkg/ufo"
)

var serviceEnvListCmd = &cobra.Command{
	Use:   "list",
	Short: "List environment variables",
	Run:   listServiceEnv,
}

func listServiceEnv(cmd *cobra.Command, args []string) {
	cfgCluster, err := cfg.getCluster(flagCluster)

	handleError(err)

	cfgService, err := cfg.getService(cfgCluster.Services, flagService)

	handleError(err)

	ufo := UFO.New(ufoCfg)

	c, err := ufo.GetCluster(cfgCluster.Name)

	handleError(err)

	s, err := ufo.GetService(c, *cfgService)

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
	serviceEnvCmd.AddCommand(serviceEnvListCmd)
}