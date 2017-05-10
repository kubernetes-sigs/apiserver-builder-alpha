# Apiserver Builder

Apiserver builder is a collection of libraries and tools to
build Kubernetes native extensions using Kubernetes apiserver cdoe.

## Installation

Instructions on installing the set of binary tools for:

- bootstrapping new apiserver code, godeps, and APIs
- generating code for APIs
- generating docs

Details [here](https://github.com/kubernetes-incubator/apiserver-builder/blob/master/docs/installing.md)

## Getting started

Instructions on how to bootstrap a new apiserver with a simple type

Details [here](https://github.com/kubernetes-incubator/apiserver-builder/blob/master/docs/creating_an_api_server.md)

## Adding a new resource

Instructions on how to add a new resource

Details [here](https://github.com/kubernetes-incubator/apiserver-builder/blob/master/docs/adding_resources.md)

## Adding validation

Instructions on how to add schema validation an existing resource

Details [here](https://github.com/kubernetes-incubator/apiserver-builder/blob/master/docs/adding_validation.md)

## Adding defaulting

Instructions on how to add field value defaulting to an existing resource

Details [here](https://github.com/kubernetes-incubator/apiserver-builder/blob/master/docs/adding_defaulting.md)

## Adding subresource

Instructions on how to add a new subresource to an existing resource

Details [here](https://github.com/kubernetes-incubator/apiserver-builder/blob/master/docs/adding_subresources.md)

## Generating code

Run:

`apiserver-boot generate --api-versions "your-group/your-version"`

- `api-versions` may be specified multiple times

## Generating docs

Run:

`apiserver-boot generate docs --server <apiserver-binary>`

## Running the apiserver

Run:

`apiserver-boot run --server <apiserver-binary>`

This will create a kubeconfig file to use with `kubectl --kubeconfig`

## Using delegated auth with minikube

Instructions on how to run an apiserver using delegated auth with a minikube cluster

Details [here](https://github.com/kubernetes-incubator/apiserver-builder/blob/master/docs/using_minikube.md)

## Using apiserver-builder libraries directly (without generating code)

TODO: Write this