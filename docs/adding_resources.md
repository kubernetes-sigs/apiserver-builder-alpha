# Adding resources

Resources live under `pkg/apis/<group>/<version>/<resource>_types.go`.
It is recommended to use `kubebuilder` to create new groups,
versions, and resources.

## Creating a resource with kubebuilder

Provide your the api group and version + the resource Kind.
The resource name will be the pluralized lowercase kind.

`kubebuilder create resource --group <group> --version <version> --kind <Kind>`

## Anatomy of a resource

A resource has a go struct which defines the *Kind* schema, and is
annotated with comment directives used by the code generator to
generate the CRD definition and go clients.

Example:

```go
// +genclient=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

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

This tells the code generator to generate the CRD.

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

## Controller

By default, a controller for your resource will also be created under
`pkg/controller/<kind>/controller.go`.  This will listen for creates
or updates to your resource and execute code in response.  You can modify
the code to also listen for changes to other resource types that your
kind manages.

```go
// +controller:group=bar,version=v1alpha1,kind=Foo,resource=foos
type FooControllerImpl struct {
	builders.DefaultControllerFns

	// lister indexes properties about Foo
	lister listers.FooLister
}

// Init initializes the controller and is called by the generated code
// Register watches for additional resource types here.
func (c *FooControllerImpl) Init(arguments sharedinformers.ControllerInitArguments) {
	// Use the lister for indexing foos labels
	c.lister = arguments.GetSharedInformers().Factory.Bar().V1alpha1().Foos().Lister()
}

// Reconcile handles enqueued messages
func (c *FooControllerImpl) Reconcile(u *v1alpha1.Foo) error {
	// Implement controller logic here
	log.Printf("Running reconcile Foo for %s\n", u.Name)
	return nil
}

func (c *FooControllerImpl) Get(namespace, name string) (*v1alpha1.Foo, error) {
	return c.lister.Foos(namespace).Get(name)
}
```

### Breakdown of example

```go
// +controller:group=bar,version=v1alpha1,kind=Foo,resource=foos
type FooControllerImpl struct {
	builders.DefaultControllerFns

	// lister indexes properties about Foo
	lister listers.FooLister
}
```

This declares a new controller that responds to events on Foo resources

```go
// Init initializes the controller and is called by the generated code
// Register watches for additional resource types here.
func (c *FooControllerImpl) Init(arguments sharedinformers.ControllerInitArguments) {
	// Use the lister for indexing foos labels
	c.lister = arguments.GetSharedInformers().Factory.Bar().V1alpha1().Foos().Lister()
}

```

This registers a new EventHandler for Add and Update events to Foo
resources and queues messages in response.  To add event handlers for
other resource types see
[watching additional resources from your controller](watching_additional_resources.md)

```go
// Reconcile handles enqueued messages
func (c *FooControllerImpl) Reconcile(u *v1alpha1.Foo) error {
	// Implement controller logic here
	log.Printf("Running reconcile Foo for %s\n", u.Name)
	return nil
}
```

This function is called when messages are dequeued.  It should read the
actual state and reconcile it with the desired state.

```go
func (c *FooControllerImpl) Get(namespace, name string) (*v1alpha1.Foo, error) {
	return c.lister.Foos(namespace).Get(name)
}
```

This function looks up a Foo object for a namespace + name.  It is
executed just before the Reconcile method to lookup the Foo object.

## Generating the wiring

To generate the REST endpoint and storage wiring for your resource,
run `kubebuilder build generated` from the go package root directory.

This will also generate go client code to read and write your resources
under `pkg/client`.