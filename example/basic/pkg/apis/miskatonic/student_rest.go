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

package miskatonic

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"sigs.k8s.io/apiserver-builder-alpha/pkg/builders"
)

//var _ rest.CreaterUpdater = NewStudentREST()
//var _ rest.Patcher = NewStudentREST()

var _ rest.CreaterUpdater = &StudentREST{}
var _ rest.Patcher = &StudentREST{}
var _ rest.Scoper = &StudentREST{}

// +k8s:deepcopy-gen=false
type StudentREST struct {
	*genericregistry.Store
}

func (r *StudentREST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	s := obj.(*Student)
	s.Spec.ID = s.Spec.ID + 1
	return s, nil
}

// Get retrieves the object from the storage. It is required to support Patch.
func (r *StudentREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return &Student{}, nil
}

// Update alters the status subset of an object.
func (r *StudentREST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	return nil, false, nil
}

func (r *StudentREST) New() runtime.Object {
	return &Student{}
}

func (r *StudentREST) NamespaceScoped() bool {
	return true
}

func (r *StudentREST) ShortNames() []string {
	return []string{"st"}
}

func (r *StudentREST) Categories() []string {
	return []string{""}
}

// Custom REST storage that delegates to the generated standard Registry
func NewStudentREST(optsGetter generic.RESTOptionsGetter) rest.Storage {
	groupResource := schema.GroupResource{
		Group:    "miskatonic.k8s.io",
		Resource: "students",
	}
	strategy := &StudentStrategy{builders.StorageStrategySingleton}
	store := &genericregistry.Store{
		NewFunc:                  func() runtime.Object { return &Student{} },
		NewListFunc:              func() runtime.Object { return &StudentList{} },
		DefaultQualifiedResource: groupResource,
		TableConvertor:           rest.NewDefaultTableConvertor(groupResource),

		CreateStrategy: strategy, // TODO: specify create strategy
		UpdateStrategy: strategy, // TODO: specify update strategy
		DeleteStrategy: strategy, // TODO: specify delete strategy
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter}
	if err := store.CompleteWithOptions(options); err != nil {
		panic(err) // TODO: Propagate error up
	}
	return &StudentREST{store}
}
