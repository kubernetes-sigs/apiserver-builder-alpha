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

package create

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/kubernetes-incubator/apiserver-builder/cmd/apiserver-boot/boot/util"
	"github.com/markbates/inflect"
	"github.com/spf13/cobra"
	"os/exec"
)

var kindName string
var resourceName string
var nonNamespacedKind bool

var createResourceCmd = &cobra.Command{
	Use:   "resource",
	Short: "Creates an API group, version and resource",
	Long:  `Creates an API group, version and resource.  Will not recreate group or resource if they already exist.`,
	Example: `# Create new resource "Bee" in the "insect" group with version "v1beta"
# Will automatically the group and version if they do not exist
apiserver-boot create group version kind --group insect --version v1beta --kind Bee`,
	Run: RunCreateResource,
}

func AddCreateResource(cmd *cobra.Command) {
	createResourceCmd.Flags().StringVar(&groupName, "group", "", "name of the API group to create")
	createResourceCmd.Flags().StringVar(&versionName, "version", "", "name of the API version to create")
	createResourceCmd.Flags().StringVar(&kindName, "kind", "", "name of the API kind to create")
	createResourceCmd.Flags().StringVar(&resourceName, "resource", "", "optional name of the API resource to create, normally the plural name of the kind in lowercase")
	createResourceCmd.Flags().BoolVar(&nonNamespacedKind, "non-namespaced", false, "if set, the API kind will be non namespaced")

	cmd.AddCommand(createResourceCmd)
}

func RunCreateResource(cmd *cobra.Command, args []string) {
	if _, err := os.Stat("pkg"); err != nil {
		log.Fatalf("could not find 'pkg' directory.  must run apiserver-boot init before creating resources")
	}

	util.GetDomain()
	if len(groupName) == 0 {
		log.Fatalf("Must specify --group")
	}
	if len(versionName) == 0 {
		log.Fatalf("Must specify --version")
	}
	if len(kindName) == 0 {
		log.Fatal("Must specify --kind")
	}
	if len(resourceName) == 0 {
		resourceName = inflect.NewDefaultRuleset().Pluralize(strings.ToLower(kindName))
	}

	if strings.ToLower(groupName) != groupName {
		log.Fatalf("--group must be lowercase was (%s)", groupName)
	}
	versionMatch := regexp.MustCompile("^v\\d+(alpha\\d+|beta\\d+)*$")
	if !versionMatch.MatchString(versionName) {
		log.Fatalf(
			"--version has bad format. must match ^v\\d+(alpha\\d+|beta\\d+)*$.  "+
				"e.g. v1alpha1,v1beta1,v1 was(%s)", versionName)
	}
	if string(kindName[0]) != strings.ToUpper(string(kindName[0])) {
		log.Fatalf("--kind must start with uppercase letter was (%s)", kindName)
	}

	cr := util.GetCopyright(copyright)

	ignoreGroupExists = true
	createGroup(cr)
	ignoreVersionExists = true
	createVersion(cr)

	createResource(cr)
}

