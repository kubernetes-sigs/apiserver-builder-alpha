This is not a Google project

**Note:** Don't use `go get` / `go install`, instead you MUST download a tar binary release or create your
own release using the release program.

## `kubebuilder`

Kubebuilder is a framework for building Kubernetes APIs.

**Note:** kubebuilder does not exist as an example to *copy-paste*, but instead provides powerful libraries and tools
to simplify building and publishing Kubernetes APIs from scratch.

## TL;DR

**First:** Download the latest [release](https://github.com/kubernetes-incubator/apiserver-builder/releases/tag/kubebuilder-v1.9-alpha.11) and
extract the tar.gz into /usr/local/kubebuilder and update your PATH to include
/usr/local/kubebuilder/bin.

```sh
# Initialize your project
kubebuilder init --domain example.com

# Create a new API and controller
kubebuilder create resource --group bar --version v1alpha1  --kind Foo

# Install and run your API into the cluster for your current kubeconfig context
kubebuilder run local
```

## Guides

**Note:** The guides are presented roughly in the order of recommended progression.

#### Installation guide

Download the latest release and install on your PATH.

[installation guide](docs/installing.md)

#### Tooling and development workflow overview

Instructions on how to use the tools packaged with kubebuilder to build APIs from scratch.

[tools guide](docs/tools_user_guide.md)

// TODO: Write these

- Creating a new project
- Creating a new resource and controller
- Configuring and running integration tests
- Running locally against a cluster
- Building container images from Docker files
- Implementing custom installation logic
- Adding examples to reference documentation
- Installing reference docs server with the APIs
- Installing using apiserver aggregation instead of CRDs
- Building using Bazel

#### Libraries and API coding guides

Instructions for how to build.

- [Creating a new resource and controller](adding_resources.md)
- [Watching other resources from your resource controller](watching_additional_resources.md)
- [Adding RBAC rules for your resource controller](declaring_rbac_rules_for_controllers.md)
- [Creating a non-namespaced resource](adding_non_namespaced_resources.md)

// TODO: Write these

- Customizing generated CRDs

## Motivation

Building Kubernetes tools and APIs involves making a lot of decisions
and writing a lot of boilerplate.

In order to facilitate easily building Kubernetes APIs and tools using
the canonical approach, this framework provides a collection of
Kubernetes development tools to minimize toil.


Kubebuilder attempts to facilitate the following developer workflow for building APIs

1. Create a new project directory
2. Create one or more resource APIs as CRDs and then add fields to the resources
3. Implement reconcile loops in controllers and watch additional resources
4. Test by running against a cluster (self-installs CRDs automatically)
5. Update bootstrapped integration tests to test new fields and business logic
6. Build and publish a container from the provided Dockerfile
7. Build and publish reference documentation for new APIs

### Scope

The current scope is focused on building APIs as CRDs or extension apiservers with controllers.

### Philosophy

- Prefer using go *interfaces* over relying on *code generation*
- Prefer using *code generation* over *1 time init* of stubs
- Prefer *1 time init* of stubs over handwritten boilerplate

### Techniques

- Provide higher level libraries on top of low level client libraries
  - Protect developers from breaking changes in low level libraries
    by providing high-level abstractions.
  - Start minimal and provide progressive discovery of functionality
  - Provide sane defaults and allow users to override when they exist
- Provide code generators to maintain common boilerplate that can't be addressed by interfaces
  - Driven off of `//+` comments
- Provide bootstrapping commands to initialize new packages
