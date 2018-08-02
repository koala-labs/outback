package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
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

func (c *Config) getConfigs() []string {
	cwd, err := os.Getwd()
	files, err := ioutil.ReadDir(cwd + "/.ufo")
	if err != nil {
		log.Fatal(err)
	}
	var configs []string
	for _, f := range files {
		name := strings.TrimSuffix(f.Name(), ".json")
		configs = append(configs, name)
	}
	return configs
}

func (c *Config) getCluster(cluster string) (*Cluster, error) {
	for _, c := range c.Clusters {
		if c.Name == cluster {
			return c, nil
		}
	}

	return nil, ErrClusterNotFound
}

func (c *Config) getClusters() []string {
	var clusters []string
	for _, cluster := range c.Clusters {
		fmt.Printf(cluster.Name)
		clusters = append(clusters, cluster.Name)
	}
	return clusters
}

func (c *Config) getServices(in string) []string {
	for _, cluster := range c.Clusters {
		if cluster.Name == in {
			return cluster.Services
		}
	}
	return []string{}
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
