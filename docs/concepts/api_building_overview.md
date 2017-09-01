# Building and using Kubernetes APIs

This document describes how Kubernetes APIs are structured and to use the apiserver-builder project
to build them.

## Vendoring the K8s libraries for building a new API server extension

First create a new go package and use `apieserver-boot init repo --domain <your domain>`.  This
will vendor the required go libraries from kubernetes and set the initial directory and package
structure including the *apiserver* and *controller* binaries responsible for storing and
reconciling objects.

## API Structure

Following is a short summary of how APIs are structured.  This may be familiar to anyone that has worked with the
Kubernetes APIs before.

Kubernetes APIs are designed such that the desired state of an object is
sent to the API, and the cluster works to reconcile the actual state with the desired state.
For example to rollout a new container image to a Deployment, only the declared Deployment image is updated
and the cluster will automatically perform a rollout of the new image.

The APIs have a level-based implementation, so they will work to the current desired state, ignoring
previous desired states that may have been set.  For example, updating a Deployment image in the middle
of an existing rollout will drive directly to the new desired state instead of completing the previous
rollout.  (e.g. it won't complete rollingout of the previous image, and will switch to rolling out the new image).

The APIs reconcile the declared desired state asynchronously, meaning the request to create a Deployment will
return before the system has tried to start any Pods.  This means many errors will not be returned to the
client as part of the initial request - e.g. if the container image is invalid the user won't find out unless they
poll or watch the Pods.  For this reason, it is important to write Status messages back to the resource as
part of reconciling the object.

Summary:

**Declarative:** Users specify desired state, not specific operations.
**Asynchronous:** Requests will return success before the system tries to reconcile the desired state
**Level based:** System drives towards current desired state and ignores previous desired states.

Next:

The API can be conceptualized in 2 parts: **Storage** and **Reconcilation**.

## Storage

Users invoke Kubernetes APIs by creating, updating or deleting resources.  The Kubernetes apiserver stores resource
objects in an etcd instance and clients can *watch* resources for changes (create, update, delete).

To define a new resource type using apiserver-builder create a new `_types.go` file in the package
`pkg/apis/<api-group>/<api-version>/`.

You can use the `apiserver-boot create group version resource` command to create a new resource definition for you.

**Note**: All types created through apiserver-boot are automatically registered with the apiserver.

### Storage structure anatomy

Resource definitions have 3 subsections.  Creating a resource with `apiserver-boot` will automatically populate
the scaffolding for you resource definition with each of these fields.

*Metadata*: Contains metadata about the resource
- Name (unique key)
- Annotations (non-queryable key-value pairs)
- Labels (queryable key-value pairs)

*Spec*: Contains the desired state
Add fields specifying the desired state here.  Used by reconciliation loops to update the cluster.

*Status*: Contains the observed state
Add fields specifying the observed state here.  Used by clients and reconciliation loops
to understand the state of the cluster.

**Note**: Resources may have different versions with different representation.  Resources
are converted between versions during storage using an "unversioned" object.

### Storage operation anatomy

During storage operations there are several opportunities to either reject the request or
modify the stored object before it is written.

#### Create

**Note**: The following operations are synchronous, however reconciling the
stored object's desired state will happen asynchronously.

*PrepareForCreate*: Perform modifications to the underlying object before it is stored.

If unspecified
in your type, a default PrepareForCreate  implementation will be provided by apiserver-builder.
The default will drop updates to Status that do not go through the Status subresource.  The
default maybe overridden by providing a function attached to `<ResourceType>Strategy` in the
`_types.go`

**Note:** This works on the *unversioned* representation of the object.

Example:

```go
func (s <ResourceType>Strategy) PrepareForCreate(ctx request.Context, obj runtime.Object) {
    // Invoke the parent implementation
    s.DefaultStorageStrategy.PrepareForCreate(ctx, obj)

    // Cast the element
    o := obj.(*<group>.<ResourceType>)

    // Your PrepareForCreate logic here
}

```

*Validate*: Perform static validation of the values set in the resource.  Reject
the creation request if the object is not valid.

If unspecified in your type, a default Validate
implementation will be provided by apiserver-builder.  The default implementation will do no
validation.  The default maybe overridden by providing a function attached to `<ResourceType>Strategy`
in the `_types.go`:

**Note:** This works on the *unversioned* representation of the object.

Example:

```go
func (<ResourceType>Strategy) Validate(ctx request.Context, obj runtime.Object) field.ErrorList {
    o := obj.(*<group>.<ResourceType>Strategy)
    errors := field.ErrorList{}

    if <some-condition>
        errors = append(errors, field.Invalid(
            field.NewPath("spec", "<some field>"), o.Spec.<SomeField>, fmt.Sprintf("Bad value %s", o.Spec.<SomeField>)))
    }
    return errors
}
```


*DefaultingFunction*: Optional fields that need to be interpreted with some value
when they are unset should be defaulted and written to the object.  Persisting
the defaulted values makes it easier to change the default values for different
versions of the same API.

If unspecified in your type, a default DefaultingFunction will be provided by
apiserver-builder.  The default implementation will do no defaulting.

**Note:** This works on the *versioned* representation of the object.

Example:

```go
func (<ResourceType>SchemeFns) DefaultingFunction(o interface{}) {
    obj := o.(*<ResourceType>)

    // Defaulting logic here
}
```

#### Update

**Note**: The following operations are synchronous, however reconciling the
stored object's desired state will happen asynchronously.

*ValidateUpdate*: A separate validate may be specified for updates that can validate
new values against old values.  This is useful for enforcing immutability of
certain fields.

If unspecified in your type, a default ValidateUpdate will be provided
by the apiserver-builder.  Can be overriden by specifying the function
`<ResourceType>Strategy.PrepareForUpdate`.

*PrepareForUpdate*: Similar to PrepareForCreate, but for update operations.

If unspecified in your type, a default PrepareForUpdate will be provided
by the apiserver-builder.  Can be overriden by specifying the function
`<ResourceType>Strategy.PrepareForUpdate`.


#### Delete

**Note**: The delete is asynchronous and the garbage collection will be
executed after the request returns.

*Finalizers*: If a finalizer is specified on an object (e.g. PrepareForCreate)
deleting an object will set the `DeletionTimestamp` field on the object with
a grace period.  This allows controllers to see the object has been deleted
and perform cleanup of resources the object may have created.

*OwnerReference*: Automatic garbage collection may be specified by setting
an OwnerReference on the object to be garbage collected.  When the owning object
is deleted, objects with the OwnerReference will automatically be deleted.
Note: This is the preferred method for garbage collection, but will not
work for cleaning up external (non-kubernetes-object) resources.  e.g.
resources provisioned through the cloud provider.  See [garbage-collection](https://kubernetes.io/docs/concepts/workloads/controllers/garbage-collection/)
for more information.

**Summary**:

- Objects have Metadata, Spec and Status fields
- Storage operations may be rejected or transformed by the API server before the object is written.
- `apiserver-boot create group version resource` will setup the basic scaffolding for a new resource type using sane overridable defaults.

### Storage operation types

Resource Status and Spec operations have separate endpoints:

- By default performing a storage operation to the resource will only update the Spec and Metadata.
- By default performing a storage operation to the status subresource will only update the Status and Metadata.

Additional subresources may be specified for polymorphic operations such as `scale`.

### Reconcilation loops

The stored object is realized in the actual cluster state through reconciliation loops.
*Controller* processes watch resource types for any updates, and then perform operations.

Operations may include:
- updating the object metadata, spec or status
- creating, updating or deleting other objects in the Kubernetes cluster (e.g. Kubernetes APIs such as Pod)
- creating, updating or deleting resources outside the Kubernetes cluster (e.g. cloud provider resources such as CloudSQL).

#### Watching for updates

Controllers watch for changes to an object using the *informer framework*.  The *informer framework*
is part of the *client-go* package, and will cache and index the objects it is watching.

**Note**: A simple *watch* will timeout, however the informer framework will continue to
watch for changes as long as the process is running.

When using `apiserver-boot create` to create a new resource type, a corresponding controller
is automatically created under the `pkg/controller/<resource-type>` package with a
function that will be invoked on any changes.  An init function for the controller
is also generated and can be modified to initialize additional dependencies and clients.

#### Defining reconciliation constrations

As part of the Spec, it may be useful to provide constraints on how a controller may reconcile the object.
For instance, when rolling out a new Deployment the controller needs to know how many Pods may be unhealthy
at a time (specified in the Spec.RolloutStrategy).  This allows clients to still control how changes are
made to the cluster while still being declarative.

### Cross-cutting initialization hooks

It is possible to specify dynamic cross-cutting hooks for initializing and validating
objects.  This allows supporting dynamic "BeforeCreate" and "Validate" semantics
for objects by installing extensions into the cluster.  For more information see
[Dynamic Admission Control](https://kubernetes.io/docs/admin/extensible-admission-controllers/)

### Cross-cutting finalization hooks

As described in the [Storage](#Storage) section, if an object has Finalizers in
its metadata it will not be deleted until a controller removes the finalizer after
completing some final action.  If the finalizer is not removed after some grace period
the object will be deleted anyway.
