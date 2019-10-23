# Adding non-namespaced resources

This document covers how to create a non-namespaced resource using
`apiserver-boot`.

## Prerequisites

- [adding resources](adding_resources.md)

## Creating a non-namespaced resource with apiserver-boot

Use the `--non-namespaced=true` flag when creating a resource:

`apiserver-boot create group version resource --non-namespaced=true --group <group> --version <version> --kind <Kind>`

## Non-namespaced resources

Non-namespaced resources have the following differences from namespaced resources:

- nonNamespaced comment directive above the type in `pkg/apis/{group}/{version}/{Kind}_types.go`
  - `// +genclient:nonNamespaced` comment under `// +genclient`
- Strategy and StatusStrategy override NamespacedScoped to false in `pkg/apis/{group}/{Kind}_strategy.go`
  - `func ({{.Kind}}Strategy) NamespaceScoped() bool { return false }`
  - `func ({{.Kind}}StatusStrategy) NamespaceScoped() bool { return false }`
- Do not provide namespace when creating the client from a clientset

Example:

File : `pkg/apis/{group}/{version}/{Kind}_types.go`
```go
// +genclient=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient:nonNamespaced

// +resource:path=foos
// +k8s:openapi-gen=true
// Foo defines some thing
type Foo struct {
...
}
```

File: `pkg/apis/{group}/{Kind}_strategy.go`
```go
...

func (FooStrategy) NamespaceScoped() bool { return false }

func (FooStatusStrategy) NamespaceScoped() bool { return false }
```
