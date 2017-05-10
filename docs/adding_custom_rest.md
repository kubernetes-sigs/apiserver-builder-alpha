# Adding resources

**Important:** Read [this doc](https://github.com/kubernetes-incubator/apiserver-builder/blob/master/docs/adding_resources.md)
first to understand how resources are added

## Create a resource with custom rest

You can implement your own REST implementation instead of using the
standard storage by providing the `rest=KindREST` parameter
and providing a `newKindREST() rest.Storage {}` function to return the
storage.

For more information on custom REST implementations, see the
[subresources doc](https://github.com/kubernetes-incubator/apiserver-builder/blob/master/docs/adding_subresources.md)

```go
// +genclient=true

// +resource:path=foos,rest=FooREST
// +k8s:openapi-gen=true
// Foo defines some thing
type Foo struct {
    // Your resource definition here
}

// Custom REST storage
func NewFooREST() rest.Storage {
    // Your rest.Storage implementation here
}
```