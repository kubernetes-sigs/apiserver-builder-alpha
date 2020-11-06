# Adding a value defaulting to a resources schema

To add server side field value defaulting for your resource implement 
the interface `var _ resourcestrategy.Defaulter = &<Kind>{}`
in the type definition file. Specifically the defaulter interface lies
in the package `sigs.k8s.io/apiserver-runtime/pkg/builder/resource/resourcestrategy`.

Example:

File: `pkg/apis/<group>/<version>/<kind>_types.go`

```go
var _ resourcestrategy.Defaulter = &Foo{}

func (in *Foo) Default() {
	if in.Spec.Field == nil {
		f := "value"
		in.Spec.Field = &f
	}
}
```

## Anatomy of defaulting

By default, the apiserver-boot won't generate defaulter related codes for 
you. You're supposed to manually add the interface assertion and implementation.

---

Update set values for fields with nil values.

```go
	if in.Spec.Field == nil {
		f := "value"
		in.Spec.Field = &f
	}
```