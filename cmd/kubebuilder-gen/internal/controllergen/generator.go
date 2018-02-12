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

package controllergen

import (
	"github.com/kubernetes-sigs/kubebuilder/cmd/kubebuilder-gen/codegen"
	"k8s.io/gengo/generator"
)

type Generator struct{}

// GenerateController returns a Generator for a controller for a specific resource e.g. pkg/controller/<resource>
func (g *Generator) GenerateController(controller codegen.Controller, filename string) generator.Generator {
	return &controllerGenerator{
		generator.DefaultGen{OptionalName: filename},
		controller,
	}
}

// GenerateControllers returns a Generator for the controller package e.g. pkg/controller
func (g *Generator) GenerateControllers(controllers []codegen.Controller, filename string) generator.Generator {
	return &controllersGenerator{
		generator.DefaultGen{OptionalName: filename},
		controllers,
	}
}

// GenerateInformers returns a Generator for the sharedinformers package e.g. pkg/controller/sharedinformers
func (g *Generator) GenerateInformers(controllers []codegen.Controller, apis *codegen.APIs, filename string) generator.Generator {
	return &informersGenerator{
		generator.DefaultGen{OptionalName: filename},
		controllers,
		apis,
	}
}
