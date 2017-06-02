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

package v1beta1

import (
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	extensionsv1beta1 "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

// Generating code from university_types.go file will generate storage and status REST endpoints for
// University.

// +genclient=true

// +k8s:openapi-gen=true
// +resource:path=universities
// +subresource:request=Scale,path=scale,rest=ScaleUniversityREST
type University struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   UniversitySpec   `json:"spec,omitempty"`
	Status UniversityStatus `json:"status,omitempty"`
}

// UniversitySpec defines the desired state of University
type UniversitySpec struct {
	// faculty_size defines the desired faculty size of the university.  Defaults to 15.
	FacultySize int `json:"faculty_size,omitempty"`

	// max_students defines the maximum number of enrolled students.  Defaults to 300.
	// +optional
	MaxStudents *int `json:"max_students,omitempty"`

	// WARNING: Using types from client-go as fields does not work outside this example
	// This example hacked the vendored client-go to add the openapi generation directives
	// to make this work
	Template *apiv1.PodSpec `json:"template,omitempty"`

	// WARNING: Using types from client-go as fields does not work outside this example
	// This example hacked the vendored client-go to add the openapi generation directives
	// to make this work
	ServiceSpec apiv1.ServiceSpec `json:"service_spec,omitempty"`

	// WARNING: Using types from client-go as fields does not work outside this example
	// This example hacked the vendored client-go to add the openapi generation directives
	// to make this work
	Rollout []extensionsv1beta1.Deployment `json:"rollout,omitempty"`
}

// UniversityStatus defines the observed state of University
type UniversityStatus struct {
	// enrolled_students is the number of currently enrolled students
	EnrolledStudents []string `json:"enrolled_students,omitempty"`

	// statusfield provides status information about University
	FacultyEmployed []string `json:"faculty_employed,omitempty"`
}

// GetDefaultingFunctions returns functions for defaulting v1beta1.University values
func (UniversitySchemeFns) DefaultingFunction(o interface{}) {
	obj := o.(*University)
	log.Printf("Defaulting University %s\n", obj.Name)
	if obj.Spec.MaxStudents == nil {
		n := 15
		obj.Spec.MaxStudents = &n
	}
}

func (UniversitySchemeFns) GetConversionFunctions() []interface{} {
	return []interface{}{
		apiv1.Convert_api_PodSpec_To_v1_PodSpec,
		apiv1.Convert_v1_PodSpec_To_api_PodSpec,
	}
}

// +genclient=true

// +subresource-request
type Scale struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Faculty int `json:"faculty,omitempty"`
}
