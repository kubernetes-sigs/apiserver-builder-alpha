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

package kingsport

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/klog"
)

func (FestivalStrategy) NamespaceScoped() bool { return false }

func (FestivalStatusStrategy) NamespaceScoped() bool { return false }

// Validate checks that an instance of Festival is well formed
func (FestivalStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	o := obj.(*Festival)
	klog.Infof("Validating fields for Festival %s", o.Name)
	errors := field.ErrorList{}

	if o.Spec.Year < 0 {
		errors = append(errors,
			field.Invalid(field.NewPath("spec", "year"), o.Spec.Year, "year must be > 0"))
	}

	// perform validation here and add to errors using field.Invalid
	return errors
}
