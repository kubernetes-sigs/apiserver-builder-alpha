package poseidon

import (
	"context"

	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	olympusv1beta1 "sigs.k8s.io/apiserver-builder-alpha/example/basic/pkg/apis/olympus/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ reconcile.Reconciler = &ReconcilePoseidon{}

// ReconcilePoseidon reconciles a Poseidon object
type ReconcilePoseidon struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
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

func (r *ReconcilePoseidon) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&olympusv1beta1.Poseidon{}).
		Complete(r)
}
