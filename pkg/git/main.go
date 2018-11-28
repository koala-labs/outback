package git

import (
	"os/exec"
	"strings"
)

// GetCommit returns the commit hash from HEAD of a git repo
func GetCommit() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")

	r, err := cmd.Output()

	if err != nil {
		return "", ErrGitError
	}

	return strings.Trim(string(r), "\n"), nil
}

// GetBranch returns the branch of a git repo
func GetBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")

	r, err := cmd.Output()

	if err != nil {
		return "", ErrGitError
	}

	return strings.Trim(string(r), "\n"), nil
}
