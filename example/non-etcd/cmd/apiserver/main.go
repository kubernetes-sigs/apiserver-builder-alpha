/*
Copyright 2020 The Kubernetes Authors.

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
	"sigs.k8s.io/apiserver-runtime/pkg/experimental/storage/filepath"

	filepathv1 "sigs.k8s.io/apiserver-builder-alpha/example/non-etcd/pkg/apis/filepath/v1"
)

func main() {
	err := builder.APIServer.
		// writes burger resources as static files under the "data" folder in the working directory.
		WithResourceAndHandler(&filepathv1.Burger{}, filepath.NewJsonFilepathStorageProvider(&filepathv1.Burger{}, "data")).
		WithOptionsFns(func(o *builder.ServerOptions) *builder.ServerOptions {
			o.RecommendedOptions.Authorization = nil
			o.RecommendedOptions.Admission = nil
			return nil
		}).
		Execute()
	if err != nil {
		klog.Fatal(err)
	}
}
