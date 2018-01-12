package git

import (
	"errors"
)

var (
	ErrGitError = errors.New("Could not read git information. Please make sure you have git installed and are in a git repository")
)
