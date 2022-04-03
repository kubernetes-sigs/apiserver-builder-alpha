
/*
Copyright 2022.

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
 	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource/resourcestrategy"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient:nonNamespaced

// Volume
// +k8s:openapi-gen=true
type Volume struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VolumeSpec   `json:"spec,omitempty"`
	Status VolumeStatus `json:"status,omitempty"`
}

// VolumeList
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type VolumeList struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Volume `json:"items"`
}

// VolumeSpec defines the desired state of Volume
type VolumeSpec struct {
}

var _ resource.Object = &Volume{}
var _ resourcestrategy.Validater = &Volume{}

func (in *Volume) GetObjectMeta() *metav1.ObjectMeta {
	return &in.ObjectMeta
}

func (in *Volume) NamespaceScoped() bool {
	return false
}

func (in *Volume) New() runtime.Object {
	return &Volume{}
}

func (in *Volume) NewList() runtime.Object {
	return &VolumeList{}
}

func (in *Volume) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "storage.sample.kubernetes.io",
		Version:  "v1",
		Resource: "volumes",
	}
}

func (in *Volume) IsStorageVersion() bool {
	return true
}

func (in *Volume) Validate(ctx context.Context) field.ErrorList {
	// TODO(user): Modify it, adding your API validation here.
	return nil
}

var _ resource.ObjectList = &VolumeList{}

func (in *VolumeList) GetListMeta() *metav1.ListMeta {
	return &in.ListMeta
}
// VolumeStatus defines the observed state of Volume
type VolumeStatus struct {
}

func (in VolumeStatus) SubResourceName() string {
	return "status"
}

// Volume implements ObjectWithStatusSubResource interface.
var _ resource.ObjectWithStatusSubResource = &Volume{}

func (in *Volume) GetStatus() resource.StatusSubResource {
	return in.Status
}

// VolumeStatus{} implements StatusSubResource interface.
var _ resource.StatusSubResource = &VolumeStatus{}

func (in VolumeStatus) CopyTo(parent resource.ObjectWithStatusSubResource) {
	parent.(*Volume).Status = in
}
