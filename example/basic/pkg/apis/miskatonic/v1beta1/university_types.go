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
	"context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/klog"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource/resourcestrategy"
)

// Generating code from university_types.go file will generate storage and status REST endpoints for
// University.

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +k8s:openapi-gen=true
// +resource:path=universities,strategy=UniversityStrategy
// +subresource:request=UniversityCampus,path=campus,kind=UniversityCampus
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

	// The unversioned struct definition for this field must be manually defined in the group package
	Manual ManualCreateUnversionedType

	// The unversioned struct definition for this field is automatically generated in the group package
	Automatic AutomaticCreateUnversionedType

	Template *corev1.PodSpec `json:"template,omitempty"`

	ServiceSpec corev1.ServiceSpec `json:"service_spec,omitempty"`

	Rollout []appsv1.Deployment `json:"rollout,omitempty"`
}

// Require that the unversioned struct is manually created.  This is *NOT* the default behavior for
// structs appearing as fields in a resource that are defined in the same package as that resource,
// but is explicitly configured through the +genregister comment.
// +genregister:unversioned=false
type ManualCreateUnversionedType struct {
	A string
	B bool
}

// Automatically create an unversioned copy of this struct by copying its definition
// This is the default behavior for structs appearing as fields in a resource and that are defined in the
// same package as that resource.
type AutomaticCreateUnversionedType struct {
	A string
	B bool
}

// UniversityStatus defines the observed state of University
type UniversityStatus struct {
	// enrolled_students is the number of currently enrolled students
	EnrolledStudents []string `json:"enrolled_students,omitempty"`

	// statusfield provides status information about University
	FacultyEmployed []string `json:"faculty_employed,omitempty"`
}

var _ resource.Object = &University{}
var _ resourcestrategy.Defaulter = &University{}
var _ resourcestrategy.Validater = &University{}
var _ resource.ObjectWithStatusSubResource = &University{}
var _ resource.ObjectList = &UniversityList{}

func (in *University) GetObjectMeta() *metav1.ObjectMeta {
	return &in.ObjectMeta
}

func (in *University) NamespaceScoped() bool {
	return true
}

func (in *University) New() runtime.Object {
	return &University{}
}

func (in *University) NewList() runtime.Object {
	return &UniversityList{}
}

func (in *University) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "miskatonic.k8s.io",
		Version:  "v1beta1",
		Resource: "universities",
	}
}

func (in *University) IsStorageVersion() bool {
	return true
}

func (in *University) Default() {
	klog.Infof("Defaulting University %s", in.Name)
	if in.Spec.MaxStudents == nil {
		n := 15
		in.Spec.MaxStudents = &n
	}
}

func (in *University) Validate(ctx context.Context) field.ErrorList {
	klog.Infof("Validating University %s\n", in.Name)
	errors := field.ErrorList{}
	if in.Spec.MaxStudents == nil || *in.Spec.MaxStudents < 1 || *in.Spec.MaxStudents > 150 {
		errors = append(errors, field.Invalid(
			field.NewPath("spec", "MaxStudents"),
			*in.Spec.MaxStudents,
			"Must be between 1 and 150"))
	}
	return errors
}

func (in *University) SetStatus(statusSubResource interface{}) {
	in.Status = statusSubResource.(UniversityStatus)
}

func (in *University) GetStatus() (statusSubResource interface{}) {
	return in.Status
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type UniversityList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []University `json:"items"`
}

func (in *UniversityList) GetListMeta() *metav1.ListMeta {
	return &in.ListMeta
}
