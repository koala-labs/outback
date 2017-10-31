package main

import (
	"os"
	"fmt"
)

const DEFAULT_CONFIG = `{
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

const UFO_DIR = ".ufo/"
const UFO_FILE = "config.json"

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
type osFS struct {}

func (osFS) Open(name string) (*os.File, error)        { return os.Open(name) }
func (osFS) Stat(name string) (os.FileInfo, error) { return os.Stat(name) }
func (osFS) Mkdir(name string, perm os.FileMode) error { return os.Mkdir(name, perm) }
func (osFS) IsNotExist(err error) bool { return os.IsNotExist(err) }
func (osFS) Create(name string) (*os.File, error) { return os.Create(name) }

func RunInitCommand(path string, fs fileSystem) error {
	if _, err := fs.Stat(path); fs.IsNotExist(err) {
		fmt.Printf("Creating directory %s\n", path)
		fs.Mkdir(UFO_DIR, 755)
	}

	if _, err := fs.Stat(UFO_CONFIG); ! fs.IsNotExist(err) {
		return ErrConfigFileAlreadyExists
	}

	fmt.Printf("Creating config file %s.\n", UFO_FILE)
	f, err := fs.Create(UFO_CONFIG)

	if err != nil {
		return ErrCouldNotCreateConfig
	}

	defer f.Close()

	fmt.Println("Writing default config to config file.")
	fmt.Fprint(f, DEFAULT_CONFIG)

	return nil
}
