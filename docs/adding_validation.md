# Adding a validation to a resources schema

To add server side schema validation for your resource override
the function `func (<group>.<Kind>Strategy) Validate(ctx request.Context, obj runtime.Object) field.ErrorList`
in the group package.

**Important:** The validation logic lives in the group package *not* the version package.

Example:

File: `pkg/apis/<group>/bar.go`

```go
// Resource Validation
func (BarStrategy) Validate(ctx request.Context, obj runtime.Object) field.ErrorList {
	bar := obj.(*Bar)
	errors := field.ErrorList{}
	if ... {
		errors = append(errors, field.Invalid(
			field.NewPath("spec", "Field"),
			*bar.Spec.Field,
			"Error message"))
	}
	return errors
}
```

## Anatomy of validation

A default `<Kind>Strategy` is generated for each resource with an embedded
default Validation function.  To specify custom validation logic,
override the embedded implementation.

Cast the object type to your resource Kind

```go
bar := obj.(*Bar)
```

---

Use the field.Invalid function to specify errors scoped to fields in the object.

```go
field.Invalid(field.NewPath("spec", "Field"), *bar.Spec.Field, "Error message")
```