func createResource(boilerplate string) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	typesFileName := fmt.Sprintf("%s_types.go", strings.ToLower(kindName))
	path := filepath.Join(dir, "pkg", "apis", groupName, versionName, typesFileName)
	a := resourceTemplateArgs{
		boilerplate,
		util.Domain,
		groupName,
		versionName,
		kindName,
		resourceName,
		util.Repo,
		inflect.NewDefaultRuleset().Pluralize(kindName),
		nonNamespacedKind,
	}

	found := false

	created := util.WriteIfNotFound(path, "resource-template", resourceTemplate, a)
	if !created {
		if !found {
			log.Printf("API group version kind %s/%s/%s already exists.",
				groupName, versionName, kindName)
			found = true
		}
	}

	exec.Command("mkdir", "-p", filepath.Join("docs", "examples")).CombinedOutput()
	docpath := filepath.Join("docs", "examples", strings.ToLower(kindName), fmt.Sprintf("%s.yaml", strings.ToLower(kindName)))
	created = util.WriteIfNotFound(docpath, "example-template", exampleTemplate, a)
	if !created {
		if !found {
			log.Printf("Example %s already exists.", docpath)
			found = true
		}
	}

	exec.Command("mkdir", "-p", filepath.Join("sample")).CombinedOutput()
	samplepath := filepath.Join("sample", fmt.Sprintf("%s.yaml", strings.ToLower(kindName)))
	created = util.WriteIfNotFound(samplepath, "sample-template", sampleTemplate, a)
	if !created {
		if !found {
			log.Printf("Sample %s already exists.", docpath)
			found = true
		}
	}

	// write the suite if it is missing
	typesFileName = fmt.Sprintf("%s_suite_test.go", strings.ToLower(versionName))
	path = filepath.Join(dir, "pkg", "apis", groupName, versionName, typesFileName)
	util.WriteIfNotFound(path, "version-suite-test-template", resourceSuiteTestTemplate, a)

	typesFileName = fmt.Sprintf("%s_types_test.go", strings.ToLower(kindName))
	path = filepath.Join(dir, "pkg", "apis", groupName, versionName, typesFileName)
	created = util.WriteIfNotFound(path, "resource-test-template", resourceTestTemplate, a)
	if !created {
		if !found {
			log.Printf("API group version kind %s/%s/%s test already exists.",
				groupName, versionName, kindName)
			found = true
		}
	}

	path = filepath.Join(dir, "pkg", "controller", strings.ToLower(kindName), "controller.go")
	created = util.WriteIfNotFound(path, "resource-controller-template", resourceControllerTemplate, a)
	if !created {
		if !found {
			log.Printf("Controller for %s/%s/%s already exists.",
				groupName, versionName, kindName)
			found = true
		}
	}

	path = filepath.Join(dir, "pkg", "controller", strings.ToLower(kindName), fmt.Sprintf("%s_suite_test.go", strings.ToLower(kindName)))
	util.WriteIfNotFound(path, "resource-controller-suite-test-template", controllerSuiteTestTemplate, a)

	path = filepath.Join(dir, "pkg", "controller", strings.ToLower(kindName), "controller_test.go")
	created = util.WriteIfNotFound(path, "controller-test-template", controllerTestTemplate, a)
	if !created {
		if !found {
			log.Printf("Controller test for %s/%s/%s already exists.",
				groupName, versionName, kindName)
			found = true
		}
	}

	if found {
		os.Exit(-1)
	}
}

type resourceTemplateArgs struct {
	BoilerPlate       string
	Domain            string
	Group             string
	Version           string
	Kind              string
	Resource          string
	Repo              string
	PluralizedKind    string
	NonNamespacedKind bool
}

var resourceTemplate = `
{{.BoilerPlate}}

package {{.Version}}

import (
	"log"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/endpoints/request"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"{{ .Repo }}/pkg/apis/{{.Group}}"
)

// +genclient=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
{{- if .NonNamespacedKind }}
// +nonNamespaced=true{{ end }}

// {{.Kind}}
// +k8s:openapi-gen=true
// +resource:path={{.Resource}},strategy={{.Kind}}Strategy
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

// Validate checks that an instance of {{.Kind}} is well formed
func ({{.Kind}}Strategy) Validate(ctx request.Context, obj runtime.Object) field.ErrorList {
	o := obj.(*{{.Group}}.{{.Kind}})
	log.Printf("Validating fields for {{.Kind}} %s\n", o.Name)
	errors := field.ErrorList{}
	// perform validation here and add to errors using field.Invalid
	return errors
}

{{- if .NonNamespacedKind }}

func ({{.Kind}}Strategy) NamespaceScoped() bool { return false }

func ({{.Kind}}StatusStrategy) NamespaceScoped() bool { return false }
{{- end }}

// DefaultingFunction sets default {{.Kind}} field values
func ({{.Kind}}SchemeFns) DefaultingFunction(o interface{}) {
	obj := o.(*{{.Kind}})
	// set default field values here
	log.Printf("Defaulting fields for {{.Kind}} %s\n", obj.Name)
}
`

