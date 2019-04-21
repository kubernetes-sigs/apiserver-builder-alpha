package suite

import "sigs.k8s.io/controller-runtime/pkg/reconcile"

var _ reconcile.Reconciler = &ReconcilerInterceptor{}

// ReconcilerInterceptor is only for testing, allowing developers to intercept the reconcile func, before and after.
type ReconcilerInterceptor struct {
	delegate reconcile.Reconciler

	BeforeReconcile func(req reconcile.Request)
	AfterReconcile  func(req reconcile.Request, err error)
}

func CreateProxyReconciler(reconciler reconcile.Reconciler) *ReconcilerInterceptor {
	return &ReconcilerInterceptor{
		delegate: reconciler,
	}
}

func (r *ReconcilerInterceptor) Reconcile(req reconcile.Request) (result reconcile.Result, err error) {
	if r.BeforeReconcile != nil {
		r.BeforeReconcile(req)
	}
	if r.AfterReconcile != nil {
		defer r.AfterReconcile(req, err)
	}
	result, err = r.delegate.Reconcile(req)
	return
}
