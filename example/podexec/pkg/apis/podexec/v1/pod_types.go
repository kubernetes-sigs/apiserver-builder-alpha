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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Pod
// +k8s:openapi-gen=true
type Pod struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PodSpec   `json:"spec,omitempty"`
	Status PodStatus `json:"status,omitempty"`
}

var _ resource.Object = &Pod{}

func (in *Pod) GetObjectMeta() *metav1.ObjectMeta {
	return &in.ObjectMeta
}

func (in *Pod) NamespaceScoped() bool {
	return true
}

func (in *Pod) New() runtime.Object {
	return &Pod{}
}

func (in *Pod) NewList() runtime.Object {
	return &PodList{}
}

func (in *Pod) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "podexec.example.com",
		Version:  "v1",
		Resource: "pods",
	}
}

func (in *Pod) IsStorageVersion() bool {
	return true
}

// PodSpec defines the desired state of Pod
type PodSpec struct {
}

// PodStatus defines the observed state of Pod
type PodStatus struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PodList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Pod `json:"items"`
}

var _ resource.ObjectWithArbitrarySubResource = &Pod{}

func (in *Pod) GetArbitrarySubResources() []resource.ArbitrarySubResource {
	return []resource.ArbitrarySubResource{
		&PodExec{},
	}
}
