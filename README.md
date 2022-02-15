# RISKEN Core

![Build Status](https://codebuild.ap-northeast-1.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoicUorMEd6RFhSZ28yb0czN0ZqQkV4eFdraU9VUE4rZUpsU2tndVQxeXNuUGE3K3UvTEJTN3FtWDdQTkp0SWMxVDRLWFQwRXlCWlQ4RnYvVS85dEdzd1pRPSIsIml2UGFyYW1ldGVyU3BlYyI6IlZnZkIwc3BnS0szd2lVazgiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=master)

`RISKEN` is a monitoring tool for your cloud platforms, web-site, source-code... 
`RISKEN Core` is the core-system that analyzes, searches, and alerts on discovered threat information.

Please check [RISKEN Documentation](https://docs.security-hub.jp/).

## Installation

### Requirements

This module requires the following modules:

- [Go](https://go.dev/doc/install)
- [Docker](https://docs.docker.com/get-docker/)
- [Protocol Buffer](https://grpc.io/docs/protoc-installation/)

### Install packages

This module is developed in the `Go language`, please run the following command after installing the `Go`.

```bash
$ make install
```

### Building

Build the containers on your machine with the following command

```bash
$ make build
```

### Running Apps

Deploy the pre-built containers to the Kubernetes environment on your local machine.

- Follow the [documentation](https://docs.security-hub.jp/admin/infra_local/#risken) to download the Kubernetes manifest sample.
- Fix the Kubernetes object specs of the manifest file as follows and deploy it.

`k8s-sample/overlays/local/core.yaml`

| service | spec                                | before (public images)                      | after (pre-build images on your machine) |
| ------- | ----------------------------------- | ------------------------------------------- | ---------------------------------------- |
| alert   | spec.template.spec.containers.image | `public.ecr.aws/risken/core/alert:latest`   | `core/alert:latest`                      |
| finding | spec.template.spec.containers.image | `public.ecr.aws/risken/core/finding:latest` | `core/finding:latest`                    |
| iam     | spec.template.spec.containers.image | `public.ecr.aws/risken/core/iam:latest`     | `core/iam:latest`                        |
| project | spec.template.spec.containers.image | `public.ecr.aws/risken/core/project:latest` | `core/project:latest`                    |
| report  | spec.template.spec.containers.image | `public.ecr.aws/risken/core/report:latest`  | `core/report:latest`                     |

## Community

Info on reporting bugs, getting help, finding roadmaps,
and more can be found in the [RISKEN Community](https://github.com/ca-risken/community).

## License

[MIT](LICENSE).
