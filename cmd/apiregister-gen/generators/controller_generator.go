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

package generators

import (
	"io"
	"strings"
	"text/template"

	"k8s.io/gengo/generator"
)

type controllerGenerator struct {
	generator.DefaultGen
	controller Controller
}

var _ generator.Generator = &controllerGenerator{}

func CreateControllerGenerator(controller Controller, filename string) generator.Generator {
	return &controllerGenerator{
		generator.DefaultGen{OptionalName: filename},
		controller,
	}
}

func (d *controllerGenerator) Imports(c *generator.Context) []string {
	im := []string{
		"github.com/golang/glog",
		"github.com/kubernetes-incubator/apiserver-builder/pkg/controller",
		"k8s.io/apimachinery/pkg/api/errors",
		"k8s.io/client-go/rest",
		"k8s.io/client-go/tools/cache",
		"k8s.io/client-go/util/workqueue",
		d.controller.Repo + "/pkg/controller/sharedinformers",
	}

	return im
}

func (d *controllerGenerator) Finalize(context *generator.Context, w io.Writer) error {
	temp := template.Must(template.New("controller-template").Funcs(
		template.FuncMap{
			"title": strings.Title,
		},
	).Parse(ControllerAPITemplate))
	return temp.Execute(w, d.controller)
}

var ControllerAPITemplate = `
// {{.Target.Kind}}Controller implements the controller.{{.Target.Kind}}Controller interface
type {{.Target.Kind}}Controller struct {
	queue *controller.QueueWorker

	// Handles messages
	controller *{{.Target.Kind}}ControllerImpl

	Name string
}

// NewController returns a new {{.Target.Kind}}Controller for responding to {{.Target.Kind}} events
func New{{.Target.Kind}}Controller(config *rest.Config, si *sharedinformers.SharedInformers) *{{.Target.Kind}}Controller {
	q := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "{{.Target.Kind}}")

	// For non-generated code to add events
	uc := &{{.Target.Kind}}ControllerImpl{}
	uc.Init(config, si, q)

	queue := &controller.QueueWorker{q, 10, "{{.Target.Kind}}", nil}
	c := &{{.Target.Kind}}Controller{queue, uc, "{{.Target.Kind}}"}
	queue.Reconcile = c.reconcile
	return c
}

func (c *{{.Target.Kind}}Controller) GetName() string {
	return c.Name
}

func (c *{{.Target.Kind}}Controller) reconcile(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	u, err := c.controller.Get(namespace, name)
	if errors.IsNotFound(err) {
		glog.Infof("Not doing work for {{.Target.Kind}} %v because it has been deleted", key)
		return nil
	}
	if err != nil {
		glog.Errorf("Unable to retrieve {{.Target.Kind}} %v from store: %v", key, err)
		return err
	}

	return c.controller.Reconcile(u)
}

func (c *{{.Target.Kind}}Controller) Run(stopCh <-chan struct{}) {
	c.queue.Run(stopCh)
}
`

type allControllerGenerator struct {
	generator.DefaultGen
	Controllers []Controller
}

var _ generator.Generator = &allControllerGenerator{}

func CreateAllControllerGenerator(controllers []Controller, filename string) generator.Generator {
	return &allControllerGenerator{
		generator.DefaultGen{OptionalName: filename},
		controllers,
	}
}

func (d *allControllerGenerator) Imports(c *generator.Context) []string {
	if len(d.Controllers) == 0 {
		return []string{}
	}

	repo := d.Controllers[0].Repo
	im := []string{
		"k8s.io/client-go/rest",
		"github.com/kubernetes-incubator/apiserver-builder/pkg/controller",
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

func (d *allControllerGenerator) Finalize(context *generator.Context, w io.Writer) error {
	temp := template.Must(template.New("all-controller-template").Funcs(
		template.FuncMap{
			"title": strings.Title,
		},
	).Parse(AllControllerAPITemplate))
	return temp.Execute(w, d)
}

var AllControllerAPITemplate = `

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

type informersGenerator struct {
	generator.DefaultGen
	Controllers []Controller
}

var _ generator.Generator = &informersGenerator{}

func CreateInformersGenerator(controllers []Controller, filename string) generator.Generator {
	return &informersGenerator{
		generator.DefaultGen{OptionalName: filename},
		controllers,
	}
}

func (d *informersGenerator) Imports(c *generator.Context) []string {
	if len(d.Controllers) == 0 {
		return []string{}
	}

	repo := d.Controllers[0].Repo
	return []string{
		"time",
		"k8s.io/client-go/rest",
		repo + "/pkg/client/clientset_generated/clientset",
		repo + "/pkg/client/informers_generated/externalversions",
	}
}

func (d *informersGenerator) Finalize(context *generator.Context, w io.Writer) error {
	temp := template.Must(template.New("informersGenerator-template").Funcs(
		template.FuncMap{
			"title": strings.Title,
		},
	).Parse(InformersTemplate))
	return temp.Execute(w, d.Controllers)
}

var InformersTemplate = `
// SharedInformers wraps all informers used by controllers so that
// they are shared across controller implementations
type SharedInformers struct {
	Factory externalversions.SharedInformerFactory
}

// newSharedInformers returns a set of started informers
func NewSharedInformers(config *rest.Config, stop <-chan struct{}) *SharedInformers {
	cs := clientset.NewForConfigOrDie(config)
	si := &SharedInformers{externalversions.NewSharedInformerFactory(cs, 10*time.Minute)}
	si.startInformers(stop)
	return si
}

// startInformers starts all of the informers
func (si *SharedInformers) startInformers(stop <-chan struct{}) {
	shutdown := make(chan struct{})
	{{ range $c := . -}}
	go si.Factory.{{title $c.Target.Group}}().{{title $c.Target.Version}}().{{title $c.Resource}}().Informer().Run(shutdown)
	{{ end -}}
	go func() {
		m := <-stop
		{{ range $c := . -}}
		shutdown <- m
		{{ end -}}
	}()
}

`
