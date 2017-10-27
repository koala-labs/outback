package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
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
	//initCommand := flag.NewFlagSet("init", flag.ExitOnError)
	//useCommand := flag.NewFlagSet("use", flag.ExitOnError)
	//listCommand := flag.NewFlagSet("list", flag.ExitOnError)

	if len(os.Args) < 2 {
		log.Println("A subcommand is required.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	subCommand := os.Args[1]

	switch subCommand {
	case "deploy":
		deployCommand.Parse(os.Args[2:])

		RunDeployCmd(LoadConfigFromFile(*deployConfig), *deployVerbose)
		// foo
	case "init":
		fallthrough
	case "use":
		fallthrough
	case "list":
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