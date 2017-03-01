# `sfdc`

A generator that can be used by `go generate` to generate a model and helper
functions for a given SFDC account.

### Configuration

Because this package requires credentials to function, you should set the
following environment variables prior to running `go generate`:

* `export SFDC_CLIENTID="YOUR SFDC CLIENT ID"`
* `export SFDC_CLIENTSECRET="YOUR SFDC CLIENT SECRET"`
* `export SFDC_USERNAME="username@domain.com"`
* `export SFDC_PASSWORD="password"`
* `export SFDC_SECURITYTOKEN="YOUR SFDC SECURITY TOKEN"`
* `export SFDC_VERSION="v36.0"` (Optional, default is `v36.0`)
* `export SFDC_ENVIRONMENT="Sandbox"` (Optional, default is `Production`)

### Usage

First, get this package and ensure that your `$GOPATH/bin` is in your `$PATH`:

```shell
go get -u github.com/joefitzgerald/sfdc
```

**Note:** this package requires Go 1.8 or later.

Create a new package (e.g. sfdc), and create a file in the package called `definition.go`
(note that this can have any name you want). In this file, place the following
contents:

```go
package sfdc

//go:generate sfdc -object=Opportunity,OpportunityLineItem,User,Account
```

Then navigate to this package (or any package that contains it) and run `go generate ./...`. Observe new generated files that you can use to interact with the Salesforce REST API :tada: .

### Flags

* `output`: (Required) A comma separated list of Salesforce objects to generate
* `prefix`: (Optional) A prefix to use for each generated file
* `suffix`: (Optional; Default: `_sfdc`) A suffix to use for each generated file
