# UFO (Universal Fuzz Orchestrator)

## Build Prerequisites
  - `$ brew install go`
  - Setup GOPATH - defaults to `$HOME/go`
  - Install [godep](https://github.com/tools/godep)
  - Install [go-bindata](https://github.com/jteeuwen/go-bindata)

## Build Backend

1. Install frontend dependencies - `cd app && yarn install`
1. Build static assets - `$ cd app && yarn run build`
1. Convert assets to go code - `$ cd backend && go-bindata ../app/dist/...`
1. Build go binary - `$ go build -o ufo`
1. Run app - `$ ./backend/ufo -profile=default -region=us-east-1`

## Build CLI

1. `cd ufo`
1. `go build -o ufo`
1. `./ufo <command> <args>`

## Installing CLI
1. Install go `brew install go`
1. Install the UFO binary `go install gitlab.fuzzhq.com/Web-Ops/ufo/ufo`
    * If you have issues pulling a private repository, see https://gist.github.com/shurcooL/6927554
