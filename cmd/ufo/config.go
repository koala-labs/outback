package cmd

import (
	"fmt"
	"os"
)

type Config struct {
	Profile  string     `mapstructure:"profile"`
	Region   string     `mapstructure:"region"`
	Repo     string     `mapstructure:"repo"`
	Clusters []*Cluster `mapstructure:"clusters"`
	Tasks    []*Task    `mapstructure:"tasks"`
}

type Cluster struct {
	Name       string   `mapstructure:"name"`
	Services   []string `mapstructure:"services"`
	Dockerfile string   `mapstructure:"dockerfile"`
}

type Task struct {
	Name    string `mapstructure:"name"`
	Command string `mapstructure:"command"`
}

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

const gitIgnoreString = `
# UFO Config
.ufo/
`

func (c *Config) getCluster(cluster string) (*Cluster, error) {
	for _, c := range c.Clusters {
		if c.Name == cluster {
			return c, nil
		}
	}

	return nil, ErrClusterNotFound
}

func (c *Config) getService(services []string, service string) (*string, error) {
	for _, s := range services {
		if s == service {
			return &s, nil
		}
	}

	return nil, ErrServiceNotFound
}

func (c *Config) getCommand(name string) (*string, error) {
	for _, t := range c.Tasks {
		if t.Name == name {
			return &t.Command, nil
		}
	}

	return nil, ErrCommandNotFound
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
