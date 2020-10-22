/*
Copyright YEAR The Kubernetes Authors.

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
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource/resourcerest"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Poseidon
// +k8s:openapi-gen=true
// +resource:path=poseidons,strategy=PoseidonStrategy
type Poseidon struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PoseidonSpec   `json:"spec,omitempty"`
	Status PoseidonStatus `json:"status,omitempty"`
}

// PoseidonSpec defines the desired state of Poseidon
type PoseidonSpec struct {
	PodSpec    v1.PodTemplate
	Deployment appsv1.Deployment
}

// PoseidonStatus defines the observed state of Poseidon
type PoseidonStatus struct {
}

var _ resource.Object = &Poseidon{}
var _ resourcerest.FieldsIndexer = &Poseidon{}
var _ resource.ObjectWithStatusSubResource = &Poseidon{}
var _ resource.ObjectList = &PoseidonList{}

func (in *Poseidon) GetObjectMeta() *metav1.ObjectMeta {
	return &in.ObjectMeta
}

func (in *Poseidon) NamespaceScoped() bool {
	return true
}

func (in *Poseidon) New() runtime.Object {
	return &Poseidon{}
}

func (in *Poseidon) NewList() runtime.Object {
	return &PoseidonList{}
}

func (in *Poseidon) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "olympus.k8s.io",
		Version:  "v1beta1",
		Resource: "poseidons",
	}
}

func (in *Poseidon) IsStorageVersion() bool {
	return true
}

func (in *Poseidon) SetStatus(statusSubResource interface{}) {
	in.Status = statusSubResource.(PoseidonStatus)
}

func (in *Poseidon) GetStatus() (statusSubResource interface{}) {
	return in.Status
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PoseidonList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Poseidon `json:"items"`
}

func (in *PoseidonList) GetListMeta() *metav1.ListMeta {
	return &in.ListMeta
}

func (in *Poseidon) IndexingFields() []string {
	return []string{
		"metadata.name",
		"metadata.namespace",
		"spec.deployment.name",
	}
}

func (in *Poseidon) GetField(fieldName string) string {
	switch fieldName {
	case "metadata.name":
		return in.Name
	case "metadata.namespace":
		return in.Namespace
	case "spec.deployment.name":
		return in.Spec.Deployment.Name
	}
	panic(fmt.Sprintf("getting field %v not supported", fieldName))
}
