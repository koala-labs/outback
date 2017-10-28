package main

import (
	"os"
	"fmt"
)

const DEFAULT_CONFIG = `{
	"profile": "fooProfile",
	"image_repository_url": "foo.dkr.ecr.us-west-1.amazonaws.com/fooRepo",
	"environments": [
		{
			"branch": "dev",
			"region": "us-west-1",
			"cluster": "api-dev",
            "service": "api",
            "dockerfile": "Dockerfile.local"
		}
	]
}
`

const UFO_DIR = ".ufo/"
const UFO_FILE = "config.json"

func RunInitCommand(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("Creating directory %s\n", path)
		os.Mkdir(UFO_DIR, 755)
	}

	// @todo if file not exists

	if _, err := os.Stat(UFO_CONFIG); ! os.IsNotExist(err) {
		HandleError(ErrConfigFileAlreadyExists)
	}

	fmt.Printf("Creating config file %s.\n", UFO_FILE)
	f, err := os.Create(UFO_CONFIG)

	if err != nil {
		HandleError(ErrCouldNotCreateConfig)
	}

	defer f.Close()

	fmt.Println("Writing default config to config file.")
	fmt.Fprint(f, DEFAULT_CONFIG)
}
