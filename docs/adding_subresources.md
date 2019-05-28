# Adding a subresource to a resource

## Easy way: Create the subresource with apiserver-boot

**Note:** Added in release v0.1-alpha.11

Create the subresource definition using apiserver-boot

To create a subresource of for the resource *Group/Version/Kind* run the following command
from the root directory (e.g. the one that containings *glide.yaml*)

```sh
apiserver-boot create subresource --subresource <subresource> --group <resource-group> --version <resource-version> --kind <resource-kind>
```

This will:

- create `pkg/apis/<group>/<version>/<subresource>_<kind>_types.go`
  - contains the subresource type definition
- create `pkg/apis/<group>/<subresource>_<kind>_rest.go`
  - contains the rest implementation
- create `pkg/apis/<group>/<version>/<subresource>_<kind>_types_test.go`
  - contains a simple test to invoke the subresource and make sure it returns 200
- update `pkg/apis/<group>/<version>/<kind>_types.go`
  - add the subresource comment directive to the resource

Next regenerate the generated code to wire up the subresource
  
```sh
apiserver-boot build generated
```

Run the tests

```sh
go test ./pkg/...
```

Look for the subresource endpoint through the discovery service:

```sh
# shell #1
apiserver-boot run local
```

```sh
# shell #2
kubectl --kubeconfig kubeconfig proxy
```

```sh
# shell #3
curl 127.0.0.1:8001/
curl 127.0.0.1:8001/apis/<group>.<domain>/<version>
```

## Hard way: Manually create the subresource

Create the subresource definition by hand

### Update a resource with the subresource

Create a resource under `pkg/apis/<group>/<version>/<resource>_types.go`

```go
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +resource:path=bars
// +subresource:request=Status,path=status,rest=BarStatusREST
type Bar struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BarSpec   `json:"spec,omitempty"`
	Status BarStatus `json:"status,omitempty"`
}

```

The following line tells the code generator to generate a subresource for this resource.

- under the path `bar/status`
- with request Kind `Status`
- implemented by the go type `BarStatusREST`

Status and BarStatusREST live in the versioned package (same as the versioned resource definition)

```go
// +subresource:request=Status,path=status,rest=BarStatusREST
```



### Create the subresource request

Define the request type in the same <kind>_types.go file

```go
// +genclient=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +subresource-request
type Status struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Faculty int `json:"faculty,omitempty"`
}
```

Note the line:

```go
// +subresource-request
```

This tells the code generator that this is a subresource type and to
register it in the wiring.

### Create the REST implementation

Create the rest implementation in the *unversioned* package.

Example:

```go
// +k8s:deepcopy-gen=false
type BarStatusREST struct {
	Registry BarRegistry
}

// Status Subresource
var _ rest.CreaterUpdater = &BarStatusREST{}
var _ rest.Patcher = &BarStatusREST{}

func (r *BarStatusREST) Create(ctx request.Context, obj runtime.Object) (runtime.Object, error) {
	status := obj.(*Status)
	b, err := r.Registry.GetBar(ctx, status.Name, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
    // Do something with b...

    // Save the udpated b
	return r.Registry.UpdateBar(ctx, b)
}

// Get retrieves the object from the storage. It is required to support Patch.
func (r *BarStatusREST) Get(ctx request.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return nil, nil
}

// Update alters the status subset of an object.
func (r *BarStatusREST) Update(ctx request.Context, name string, objInfo rest.UpdatedObjectInfo) (runtime.Object, bool, error) {
	return nil, false, nil
}

func (r *BarStatusREST) New() runtime.Object {
	return &Status{}
}

```


## Anatomy of a REST implementation

Define the struct type implementing the REST api.  The Registry
field is required, and provides a type safe library to read / write
instances of Bar from the storage.


```go
// +k8s:deepcopy-gen=false
type BarStatusREST struct {
	Registry BarRegistry
}
```


---

Enforce local compile time checks that the struct implements
the needed REST methods

```go
// Status Subresource
var _ rest.CreaterUpdater = &BarStatusREST{}
var _ rest.Patcher = &BarStatusREST{}
```


---

Implement create and update methods using the Registry to update the parent resource.

```go
func (r *BarStatusREST) Create(ctx request.Context, obj runtime.Object) (runtime.Object, error) {
    ...
}

// Update alters the status subset of an object.
func (r *BarStatusREST) Update(ctx request.Context, name string, objInfo rest.UpdatedObjectInfo) (runtime.Object, bool, error) {
	...
}
```

---

Implement a read method using the Registry to read the parent resource.


```go
// Get retrieves the object from the storage. It is required to support Patch.
func (r *BarStatusREST) Get(ctx request.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	...
}
```

---

Implement a method that creates new instance of the request.

```go
func (r *BarStatusREST) New() runtime.Object {
	return &Status{}
}
```

__NOTE__: if you're trying extending streaming long-running subresource, please 
refer to the examples at `miskatonic.students/pencil`, or `miskatonic.students/book`.

## Generate the code for your subresource

Run the code generation command to generate the wiring for your subresource.

`apiserver-boot build generated`

## Invoke your subresource from a test

Use the RESTClient to call your subresource.  Client go is not generated
for subresources, so you will need to manually invoke the subresource.

```
client.RESTClient()
	err := restClient.Post().Namespace("default").
		Name("name").
		Resource("bars").
		SubResource("status").
		Body(status).Do().Error()
	...
```

