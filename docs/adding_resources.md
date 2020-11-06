# Adding resources

Resources live under `pkg/apis/<group>/<version>/<resource>_types.go`.
It is recommended to use `apiserver-boot` to create new groups,
versions, and resources.

## Creating a resource with apiserver-boot

Provide your domain + the api group and version + the resource Kind.
The resource name will be the pluralized lowercased kind.

`apiserver-boot create group version resource --group <group> --version <version> --kind <Kind>`

Add `--with-status-subresource=false` option, if your resource is stateless.

## Anatomy of a resource

A resource has a go struct which defines the *Kind* schema, and is
annotated with comment directives used by the code generator to
wire the storage and REST endpoints.

Example:

```go
var _ resource.Object = &Foo{}
var _ resourcestrategy.Validater = &Foo{}

// +genclient=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

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

var _ resource.ObjectWithStatusSubResource = &Foo{}
var _ resource.StatusSubResource = &FooStatus{}

// FooStatus records the observed state of Foo
type FooStatus struct {
	// some_status_field records some observed state about Foo
	SomeStatusField int `json:"some_status_field,omitempty"`
}
```

### Breakdown of example

```go
var _ resource.Object = &Foo{}
```

This line ensures the resource implements `resource.Object` from apiserver-runtime
so that it can be served in the apiserver.

```go
var _ resourcestrategy.Validater = &Foo{}
```

This tells the apiserver-runtime the resource has a custom validation function
which is called before we writing the resource into storage.

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

```go
var _ resource.ObjectWithStatusSubResource = &Foo{}
var _ resource.StatusSubResource = &FooStatus{}
```

These lines ensure that the resource its status subresource, and defines the
behaviod of the status subresource.

## Controller

The controller reuses the kubebuilder's scaffolding, read [this doc](https://book.kubebuilder.io/cronjob-tutorial/controller-overview.html)
for more detailed information.
