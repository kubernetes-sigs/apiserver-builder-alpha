/*
Copyright 2020 The Kubernetes Authors.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Burger
// +k8s:openapi-gen=true
// +resource:path=burgers,strategy=BurgerStrategy,rest=BurgerREST
type Burger struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec BurgerSpec `json:"spec,omitempty"`
}

// BurgerSpec defines the desired state of Burger
type BurgerSpec struct {
}

var _ resource.Object = &Burger{}

func (b *Burger) GetObjectMeta() *metav1.ObjectMeta {
	return &b.ObjectMeta
}

func (b *Burger) NamespaceScoped() bool {
	return true
}

func (b *Burger) New() runtime.Object {
	return &Burger{}
}

func (b *Burger) NewList() runtime.Object {
	return &BurgerList{}
}

func (b Burger) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "filepath.k8s.io",
		Version:  "v1",
		Resource: "burgers",
	}
}

func (b Burger) IsStorageVersion() bool {
	return true
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type BurgerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Burger `json:"items"`
}

var _ resource.ObjectList = &BurgerList{}

func (in *BurgerList) GetListMeta() *metav1.ListMeta {
	return &in.ListMeta
}
