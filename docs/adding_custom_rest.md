# Adding resources

**Important:** Read [this doc](https://sigs.k8s.io/apiserver-builder-alpha/docs/adding_resources.md)
first to understand how resources are added.

## Create a resource with custom rest

You can implement your own REST implementation instead of using the
standard storage by any one of `Getter`, `Lister`, `Creator`, `Updater`
from `sigs.k8s.io/apiserver-runtime/pkg/builder/resource/resourcerest`
package.

```go
// +genclient=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +resource:path=foos,rest=FooREST
// +k8s:openapi-gen=true
// Foo defines some thing
type Foo struct {
    // Your resource definition here
}

// Foo resource supports "get", "create", "update" verbs. Hence you can't invoke 
// "list" upon Foo resource.
var _ resourcerest.Getter = &Foo{}
var _ resourcerest.Creator = &Foo{}
var _ resourcerest.Updater = &Foo{}

// Your rest.Storage implementation below
// ...
```
