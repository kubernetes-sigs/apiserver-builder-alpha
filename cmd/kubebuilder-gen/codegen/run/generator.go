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

package run

import (
	"github.com/golang/glog"
	"github.com/najena/kubebuilder/cmd/kubebuilder-gen/codegen"
	"github.com/najena/kubebuilder/cmd/kubebuilder-gen/codegen/internal/parse"
	"k8s.io/gengo/args"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"path/filepath"
)

// CodeGenerator generates code for Kubernetes resources and controllers
type CodeGenerator struct {
	resourceGenerators   []codegen.ResourceGenerator
	controllerGenerators []codegen.ControllerGenerator

	// OutputFileBaseName is the base name used for output files
	OutputFileBaseName string
}

// AddControllerGenerator adds a controller generator that will be called with parsed controllers
func (g *CodeGenerator) AddControllerGenerator(generator codegen.ControllerGenerator) *CodeGenerator {
	g.controllerGenerators = append(g.controllerGenerators, generator)
	return g
}

// AddResourceGenerator adds a resource generator that will be called with parsed resources
func (g *CodeGenerator) AddResourceGenerator(generator codegen.ResourceGenerator) *CodeGenerator {
	g.resourceGenerators = append(g.resourceGenerators, generator)
	return g
}

type customArgs struct{}

// Execute parses packages and executes the code generators against the resource and controller packages
func (g *CodeGenerator) Execute() error {
	arguments := args.Default()
	arguments.CustomArgs = &customArgs{}
	arguments.OutputFileBaseName = g.OutputFileBaseName

	err := arguments.Execute(nameSystems(), defaultNameSystem(), g.packages)
	if err != nil {
		glog.Fatalf("Error: %v", err)
	}
	return nil
}

// packages parses the observed packages and creates code generators
func (g *CodeGenerator) packages(context *generator.Context, arguments *args.GeneratorArgs) generator.Packages {
	p := packages{}

	b := parse.NewAPIs(context, arguments)

	// Add resource generators
	for _, apigroup := range b.APIs.Groups {
		for _, apiversion := range apigroup.Versions {
			// Do the versioned resource packages
			for _, r := range g.resourceGenerators {
				g := r.GenerateVersionedResource(apiversion, apigroup, arguments.OutputFileBaseName)
				if g != nil {
					p.add(apiversion.Pkg.Path, g)
				}
			}
		}
		// Do the unversioned packages
		for _, r := range g.resourceGenerators {
			g := r.GenerateUnversionedResource(apigroup, arguments.OutputFileBaseName)
			if g != nil {
				p.add(apigroup.Pkg.Path, g)
			}
		}
		// Do the install generators
		for _, r := range g.resourceGenerators {
			g := r.GenerateInstall(apigroup, arguments.OutputFileBaseName)
			if g != nil {
				p.add(filepath.Join(apigroup.Pkg.Path, "install"), g)
			}
		}
	}
	// Do apis package
	for _, r := range g.resourceGenerators {
		g := r.GenerateAPIs(b.APIs, arguments.OutputFileBaseName)
		if g != nil {
			p.add(b.APIs.Pkg.Path, g)
		}
	}

	// Add generators for Controllers.
	// Do each controller package
	repo := ""
	for _, c := range b.Controllers {
		repo = c.Repo
		for _, cg := range g.controllerGenerators {
			g := cg.GenerateController(c, arguments.OutputFileBaseName)
			if g != nil {
				p.add(c.Pkg.Path, g)
			}

		}
	}
	// Do controller package
	if len(b.Controllers) > 0 {
		for _, cg := range g.controllerGenerators {
			g := cg.GenerateControllers(b.Controllers, arguments.OutputFileBaseName)
			if g != nil {
				p.add(context.Universe[repo+"/pkg/controller"].Path, g)
			}
		}
	}
	// Do sharedinformer package
	if len(b.Controllers) > 0 {
		for _, cg := range g.controllerGenerators {
			g := cg.GenerateInformers(b.Controllers, b.APIs, arguments.OutputFileBaseName)
			if g != nil {
				p.add(context.Universe[repo+"/pkg/controller/sharedinformers"].Path, g)
			}
		}
	}
	return p.value
}

// defaultNameSystem returns the default name system for ordering the types to be
// processed by the generators in this package.
// e.g. black-magic
func defaultNameSystem() string {
	return "public"
}

// nameSystems returns the name system used by the generators in this package.
// e.g. black-magic
func nameSystems() namer.NameSystems {
	return namer.NameSystems{
		"public": namer.NewPublicNamer(1),
		"raw":    namer.NewRawNamer("", nil),
	}
}
