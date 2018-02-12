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

package codegen

import (
	"k8s.io/gengo/generator"
)

// ResourceGenerator provides a code generator that takes a package of an API GroupVersion
// and generates a file
type ResourceGenerator interface {
	// Returns a Generator for a versioned resource package e.g. pkg/apis/<group>/<version>
	GenerateVersionedResource(
		apiversion *APIVersion, apigroup *APIGroup, filename string) generator.Generator

	// GenerateUnversionedResource returns a Generator for an unversioned resource package e.g. pkg/apis/<group>
	GenerateUnversionedResource(apigroup *APIGroup, filename string) generator.Generator

	GenerateInstall(apigroup *APIGroup, filename string) generator.Generator

	// GenerateAPIs returns a Generator for the apis package e.g. pkg/apis
	GenerateAPIs(apis *APIs, filename string) generator.Generator
}

// ControllerGenerator provides a code generator that takes a package of a controller
// and generates a file
type ControllerGenerator interface {
	// GenerateController returns a Generator for a controller for a specific resource e.g. pkg/controller/<resource>
	GenerateController(controller Controller, filename string) generator.Generator

	// GenerateControllers returns a Generator for the controller package e.g. pkg/controller
	GenerateControllers(controllers []Controller, filename string) generator.Generator

	// GenerateInformers returns a Generator for the sharedinformers package e.g. pkg/controller/sharedinformers
	GenerateInformers(controllers []Controller, apis *APIs, filename string) generator.Generator
}
