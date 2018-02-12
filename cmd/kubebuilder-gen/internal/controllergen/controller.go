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

type controllerGenerator struct {
	generator.DefaultGen
	controller codegen.Controller
}

var _ generator.Generator = &controllerGenerator{}

func (d *controllerGenerator) Imports(c *generator.Context) []string {
	im := []string{
		"github.com/golang/glog",
		"github.com/najena/kubebuilder/pkg/controller",
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
			"title":  strings.Title,
			"plural": inflect.NewDefaultRuleset().Pluralize,
		},
	).Parse(controllerAPITemplate))
	return temp.Execute(w, d.controller)
}

var controllerAPITemplate = `
// {{.Target.Kind}}Controller implements the controller.{{.Target.Kind}}Controller interface
type {{.Target.Kind}}Controller struct {
    queue *controller.QueueWorker

    // Handles messages
    controller *{{.Target.Kind}}ControllerImpl

    Name string

    BeforeReconcile func(key string)
    AfterReconcile  func(key string, err error)

    Informers *sharedinformers.SharedInformers
}

// NewController returns a new {{.Target.Kind}}Controller for responding to {{.Target.Kind}} events
func New{{.Target.Kind}}Controller(config *rest.Config, si *sharedinformers.SharedInformers) *{{.Target.Kind}}Controller {
    q := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "{{.Target.Kind}}")

    queue := &controller.QueueWorker{q, 10, "{{.Target.Kind}}", nil}
    c := &{{.Target.Kind}}Controller{queue, nil, "{{.Target.Kind}}", nil, nil, si}

    // For non-generated code to add events
    uc := &{{.Target.Kind}}ControllerImpl{}
    var ci sharedinformers.Controller = uc

    if i, ok := ci.(sharedinformers.ControllerInit); ok {
        i.Init(&sharedinformers.ControllerInitArgumentsImpl{si, config, c.LookupAndReconcile})
    }

    c.controller = uc

    queue.Reconcile = c.LookupAndReconcile
    if c.Informers.WorkerQueues == nil {
        c.Informers.WorkerQueues = map[string]*controller.QueueWorker{}
    }
    c.Informers.WorkerQueues["{{.Target.Kind}}"] = queue
    si.Factory.{{title .Target.Group}}().{{title .Target.Version}}().{{plural .Target.Kind }}().Informer().
        AddEventHandler(&controller.QueueingEventHandler{q, nil, false})
    return c
}

func (c *{{.Target.Kind}}Controller) GetName() string {
    return c.Name
}

func (c *{{.Target.Kind}}Controller) LookupAndReconcile(key string) (err error) {
    var namespace, name string

    if c.BeforeReconcile != nil {
        c.BeforeReconcile(key)
    }
    if c.AfterReconcile != nil {
        // Wrap in a function so err is evaluated after it is set
        defer func() { c.AfterReconcile(key, err) }()
    }

    namespace, name, err = cache.SplitMetaNamespaceKey(key)
    if err != nil {
        return
    }

    u, err := c.controller.Get(namespace, name)
    if errors.IsNotFound(err) {
        glog.Infof("Not doing work for {{.Target.Kind}} %v because it has been deleted", key)
        // Set error so it is picked up by AfterReconcile and the return function
        err = nil
        return
    }
    if err != nil {
        glog.Errorf("Unable to retrieve {{.Target.Kind}} %v from store: %v", key, err)
        return
    }

    // Set error so it is picked up by AfterReconcile and the return function
    err = c.controller.Reconcile(u)

    return
}

func (c *{{.Target.Kind}}Controller) Run(stopCh <-chan struct{}) {
    for _, q := range c.Informers.WorkerQueues {
        q.Run(stopCh)
    }
    controller.GetDefaults(c.controller).Run(stopCh)
}
`
