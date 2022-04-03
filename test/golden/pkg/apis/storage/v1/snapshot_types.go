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

// Snapshot
// +k8s:openapi-gen=true
type Snapshot struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SnapshotSpec   `json:"spec,omitempty"`
	Status SnapshotStatus `json:"status,omitempty"`
}

// SnapshotList
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SnapshotList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Snapshot `json:"items"`
}

// SnapshotSpec defines the desired state of Snapshot
type SnapshotSpec struct {
}

var _ resource.Object = &Snapshot{}
var _ resourcestrategy.Validater = &Snapshot{}

func (in *Snapshot) GetObjectMeta() *metav1.ObjectMeta {
	return &in.ObjectMeta
}

func (in *Snapshot) NamespaceScoped() bool {
	return false
}

func (in *Snapshot) New() runtime.Object {
	return &Snapshot{}
}

func (in *Snapshot) NewList() runtime.Object {
	return &SnapshotList{}
}

func (in *Snapshot) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "storage.sample.kubernetes.io",
		Version:  "v1",
		Resource: "snapshots",
	}
}

func (in *Snapshot) IsStorageVersion() bool {
	return true
}

func (in *Snapshot) Validate(ctx context.Context) field.ErrorList {
	// TODO(user): Modify it, adding your API validation here.
	return nil
}

var _ resource.ObjectList = &SnapshotList{}

func (in *SnapshotList) GetListMeta() *metav1.ListMeta {
	return &in.ListMeta
}

// SnapshotStatus defines the observed state of Snapshot
type SnapshotStatus struct {
}

func (in SnapshotStatus) SubResourceName() string {
	return "status"
}

// Snapshot implements ObjectWithStatusSubResource interface.
var _ resource.ObjectWithStatusSubResource = &Snapshot{}

func (in *Snapshot) GetStatus() resource.StatusSubResource {
	return in.Status
}

// SnapshotStatus{} implements StatusSubResource interface.
var _ resource.StatusSubResource = &SnapshotStatus{}

func (in SnapshotStatus) CopyTo(parent resource.ObjectWithStatusSubResource) {
	parent.(*Snapshot).Status = in
}

var _ resource.ObjectWithArbitrarySubResource = &Snapshot{}

func (in *Snapshot) GetArbitrarySubResources() []resource.ArbitrarySubResource {
	return []resource.ArbitrarySubResource{
		// +kubebuilder:scaffold:subresource
		&SnapshotBar{},
		&SnapshotFoo{},
	}
}
