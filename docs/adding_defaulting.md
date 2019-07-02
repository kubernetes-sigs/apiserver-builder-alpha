# Adding a value defaulting to a resources schema

To add server side field value defaulting for your resource override
the function `func SetDefaults_<Kind>(obj <Kind>)`
in the group package. And the defaulter-gen will do all the rest for you.

**Important:** The validation logic lives in the unversioned package *not* the group package.

Example:

File: `pkg/apis/<group>/<version>/defaults.go`

```go
func SetDefaults_<Kind>(o <Kind>) {
	obj := o.(*<Kind>)
	if obj.Spec.Field == nil {
		f := "value"
		obj.Spec.Field = &f
	}
}
```

## Anatomy of defaulting

You're supposed to create your own "defaults.go" and do the coding. To specify custom defaulting logic,
override the embedded implementation.

Cast the object type to your resource Kind

```go
bar := obj.(*Bar)
```

---

Update set values for fields with nil values.

```go
	if obj.Spec.Field == nil {
		f := "value"
		obj.Spec.Field = &f
	}
```