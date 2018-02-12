# Getting started

This document covers building an API using CRDs and a controller
`kubebuilder`.  It is focused on how to use the most basic aspects of
the tooling to be productive quickly.

For information on the libraries, see the [libraries user guide](libraries_user_guide.md)

New API workflow:

- Bootstrap go vendor + initialize required directory structure and
  go packages
- Create an API group, version, resource + controller
- Build and run against a Kubernetes cluster
- Run tests

## Download the latest release

Make sure you downloaded and installed the latest release as described
[here](https://github.com/kubernetes-sigs/kubebuilder/blob/master/docs/installing.md)

## Create your Go project

Create a Go project under GOPATH/src/

For example

> GOPATH/src/github.com/my-org/my-project

## Create a copyright header

Create a file called `boilerplate.go.txt`.  This file will contain the
copyright boilerplate appearing at the top of all generated files.

Under GOPATH/src/github.com/my-org/my-project:

- `boilerplate.go.txt`

e.g.

```go
/*
Copyright 2018 The Kubernetes Authors.

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

This will setup the initial file structure for your API, including
vendored go libraries pinned to a set of versions tested and known
to work together.  Vendored libraries are distributed with the
kubebuilder release and extracted from the installation directory.

Flags:

- your-domain: unique namespace for your API groups

At the root of your go package under your GOPATH run the following command:

```sh
kubebuilder init --domain <your-domain>
```

## Create an API

An API resource provides REST endpoints for CRUD operations on a resource
type.  This is what will be used by clients to read and store instances
of the resource kind.

API resources are defined by a group (like a package),
a version (v1alpha1, v1beta1, v1), and a Kind (the type)
Running the `kubebuilder create resource` command will create the
api group, version and Kind for you.

Files created under GOPATH/src/github.com/my-org/my-project:

- `pkg/apis/your-group/your-version/your-kind_types.go`
  - See the [libraries user guide](libraries_user_guide.md) for addtional information
- `pkg/apis/your-group/your-version/your-kind_types_test.go`
  - type integration test - basic storage read / write test
- `pkg/controller/your-kind/controller.go`
  - controller implementation - empty control loop created
- `pkg/controller/your-kind/controller_test.go`
  - controller integration test - basic control loop runs on create test
- `docs/examples/your-kind/your-kind.yaml`
  - example to show in the reference documentation - empty example
- `samples/your-kind.yaml`
  - sample config for testing your resource in your cluster - empty sample

Flags:

- your-group: name of the API group e.g. `batch`
- your-version: name of the API version e.g. `v1beta1` or `v1`
- your-kind: **Upper CamelCase** name of the type e.g. `MyKind`

At the root of your go package under your GOPATH run the following command:

```sh
kubebuilder create resource --group <yourgroup> --version <yourversion> --kind <YourKind>
```

> **Note:** To skip creating the control, pass the `--controller=false` flag.

## Setup the CRD + controller against a remote cluster (run locally)

Run the controller manager locally.

```sh
kubebuilder run local
```

> **Note:** the controller manager will install or update the CRDs in the cluster as needed.

Code generates and building executables maybe run separate using
`kubebuilder build generated` or `kubebuilder build executables`.

> **Note:** The generators must be rerun after fields are added or removed from your resources

## Create a new instance of your CRD

```sh
kubectl create -f sample/<type>.yaml
kubectl get <type>s
```

Look at the controller logs to see the reconcile loop print a message

## Run the tests

A placeholder test was created for your resource.  The test will
start a Kubernetes apiserver + etcd, install your CRDs and in some cases start
your controller.

```sh
# tell the test framework where the control plane binaries live
export TEST_ASSET_KUBECTL=/usr/local/kubebuilder/bin/kubectl
export TEST_ASSET_KUBE_APISERVER=/usr/local/kubebuilder/bin/kube-apiserver
export TEST_ASSET_ETCD=/usr/local/kubebuilder/bin/etcd

go test ./pkg/...
```

## Build reference documentation

You can build the reference documentation for your APIs by running

```sh
kubebuilder build docs
```

The output will be written to docs/build/index.html

> **Note:** To update the examples in the reference doc, modify the examples under `docs/examples`.

## Build and run an image for your CRD and Controller

A `Dockerfile` was created for you as part of `kubebuilder init`.
This Dockerfile will build the controller-manager from source and
run the tests under `./pkg/...` and `./cmd/...`

```sh
docker build . -f Dockerfile.install -t <install-image>:<version>
docker build . -f Dockerfile.controller -t <controller-image>:<version>
docker build . -f Dockerfile.docs -t <docs-image>:<version>
kubectl run <my-image> --image <my-image>
```

## Build and run the controller using Bazel

Bazel offers faster build times for development.

```sh
kubebuilder build generated
bazel run gazelle
bazel run cmd/controller-manager:controller-manager -- --kubeconfig ~/.kube/config
```