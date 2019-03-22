package miskatonic

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/registry/generic"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ rest.CreaterUpdater = &StudentComputerREST{}
var _ rest.Patcher = &StudentComputerREST{}

// +k8s:deepcopy-gen=false
type StudentComputerREST struct {
	Registry StudentRegistry
}

func (r *StudentComputerREST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	sub := obj.(*StudentComputer)
	return sub, nil
}

// Get retrieves the object from the storage. It is required to support Patch.
func (r *StudentComputerREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return nil, nil
}

// Update alters the status subset of an object.
func (r *StudentComputerREST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	return nil, false, nil
}

func (r *StudentComputerREST) New() runtime.Object {
	return &StudentComputer{}
}

// Custom REST storage that delegates to the generated standard Registry
func NewStudentComputerREST(getter generic.RESTOptionsGetter) rest.Storage {
	return &StudentComputerREST{NewStudentRegistry(nil)}
}
