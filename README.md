# Outback CLI 🦘

## About The Project

Outback is CLI tool written in Go to help streamline the process of deploying containerized applications to AWS Elastic Container Service.

Outback automates the process of building a Docker image, pushing the image to AWS's container registry, creating an updated ECS Task Definition pointing to the new image, and finally updating the ECS service to use the new image.

## Getting Started

### Prerequisites

1. Install [Go](https://golang.org/doc/install)

   `brew install go`

2. Setup make sure that [Go compiled binaries are included in your $PATH](https://golang.org/doc/tutorial/compile-install)

3. Install [Docker](https://docs.docker.com/install/)

4. Install [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html)

5. Create an [AWS access key](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_access-keys.html#Using_CreateAccessKey)

6. Add AWS access key to the file `~/.aws/credentials` to match the profile in the config of Outback `.outback/config.json`.
   ```
   [default]
   aws_access_key_id = KEYVALUE
   aws_secret_access_key = ACCESSKEYVALUE
   ```

### Installing Outback

1. Install the Outback binary `go install github.com/koala-labs/outback@latest`

## Usage

### Configuration

On first run, if there is not a `.outback/config.json` config present in your current working directory, Outback will create this config for you with the default configuration below.

```json
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

Outback will relies on this config to run its operations.

### Commands

- outback deploy
- outback service
- outback task
- outback rollback

#### Global Flags

| Flag      | Shorthand | Default | Description      |
| --------- | --------- | ------- | ---------------- |
| --cluster | -c        | dev     | ECS cluster name |
| --service | -s        | api     | ECS service name |

#### Deployments

A deployment consists of 5 steps necessary to update an AWS ECS Service.

1. It builds a docker image
2. It tags a docker image with the current short git commit hash
3. It pushes a docker image to AWS ECR
4. It creates a new task definition revision, only replacing its image with the newly tagged one
5. It updates a service on ecs to use the newly created task definition

- [deploy](#outback-deploy)

##### `outback deploy`

```console
outback deploy --cluster --verbose --login
```

Create a deployment

A cluster must be specified via the --cluster flag. The --verbose flag can be input to enable verbose output. The --login flag can be input to login to AWS ECR.

Docker build arguments

Outback can use `--build-arg` or `-b` to pass arguments during the docker build phase. Multiple build arguments can be passed, see example below.

```console
outback deploy --cluster dev --build-arg NODE_ENV=dev --build-arg CAT=lazy
```

Docker build arguments can also be passed though the `.outback/config.json` and coexist with the `--build-arg` command option.

```json
{
  "profile": "default",
  "region": "us-east-1",
  "repo": "default.dkr.ecr.us-west-1.amazonaws.com/default",
  "clusters": [
    {
      "name": "dev",
      "services": ["api"],
      "dockerfile": "Dockerfile",
      "build-args": ["NODE_ENV=dev", "CAT=lazy"]
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

#### Services

Services manage long-lived instances of your containers that are run on AWS
ECS. If your container exits for any reason, the service scheduler will
restart your containers and ensure your service has the desired number of
tasks running. Services can be used in concert with a load balancer to
distribute traffic amongst the tasks in your service.

- [env add](#outback-service-env-add)
- [env rm](#outback-service-env-rm)
- [env list](#outback-service-env-list)

##### `outback service env add`

```console
outback service env add --env <key=value>
```

Add/Update environment variables

At least one environment variable must be specified via the --env flag. Specify
--env with a key=value parameter multiple times to add multiple variables.

##### `outback service env rm`

```console
outback service env rm --key <key-name>
```

Remove environment variables

Removes the environment variable specified via the --key flag. Specify --key with
a key name multiple times to unset multiple variables.

##### `outback service env list`

```console
outback service env list
```

List environment variables

##### `outback service list`

```console
outback service info --cluster dev --service frontend
```

To list the status of service on a cluster, in this example: cluster dev, frontend service.

#### Tasks

Tasks are one-time executions of your container. Instances of your task are run
until you manually stop them either through AWS APIs, the AWS Management
Console, or until they are interrupted for any reason.

- [run](#outback-task-run)

##### `outback task run`

```console
outback task run --command "<command>"
```

Run a one off tasks

You must specify a cluster, service, and command to run. The command will use the image described in the task definition for the service that is specified. When specifying a command, the task definitions current command will be overridden with the one specified.

There is also an option of creating command aliases in `.outback/config.json`. Once a command alias is in the outback config, specifying that alias via the --command flag will run the configured command.

If the awslogs driver is configured for the service in which you base your task. Logs for that task will be sent to cloudwatch under the same log group and prefix as described in the task definition.

##### `outback rollback`

The rollback option will update the ECS service revision number to the desired task number. If the need is to rollback to the previous deploy, use:

```console
outback rollback --cluster dev
```

Revision Number

Rollback can use `--revision` or `-r` to pass the revision number that is desired for the ECS service to run:

```console
outback rollback --cluster dev --revision 123
```

## Tests

Use the following command to run the tests and output function-level code coverage

```sh
go test ./... -coverprofile coverage.txt && go tool cover -func coverage.txt
```

## Contributing

Contributions are what make the open source community such an amazing place to be learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Run the test suite (`go test ./...`)
4. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
5. Push to the Branch (`git push origin feature/AmazingFeature`)
6. Open a Pull Request

## License

Distributed under the MIT License. See `LICENSE` for more information.

## Contact

Koala Labs - [@koala_labs](https://twitter.com/koala_labs) - engineering@koala.io
