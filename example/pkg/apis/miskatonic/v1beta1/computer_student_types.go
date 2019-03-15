/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	"context"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/rest"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kubernetes-incubator/apiserver-builder-alpha/example/pkg/apis/miskatonic"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +subresource-request
type StudentComputer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
}

var _ rest.CreaterUpdater = &StudentComputerREST{}
var _ rest.Patcher = &StudentComputerREST{}

// +k8s:deepcopy-gen=false
type StudentComputerREST struct {
	Registry miskatonic.StudentRegistry
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
func NewStudentComputerREST() rest.Storage {
	return &StudentComputerREST{miskatonic.NewStudentRegistry(nil)}
}
