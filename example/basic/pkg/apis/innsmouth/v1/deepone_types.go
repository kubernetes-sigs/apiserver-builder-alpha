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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/apiserver-builder-alpha/example/basic/pkg/apis/innsmouth/common"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +k8s:openapi-gen=true
// +resource:path=deepones
// DeepOne defines a resident of innsmouth
type DeepOne struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DeepOneSpec   `json:"spec,omitempty"`
	Status DeepOneStatus `json:"status,omitempty"`
}

type SamplePrimitiveAlias int64

// DeepOnesSpec defines the desired state of DeepOne
type DeepOneSpec struct {
	// fish_required defines the number of fish required by the DeepOne.
	FishRequired int `json:"fish_required,omitempty"`

	Sample               SampleElem                       `json:"sample,omitempty"`
	SamplePointer        *SamplePointerElem               `json:"sample_pointer,omitempty"`
	SampleList           []SampleListElem                 `json:"sample_list,omitempty"`
	SamplePointerList    []*SampleListPointerElem         `json:"sample_pointer_list,omitempty"`
	SampleMap            map[string]SampleMapElem         `json:"sample_map,omitempty"`
	SamplePointerMap     map[string]*SampleMapPointerElem `json:"sample_pointer_map,omitempty"`
	SamplePrimitiveAlias SamplePrimitiveAlias

	// Example of using a constant
	Const      common.CustomType            `json:"const,omitempty"`
	ConstPtr   *common.CustomType           `json:"constPtr,omitempty"`
	ConstSlice []common.CustomType          `json:"constSlice,omitempty"`
	ConstMap   map[string]common.CustomType `json:"constMap,omitempty"`

	// TODO: Fix issues with deep copy to make these work
	//ConstSlicePtr []*common.CustomType          `json:"constSlicePtr,omitempty"`
	//ConstMapPtr map[string]*common.CustomType `json:"constMapPtr,omitempty"`
}

type SampleListElem struct {
	Sub []SampleListSubElem `json:"sub,omitempty"`
}

type SampleListSubElem struct {
	Foo string `json:"foo,omitempty"`
}

type SampleListPointerElem struct {
	Sub []*SampleListPointerSubElem `json:"sub,omitempty"`
}

type SampleListPointerSubElem struct {
	Foo string `json:"foo,omitempty"`
}

type SampleMapElem struct {
	Sub map[string]SampleMapSubElem `json:"sub,omitempty"`
}

type SampleMapSubElem struct {
	Foo string `json:"foo,omitempty"`
}

type SampleMapPointerElem struct {
	Sub map[string]*SampleMapPointerSubElem `json:"sub,omitempty"`
}

type SampleMapPointerSubElem struct {
	Foo string `json:"foo,omitempty"`
}

type SamplePointerElem struct {
	Sub *SamplePointerSubElem `json:"sub,omitempty"`
}

type SamplePointerSubElem struct {
	Foo string `json:"foo,omitempty"`
}

type SampleElem struct {
	Sub SampleSubElem `json:"sub,omitempty"`
}

type SampleSubElem struct {
	Foo string `json:"foo,omitempty"`
}

// DeepOneStatus defines the observed state of DeepOne
type DeepOneStatus struct {
	// actual_fish defines the number of fish caught by the DeepOne.
	ActualFish int `json:"actual_fish,omitempty"`
}

var _ resource.Object = &DeepOne{}
var _ resource.ObjectWithStatusSubResource = &DeepOne{}

func (in *DeepOne) GetObjectMeta() *metav1.ObjectMeta {
	return &in.ObjectMeta
}

func (in *DeepOne) NamespaceScoped() bool {
	return true
}

func (in *DeepOne) New() runtime.Object {
	return &DeepOne{}
}

func (in *DeepOne) NewList() runtime.Object {
	return &DeepOneList{}
}

func (in *DeepOne) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "innsmouth.k8s.io",
		Version:  "v1",
		Resource: "deepones",
	}
}

func (in *DeepOne) IsStorageVersion() bool {
	return true
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DeepOneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []DeepOne `json:"items"`
}

var _ resource.ObjectList = &DeepOneList{}

func (in *DeepOneList) GetListMeta() *metav1.ListMeta {
	return &in.ListMeta
}

var _ resource.ObjectWithStatusSubResource = &DeepOne{}
var _ resource.StatusSubResource = &DeepOneStatus{}

func (in DeepOneStatus) SubResourceName() string {
	return "status"
}

func (in *DeepOne) GetStatus() resource.StatusSubResource {
	return in.Status
}

func (in DeepOneStatus) CopyTo(parent resource.ObjectWithStatusSubResource) {
	parent.(*DeepOne).Status = in
}
