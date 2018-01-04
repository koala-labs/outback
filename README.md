# UFO (Universal Fuzz Orchestrator)

## Build Prerequisites
  - `$ brew install go`
  - Setup GOPATH - defaults to `$HOME/go`
  - Install [glide](https://github.com/Masterminds/glide)

## Build CLI

1. `cd ufo`
1. `go build -o ufo`
1. `./ufo <command> <args>`

## Installing CLI
1. Install go `brew install go`
1. Install the UFO binary `go install gitlab.fuzzhq.com/Web-Ops/ufo/...`
    * If you have issues pulling a private repository, see https://gist.github.com/shurcooL/6927554
