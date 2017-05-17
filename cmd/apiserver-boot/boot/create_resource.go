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

package boot

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var createResourceCmd = &cobra.Command{
	Use:   "create-resource",
	Short: "Creates an API resource",
	Long:  `Creates an API resource`,
	Run:   RunCreateResource,
}

func AddCreateResource(cmd *cobra.Command) {
	createResourceCmd.Flags().StringVar(&groupName, "group", "", "name of the API group")
	createResourceCmd.Flags().StringVar(&versionName, "version", "", "name of the API version")
	createResourceCmd.Flags().StringVar(&kindName, "kind", "", "name of the API kind to create")
	createResourceCmd.Flags().StringVar(&resourceName, "resource", "", "name of the API resource to create, plural name of the kind")
	createResourceCmd.Flags().StringVar(&copyright, "copyright", "", "path to copyright file.  defaults to boilerplate.go.txt")
	createResourceCmd.Flags().StringVar(&domain, "domain", "", "domain the api group lives under")
	cmd.AddCommand(createResourceCmd)
}

func RunCreateResource(cmd *cobra.Command, args []string) {
	if len(domain) == 0 {
		fmt.Fprintf(os.Stderr, "apiserver-boot create-resource requires the --domain flag\n")
		os.Exit(-1)
	}
	if len(groupName) == 0 {
		fmt.Fprintf(os.Stderr, "apiserver-boot create-resource requires the --group flag\n")
		os.Exit(-1)
	}
	if len(versionName) == 0 {
		fmt.Fprintf(os.Stderr, "apiserver-boot create-resource requires the --version flag\n")
		os.Exit(-1)
	}
	if len(kindName) == 0 {
		fmt.Fprintf(os.Stderr, "apiserver-boot create-resource requires the --kind flag\n")
		os.Exit(-1)
	}
	if len(resourceName) == 0 {
		fmt.Fprintf(os.Stderr, "apiserver-boot create-resource requires the --resource flag\n")
		os.Exit(-1)
	}

	cr := getCopyright()

	ignoreExists = true
	createGroup(cr)
	createVersion(cr)

	ignoreExists = false
	createResource(cr)
}

func createResource(boilerplate string) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(-1)
	}
	typesFileName := fmt.Sprintf("%s_types.go", strings.ToLower(kindName))
	path := filepath.Join(dir, "pkg", "apis", groupName, versionName, typesFileName)
	a := resourceTemplateArgs{
		boilerplate,
		domain,
		groupName,
		versionName,
		kindName,
		resourceName,
		Repo,
	}

	found := false

	created := writeIfNotFound(path, "resource-template", resourceTemplate, a)
	if !created {
		fmt.Fprintf(os.Stderr,
			"API group version kind %s/%s/%s already exists.\n", groupName, versionName, kindName)
		found = true
	}

	typesFileName = fmt.Sprintf("%s_types_test.go", strings.ToLower(kindName))
	path = filepath.Join(dir, "pkg", "apis", groupName, versionName, typesFileName)
	created = writeIfNotFound(path, "resource-test-template", resourceTestTemplate, a)
	if !created {
		fmt.Fprintf(os.Stderr,
			"API group version kind %s/%s/%s test already exists.\n", groupName, versionName, kindName)
		found = true
	}

	path = filepath.Join(dir, "pkg", "controller", strings.ToLower(kindName), "controller.go")
	created = writeIfNotFound(path, "resource-controller-template", resourceControllerTemplate, a)
	if !created {
		fmt.Fprintf(os.Stderr,
			"Controller for %s/%s/%s already exists.\n", groupName, versionName, kindName)
		found = true
	}

	if found {
		os.Exit(-1)
	}
}

type resourceTemplateArgs struct {
	BoilerPlate string
	Domain      string
	Group       string
	Version     string
	Kind        string
	Resource    string
	Repo        string
}

