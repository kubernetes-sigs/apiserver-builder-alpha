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
	"io"
	"strings"
	"text/template"

	"github.com/markbates/inflect"
	"github.com/najena/kubebuilder/cmd/kubebuilder-gen/codegen"
	"k8s.io/gengo/generator"
)

type controllersGenerator struct {
	generator.DefaultGen
	Controllers []codegen.Controller
}

var _ generator.Generator = &controllersGenerator{}

func (d *controllersGenerator) Imports(c *generator.Context) []string {
	if len(d.Controllers) == 0 {
		return []string{}
	}

	repo := d.Controllers[0].Repo
	im := []string{
		"k8s.io/client-go/rest",
		"github.com/najena/kubebuilder/pkg/controller",
		repo + "/pkg/controller/sharedinformers",
	}

	// Import package for each controller
	repos := map[string]string{}
	for _, c := range d.Controllers {
		repos[c.Pkg.Path] = ""
	}
	for k, _ := range repos {
		im = append(im, k)
	}

	return im
}

func (d *controllersGenerator) Finalize(context *generator.Context, w io.Writer) error {
	temp := template.Must(template.New("all-controller-template").Funcs(
		template.FuncMap{
			"title":  strings.Title,
			"plural": inflect.NewDefaultRuleset().Pluralize,
		},
	).Parse(controllersAPITemplate))
	return temp.Execute(w, d)
}

var controllersAPITemplate = `

func GetAllControllers(config *rest.Config) ([]controller.Controller, chan struct{}) {
    shutdown := make(chan struct{})
    si := sharedinformers.NewSharedInformers(config, shutdown)
    return []controller.Controller{
        {{ range $c := .Controllers -}}
        {{ $c.Pkg.Name }}.New{{ $c.Target.Kind }}Controller(config, si),
        {{ end -}}
    }, shutdown
}

`
