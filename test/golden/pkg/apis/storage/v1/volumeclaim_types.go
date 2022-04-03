
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

// VolumeClaim
// +k8s:openapi-gen=true
type VolumeClaim struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VolumeClaimSpec   `json:"spec,omitempty"`
	Status VolumeClaimStatus `json:"status,omitempty"`
}

// VolumeClaimList
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type VolumeClaimList struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []VolumeClaim `json:"items"`
}

// VolumeClaimSpec defines the desired state of VolumeClaim
type VolumeClaimSpec struct {
}

var _ resource.Object = &VolumeClaim{}
var _ resourcestrategy.Validater = &VolumeClaim{}

func (in *VolumeClaim) GetObjectMeta() *metav1.ObjectMeta {
	return &in.ObjectMeta
}

func (in *VolumeClaim) NamespaceScoped() bool {
	return false
}

func (in *VolumeClaim) New() runtime.Object {
	return &VolumeClaim{}
}

func (in *VolumeClaim) NewList() runtime.Object {
	return &VolumeClaimList{}
}

func (in *VolumeClaim) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "storage.sample.kubernetes.io",
		Version:  "v1",
		Resource: "volumeclaims",
	}
}

func (in *VolumeClaim) IsStorageVersion() bool {
	return true
}

func (in *VolumeClaim) Validate(ctx context.Context) field.ErrorList {
	// TODO(user): Modify it, adding your API validation here.
	return nil
}

var _ resource.ObjectList = &VolumeClaimList{}

func (in *VolumeClaimList) GetListMeta() *metav1.ListMeta {
	return &in.ListMeta
}
// VolumeClaimStatus defines the observed state of VolumeClaim
type VolumeClaimStatus struct {
}

func (in VolumeClaimStatus) SubResourceName() string {
	return "status"
}

// VolumeClaim implements ObjectWithStatusSubResource interface.
var _ resource.ObjectWithStatusSubResource = &VolumeClaim{}

func (in *VolumeClaim) GetStatus() resource.StatusSubResource {
	return in.Status
}

// VolumeClaimStatus{} implements StatusSubResource interface.
var _ resource.StatusSubResource = &VolumeClaimStatus{}

func (in VolumeClaimStatus) CopyTo(parent resource.ObjectWithStatusSubResource) {
	parent.(*VolumeClaim).Status = in
}
