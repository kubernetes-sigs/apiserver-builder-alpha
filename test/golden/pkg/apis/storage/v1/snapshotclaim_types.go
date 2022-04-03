
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

// SnapshotClaim
// +k8s:openapi-gen=true
type SnapshotClaim struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SnapshotClaimSpec   `json:"spec,omitempty"`
	Status SnapshotClaimStatus `json:"status,omitempty"`
}

// SnapshotClaimList
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SnapshotClaimList struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []SnapshotClaim `json:"items"`
}

// SnapshotClaimSpec defines the desired state of SnapshotClaim
type SnapshotClaimSpec struct {
}

var _ resource.Object = &SnapshotClaim{}
var _ resourcestrategy.Validater = &SnapshotClaim{}

func (in *SnapshotClaim) GetObjectMeta() *metav1.ObjectMeta {
	return &in.ObjectMeta
}

func (in *SnapshotClaim) NamespaceScoped() bool {
	return false
}

func (in *SnapshotClaim) New() runtime.Object {
	return &SnapshotClaim{}
}

func (in *SnapshotClaim) NewList() runtime.Object {
	return &SnapshotClaimList{}
}

func (in *SnapshotClaim) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "storage.sample.kubernetes.io",
		Version:  "v1",
		Resource: "snapshotclaims",
	}
}

func (in *SnapshotClaim) IsStorageVersion() bool {
	return true
}

func (in *SnapshotClaim) Validate(ctx context.Context) field.ErrorList {
	// TODO(user): Modify it, adding your API validation here.
	return nil
}

var _ resource.ObjectList = &SnapshotClaimList{}

func (in *SnapshotClaimList) GetListMeta() *metav1.ListMeta {
	return &in.ListMeta
}
// SnapshotClaimStatus defines the observed state of SnapshotClaim
type SnapshotClaimStatus struct {
}

func (in SnapshotClaimStatus) SubResourceName() string {
	return "status"
}

// SnapshotClaim implements ObjectWithStatusSubResource interface.
var _ resource.ObjectWithStatusSubResource = &SnapshotClaim{}

func (in *SnapshotClaim) GetStatus() resource.StatusSubResource {
	return in.Status
}

// SnapshotClaimStatus{} implements StatusSubResource interface.
var _ resource.StatusSubResource = &SnapshotClaimStatus{}

func (in SnapshotClaimStatus) CopyTo(parent resource.ObjectWithStatusSubResource) {
	parent.(*SnapshotClaim).Status = in
}
