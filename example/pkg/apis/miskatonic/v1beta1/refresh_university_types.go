/*
Copyright 2017 The Kubernetes Authors.

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
type UniversityRefresh struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Faculty int `json:"faculty,omitempty"`
}

var _ rest.CreaterUpdater = &UniversityRefreshREST{}
var _ rest.Patcher = &UniversityRefreshREST{}

// +k8s:deepcopy-gen=false
type UniversityRefreshREST struct {
	Registry miskatonic.UniversityRegistry
}

func (r *UniversityRefreshREST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	sub := obj.(*UniversityRefresh)
	rec, err := r.Registry.GetUniversity(ctx, sub.Name, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	// Modify rec in someway before writing it back to storage
	rec.Spec.FacultySize = sub.Faculty
	r.Registry.UpdateUniversity(ctx, rec)
	return rec, nil
}

// Get retrieves the object from the storage. It is required to support Patch.
func (r *UniversityRefreshREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return nil, nil
}

// Update alters the status subset of an object.
func (r *UniversityRefreshREST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	return nil, false, nil
}

func (r *UniversityRefreshREST) New() runtime.Object {
	return &UniversityRefresh{}
}
