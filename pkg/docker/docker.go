package docker

import (
	"fmt"
	"os/exec"

	"github.com/fuzz-productions/ufo/pkg/term"
)

// ImageBuild builds a docker image based on the configured dockerfile for
// the cluster you are deploying to and tags the image with the vcs head
func ImageBuild(repo string, tag string, dockerfile string, buildArgs []string, configBuildArgs []string) error {
	image := fmt.Sprintf("%s:%s", repo, tag)

	dockerConfigBuildArgs := make([]string, 2*len(configBuildArgs), 2*len(configBuildArgs))
	dockerBuildArgs := make([]string, 2*len(buildArgs), 2*len(buildArgs))

	for i, v := range configBuildArgs {

		dockerConfigBuildArgs[i*2] = "--build-arg"
		dockerConfigBuildArgs[i*2+1] = v

	}

	for i, v := range buildArgs {

		dockerBuildArgs[i*2] = "--build-arg"
		dockerBuildArgs[i*2+1] = v

	}

	dockerCmd := "docker"
	dockerCmdArgs := []string{"build", "-f", dockerfile, "-t", image, "."}
	dockerCmdBuildArgs := append([]string(dockerConfigBuildArgs), []string(dockerBuildArgs)...)
	dockerCmdFullArgs := append([]string(dockerCmdArgs), []string(dockerCmdBuildArgs)...)

	cmd := exec.Command(dockerCmd, dockerCmdFullArgs...)

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
