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
	"sigs.k8s.io/apiserver-builder-alpha/example/pkg/apis"
	"sigs.k8s.io/apiserver-builder-alpha/example/pkg/openapi"
	"sigs.k8s.io/apiserver-builder-alpha/pkg/cmd/server"

	// import the package to install custom admission controllers and custom admission initializers
	_ "sigs.k8s.io/apiserver-builder-alpha/example/plugin/admission/install"
)

func main() {
	err := server.StartApiServerWithOptions(&server.StartOptions{
		EtcdPath:    "/registry/sample.kubernetes.io",
		Apis:        apis.GetAllApiBuilders(),
		Openapidefs: openapi.GetOpenAPIDefinitions,
		Title:       "Api",
		Version:     "v0",

		// TweakConfigFuncs []func(apiServer *apiserver.Config) error
		// FlagConfigFuncs []func(*cobra.Command) error
	})

	if err != nil {
		panic(err)
	}
}
