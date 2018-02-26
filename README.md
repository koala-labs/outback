# UFO CLI

### Install Prerequisites
  1.  `$ brew install go`
  1. Setup your GOPATH
			
		export GOPATH="$HOME/go"

  1. Add go compiled binaries to your PATH

		export PATH=$GOPATH/bin:$PATH 

### Installing UFO
1. Install the UFO binary `go get gitlab.fuzzhq.com/Web-Ops/ufo/...`
    * If you have issues pulling a private repository, see https://gist.github.com/shurcooL/6927554

## Usage

### Configuration

On first run, if there is not a `.ufo/config.json` config present in your current working directory, UFO will create this config for you with the default configuration below.

```
{
	"profile": "default",
	"region": "us-east-1",
	"repo": "default.dkr.ecr.us-west-1.amazonaws.com/default",
	"clusters": [
		{
			"name": "dev",
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

### Commands

- ufo deploy
- ufo service
- ufo task

#### Global Flags

| Flag | Shorthand | Default | Description |
| --- | --- | --- | --- |
| --cluster | -c | dev | ECS cluster name |
| --service | -s | api | ECS service name |

#### Deployments

A deployment consists of 5 steps necessary to update an AWS ECS Service.

1. It builds a docker image
2. It tags a docker image with the current short git commit hash
3. It pushes a docker image to AWS ECR
4. It creates a new task definition revision, only replacing its image with the newly tagged one
5. It updates a service on ecs to use the newly created task definition

- [deploy](#ufo-deploy)

##### ufo deploy

```console
ufo deploy --cluster --verbose --login
```

Create a deployment

A cluster must be specified via the --cluster flag. The --verbose flag can be input to enable verbose output. The --login flag can be input to login to AWS ECR.

#### Services

Services manage long-lived instances of your containers that are run on AWS
ECS. If your container exits for any reason, the service scheduler will
restart your containers and ensure your service has the desired number of
tasks running. Services can be used in concert with a load balancer to
distribute traffic amongst the tasks in your service.

- [env add](#ufo-service-env-add)
- [env rm](#ufo-service-env-rm)
- [env list](#ufo-service-env-list)

##### ufo service env add

```console
ufo service env add --env <key=value>
```

Add/Update environment variables

At least one environment variable must be specified via the --env flag. Specify
--env with a key=value parameter multiple times to add multiple variables.

##### ufo service env rm

```console
ufo service env rm --key <key-name>
```

Remove environment variables

Removes the environment variable specified via the --key flag. Specify --key with
a key name multiple times to unset multiple variables.


##### ufo service env list

```console
ufo service env list
```

List environment variables

#### Tasks

Tasks are one-time executions of your container. Instances of your task are run
until you manually stop them either through AWS APIs, the AWS Management
Console, or until they are interrupted for any reason.

- [run](#ufo-task-run)

##### ufo task run

```console
ufo task run --command "<command>"
```

Run a one off tasks

You must specify a cluster, service, and command to run. The command will use the image described in the task definition for the service that is specified. When specifying a command, the task definitions current command will be overriden with the one specified. 

There is also an option of creating command aliases in `.ufo/config.json`. Once a command alias is in the ufo config, specifying that alias via the --command flag will run the configured command.

If the awslogs driver is configured for the service in which you base your task. Logs for that task will be sent to cloudwatch under the same log group and prefix as described in the task definition.

