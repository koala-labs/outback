package main

import (
	"os"
	"fmt"
)

const DEFAULT_CONFIG = `
{
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
		os.Mkdir(UFO_DIR, 755)
	}

	f, err := os.Create(UFO_DIR + UFO_FILE)

	if err != nil {
		HandleError(ErrCouldNotCreateConfig)
	}

	defer f.Close()

	fmt.Fprint(f, DEFAULT_CONFIG)
}
