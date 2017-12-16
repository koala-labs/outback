package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const defaultConfig = `{
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
}
`

const gitIgnoreConfig = `
/* UFO Config */
.ufo/
`

const UFOConfig = "/.ufo/config.json"
const UFODir = "/.ufo/"
const UFOFile = "/config.json"

var fs fileSystem = osFS{}

type fileSystem interface {
	//Open(name string) (file, error)
	Stat(name string) (os.FileInfo, error)
	Mkdir(name string, perm os.FileMode) error
	IsNotExist(err error) bool
	Create(name string) (*os.File, error)
}

//type file interface {
//	io.Closer
//	io.Reader
//	io.ReaderAt
//	io.Seeker
//	Stat() (os.FileInfo, error)
//}

// osFS implements fileSystem using the local disk.
type osFS struct{}

func (osFS) Open(name string) (*os.File, error)        { return os.Open(name) }
func (osFS) Stat(name string) (os.FileInfo, error)     { return os.Stat(name) }
func (osFS) Mkdir(name string, perm os.FileMode) error { return os.Mkdir(name, perm) }
func (osFS) IsNotExist(err error) bool                 { return os.IsNotExist(err) }
func (osFS) Create(name string) (*os.File, error)      { return os.Create(name) }

func RunInitCommand(path string, fs fileSystem) error {
	createUFODirectory(path+UFODir, fs)

	f, err := createConfigFile(path+UFOConfig, fs)

	if err == nil {
		defer f.Close()

		fmt.Println("Writing default config to config file.")
		fmt.Fprint(f, defaultConfig)
	}

	addUFOToGitignore(path)

	return nil
}

func createUFODirectory(path string, fs fileSystem) {
	if _, err := fs.Stat(path); fs.IsNotExist(err) {
		fmt.Printf("Creating directory %s\n", path)
		fs.Mkdir(path, os.ModePerm)
	}
}

func createConfigFile(path string, fs fileSystem) (*os.File, error) {
	if _, err := fs.Stat(path); !fs.IsNotExist(err) {
		return nil, ErrConfigFileAlreadyExists
	}

	fmt.Printf("Creating config file %s.\n", UFOFile)
	f, err := fs.Create(path)

	if err != nil {
		return nil, ErrCouldNotCreateConfig
	}

	return f, nil
}

func addUFOToGitignore(path string) error {
	gitIgnore := path + "/.gitignore"

	if _, err := fs.Stat(gitIgnore); fs.IsNotExist(err) {
		return ErrNoGitIgnore
	}

	f, err := os.OpenFile(gitIgnore, os.O_APPEND|os.O_WRONLY, 0600)

	if err != nil {
		return err
	}

	file, err := ioutil.ReadFile(gitIgnore)

	if strings.Contains(string(file), gitIgnoreConfig) {
		fmt.Println("UFO .gitignore already set.")
		return nil
	}

	defer f.Close()

	fmt.Println("Adding UFO config to .gitignore.")
	_, err = f.WriteString(gitIgnore)

	return err
}