var resourceTemplate = `
{{.BoilerPlate}}

package {{.Version}}

import (
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient=true

// {{.Kind}}
// +k8s:openapi-gen=true
// +resource:path={{.Resource}}
type {{.Kind}} struct {
	metav1.TypeMeta   ` + "`" + `json:",inline"` + "`" + `
	metav1.ObjectMeta ` + "`" + `json:"metadata,omitempty"` + "`" + `

	Spec   {{.Kind}}Spec   ` + "`" + `json:"spec,omitempty"` + "`" + `
	Status {{.Kind}}Status ` + "`" + `json:"status,omitempty"` + "`" + `
}

// {{.Kind}}Spec defines the desired state of {{.Kind}}
type {{.Kind}}Spec struct {
}

// {{.Kind}}Status defines the observed state of {{.Kind}}
type {{.Kind}}Status struct {
}

// DefaultingFunction sets default {{.Kind}} field values
func ({{.Kind}}SchemeFns) DefaultingFunction(o interface{}) {
	obj := o.(*{{.Kind}})
	// Set default field values here
	log.Printf("Defaulting fields for {{.Kind}} %s\n", obj.Name)
}

`

var resourceTestTemplate = `
{{.BoilerPlate}}

package {{.Version}}_test

import (
	"os"
	"testing"

	"k8s.io/client-go/rest"
	"github.com/kubernetes-incubator/apiserver-builder/pkg/test"

	"{{.Repo}}/pkg/apis"
	"{{.Repo}}/pkg/client/clientset_generated/clientset"
	"{{.Repo}}/pkg/openapi"
)

var testenv *test.TestEnvironment
var config *rest.Config
var client *clientset.Clientset

// Do Test Suite setup / teardown
func TestMain(m *testing.M) {
	testenv = test.NewTestEnvironment()
	config = testenv.Start(apis.GetAllApiBuilders(), openapi.GetOpenAPIDefinitions)
	client = clientset.NewForConfigOrDie(config)
	retCode := m.Run()
	testenv.Stop()
	os.Exit(retCode)
}

func TestCreateDelete{{.Kind}}(t *testing.T) {
}
`

var resourceControllerTemplate = `
{{.BoilerPlate}}

package controller

import (
	"log"

	"github.com/kubernetes-incubator/apiserver-builder/pkg/controller"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	"{{.Repo}}/pkg/apis/{{.Group}}/{{.Version}}"
	"{{.Repo}}/pkg/controller/sharedinformers"
	listers "{{.Repo}}/pkg/client/listers_generated/{{.Group}}/{{.Version}}"
)

// +controller:group={{ .Group }},version={{ .Version }},kind={{ .Kind}},resource={{ .Resource }}
type {{.Kind}}ControllerImpl struct {
	// informer listens for events about {{.Kind}}
	informer cache.SharedIndexInformer

	// lister indexes properties about {{.Kind}}
	lister listers.{{.Kind}}Lister
}

// Init initializes the controller and is called by the generated code
// Registers eventhandlers to enqueue events
// config - client configuration for talking to the apiserver
// si - informer factory shared across all controllers for listening to events and indexing resource properties
// queue - message queue for handling new events.  unique to this controller.
func (c *{{.Kind}}ControllerImpl) Init(
	config *rest.Config,
	si *sharedinformers.SharedInformers,
	queue workqueue.RateLimitingInterface) {

	// Set the informer and lister for subscribing to events and indexing {{.Resource}} labels
	i := si.Factory.{{title .Group}}().{{title .Version}}().{{title .Resource}}()
	c.informer = i.Informer()
	c.lister = i.Lister()

	// Add an event handler to enqueue a message for {{.Resource}} adds / updates
	c.informer.AddEventHandler(&controller.QueueingEventHandler{queue})
}

// Reconcile handles enqueued messages
func (c *{{.Kind}}ControllerImpl) Reconcile(u *{{.Version}}.{{.Kind}}) error {
	// Implement controller logic here
	log.Printf("Running reconcile {{.Kind}} for %s\n", u.Name)
	return nil
}
`
