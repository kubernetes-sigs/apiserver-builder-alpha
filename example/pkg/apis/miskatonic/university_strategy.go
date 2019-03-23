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

package miskatonic

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"context"
	"log"
)

// Resource Validation
func (UniversityStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	university := obj.(*University)
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
