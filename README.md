# UFO CLI

## Build Prerequisites
  - `$ brew install go`
  - Setup GOPATH - defaults to `$HOME/go`
  - Install [glide](https://github.com/Masterminds/glide)

## Installing UFO
1. Install go `brew install go`
1. Install the UFO binary `go install gitlab.fuzzhq.com/Web-Ops/ufo/...`
    * If you have issues pulling a private repository, see https://gist.github.com/shurcooL/6927554

## Configuration

On first run, if there is not a `.ufo/config.json` config present in your current working directory, UFO will create this config for you with the default configuration below.

```
{
	"profile": "default",
	"region": "us-east-1",
	"repo": "default.dkr.ecr.us-west-1.amazonaws.com/default",
	"clusters": [
		{
			"name": "dev",
			"branch": "dev",
			"services": ["api"],
			"dockerfile": "Dockerfile"
		}
	],
	"tasks": [
		{
			"name": "migrate",
			"command": "php artisan migrate"
		}
	]
}
```

UFO will relies on this config to run its operations.

## Commands

- ufo service
- ufo task
- ufo deploy
