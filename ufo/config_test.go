package main

import (
	"errors"
	"fmt"
	"testing"
)

func TestItCanLoadConfigFromBytes(t *testing.T) {
	json := []byte(`{
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
}`)
	_, err := LoadConfig(json)

	if err != nil {
		t.Fatal(err)
	}
}

func TestItErrorsIfJsonInvalid(t *testing.T) {
	json := []byte(`{
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
	],
}`)
	_, err := LoadConfig(json)

	if err == nil {
		t.Fatal(errors.New("Expected failure case due to invalid json"))
	}
}

func TestItReturnsErrIfEnvironmentNotFoundForBranch(t *testing.T) {
	json := []byte(`{
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
}`)
	c, err := LoadConfig(json)

	if err != nil {
		t.Fatal(err)
	}

	_, err = c.GetEnvironmentByBranch("not_exists")

	if err == nil {
		t.Fatal(errors.New("Expected failure case due to not found branch"))
	}
}

func TestItCanFindConfigForBranch(t *testing.T) {
	json := []byte(`{
	"profile": "fooProfile",
	"image_repository_url": "foo.dkr.ecr.us-west-1.amazonaws.com/fooRepo",
	"environments": [
		{
			"branch": "dev",
			"region": "us-west-1",
			"cluster": "api-dev",
			"service": "api",
			"dockerfile": "Dockerfile.local"
		},
		{
			"branch": "staging",
			"region": "us-west-1",
			"cluster": "api-dev",
			"service": "api",
			"dockerfile": "Dockerfile.local"
		},
		{
			"branch": "production",
			"region": "us-west-1",
			"cluster": "api-dev",
			"service": "api",
			"dockerfile": "Dockerfile.local"
		}
	]
}`)
	c, err := LoadConfig(json)

	if err != nil {
		t.Fatal(err)
	}

	e, err := c.GetEnvironmentByBranch("staging")

	if err != nil {
		t.Fatal(err)
	}

	if e.Branch != "staging" {
		t.Fatalf("Expected %s, got %s", "staging", c.Env[0].Branch)
	}
}

func TestItRequiresProfile(t *testing.T) {
	require := "profile"

	json := []byte(`{
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
}`)

	c, err := LoadConfig(json)

	if err != nil {
		t.Fatal(err)
	}

	err = c.validate()

	if err.Error() != fmt.Sprintf("Missing required attribute: %s", require) {
		t.Fatalf("Failed to require %s.", require)
	}
}

func TestItRequiresImageRepo(t *testing.T) {
	require := "image_repository_url"

	json := []byte(`{
	"profile": "fooProfile",
	"environments": [
		{
			"branch": "dev",
			"region": "us-west-1",
			"cluster": "api-dev",
			"service": "api",
			"dockerfile": "Dockerfile.local"
		}
	]
}`)

	c, err := LoadConfig(json)

	if err != nil {
		t.Fatal(err)
	}

	err = c.validate()

	if err.Error() != fmt.Sprintf("Missing required attribute: %s", require) {
		t.Fatalf("Failed to require %s.", require)
	}
}

func TestItRequiresBranch(t *testing.T) {
	require := "branch"

	json := []byte(`{
	"profile": "fooProfile",
	"image_repository_url": "foo.dkr.ecr.us-west-1.amazonaws.com/fooRepo",
	"environments": [
		{
			"region": "us-west-1",
			"cluster": "api-dev",
			"service": "api",
			"dockerfile": "Dockerfile.local"
		}
	]
}`)

	c, err := LoadConfig(json)

	if err != nil {
		t.Fatal(err)
	}

	err = c.validate()

	if err.Error() != fmt.Sprintf("Missing required attribute %s under environment ", require) {
		t.Fatalf("Failed to require %s.", require)
	}
}

func TestItRequiresRegion(t *testing.T) {
	require := "region"

	json := []byte(`{
	"profile": "fooProfile",
	"image_repository_url": "foo.dkr.ecr.us-west-1.amazonaws.com/fooRepo",
	"environments": [
		{
			"branch": "dev",
			"cluster": "api-dev",
			"service": "api",
			"dockerfile": "Dockerfile.local"
		}
	]
}`)

	c, err := LoadConfig(json)

	if err != nil {
		t.Fatal(err)
	}

	err = c.validate()

	if err.Error() != fmt.Sprintf("Missing required attribute %s under environment dev", require) {
		t.Fatalf("Failed to require %s.", require)
	}
}

func TestItRequiresCluster(t *testing.T) {
	require := "cluster"

	json := []byte(`{
	"profile": "fooProfile",
	"image_repository_url": "foo.dkr.ecr.us-west-1.amazonaws.com/fooRepo",
	"environments": [
		{
			"branch": "dev",
			"region": "us-west-1",
			"service": "api",
			"dockerfile": "Dockerfile.local"
		}
	]
}`)

	c, err := LoadConfig(json)

	if err != nil {
		t.Fatal(err)
	}

	err = c.validate()

	if err.Error() != fmt.Sprintf("Missing required attribute %s under environment dev", require) {
		t.Fatalf("Failed to require %s.", require)
	}
}

func TestItRequiresService(t *testing.T) {
	require := "service"

	json := []byte(`{
	"profile": "fooProfile",
	"image_repository_url": "foo.dkr.ecr.us-west-1.amazonaws.com/fooRepo",
	"environments": [
		{
			"branch": "dev",
			"region": "us-west-1",
			"cluster": "api-dev",
			"dockerfile": "Dockerfile.local"
		}
	]
}`)

	c, err := LoadConfig(json)

	if err != nil {
		t.Fatal(err)
	}

	err = c.validate()

	if err.Error() != fmt.Sprintf("Missing required attribute %s under environment dev", require) {
		t.Fatalf("Failed to require %s.", require)
	}
}

func TestItDefaultsDockerfile(t *testing.T) {
	json := []byte(`{
	"profile": "fooProfile",
	"image_repository_url": "foo.dkr.ecr.us-west-1.amazonaws.com/fooRepo",
	"environments": [
		{
			"branch": "dev",
			"region": "us-west-1",
			"cluster": "api-dev",
			"service": "api"
		}
	]
}`)

	c, err := LoadConfig(json)

	if err != nil {
		t.Fatal(err)
	}

	err = c.validate()

	if err != nil {
		t.Fatal(err)
	}

	if c.Env[0].Dockerfile != "Dockerfile" {
		t.Fatal("Failed to default Dockerfile")
	}
}

func TestItRequiresAtLeastOneEnvironment(t *testing.T) {
	require := "environments"

	json := []byte(`{
	"profile": "fooProfile",
	"image_repository_url": "foo.dkr.ecr.us-west-1.amazonaws.com/fooRepo",
	"environments": []
}`)

	c, err := LoadConfig(json)

	if err != nil {
		t.Fatal(err)
	}

	err = c.validate()

	if err != ErrNoEnvironments {
		t.Fatalf("Failed to require %s", require)
	}
}
