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

package olympus

import (
	"context"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/klog"
)

// Validate checks that an instance of Poseidon is well formed
func (PoseidonStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	o := obj.(*Poseidon)
	klog.Infof("Validating fields for Poseidon %s\n", o.Name)
	errors := field.ErrorList{}
	// perform validation here and add to errors using field.Invalid
	return errors
}

func (b PoseidonStrategy) GetTriggerFuncs() storage.IndexerFuncs {
	// Change this function to override the trigger fn that is used
	value := b.DefaultStorageStrategy.GetTriggerFuncs()
	return value
}

func (b PoseidonStrategy) BasicMatch(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:       label,
		Field:       field,
		GetAttrs:    b.GetAttrs,
		IndexFields: []string{"spec.deployment.name"},
	}
}

// The following functions allow spec.deployment.name to be selected when listing
// or watching resources
func (b PoseidonStrategy) GetAttrs(o runtime.Object) (labels.Set, fields.Set, error) {
	// Change this function to override the attributes that are matched
	l, _, e := b.DefaultStorageStrategy.GetAttrs(o)
	obj := o.(*Poseidon)

	fs := fields.Set{"spec.deployment.name": obj.Spec.Deployment.Name}
	fs = generic.AddObjectMetaFieldsSet(fs, &obj.ObjectMeta, true)
	return l, fs, e
}
