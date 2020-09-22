package v1beta1

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/rest"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource/resourcerest"
)

var _ resourcerest.Getter = &Student{}
var _ resourcerest.Creator = &Student{}
var _ resourcerest.Updater = &Student{}

func (in *Student) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	s := obj.(*Student)
	s.Spec.ID = s.Spec.ID + 1
	return s, nil
}

// Get retrieves the object from the storage. It is required to support Patch.
func (in *Student) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return &Student{}, nil
}

// Update alters the status subset of an object.
func (in *Student) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	return nil, false, nil
}
