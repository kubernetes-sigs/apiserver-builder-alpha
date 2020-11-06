# Adding a validation to a resources schema

**Note:** By default, when creating a resource with `apiserver-boot create group version resource` a validation
function will be created in the versioned `<kind>_types.go` file.

```go
func (in *Foo) Validate(ctx context.Context) field.ErrorList {
	// TODO(user): Modify it, adding your API validation here.
	return nil
}
```

To add server side validation for your resource fill the `Validate` function with
your implementation.

Example:

File: `pkg/apis/<group>/<version>/foo_types.go`

```go
// Resource Validation
func (in *Foo) Validate(ctx request.Context) field.ErrorList {
	foo := in.(*Foo)
	errors := field.ErrorList{}
	if ... {
		errors = append(errors, field.Invalid(
			field.NewPath("spec", "Field"),
			*foo.Spec.Field,
			"Error message"))
	}
	return errors
}
```

## Anatomy of validation

Use the field.Invalid function to specify errors scoped to fields in the object.

```go
field.Invalid(field.NewPath("spec", "Field"), *bar.Spec.Field, "Error message")
```
