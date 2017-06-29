## `apiserver-builder`

[![Build Status](https://travis-ci.org/kubernetes-incubator/apiserver-builder.svg?branch=master)](https://travis-ci.org/kubernetes-incubator/apiserver-builder "Travis")

Apiserver Builder is a collection of libraries and tools to build native
Kubernetes extensions using Kubernetes apiserver code.

## Motivation

Standing up apiservers from scratch and adding apis requires non-trivial boilerplate
code that must be maintained and rebased against master. The goal of this project is
to make building apiservers in go simple and accessible to everyone in the
Kubernetes community.

The project aims to provide libraries, code generators, and tooling to make it possible to build
and run a basic apiserver in an afternoon, while providing all of the hooks to offer the
same capabilities when building from scratch.

## Installation

Download the latest release and install on your PATH. Details [here](https://github.com/kubernetes-incubator/apiserver-builder/blob/master/docs/installing.md)

## Getting started guide

**Note:** Go through this guide first.

Instructions on how to bootstrap a new apiserver with a simple type

Details [here](https://github.com/kubernetes-incubator/apiserver-builder/blob/master/docs/getting_started.md)

## User guide

### Adding a new resource with a controller

Instructions on how to add a new resource

Details [here](https://github.com/kubernetes-incubator/apiserver-builder/blob/master/docs/adding_resources.md)

### Adding validation

Instructions on how to add schema validation an existing resource

Details [here](https://github.com/kubernetes-incubator/apiserver-builder/blob/master/docs/adding_validation.md)

### Adding defaulting

Instructions on how to add field value defaulting to an existing resource

Details [here](https://github.com/kubernetes-incubator/apiserver-builder/blob/master/docs/adding_defaulting.md)

### Adding subresource

Instructions on how to add a new subresource to an existing resource

Details [here](https://github.com/kubernetes-incubator/apiserver-builder/blob/master/docs/adding_subresources.md)

### Defining custom REST handlers

Instructions on how to Overriding the default resource storage with
custom REST handlers

Details [here](https://github.com/kubernetes-incubator/apiserver-builder/blob/master/docs/adding_custom_rest.md)

## Using delegated auth with minikube

Instructions on how to run an apiserver using delegated auth with a minikube cluster

Details [here](https://github.com/kubernetes-incubator/apiserver-builder/blob/master/docs/using_minikube.md)
