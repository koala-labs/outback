package main

import "fmt"

func RunListCommand(c *Config) {
	PrintConfigInfo(c)
	PrintInfoForAllEnvironments(c)
}

func PrintConfigInfo(c *Config) {
	fmt.Printf("AWS Profile:            %s\n", c.Profile)
	fmt.Printf("ECR Repository:         %s\n", c.ImageRepositoryUrl)
	fmt.Printf("Number of environments: %d\n", len(c.Env))
}

func PrintInfoForAllEnvironments(c *Config) {
	for _, env := range c.Env {
		fmt.Println("================================")
		PrintInfoForEnvironment(env)
		fmt.Println("================================\n")
	}
}

func PrintInfoForEnvironment(e *Environment) {
	fmt.Printf("Branch:     %s\n", e.Branch)
	fmt.Printf("Region:     %s\n", e.Region)
	fmt.Printf("Cluster:    %s\n", e.Cluster)
	fmt.Printf("Service:    %s\n", e.Service)
	fmt.Printf("Dockerfile: %s\n", e.Dockerfile)
}