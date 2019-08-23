package poseidon

import (
	"context"

	olympusv1beta1 "sigs.k8s.io/apiserver-builder-alpha/example/basic/pkg/apis/olympus/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Poseidon Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcilePoseidon{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("poseidon-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to Poseidon
	err = c.Watch(&source.Kind{Type: &olympusv1beta1.Poseidon{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	/*
		// TODO(user): Modify this to be the types you create
		// Uncomment watch a Deployment created by Poseidon - change this for objects you create
		err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &olympusv1beta1.Poseidon{},
		})
		if err != nil {
			return err
		}
	*/

	return nil
}

var _ reconcile.Reconciler = &ReconcilePoseidon{}

// ReconcilePoseidon reconciles a Poseidon object
type ReconcilePoseidon struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Poseidon object and makes changes based on the state read
// and what is in the Poseidon.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  The scaffolding writes
// a Deployment as an example
// +kubebuilder:rbac:groups=olympus.k8s.io,resources=poseidons,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=olympus.k8s.io,resources=poseidons/status,verbs=get;update;patch
func (r *ReconcilePoseidon) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the Poseidon instance
	instance := &olympusv1beta1.Poseidon{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
