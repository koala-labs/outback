package docker

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/koala-labs/outback/pkg/term"
)

// ImageBuild builds a docker image based on the configured dockerfile for
// the cluster you are deploying to and tags the image with the vcs head
func ImageBuild(repo string, tag string, dockerfile string, buildArgs []string, configBuildArgs []string, cacheFrom []string) error {
	image := fmt.Sprintf("%s:%s", repo, tag)
	combinedBuildArgs := append(configBuildArgs, buildArgs...)
	// ensure built images are optimized for remote caching
	combinedBuildArgs = append(combinedBuildArgs, "BUILDKIT_INLINE_CACHE=1")

	dockerCmdBuildArgs := make([]string, 2*(len(combinedBuildArgs)), 2*(len(combinedBuildArgs)))
	dockerCacheFromArgs := make([]string, 2*len(cacheFrom), 2*len(cacheFrom))

	for i, v := range combinedBuildArgs {
		dockerCmdBuildArgs[i*2] = "--build-arg"
		dockerCmdBuildArgs[i*2+1] = v
	}

	for i, v := range cacheFrom {
		dockerCacheFromArgs[i*2] = "--cache-from"
		dockerCacheFromArgs[i*2+1] = v
	}

	dockerCmd := "docker"
	dockerCmdArgs := []string{"build", "-f", dockerfile, "-t", image, "."}
	dockerCmdFullArgs := append(dockerCmdArgs, dockerCmdBuildArgs...)
	dockerCmdFullArgs = append(dockerCmdFullArgs, dockerCacheFromArgs...)

	cmd := exec.Command(dockerCmd, dockerCmdFullArgs...)
	// enable BuildKit: https://docs.docker.com/engine/reference/builder/#buildkit
	cmd.Env = append(os.Environ(), "DOCKER_BUILDKIT=1")

	if err := term.PrintStdout(cmd); err != nil {
		return ErrImageBuild
	}

	return nil
}

// ImagePush pushes the image built from buildImage to the configured repository
func ImagePush(repo string, tag string) error {
	image := fmt.Sprintf("%s:%s", repo, tag)

	cmd := exec.Command("docker", "push", image)

	if err := term.PrintStdout(cmd); err != nil {
		return ErrImagePush
	}

	return nil
}

// ImagePull pulls an image with a specific tag from the configured repository
func ImagePull(repo string, tag string) error {
	image := fmt.Sprintf("%s:%s", repo, tag)

	cmd := exec.Command("docker", "pull", image)

	if err := term.PrintStdout(cmd); err != nil {
		return ErrImagePull
	}

	return nil
}
