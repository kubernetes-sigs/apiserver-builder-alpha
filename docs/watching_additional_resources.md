# Watching Additional Resources

This document describes how to watch multiple resources from a
controller.

## Triggering a controller reconcile for additional resources that you've created

**Step 1:** Map the additional resource instance to the key of an
instance of your resource

When receiving a notification for the additional resource
(e.g. Baz) you want to run the `Reconcile` loop for your
resource instance that *owns* the resource we were notified about.
If you write an `OwnerReference` for each of the additional resources
created by your resource, you may look at that field to find the key
for your resource.

- Cast the argument to the type
  - **Important:** double check the group/version are consistent between
    - the informer you started in `pkg/controller/sharedinformers/informers.go`
    - the type you cast the argument to
    - the informer with which you register your reconcile loop

```go
// BazToFoo takes an instance of *v1beta1.Baz and returns the
// namespace/name key of the Foo instance to call the reconcilation
// loop for
func (c *FooControllerImpl) BazToFoo(i interface{}) (string, error) {
	d, _ := i.(*v1beta1.Baz)
	if len(d.OwnerReferences) == 1 && d.OwnerReferences[0].Kind == "Foo" {
		return d.Namespace + "/" + d.OwnerReferences[0].Name, nil
	} else {
		// Not owned
		return "", nil
	}
}
```

**Step 2:** Register your resource Reconcile loop with the informer
so it gets called in response to changes to the resource you want to watch

In your controller init function tie it all together:
- The sharedinformer you started that watches for events
- The conversion from a Deployment to the key of your resource
- The reconcile function that takes the key of one of your resources

```go
func (c *FooControllerImpl) Init(arguments sharedinformers.ControllerInitArguments) {
    ...
    // Call the Reconcile(Foo) function for observed updates to
    // Baz objects.  Use the BazToFoo function to map the changed
    // Baz to the Foo instance to reconcile
	arguments.Watch(
	    "FooBaz",
	    arguments.GetSharedInformers().Factory.Bar().V1beta1().Bazs().Informer(),
	    c.BazToFoo)
}
```

### CRUD the resource from your controller reconcile loop

Add a PodSpec to your resource Spec.

Use the si.KubernetesClientSet from within your controllers
`Reconcile` function to update Kubernetes objects.

**Note**: Consider using a `Lister` for reading and indexing cached
objects to reduce load on the apiserver.


## Triggering a controller reconcile for Kubernetes resources

Enable watching Kubernetes types and initializing the ClientSet

[Example](https://github.com/kubernetes-incubator/apiserver-builder/blob/master/example/pkg/controller/sharedinformers/informers.go)

**Step 1:** Modify `pkg/controller/sharedinformer/informers.go` to have
`SetupKubernetesTypes` return `true`

```go
// SetupKubernetesTypes registers the config for watching Kubernetes types
func (si *SharedInformers) SetupKubernetesTypes() bool {
	return true
}
```

This will also allow your controller to read/write to the Kubernetes
apiserver using the `ClientSet` passed into your controller's
`Init` function through `si.KubernetesClientSet`.

**Step 2:** Modify `pkg/controller/sharedinformer/informers.go` to
start the informers for the additional resources you want to watch

For each Kind you that want to get notified about when it is
created/updated/deleted, `Run` a corresponding Informer in
`StartAdditionalInformers`.

**Note:** If you want to watch Deployments, you do *not* need to start
informers for all group/versions.  You only need to watch 1
group/version and other Deployment group/versions will be converted to
the group/version you are watching.

```go
// StartAdditionalInformers starts watching Deployments
func (si *SharedInformers) StartAdditionalInformers(shutdown <-chan struct{}) {
	go si.KubernetesFactory.Apps().V1beta1().Deployments().Informer().Run(shutdown)
}
```

**Step 3:** Map the Kubernetes resource instance to the key of an
instance of your resource

When receiving a notification for the Kubernetes resource
(e.g. Deployment) you want to run the `Reconcile` loop for your
resource instance that *owns* the resource we were notified about.
If you write an `OwnerReference` for each of the Kubernetes resources
created by your resource, you may look at that field to find the key
for your resource.

- Cast the argument to the type
  - **Important:** double check the group/version are consistent between
    - the informer you started in `pkg/controller/sharedinformers/informers.go`
    - the type you cast the argument to
    - the informer with which you register your reconcile loop

```go
func (c *FooControllerImpl) DeploymentToFoo(i interface{}) (string, error) {
	d, _ := i.(*v1beta1.Deployment)
	if len(d.OwnerReferences) == 1 && d.OwnerReferences[0].Kind == "Foo" {
		return d.Namespace + "/" + d.OwnerReferences[0].Name, nil
	} else {
		// Not owned
		return "", nil
	}
}
```

**Step 4:** Register your resource Reconcile loop with the informer
so it gets called in response to changes to the resource you want to watch

In your controller init function tie it all together:
- The sharedinformer you started that watches for events
- The conversion from a Deployment to the key of your resource
- The reconcile function that takes the key of one of your resources

```go
func (c *FooControllerImpl) Init(arguments sharedinformers.ControllerInitArguments) error) {
    ...
	arguments.Watch(
	    "FooDeployment",
	    arguments.GetSharedInformers().KubernetesFactory.Extensions().V1beta1().Deployments().Informer(),
	    c.DeploymentToFoo)
}
```
