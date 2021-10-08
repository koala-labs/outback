package cmd

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/ecs"
	UFO "github.com/fuzz-productions/ufo/pkg/ufo"
	"github.com/spf13/cobra"
)

var serviceInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "List information of currently deploy service",
	Run:   listInfo,
}

func listInfo(cmd *cobra.Command, args []string) {
	cfgCluster, err := cfg.getCluster(flagCluster)

	handleError(err)

	cfgService, err := cfg.getService(cfgCluster.Services, flagService)

	handleError(err)

	ufo := UFO.New(awsConfig)

	c, err := ufo.GetCluster(cfgCluster.Name)

	handleError(err)

	s, err := ufo.GetService(c, *cfgService)

	handleError(err)

	t, err := ufo.GetTaskDefinition(c, s)

	handleError(err)

	printServiceInfoTable(t)

	fmt.Printf("\n")

	runningTasks, err := ufo.RunningTasks(c, s)

	handleError(err)

	tasks, err := ufo.GetTasks(c, runningTasks)

	handleError(err)

	printRunningTaskTable(tasks)
}

func printServiceInfoTable(t *ecs.TaskDefinition) {
	longestName := 20
	longestValue := 100
	nameDashes := strings.Repeat("-", longestName+2) // Adding two because of the table padding
	valueDashes := strings.Repeat("-", longestValue+2)
	fmt.Printf("Active Task\n")
	fmt.Printf("+%s+%s+\n", nameDashes, valueDashes)

	for _, containerDefinition := range t.ContainerDefinitions {
		imageName := "Image"
		imageValue := *containerDefinition.Image
		imageNameSpaces := longestName - len(imageName)
		imageValueSpaces := longestValue - len(imageValue)
		imageSpacesForName := strings.Repeat(" ", imageNameSpaces)
		imageSpacesForValue := strings.Repeat(" ", imageValueSpaces)
		fmt.Printf("| %s%s | %s%s |\n", imageName, imageSpacesForName, imageValue, imageSpacesForValue)
		fmt.Printf("+%s+%s+\n", nameDashes, valueDashes)
	}

	revisionName := "Revision"
	revisionValue := string(*t.TaskDefinitionArn)
	revisionNameSpaces := longestName - len(revisionName)
	revisionValueSpaces := longestValue - len(revisionValue)
	revisionSpacesForName := strings.Repeat(" ", revisionNameSpaces)
	revisionSpacesForValue := strings.Repeat(" ", revisionValueSpaces)
	fmt.Printf("| %s%s | %s%s |\n", revisionName, revisionSpacesForName, revisionValue, revisionSpacesForValue)
	fmt.Printf("+%s+%s+\n", nameDashes, valueDashes)

	statusName := "Status"
	statusValue := *t.Status
	statusNameSpaces := longestName - len(statusName)
	statusValueSpaces := longestValue - len(statusValue)
	statusSpacesForName := strings.Repeat(" ", statusNameSpaces)
	statusSpacesForValue := strings.Repeat(" ", statusValueSpaces)
	fmt.Printf("| %s%s | %s%s |\n", statusName, statusSpacesForName, statusValue, statusSpacesForValue)
	fmt.Printf("+%s+%s+\n", nameDashes, valueDashes)
}

func printRunningTaskTable(tasks []*ecs.Task) {

	longestName := 20
	longestValue := 100
	nameDashes := strings.Repeat("-", longestName+2) // Adding two because of the table padding
	valueDashes := strings.Repeat("-", longestValue+2)
	fmt.Printf("Running Tasks\n")
	fmt.Printf("+%s+%s+\n", nameDashes, valueDashes)

	for _, task := range tasks {
		imageName := "Task Definition"
		imageValue := *task.TaskDefinitionArn
		imageNameSpaces := longestName - len(imageName)
		imageValueSpaces := longestValue - len(imageValue)
		imageSpacesForName := strings.Repeat(" ", imageNameSpaces)
		imageSpacesForValue := strings.Repeat(" ", imageValueSpaces)
		fmt.Printf("| %s%s | %s%s |\n", imageName, imageSpacesForName, imageValue, imageSpacesForValue)
		fmt.Printf("+%s+%s+\n", nameDashes, valueDashes)
	}
}

func init() {
	serviceCmd.AddCommand(serviceInfoCmd)
}
