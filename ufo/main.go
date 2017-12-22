package main

import (
	"flag"
	"os"
	"fmt"
	"os/signal"
	"syscall"
)

const UFO_VERSION = "0.1"

func registerSigHandler() {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT)

	go func() {
		<-sigs
	}()
}

func main() {
	HandleError(AssertCurrentVersion(UFO_VERSION))
	registerSigHandler()

	CWD, err := os.Getwd()

	if err != nil {
		HandleError(ErrNoWorkingDirectory)
	}

	// Deploy command setup
	deployCommand := flag.NewFlagSet("deploy", flag.ExitOnError)
	deployConfig := deployCommand.String("c", CWD+UFOConfig, "Path to ufo config.json, ./.ufo/config.json by default.")
	deployVerbose := deployCommand.Bool("v", false, "Verbose.")
	deployBranch := deployCommand.String("b", EmptyValue, "Branch to deploy.")

	// Init command setup
	initCommand := flag.NewFlagSet("init", flag.ExitOnError)
	initLocation := initCommand.String("p", CWD, "Path to create UFO config directory.")

	// List command setup
	listCommand := flag.NewFlagSet("list", flag.ExitOnError)
	listConfig := listCommand.String("c", CWD+UFOConfig, "Path to ufo config.json, ./.ufo/config.json by default.")
	listEnvs := listCommand.Bool("e", false, "List environment variables")

	// Run task setup
	runTaskCommand := flag.NewFlagSet("run-task", flag.ExitOnError)
	runTaskCommandOverride := runTaskCommand.String("o", "echo", "Command to run as a one off task. e.g. -o 'echo foo'")
	runTaskCommandName := runTaskCommand.String("n", EmptyValue, "Run-Task command name to run.")
	runTaskBranch := runTaskCommand.String("b", EmptyValue, "Branch respresentative of environment to run task on")
	runTaskConfig := runTaskCommand.String("c", CWD+UFOConfig, "Path to ufo config.json, ./.ufo/config.json by default.")

	commands := map[string]*flag.FlagSet{
		"deploy":   deployCommand,
		"init":     initCommand,
		"list":     listCommand,
		"run-task": runTaskCommand,
	}

	if len(os.Args) < 2 {
		fmt.Println("A subcommand is required.")
		fmt.Println("Usage:")
		fmt.Println("ufo command [arg1] [arg2] ...")

		for cmd, set := range commands {
			fmt.Printf("Usage for `%s`:\n", cmd)
			set.PrintDefaults()
			fmt.Println()
		}

		os.Exit(1)
	}

	subCommand := os.Args[1]

	switch subCommand {
	case "deploy":
		deployCommand.Parse(os.Args[2:])

		config, err := LoadConfigFromFile(*deployConfig)

		if err != nil {
			HandleError(err)
		}

		err = RunDeployCmd(config, DeployOptions{
			Verbose:        *deployVerbose,
			OverrideBranch: *deployBranch,
		})

		HandleError(err)
		// foo
	case "run-task":
		runTaskCommand.Parse(os.Args[2:])

		config, err := LoadConfigFromFile(*runTaskConfig)

		if err != nil {
			HandleError(err)
		}

		err = RunTask(config, TaskOptions{
			Command:        *runTaskCommandOverride,
			OverrideBranch: *runTaskBranch,
			CommandName:    *runTaskCommandName,
		})

		HandleError(err)
	case "init":
		HandleError(RunInitCommand(*initLocation, osFS{}))
	case "list":
		listCommand.Parse(os.Args[2:])

		config, err := LoadConfigFromFile(*listConfig)

		if err != nil {
			HandleError(err)
		}

		ListCommand(config, listOptions{
			ListEnvs: *listEnvs,
		})
	case "use":
		fallthrough
	default:
		fmt.Println("Not supported yet.")
		os.Exit(1)
	}
}

// HandleError is intended to be called with which error return to simplify error handling
// Usage:
// foo, err := GetFoo()
// HandleError(err)
// DoSomethingBecauseNoError()
func HandleError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\nEncountered an error: %s\n", err.Error())

	os.Exit(1)
}
