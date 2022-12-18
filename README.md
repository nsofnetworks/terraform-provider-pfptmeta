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

### Creating a resource and data-source

#### Client package:

In the [client](internal%2Fclient) package, add the module with the following:
- Private endpoint for the resource(i.e `const accessBridgeEndpoint = "v1/access_bridges"`)
- Public struct that defines the resource (should be similar to how it's represented in the mgmt API)
- Public constructor for that struct.
- Public CRUD functions for the resource.
- **Note**:all CRUD functions **must** take a [context](https://pkg.go.dev/context) object as an argument with which terraform enforces deadlines.


#### Provider folder:

Create a new package named after the created resource.
This package should contain three files:
- `resource.go` - contains the resource's schema definition
- `data_source.go` - contains the data source's definition
- `common.go` - contains common descriptions and functions for both resource and data source.

#### provider.go file:

Adding the resource and data-source to the [`provider.go`](internal%2Fprovider%2Fprovider.go) file to the resource and data source mapping.

#### Acceptance tests:

Under [acc_tests](internal%2Fprovider%2Facc_tests) folder and a file named `<resource>_test.go`.
In that test make sure you're covering all CRUD operations of the resource, and be sure to sanity check the data-source as well.
In order to run the tests you can either trigger a jenkins build, or run it locally using:
```shell
make acc_tests TEST_KW=<test_pattern>
```
Note: In order to run it locally you'll have to configure credentials for the TF provider. See [index.md](docs%2Findex.md) for details.

Note [2]: The default base url is configured to be `https://api.access.proofpoint.com`, if you want to run acceptance tests with other base url use the `PFPTMETA_BASE_URL` envar.
#### Documentation:

The documentation is auto generated using [this](github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs) plugin.
All the files in the [docs](docs) folder are auto generated from the following:
- resource and data source `Description` attribute - This comes from the schema (make sure you document both resource itself and its attributes)
- terraform examples in [examples](examples) folder - Add a resource and data-source folders with valid examples (see other resources for reference).
- Templates - The templates are important mostly for the subcategory attribute that's decides the hierarchy of the resource in the terraform doc page. Make sure you add a resource and data source templates in the [templates](templates) folder.

#### Generate the docs and commit:

Once all the development is done, you need to build the doc fmt your tf and go code.
In order to do that run:
```shell
make verify_clean
```
