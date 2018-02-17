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

	"github.com/kubernetes-sigs/kubebuilder/cmd/kubebuilder-gen/codegen"
	"github.com/markbates/inflect"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/gengo/generator"
)

type informersGenerator struct {
	generator.DefaultGen
	Controllers []codegen.Controller
	apis        *codegen.APIs
}

var _ generator.Generator = &informersGenerator{}

func (d *informersGenerator) Imports(c *generator.Context) []string {
	if len(d.Controllers) == 0 {
		return []string{}
	}

	repo := d.Controllers[0].Repo
	return []string{
		"time",
		"github.com/kubernetes-sigs/kubebuilder/pkg/controller",
		"k8s.io/client-go/rest",
		repo + "/pkg/client/clientset_generated/clientset",
		repo + "/pkg/client/informers_generated/externalversions",
		"k8s.io/client-go/tools/cache",
	}
}

func (d *informersGenerator) Finalize(context *generator.Context, w io.Writer) error {
	temp := template.Must(template.New("informersGenerator-template").Funcs(
		template.FuncMap{
			"title":  strings.Title,
			"plural": inflect.NewDefaultRuleset().Pluralize,
		},
	).Parse(informersTemplate))

	gvks := []schema.GroupVersionKind{}
	for _, g := range d.apis.Groups {
		for _, v := range g.Versions {
			for _, r := range v.Resources {
				gvks = append(gvks, schema.GroupVersionKind{
					Group:   r.Group,
					Version: r.Version,
					Kind:    r.Kind,
				})
			}
		}
	}
	return temp.Execute(w, gvks)
}

var informersTemplate = `
// SharedInformers wraps all informers used by controllers so that
// they are shared across controller implementations
type SharedInformers struct {
    controller.SharedInformersDefaults
    Factory           externalversions.SharedInformerFactory
}

// newSharedInformers returns a set of started informers
func NewSharedInformers(config *rest.Config, shutdown <-chan struct{}) *SharedInformers {
    si := &SharedInformers{
        controller.SharedInformersDefaults{},
        externalversions.NewSharedInformerFactory(clientset.NewForConfigOrDie(config), 10*time.Minute),
    }
    if si.SetupKubernetesTypes() {
        si.InitKubernetesInformers(config)
    }
    si.Init()
    si.startInformers(shutdown)
    si.StartAdditionalInformers(shutdown)
    return si
}

// startInformers starts all of the informers
func (si *SharedInformers) startInformers(shutdown <-chan struct{}) {
    {{ range $c := . -}}
    go si.Factory.{{title $c.Group}}().{{title $c.Version}}().{{plural $c.Kind}}().Informer().Run(shutdown)
    {{ end -}}
}

// ControllerInitArguments are arguments provided to the Init function for a new controller.
type ControllerInitArguments interface {
    // GetSharedInformers returns the SharedInformers that can be used to access
    // informers and listers for watching and indexing Kubernetes Resources
    GetSharedInformers() *SharedInformers

    // GetRestConfig returns the Config to create new client-go clients
    GetRestConfig() *rest.Config

    // Watch uses resourceInformer to watch a resource.  When create, update, or deletes
    // to the resource type are encountered, watch uses watchResourceToReconcileResourceKey
    // to lookup the key for the resource reconciled by the controller (maybe a different type
    // than the watched resource), and enqueue it to be reconciled.
    // watchName: name of the informer.  may appear in logs
    // resourceInformer: gotten from the SharedInformer.  controls which resource type is watched
    // getReconcileKeys: takes an instance of the watched resource and returns
    //                   a slice of keys for the reconciled resource type to enqueue.
    Watch(watchName string, resourceInformer cache.SharedIndexInformer,
            getReconcileKeys func(interface{}) ([]string, error))
}

type ControllerInitArgumentsImpl struct {
    Si *SharedInformers
    Rc *rest.Config
    Rk func(key string) error
}

func (c ControllerInitArgumentsImpl) GetSharedInformers() *SharedInformers {
  return c.Si
}

func (c ControllerInitArgumentsImpl) GetRestConfig() *rest.Config {
  return c.Rc
}

// Watch uses resourceInformer to watch a resource.  When create, update, or deletes
// to the resource type are encountered, watch uses watchResourceToReconcileResourceKey
// to lookup the key for the resource reconciled by the controller (maybe a different type
// than the watched resource), and enqueue it to be reconciled.
// watchName: name of the informer.  may appear in logs
// resourceInformer: gotten from the SharedInformer.  controls which resource type is watched
// getReconcileKey: takes an instance of the watched resource and returns
//                                      a key for the reconciled resource type to enqueue.
func (c ControllerInitArgumentsImpl) Watch(
    watchName string, resourceInformer cache.SharedIndexInformer,
    getReconcileKey func(interface{}) ([]string, error)) {
    c.Si.Watch(watchName, resourceInformer, getReconcileKey, c.Rk)
}

type Controller interface {}

// ControllerInit new controllers should implement this.  It is more flexible in
// allowing additional options to be passed in
type ControllerInit interface {
    Init(args ControllerInitArguments)
}
`
