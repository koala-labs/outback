package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/spf13/cobra"
	UFO "gitlab.fuzzhq.com/Web-Ops/ufo/ufo"
)

var (
	flagServiceEnvsCluster string
	flagServiceEnvsService string
)

var serviceEnvsCmd = &cobra.Command{
	Use:   "envs",
	Short: "List envs for services in your cluster",
	Long:  `envs`,
	Run:   envsRun,
}

func envsRun(cmd *cobra.Command, args []string) {
	e, err := getSelectedEnv(flagServiceEnvsCluster)

	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	s, err := getSelectedService(e.Services, flagServiceEnvsService)

	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	printEnvsForService(e.Cluster, *s)
}

func printEnvsForService(cluster string, service string) {
	ufo := UFO.New(ufoCfg)

	c, err := ufo.GetCluster(cluster)

	if err != nil {
		fmt.Printf("Erorr: %v", err)
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	s, err := ufo.GetService(c, service)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	t, err := ufo.GetTaskDefinition(c, s)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

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

func longNameAndValue(e []*ecs.KeyValuePair) (longName int, longVal int) {
	for _, value := range e {
		name := *value.Name
		value := *value.Value
		nameLength := len(name)
		valueLength := len(value)
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

	serviceEnvsCmd.Flags().StringVarP(&flagServiceEnvsCluster, "cluster", "c", "", "Cluster where your services are running")
	serviceEnvsCmd.Flags().StringVarP(&flagServiceEnvsService, "service", "s", "", "Service to list envs for")
}
