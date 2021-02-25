/*
Copyright 2016 The Kubernetes Authors.

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

package main

import (
	"k8s.io/klog"
	"sigs.k8s.io/apiserver-runtime/pkg/builder"

	innsmouthv1 "sigs.k8s.io/apiserver-builder-alpha/example/basic/pkg/apis/innsmouth/v1"
	kingsportv1 "sigs.k8s.io/apiserver-builder-alpha/example/basic/pkg/apis/kingsport/v1"
	miskatonicv1beta1 "sigs.k8s.io/apiserver-builder-alpha/example/basic/pkg/apis/miskatonic/v1beta1"
	olympusv1beta1 "sigs.k8s.io/apiserver-builder-alpha/example/basic/pkg/apis/olympus/v1beta1"
)

func main() {
	err := builder.APIServer.
		WithResource(&innsmouthv1.DeepOne{}).          // namespaced resource
		WithResource(&kingsportv1.Festival{}).         // cluster-scoped resource
		WithResource(&miskatonicv1beta1.Student{}).    // resource with arbitrary subresource and custom storage
		WithResource(&miskatonicv1beta1.University{}). // resource with arbitrary subresource
		WithResource(&olympusv1beta1.Poseidon{}).      // resource with custom storage indexers
		WithLocalDebugExtension().
		WithOptionsFns(func(options *builder.ServerOptions) *builder.ServerOptions {
			options.RecommendedOptions.CoreAPI = nil
			options.RecommendedOptions.Admission = nil
			return options
		}).
		Execute()

	if err != nil {
		klog.Fatal(err)
	}
}
