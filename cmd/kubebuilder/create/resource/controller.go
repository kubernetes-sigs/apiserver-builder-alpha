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

package resource

import (
	"fmt"
	"path/filepath"
	"strings"

	createutil "github.com/najena/kubebuilder/cmd/kubebuilder/create/util"
	"github.com/najena/kubebuilder/cmd/kubebuilder/util"
)

func doController(dir string, args resourceTemplateArgs) bool {
	path := filepath.Join(dir, "pkg", "controller", strings.ToLower(createutil.KindName), "controller.go")
	fmt.Printf("\t%s\n", filepath.Join(
		"pkg", "controller", strings.ToLower(createutil.KindName), "controller.go"))
	return util.WriteIfNotFound(path, "resource-controller-template", resourceControllerTemplate, args)
}

var resourceControllerTemplate = `
{{.BoilerPlate}}

package {{ lower .Kind }}

import (
    "log"

    "github.com/najena/kubebuilder/pkg/builders"

    "{{.Repo}}/pkg/apis/{{.Group}}/{{.Version}}"
    "{{.Repo}}/pkg/controller/sharedinformers"
    listers "{{.Repo}}/pkg/client/listers_generated/{{.Group}}/{{.Version}}"
)

// EDIT THIS FILE!
// Created by "kubebuilder create resource" for you to implement controller logic for the {{.Kind}} resource API

// Reconcile handles enqueued messages
func (c *{{.Kind}}ControllerImpl) Reconcile(u *{{.Version}}.{{.Kind}}) error {
    // INSERT YOUR CODE HERE - implement controller logic to reconcile observed and desired state of the object
    log.Printf("Running reconcile {{.Kind}} for %s\n", u.Name)
    return nil
}

// +controller:group={{ .Group }},version={{ .Version }},kind={{ .Kind}},resource={{ .Resource }}
type {{.Kind}}ControllerImpl struct {
    builders.DefaultControllerFns

    // lister indexes properties about {{.Kind}}
    lister listers.{{.Kind}}Lister
}

// Init initializes the controller and is called by the generated code
// Register watches for additional resource types here.
func (c *{{.Kind}}ControllerImpl) Init(arguments sharedinformers.ControllerInitArguments) {
    // INSERT YOUR CODE HERE - add logic for initializing the controller as needed

    // Use the lister for indexing {{.Resource}} labels
    c.lister = arguments.GetSharedInformers().Factory.{{title .Group}}().{{title .Version}}().{{plural .Kind}}().Lister()

    // To watch other resource types, uncomment this function and replace Foo with the resource name to watch.
    // Must define the func FooTo{{.Kind}}(i interface{}) (string, error) {} that returns the {{ .Kind }}
    // "namespace/name"" to reconcile in response to the updated Foo
    // Note: To watch Kubernetes resources, you must also update the StartAdditionalInformers function in
    // pkg/controllers/sharedinformers/informers.go
    // 
    // arguments.Watch("{{.Kind}}Foo",
    //     arguments.GetSharedInformers().Factory.Bar().V1beta1().Bars().Informer(),
    //     c.FooTo{{.Kind}})
}

func (c *{{.Kind}}ControllerImpl) Get(namespace, name string) (*{{.Version}}.{{.Kind}}, error) {
    return c.lister.{{ if not .NonNamespacedKind }}{{plural .Kind}}(namespace).{{ end }}Get(name)
}
`
