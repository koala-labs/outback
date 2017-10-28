package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"fmt"
)

const UFO_CONFIG = ".ufo/config.json"

func main() {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT)

	go func() {
		<-sigs
	}()

	deployCommand := flag.NewFlagSet("deploy", flag.ExitOnError)
	deployConfig := deployCommand.String("c", UFO_CONFIG, "Path to ufo config.json, ./.ufo/config.json by default.")
	deployVerbose := deployCommand.Bool("v", false, "Verbose.")

	// @todo support
	initCommand := flag.NewFlagSet("init", flag.ExitOnError)
	initLocation := initCommand.String("p", UFO_DIR, "Path to create UFO config directory.")
	listCommand := flag.NewFlagSet("list", flag.ExitOnError)
	listConfig := listCommand.String("c", UFO_CONFIG, "Path to ufo config.json, ./.ufo/config.json by default.")
	//useCommand := flag.NewFlagSet("use", flag.ExitOnError)

	commands := map[string]*flag.FlagSet{
		"deploy": deployCommand,
		"init": initCommand,
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

		RunDeployCmd(LoadConfigFromFile(*deployConfig), *deployVerbose)
		// foo
	case "init":
		RunInitCommand(*initLocation)
	case "list":
		RunListCommand(LoadConfigFromFile(*listConfig))
	case "use":
		fallthrough
	default:
		log.Fatalln("Not supported yet.")
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

	log.Printf("\nEncountered an error: %s\n", err.Error())

	os.Exit(1)
}