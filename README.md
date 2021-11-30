# Terraform Provider for Proofpoint - Meta Networks

See the [Proofpoint - Meta Networks Provider documentation](docs/index.md) to get started using the provider.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.12.x or higher
- [Go](https://golang.org/doc/install) >= 1.13
- [Goreleaser](https://goreleaser.com/install/)

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go `install` command:
```sh
$ go install
```

## Developing the Provider

### Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up-to-date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```sh
$ go get github.com/author/dependency
$ make mod-tidy
```

Then commit the changes to `go.mod` and `go.sum`.

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile and run the provider locally run:
```sh
$ make debug-mode
```
And then export the printed `TF_REATTACH_PROVIDERS` as env variable.

This will compile the provider in debug mode and will allow you to execute terraform command with it.

### Documentation

In order to generate or update documentation, run `make generate`.

The documentation are auto-generated using the [docs generation tool](github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs)
and are based on the resources descriptions in the schema and on the files in the `examples` folder

### Testing

There are two types of tests that can be executed:
- unit-tests
- acceptance-tests

Acceptance tests are expressed in terms of Test Cases, each using one or more Terraform configurations designed to create a set of resources under test, and then verify the actual infrastructure created.

*Note:* Acceptance tests create real resources, see more about it [here](https://www.terraform.io/docs/extend/testing/acceptance-tests/testcase.html)

In order to run the full suite of Acceptance tests:
- Configure the authentication env variables or configuration file specified in [provider documentation](docs/index.md).
- Execute:
```sh
$ make acc_tests
```

In order to tun unit-tests only run
```shell
$ make unittest
```

