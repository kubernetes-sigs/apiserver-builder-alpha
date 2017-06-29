## `apiserver-builder`

[![Build Status](https://travis-ci.org/kubernetes-incubator/apiserver-builder.svg?branch=master)](https://travis-ci.org/kubernetes-incubator/apiserver-builder "Travis")

**Note**: This project is still only a proof of concept, and is not production ready.

Apiserver Builder is a collection of libraries and tools to build native
Kubernetes extensions using Kubernetes apiserver code.

## Motivation

*Addon apiservers* are a Kubernetes extension point allowing fully featured Kubernetes
APIs to be developed on the same api-machinery used to build the core Kubernetes APIS,
but distributed and installed into clusters.

Building addon apiservers from using the raw api-machinery requires non-trivial
code that must be maintained and rebased against master. The goal of this project is
to make building apiservers in go simple and accessible to everyone in the
Kubernetes community.

The project provides libraries, code generators, and tooling to make it possible to build
and run a basic apiserver in an afternoon, while providing all of the hooks to offer the
same capabilities when building from scratch.

## Guides

#### Installation guide

Download the latest release and install on your PATH.

[installation guide](docs/installing.md)

#### Tools user guide

**Note:** Go through this guide first.

Instructions on how to use the tools packaged with apiserver-builder to build a new apiserver
containing a simple type.

[tools guide](docs/tools_user_guide.md)

#### Coding and libraries user guide

Instructions for how to complete various tasks using the apiserver-builder libraries.

[libraries guide](docs/libraries_user_guide.md)

#### Concept guides

Conceptual information on how to run addon apiservers

[auth](docs/concepts/auth.md)


## Additional material

##### Using delegated auth with minikube

Instructions on how to run an apiserver using delegated auth with a minikube cluster

Details [here](https://github.com/kubernetes-incubator/apiserver-builder/blob/master/docs/using_minikube.md)
