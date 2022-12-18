# Terraform Provider for Proofpoint - Meta Networks

See the [Proofpoint - Meta Networks Provider documentation](docs/index.md) to get started using the provider.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.12.x or higher
- [Go](https://golang.org/doc/install) >= 1.18

## Developing the Provider

To compile and run the provider locally run:
```sh
$ make debug-build
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

In order to run unit-tests only run
```shell
$ make unittest
```

