# Apiserver Builder

Apiserver builder is a collection of libraries and tools to
build Kubernetes native extensions using Kubernetes apiserver cdoe.

## Quick start

### Installation

Instructions on installing the set of binary tools for:

- bootstrapping new apiserver code, godeps, and APIs
- generating code for APIs
- generating docs

Details [here](https://github.com/kubernetes-incubator/apiserver-builder/blob/master/docs/installing.md)

### Getting started

Instructions on how to bootstrap a new apiserver with a simple type

Details [here](https://github.com/kubernetes-incubator/apiserver-builder/blob/master/docs/creating_an_api_server.md)

### Adding a new resource

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

### Generating code

Run:

`apiserver-boot generate --api-versions "your-group/your-version"`

- `api-versions` may be specified multiple times

### Generating docs

Run:

`apiserver-boot generate docs --server <apiserver-binary>`

### Running the apiserver

Run:

`apiserver-boot run --server <apiserver-binary>`

This will create a kubeconfig file to use with `kubectl --kubeconfig`

### Using delegated auth with minikube

Instructions on how to run an apiserver using delegated auth with a minikube cluster

Details [here](https://github.com/kubernetes-incubator/apiserver-builder/blob/master/docs/using_minikube.md)


## Motivation

Standing up apiservers from scratch and adding apis requires 100's of lines of boilerplate
code that must be understood and maintained (rebased against master).  There are few defaults,
requiring the common case configuration to be repeated for each new apiserver and resource.
Apiservers rely heavily on code generation to build libraries used by the apiserver, putting a
steep learning curve on Kubernetes community members that want to implement a native api.
Frameworks like Ruby on Rails and Spring have made standing up REST apis trivial by eliminating
boilerplate and defaulting common values, allowing developers to focus on creating
implementing the business logic of their component.

## Goals

- Working hello-world apiserver in ~5 lines
- Declaring new resource types only requires defining the struct definition
  and taging it as a resource
- Adding sub-resources only requires defining the request-type struct definition,
  implementing the REST implementation, and tagging the parent resource.
- Adding validation / defaulting to a type only requires defining the validation / defaulting method
  as a function of the appropriate struct type.
- All necessary generated code can be generated running a single command, passing in repo root.


### Binary distribution of build tools

- Distribute binaries for all of the code-generators
- Write porcelian wrapper for code-generators that is able to detect
  the appropriate arguments for each from the go PATH and types.go files

### Helper libraries

- Implement common-case defaults for create/update strategies
  - Define implementable interfaces for default actions requiring
    type specific knowledge - e.g. HasStatus - how to set and get Status
- Implement libraries for registering types and setting up strategies
  - Implement structs to defining wiring semantics instead of linking
    directly to package variables for declarations
- Implement libraries for registering subresources

### Generate code for common defaults that require type or variable declarations

- Implementations for "unversioned" types
- Implementations for "List" types
- Package variables used by code generation
- Generate invocations of helper libraries from observered types.go types

### Support hooks for overriding defaults

- Try to support 100% of the flexibility of manually writing the boilerplate by 
  providing hooks.
  - Implement functions that can be invoked to register overrides
  - Use type embeding to inherit defaults but allow new functions to override the defaults
  
### Support for generating reference documentation

- Generate k8s.io style reference documentation for declared types
  - Support for request / response examples and manual edits

### Thorough documentation and examples for how to use the framework

- Hello-world example
- How to override each default
- Build tools
- How to use libraries directly (without relying on code generation)
