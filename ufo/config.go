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
	Service    string `json:"service"`
	Dockerfile string `json:"dockerfile"`
}

type Config struct {
	Profile            string
	ImageRepositoryURL string         `json:"image_repository_url"`
	Env                []*Environment `json:"environments"`
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

func (c *Config) GetEnvironmentByBranch(branch string) (*Environment, error) {
	fmt.Println(branch)
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
			"service": env.Service,
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
