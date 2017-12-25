package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const EmptyValue = ""

type Environment struct {
	Branch     string `json:"branch"`
	Region     string `json:"region"`
	Cluster    string `json:"cluster"`
	Services   []string `json:"services"`
	Dockerfile string `json:"dockerfile"`
}

type RunTaskConfiguration struct {
	Name    string `json:"name"`
	Command string `json:"command"`
}

type Config struct {
	Profile            string
	ImageRepositoryURL string         `json:"image_repository_url"`
	Env                []*Environment `json:"environments"`
	RunTaskConfigs     []*RunTaskConfiguration `json:"run_tasks"`
}

func LoadConfigFromFile(path string) (*Config, error) {
	dat, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, ErrCouldNotLoadConfig
	}

	c, err := LoadConfig(dat)

	if err != nil {
		return nil, err
	}

	err = c.validate()

	if err != nil {
		return nil, err
	}

	return c, nil
}

func LoadConfig(config []byte) (*Config, error) {
	c := &Config{}

	err := json.Unmarshal(config, c)

	if err != nil {
		return nil, ErrCouldNotLoadConfig
	}

	return c, nil
}

func (c *Config) GetCommandForName(name string) (*RunTaskConfiguration, error) {
	for _, runTaskConfig := range c.RunTaskConfigs {
		if runTaskConfig.Name == name {
			return runTaskConfig, nil
		}
	}

	return nil, fmt.Errorf("Could not find RunTask configuration for: %s", name)
}

func (c *Config) GetEnvironmentByBranch(branch string) (*Environment, error) {
	for _, env := range c.Env {
		if env.Branch == branch {
			return env, nil
		}
	}

	return nil, fmt.Errorf("Could not find environment for chosen branch: %s", branch)
}

func (c *Config) validate() error {
	req := map[string]string{
		"profile":              c.Profile,
		"image_repository_url": c.ImageRepositoryURL,
	}

	if len(c.Env) < 1 {
		return ErrNoEnvironments
	}

	for key, val := range req {
		if val == EmptyValue {
			return fmt.Errorf("Missing required attribute: %s", key)
		}
	}

	for _, env := range c.Env {
		envReqs := map[string]string{
			"branch":  env.Branch,
			"region":  env.Region,
			"cluster": env.Cluster,
		}

		if len(env.Services) < 1 {
			return fmt.Errorf("At least one service is required for environment %s.", env.Branch)
		}

		for k, v := range envReqs {
			if v == EmptyValue {
				return fmt.Errorf("Missing required attribute %s under environment %s", k, env.Branch)
			}
		}

		if env.Dockerfile == EmptyValue {
			env.Dockerfile = "Dockerfile"
		}
	}

	return nil
}
