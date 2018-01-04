package git

import (
	"os/exec"
	"strings"
)

func GetCurrentHead() (string, error) {
	command := exec.Command("git", "rev-parse", "--short", "HEAD")

	r, err := command.Output()

	if err != nil {
		return "", ErrGitError
	}

	return strings.Trim(string(r), "\n"), nil
}

func GetCurrentBranch() (string, error) {
	command := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")

	r, err := command.Output()

	if err != nil {
		return "", ErrGitError
	}

	return strings.Trim(string(r), "\n"), nil
}
