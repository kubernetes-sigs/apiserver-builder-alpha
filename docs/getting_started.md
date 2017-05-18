# Getting started

This document covers building your first apiserver from scratch:

- Bootstrapping your go dependencies
- Initialize your project directory structure and go packages
- Create an API group, version, resource
- Build the apiserver command
- Write an automated test for the API resource

## Download the latest release

Make sure you downloaded and installed the latest release as described
[here](https://github.com/kubernetes-incubator/apiserver-builder/blob/master/docs/installing.md)

## Create your Go project

Create a Go project under GOPATH/src/

For example

> GOPATH/src/github.com/my-org/my-project

## Install the apiserver-builder go libraries as vendored deps

Using `apiserver-boot` to install the vendored go libraries will
make sure that your project is bootstrapped with a known set of
good libraries.

The will bootstrap will copy the files included in the apiserver-builder
binary distribution.

Files created under GOPATH/src/github.com/my-org/my-project:

- vendor
- glide.yaml
- glide.lock

At the root of your go package under your GOPATH run the following command.

```sh
apiserver-boot glide-install
```

## Create a copyright header

Create a file called `boilerplate.go.txt` that contains the copyright
you want to appear at the top of generated files.

Under GOPATH/src/github.com/my-org/my-project:

- `boilerplate.go.txt`

e.g.

```go
/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
```

## Initialize your project

This will setup the initial file structure for your apiserver, including:

Files created under GOPATH/src/github.com/my-org/my-project:

- `pkg/apis/doc.go`
- `pkg/openapi/doc.go`
- docs/...
- main.go

Flags:

- your-domain: unique namespace for your API groups

At the root of your go package under your GOPATH run the following command:

```sh
apiserver-boot init --domain <your-domain>
```

## Create an API group

**Note:** This step is optional, `create-resource` will automatically do this for you

An API group contains one or more related API versions.  It is similar to
a package in go or Java.

Files created under GOPATH/src/github.com/my-org/my-project:

- pkg/apis/your-group/doc.go

Flags:

- your-group: name of the API group e.g. `cicd` or `apps`

At the root of your go package under your GOPATH run the following command:

```sh
apiserver-boot create-group --domain <your-domain> --group <your-group>
```

This will create a new API group under pkg/apis/<your-group>

## Create an API version

**Note:** This step is optional, `create-resource` will automatically do this for you

An API version contains one or more APIs.  The version is used
to support introducing changes to APIs without breaking backwards
compatibility.

Files created under GOPATH/src/github.com/my-org/my-project:

- pkg/apis/your-group/your-version/doc.go

Flags:

- your-version: name of the API version e.g. `v1beta1` or `v1`

At the root of your go package under your GOPATH run the following command:

```sh
apiserver-boot create-version --domain <your-domain> --group <your-group> --version <your-version>
```

This will create a new API version under pkg/apis/<your-group>/<your-version>

## Create an API resource

**Note:** This will invoke `create-group` and `create-version` if they have not already been run.


An API resource provides REST endpoints for CRUD operations on a resource
type.  This is what will be used by clients to read and store instances
of the resource kind.

Files created under GOPATH/src/github.com/my-org/my-project:

- pkg/apis/your-group/your-version/your-kind_types.go
- pkg/apis/your-group/your-version/your-kind_types_test.go

Flags:

- your-kind: camelcase name of the type e.g. `MyKind`
- your-resource: lowercase pluralization of the kind e.g. `mykinds`

At the root of your go package under your GOPATH run the following command:

```sh
apiserver-boot create-resource --domain <your-domain> --group <your-group> --version <your-version> --kind <your-kind>
```

## Generate the code

The following command will generate the wiring to register your API resources.

**Note:** It must be rerun any time new fields are added to your resources

```sh
apiserver-boot generate
```

## Build and run the apiserver

Build the apiserver binary

```sh
apiserver-boot generate
```

Run an etcd instance and the apiserver.

**Note:** must have etcd on your PATH

```sh
apiserver-boot run
```

Test with kubectl

```sh
kubectl --kubeconfig kubeconfig version
```

## Run a test

A placehold test was created for your resource.  The test will
start your apiserver in memory, and allow you to create, read, and write
your resource types.

This is a good way to test validation and defaulting of your types.

```sh
go test pkg/apis/your-group/your-version/your-kind_types_test.go
```
