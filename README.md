# UFO (Universal Fuzz Orchestrator)

## Prerequisites
  - `brew install go`
  - Setup GOPATH - defaults to `$HOME/go`
  - Install [Godep](https://github.com/tools/godep)

## Build & Run

1. Build static assets - cd app && yarn run build
1. Build go binary - cd backend && go build -o ufo
1. This will read your default profile in `~/.aws/credentials`
1. Run app - `./backend/ufo`