var resourceSuiteTestTemplate = `
{{.BoilerPlate}}

package {{.Version}}_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/kubernetes-incubator/apiserver-builder/pkg/test"
	"k8s.io/client-go/rest"

	"{{ .Repo }}/pkg/apis"
	"{{ .Repo }}/pkg/client/clientset_generated/clientset"
	"{{ .Repo }}/pkg/openapi"
)

var testenv *test.TestEnvironment
var config *rest.Config
var cs *clientset.Clientset

func Test{{title .Version}}(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "v1 Suite", []Reporter{test.NewlineReporter{}})
}

var _ = BeforeSuite(func() {
	testenv = test.NewTestEnvironment()
	config = testenv.Start(apis.GetAllApiBuilders(), openapi.GetOpenAPIDefinitions)
	cs = clientset.NewForConfigOrDie(config)
})

var _ = AfterSuite(func() {
	testenv.Stop()
})
`

var resourceTestTemplate = `
{{.BoilerPlate}}

package {{.Version}}_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "{{.Repo}}/pkg/apis/{{.Group}}/{{.Version}}"
	. "{{.Repo}}/pkg/client/clientset_generated/clientset/typed/{{.Group}}/{{.Version}}"
)

var _ = Describe("{{.Kind}}", func() {
	var instance {{ .Kind}}
	var expected {{ .Kind}}
	var client {{ .Kind}}Interface

	BeforeEach(func() {
		instance = {{ .Kind}}{}
		instance.Name = "instance-1"

		expected = instance
	})

	AfterEach(func() {
		client.Delete(instance.Name, &metav1.DeleteOptions{})
	})

	Describe("when sending a storage request", func() {
		Context("for a valid config", func() {
			It("should provide CRUD access to the object", func() {
				client = cs.{{ title .Group}}{{title .Version}}Client.{{plural .Kind}}({{ if not .NonNamespacedKind }}"{{lower .Kind}}-test-valid"{{ end }})

				By("returning success from the create request")
				actual, err := client.Create(&instance)
				Expect(err).ShouldNot(HaveOccurred())

				By("defaulting the expected fields")
				Expect(actual.Spec).To(Equal(expected.Spec))

				By("returning the item for list requests")
				result, err := client.List(metav1.ListOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Items).To(HaveLen(1))
				Expect(result.Items[0].Spec).To(Equal(expected.Spec))

				By("returning the item for get requests")
				actual, err = client.Get(instance.Name, metav1.GetOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(actual.Spec).To(Equal(expected.Spec))

				By("deleting the item for delete requests")
				err = client.Delete(instance.Name, &metav1.DeleteOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				result, err = client.List(metav1.ListOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Items).To(HaveLen(0))
			})
		})
	})
})
`

var resourceControllerTemplate = `
{{.BoilerPlate}}

package {{ lower .Kind }}

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
	i := si.Factory.{{title .Group}}().{{title .Version}}().{{plural .Kind}}()
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

func (c *{{.Kind}}ControllerImpl) Get(namespace, name string) (*{{.Version}}.{{.Kind}}, error) {
	return c.lister.{{ if not .NonNamespacedKind }}{{ title .Resource }}(namespace).{{ end }}Get(name)
}
`

