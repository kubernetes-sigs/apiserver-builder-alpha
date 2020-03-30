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
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

