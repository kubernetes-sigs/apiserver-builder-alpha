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

	"github.com/kubernetes-incubator/apiserver-builder/example/pkg/apis/miskatonic"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
	//extensionsv1beta1 "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

// Generating code from university_types.go file will generate storage and status REST endpoints for
// University.

// +genclient=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +k8s:openapi-gen=true
// +resource:path=universities,strategy=UniversityStrategy
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

	// The unversioned struct definition for this field must be manually defined in the group package
	Manual ManualCreateUnversionedType

	// The unversioned struct definition for this field is automatically generated in the group package
	Automatic AutomaticCreateUnversionedType

	//// WARNING: Using types from client-go as fields does not work outside this example
	//// This example hacked the vendored client-go to add the openapi generation directives
	//// to make this work
	//Template *apiv1.PodSpec `json:"template,omitempty"`
	//
	//// WARNING: Using types from client-go as fields does not work outside this example
	//// This example hacked the vendored client-go to add the openapi generation directives
	//// to make this work
	//ServiceSpec apiv1.ServiceSpec `json:"service_spec,omitempty"`
	//
	//// WARNING: Using types from client-go as fields does not work outside this example
	//// This example hacked the vendored client-go to add the openapi generation directives
	//// to make this work
	//Rollout []extensionsv1beta1.Deployment `json:"rollout,omitempty"`
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

// Resource Validation
func (UniversityStrategy) Validate(ctx request.Context, obj runtime.Object) field.ErrorList {
	university := obj.(*miskatonic.University)
	log.Printf("Validating University %s\n", university.Name)
	errors := field.ErrorList{}
	if university.Spec.MaxStudents == nil || *university.Spec.MaxStudents < 1 || *university.Spec.MaxStudents > 150 {
		errors = append(errors, field.Invalid(
			field.NewPath("spec", "MaxStudents"),
			*university.Spec.MaxStudents,
			"Must be between 1 and 150"))
	}
	return errors
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

// GetConversionFunctions returns functions for converting resource versions to override the
// conversion functions
func (UniversitySchemeFns) GetConversionFunctions() []interface{} {
	return []interface{}{}
}

// +genclient=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +subresource-request
type Scale struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Faculty int `json:"faculty,omitempty"`
}

var _ rest.CreaterUpdater = &ScaleUniversityREST{}
var _ rest.Patcher = &ScaleUniversityREST{}

// +k8s:deepcopy-gen=false
type ScaleUniversityREST struct {
	Registry miskatonic.UniversityRegistry
}

func (r *ScaleUniversityREST) Create(ctx request.Context, obj runtime.Object, includeUninitialized bool) (runtime.Object, error) {
	scale := obj.(*Scale)
	u, err := r.Registry.GetUniversity(ctx, scale.Name, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	u.Spec.FacultySize = scale.Faculty
	r.Registry.UpdateUniversity(ctx, u)
	return u, nil
}

// Get retrieves the object from the storage. It is required to support Patch.
func (r *ScaleUniversityREST) Get(ctx request.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return nil, nil
}

// Update alters the status subset of an object.
func (r *ScaleUniversityREST) Update(ctx request.Context, name string, objInfo rest.UpdatedObjectInfo) (runtime.Object, bool, error) {
	return nil, false, nil
}

func (r *ScaleUniversityREST) New() runtime.Object {
	return &Scale{}
}
