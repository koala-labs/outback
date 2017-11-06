package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"fmt"
)

func registerSigHandler() {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT)

	go func() {
		<-sigs
	}()
}

func main() {
	registerSigHandler()

	CWD, err := os.Getwd()

	if err != nil {
		HandleError(ErrNoWorkingDirectory)
	}

	// Deploy command setup
	deployCommand := flag.NewFlagSet("deploy", flag.ExitOnError)
	deployConfig := deployCommand.String("c", CWD + UFO_CONFIG, "Path to ufo config.json, ./.ufo/config.json by default.")
	deployVerbose := deployCommand.Bool("v", false, "Verbose.")
	deployBranch := deployCommand.String("b", EMPTY_VALUE, "Branch to deploy.")

	// Init command setup
	initCommand := flag.NewFlagSet("init", flag.ExitOnError)
	initLocation := initCommand.String("p", CWD, "Path to create UFO config directory.")

	// List command setup
	listCommand := flag.NewFlagSet("list", flag.ExitOnError)
	listConfig := listCommand.String("c", CWD + UFO_CONFIG, "Path to ufo config.json, ./.ufo/config.json by default.")

	commands := map[string]*flag.FlagSet{
		"deploy": deployCommand,
		"init":   initCommand,
		"list":   listCommand,
	}

	if len(os.Args) < 2 {
		fmt.Println("A subcommand is required.\n")
		fmt.Println("Usage:\n")
		fmt.Println("ufo command [arg1] [arg2] ...\n")

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
	case "init":
		HandleError(RunInitCommand(*initLocation, osFS{}))
	case "list":
		config, err := LoadConfigFromFile(*listConfig)

		if err != nil {
			HandleError(err)
		}

		RunListCommand(config)
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
