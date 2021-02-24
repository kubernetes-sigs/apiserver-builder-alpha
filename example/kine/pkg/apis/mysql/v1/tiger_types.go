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

// Tiger
// +k8s:openapi-gen=true
// +resource:path=tigers,strategy=TigerStrategy,rest=TigerREST
type Tiger struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TigerSpec   `json:"spec,omitempty"`
	Status TigerStatus `json:"status,omitempty"`
}

// TigerSpec defines the desired state of Tiger
type TigerSpec struct {
}

// TigerStatus defines the observed state of Tiger
type TigerStatus struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type TigerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Tiger `json:"items"`
}

func (in *TigerList) GetListMeta() *metav1.ListMeta {
	return &in.ListMeta
}

var _ runtime.Object = &Tiger{}
var _ resource.Object = &Tiger{}
var _ resource.ObjectList = &TigerList{}

func (t *Tiger) GetObjectMeta() *metav1.ObjectMeta {
	return &t.ObjectMeta
}

func (t *Tiger) NamespaceScoped() bool {
	return true
}

func (t *Tiger) New() runtime.Object {
	return &Tiger{}
}

func (t *Tiger) NewList() runtime.Object {
	return &TigerList{}
}

func (t *Tiger) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "mysql.example.com",
		Version:  "v1",
		Resource: "tigers",
	}
}

func (t *Tiger) IsStorageVersion() bool {
	return true
}

var _ resource.ObjectWithStatusSubResource = &Tiger{}
var _ resource.StatusSubResource = &TigerStatus{}

func (in TigerStatus) SubResourceName() string {
	return "status"
}

func (t *Tiger) GetStatus() (statusSubResource resource.StatusSubResource) {
	return t.Status
}

func (in TigerStatus) CopyTo(parent resource.ObjectWithStatusSubResource) {
	parent.(*Tiger).Status = in
}
