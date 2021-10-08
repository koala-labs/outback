package cmd

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/ecs"
	Outback "github.com/koala-labs/outback/pkg/outback"
	"github.com/spf13/cobra"
)

var serviceListEnvCmd = &cobra.Command{
	Use:   "list",
	Short: "List environment variables",
	Run:   listEnv,
}

func listEnv(cmd *cobra.Command, args []string) {
	cfgCluster, err := cfg.getCluster(flagCluster)

	handleError(err)

	cfgService, err := cfg.getService(cfgCluster.Services, flagService)

	handleError(err)

	outback := Outback.New(awsConfig)

	c, err := outback.GetCluster(cfgCluster.Name)

	handleError(err)

	s, err := outback.GetService(c, *cfgService)

	handleError(err)

	t, err := outback.GetTaskDefinition(c, s)

	handleError(err)

	printEnvTable(t)
}

func printEnvTable(t *ecs.TaskDefinition) {
	for _, containerDefinition := range t.ContainerDefinitions {
		longestName, longestValue := longestNameAndValue(containerDefinition.Environment)
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

// longestNameAndValue returns the length of the longest Name and Value
func longestNameAndValue(e []*ecs.KeyValuePair) (longName int, longVal int) {
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
	serviceEnvCmd.AddCommand(serviceListEnvCmd)
}
