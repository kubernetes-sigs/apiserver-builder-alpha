
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
