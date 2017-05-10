# Adding resources

Resources live under `pkg/apis/<group>/<version>/<resource>_types.go`.
It is recommended to use `apiserver-boot` to create new groups,
versions, and resources.

## Creating a resource with apiserver-boot

1. Create the group the resource will live under.

`apiserver-boot create group --domain <domain> --group <group>`

2. Create the version the resource will live under

`apiserver-boot create group --domain <domain> --group <group> --version <version>`

3. Create the resource

The resource name should be all lower case and plural.  e.g. `deployments`.
The resource also has *Kind*, which is the CamelCase singular name of
your resource.  e.g. `Deployment`

`apiserver-boot create group --domain <domain> --group <group> --version <version> --kind <Kind> --resource <resource>`

## Anatomy of a resource

A resource has a go struct which defines the *Kind* schema, and is
annotated with comment directives used by the code generator to
wire the storage and REST endpoints.

Example:

```go
// +genclient=true

// +resource:path=foos
// +k8s:openapi-gen=true
// Foo defines some thing
type Foo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

    // spec defines the desired state of Foo
	Spec   FooSpec   `json:"spec,omitempty"`

    // status records the observed state of Foo
	Status FooStatus `json:"status,omitempty"`
}

// FooSpec defines the desired state of Foo
type FooSpec struct {
    // some_spec_field defines some desired state about Foo
	SomeSpecField int `json:"some_spec_field,omitempty"`
}

// FooStatus records the observed state of Foo
type FooStatus struct {
	// some_status_field records some observed state about Foo
	SomeStatusField int `json:"some_status_field,omitempty"`
}
```

### Breakdown of example

```go
// +resource:path=foos
```

This tells the code generator to generate the REST
storage endpoints for this resource.

```go
// +k8s:openapi-gen=true
```

This tells the code generator to include this
resource in the openapi spec published by the apiserver

```go
// Foo defines some thing
```

This will appear in the openapi spec and the
generated reference docs as the description of the resource.

```go
type Foo struct {...}
```

This block defines the resource schema

```go
metav1.TypeMeta   `json:",inline"`
metav1.ObjectMeta `json:"metadata,omitempty"`
```

These define metadata common to most resources - such as
the name, group/version/kind, annotations, and labels.

```go
// spec defines the desired state of Foo
Spec   FooSpec   `json:"spec,omitempty"`
```

This field defines the desired state of Foo that the controller loops
will work towards.

```go
// status records the observed state of Foo
Status FooStatus `json:"status,omitempty"`
```

This field records the state of Foo observed by the controller loops
for clients to read.

```go
// FooSpec defines the desired state of Foo
type FooSpec struct {
    // some_spec_field defines some desired state about Foo
	SomeSpecField int `json:"some_spec_field,omitempty"`
}

// FooStatus records the observed state of Foo
type FooStatus struct {
	// some_status_field records some observed state about Foo
	SomeStatusField int `json:"some_status_field,omitempty"`
}
```

These structures define the schema for the desired and observed
state.

## Generating the wiring

To generate the REST endpoint and storage wiring for your resource,
run `apiserver-boot generate --api-versions "your-group/your-version"`
from the go package root directory.

This will also generate go client code to read and write your resources under `pkg/client`.