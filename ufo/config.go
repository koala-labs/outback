package main

import (
	"encoding/json"
	"io/ioutil"
	"errors"
	"fmt"
)

const EMPTY_VALUE = ""

type Environment struct {
	Branch string `json:"branch"`
	Region string `json:"region"`
	Cluster string `json:"cluster"`
	Service string `json:"service"`
	Dockerfile string `json:"dockerfile"`
}

type Config struct {
	Profile string
	ImageRepositoryUrl string `json:"image_repository_url"`
	Env []*Environment `json:"environments"`
}

func LoadConfigFromFile(path string) *Config {
	dat, err := ioutil.ReadFile(path)

	if err != nil {
		HandleError(ErrCouldNotLoadConfig)
	}

	c := LoadConfig(dat)

	err = c.validate()
	HandleError(err)

	return c
}

func LoadConfig(config []byte) *Config {
	c := &Config{}

	err := json.Unmarshal(config, c)

	if err != nil {
		HandleError(ErrCouldNotLoadConfig)
	}

	return c
}

func (c *Config) GetEnvironmentByBranch(branch string) (*Environment, error) {
	for _, env := range c.Env {
		if env.Branch == branch {
			return env, nil
		}
	}

	return nil, ErrEnvironmentForBranchDoesNotExist
}

func (c *Config) validate() error {
	req := map[string]string{
		"profile": c.Profile,
		"image_repository_url": c.ImageRepositoryUrl,
	}

	if len(c.Env) < 1 {
		return ErrNoEnvironments
	}

	for key, val := range req {
		if val == EMPTY_VALUE {
			return errors.New(fmt.Sprintf("Missing required attribute: %s.", key))
		}
	}

	for _, env := range c.Env {
		envReqs := map[string]string{
			"branch": env.Branch,
			"region": env.Region,
			"cluster": env.Cluster,
			"service": env.Service,
		}

		for k, v := range envReqs {
			if v == EMPTY_VALUE {
				return errors.New(fmt.Sprintf("Missing required attribute %s under environment %s.", k, env.Branch))
			}
		}

		if env.Dockerfile == EMPTY_VALUE {
			env.Dockerfile = "Dockerfile"
		}
	}

	return nil
}