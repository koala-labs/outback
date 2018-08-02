package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// Config
const (
	configPath = "/.ufo/config.json"
	configDir  = "/.ufo/"
)

const configTemplate = `{
	"profile": "default",
	"region": "us-east-1",
	"repo": "default.dkr.ecr.us-west-1.amazonaws.com/default",
	"clusters": [
		{
			"name": "dev",
			"services": ["api"],
			"dockerfile": "Dockerfile"
		}
	],
	"tasks": [
		{
			"name": "migrate",
			"command": "php artisan migrate"
		}
	]
}
`

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a UFO config",
	RunE:  runInit,
}

func runInit(cmd *cobra.Command, args []string) error {
	return initConfig()
}

func initConfig() error {
	cwd, err := os.Getwd()

	fmt.Println("Initializing ufo config...")
	createDirectory(filepath.Join(cwd, configDir))

	f, err := createConfig(filepath.Join(cwd, configPath))

	if err != nil {
		return err
	}

	defer f.Close()

	fmt.Fprint(f, configTemplate)

	fmt.Println("ufo config initialized")

	return nil
}

func createDirectory(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("Creating %s\n", path)
		os.Mkdir(path, os.ModePerm)
	}
}

func createConfig(path string) (*os.File, error) {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return nil, ErrConfigFileAlreadyExists
	}

	fmt.Printf("Creating %s\n", path)

	f, err := os.Create(path)

	if err != nil {
		return nil, ErrCouldNotCreateConfig
	}

	return f, nil
}

func init() {
	rootCmd.AddCommand(initCmd)
}
