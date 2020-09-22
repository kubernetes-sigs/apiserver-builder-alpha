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
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/klog/v2"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource/resourcestrategy"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Festival
// +k8s:openapi-gen=true
// +resource:path=festivals,strategy=FestivalStrategy,shortname=fs
type Festival struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FestivalSpec   `json:"spec,omitempty"`
	Status FestivalStatus `json:"status,omitempty"`
}

// FestivalSpec defines the desired state of Festival
type FestivalSpec struct {
	// Year when the festival was held, may be negative (BC)
	Year int `json:"year,omitempty"`
	// Invited holds the number of invited attendees
	Invited uint `json:"invited,omitempty"`
}

// FestivalStatus defines the observed state of Festival
type FestivalStatus struct {
	// Attended holds the actual number of attendees
	Attended uint `json:"attended,omitempty"`
}

var _ resource.Object = &Festival{}
var _ resource.ObjectWithStatusSubResource = &Festival{}
var _ resourcestrategy.Validater = &Festival{}
var _ resource.ObjectList = &FestivalList{}

func (in *Festival) GetObjectMeta() *metav1.ObjectMeta {
	return &in.ObjectMeta
}

func (in *Festival) NamespaceScoped() bool {
	return false
}

func (in *Festival) New() runtime.Object {
	return &Festival{}
}

func (in *Festival) NewList() runtime.Object {
	return &FestivalList{}
}

func (in *Festival) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "kingsport.k8s.io",
		Version:  "v1",
		Resource: "festivals",
	}
}

func (in *Festival) IsStorageVersion() bool {
	return true
}

func (in *Festival) Validate(ctx context.Context) field.ErrorList {
	klog.Infof("Validating fields for Festival %s", in.Name)
	errors := field.ErrorList{}

	if in.Spec.Year < 0 {
		errors = append(errors,
			field.Invalid(field.NewPath("spec", "year"), in.Spec.Year, "year must be > 0"))
	}

	// perform validation here and add to errors using field.Invalid
	return errors
}

func (in *Festival) SetStatus(statusSubResource interface{}) {
	in.Status = statusSubResource.(FestivalStatus)
}

func (in *Festival) GetStatus() interface{} {
	return in.Status
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type FestivalList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Festival `json:"items"`
}

func (in *FestivalList) GetListMeta() *metav1.ListMeta {
	return &in.ListMeta
}