var controllerSuiteTestTemplate = `
{{.BoilerPlate}}

package {{lower .Kind}}_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/rest"
	"github.com/kubernetes-incubator/apiserver-builder/pkg/test"

	"{{ .Repo }}/pkg/apis"
	"{{ .Repo }}/pkg/client/clientset_generated/clientset"
	"{{ .Repo }}/pkg/openapi"
	"{{ .Repo }}/pkg/controller/sharedinformers"
	"{{ .Repo }}/pkg/controller/{{lower .Kind}}"
)

var testenv *test.TestEnvironment
var config *rest.Config
var cs *clientset.Clientset
var shutdown chan struct{}
var controller *{{ lower .Kind }}.{{ .Kind }}Controller
var si *sharedinformers.SharedInformers

func Test{{.Kind}}(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "{{ .Kind }} Suite", []Reporter{test.NewlineReporter{}})
}

var _ = BeforeSuite(func() {
	testenv = test.NewTestEnvironment()
	config = testenv.Start(apis.GetAllApiBuilders(), openapi.GetOpenAPIDefinitions)
	cs = clientset.NewForConfigOrDie(config)

	shutdown = make(chan struct{})
	si = sharedinformers.NewSharedInformers(config, shutdown)
	controller = {{ lower .Kind }}.New{{ .Kind}}Controller(config, si)
	controller.Run(shutdown)
})

var _ = AfterSuite(func() {
	close(shutdown)
	testenv.Stop()
})
`

var controllerTestTemplate = `
{{.BoilerPlate}}

package {{ lower .Kind }}_test

import (
	"time"

	. "{{ .Repo }}/pkg/apis/{{ .Group }}/{{ .Version }}"
	. "{{ .Repo }}/pkg/client/clientset_generated/clientset/typed/{{ .Group }}/{{ .Version }}"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("{{ .Kind }} controller", func() {
	var instance {{ .Kind }}
	var expectedKey string
	var client {{ .Kind }}Interface
	var before chan struct{}
	var after chan struct{}

	BeforeEach(func() {
		instance = {{ .Kind }}{}
		instance.Name = "instance-1"
		expectedKey = "{{ if not .NonNamespacedKind }}{{lower .Kind }}-controller-test-handler/{{ end }}instance-1"
	})

	AfterEach(func() {
		client.Delete(instance.Name, &metav1.DeleteOptions{})
	})

	Describe("when creating a new object", func() {
		It("invoke the reconcile method", func() {
			client = cs.{{title .Group}}{{title .Version}}Client.{{ plural .Kind }}({{ if not .NonNamespacedKind }}"{{lower .Kind }}-controller-test-handler"{{ end }})
			before = make(chan struct{})
			after = make(chan struct{})

			actualKey := ""
			var actualErr error = nil

			// Setup test callbacks to be called when the message is reconciled
			controller.BeforeReconcile = func(key string) {
				defer close(before)
				actualKey = key
			}
			controller.AfterReconcile = func(key string, err error) {
				defer close(after)
				actualKey = key
				actualErr = err
			}

			// Create an instance
			_, err := client.Create(&instance)
			Expect(err).ShouldNot(HaveOccurred())

			// Verify reconcile function is called against the correct key
			select {
			case <-before:
				Expect(actualKey).To(Equal(expectedKey))
				Expect(actualErr).ShouldNot(HaveOccurred())
			case <-time.After(time.Second * 2):
				Fail("reconcile never called")
			}

			select {
			case <-after:
				Expect(actualKey).To(Equal(expectedKey))
				Expect(actualErr).ShouldNot(HaveOccurred())
			case <-time.After(time.Second * 2):
				Fail("reconcile never finished")
			}
		})
	})
})
`

var exampleTemplate = `note: {{ .Kind }} Example
sample: |
  apiVersion: {{ .Group }}.{{ .Domain }}/{{ .Version }}
  kind: {{ .Kind }}
  metadata:
    name: {{ lower .Kind }}-example
  spec:
`

var sampleTemplate = `apiVersion: {{ .Group }}.{{ .Domain }}/{{ .Version }}
kind: {{ .Kind }}
metadata:
  name: {{ lower .Kind }}-example
spec:
`
