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
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/markbates/inflect"
	"github.com/spf13/cobra"
	"k8s.io/klog"
	"sigs.k8s.io/apiserver-builder-alpha/cmd/apiserver-boot/boot/util"
	"sigs.k8s.io/kubebuilder/pkg/scaffold"
	"sigs.k8s.io/kubebuilder/pkg/scaffold/controller"
	"sigs.k8s.io/kubebuilder/pkg/scaffold/input"
	"sigs.k8s.io/kubebuilder/pkg/scaffold/manager"
	"sigs.k8s.io/kubebuilder/pkg/scaffold/resource"
)

var kindName string
var resourceName string
var shortName string
var nonNamespacedKind bool
var skipGenerateAdmissionController bool
var skipGenerateResource bool
var skipGenerateController bool

var createResourceCmd = &cobra.Command{
	Use:   "resource",
	Short: "Creates an API group, version and resource",
	Long:  `Creates an API group, version and resource.  Will not recreate group or resource if they already exist.  Creates file pkg/apis/<group>/<version>/<kind>_types.go`,
	Example: `# Create new resource "Bee" in the "insect" group with version "v1beta1"
# Will automatically the group and version if they do not exist
apiserver-boot create group version resource --group insect --version v1beta1 --kind Bee`,
	Run: RunCreateResource,
}

func AddCreateResource(cmd *cobra.Command) {
	RegisterResourceFlags(createResourceCmd)

	createResourceCmd.Flags().StringVar(&shortName, "short-name", "", "if set, add a short name for the resource. It must be all lowercase.")
	createResourceCmd.Flags().BoolVar(&nonNamespacedKind, "non-namespaced", false, "if set, the API kind will be non namespaced")

	createResourceCmd.Flags().BoolVar(&skipGenerateResource, "skip-resource", false, "if set, the resources will not be generated")
	createResourceCmd.Flags().BoolVar(&skipGenerateController, "skip-controller", false, "if set, the controller will not be generated")
	createResourceCmd.Flags().BoolVar(&skipGenerateAdmissionController, "skip-admission-controller", false, "if set, the admission controller will not be generated")

	cmd.AddCommand(createResourceCmd)
}

func RunCreateResource(cmd *cobra.Command, args []string) {
	if _, err := os.Stat("pkg"); err != nil {
		klog.Fatalf("could not find 'pkg' directory.  must run apiserver-boot init before creating resources")
	}

	util.GetDomain()
	ValidateResourceFlags()

	reader := bufio.NewReader(os.Stdin)

	if !cmd.Flag("skip-resource").Changed {
		fmt.Println("Create Resource [y/n]")
		skipGenerateResource = !Yesno(reader)
	}

	if !cmd.Flag("skip-controller").Changed {
		fmt.Println("Create Controller [y/n]")
		skipGenerateController = !Yesno(reader)
	}

	if !cmd.Flag("skip-admission-controller").Changed {
		fmt.Println("Create Admission Controller [y/n]")
		skipGenerateAdmissionController = !Yesno(reader)
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
		klog.Fatal(err)
	}

	//
	a := resourceTemplateArgs{
		boilerplate,
		util.Domain,
		groupName,
		versionName,
		kindName,
		resourceName,
		shortName,
		util.Repo,
		inflect.NewDefaultRuleset().Pluralize(kindName),
		nonNamespacedKind,
	}

	found := false

	if !skipGenerateResource {
		strategyFileName := fmt.Sprintf("%s_strategy.go", strings.ToLower(kindName))
		unversionedPath := filepath.Join(dir, "pkg", "apis", groupName, strategyFileName)
		created := util.WriteIfNotFound(unversionedPath, "unversioned-strategy-template", unversionedStrategyTemplate, a)
		if !created {
			if !found {
				klog.Infof("API group version kind %s/%s/%s already exists.",
					groupName, versionName, kindName)
				found = true
			}
		}

		typesFileName := fmt.Sprintf("%s_types.go", strings.ToLower(kindName))
		path := filepath.Join(dir, "pkg", "apis", groupName, versionName, typesFileName)
		created = util.WriteIfNotFound(path, "versioned-resource-template", versionedResourceTemplate, a)
		if !created {
			if !found {
				klog.Infof("API group version kind %s/%s/%s already exists.",
					groupName, versionName, kindName)
				found = true
			}
		}

		os.MkdirAll(filepath.Join("docs", "examples"), 0700)
		docpath := filepath.Join("docs", "examples", strings.ToLower(kindName), fmt.Sprintf("%s.yaml", strings.ToLower(kindName)))
		created = util.WriteIfNotFound(docpath, "example-template", exampleTemplate, a)
		if !created {
			if !found {
				klog.Infof("Example %s already exists.", docpath)
				found = true
			}
		}

		os.MkdirAll("sample", 0700)
		samplepath := filepath.Join("sample", fmt.Sprintf("%s.yaml", strings.ToLower(kindName)))
		created = util.WriteIfNotFound(samplepath, "sample-template", sampleTemplate, a)
		if !created {
			if !found {
				klog.Infof("Sample %s already exists.", docpath)
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
				klog.Infof("API group version kind %s/%s/%s test already exists.",
					groupName, versionName, kindName)
				found = true
			}
		}
	}

	if !skipGenerateAdmissionController {
		// write the admission-controller initializer if it is missing
		os.MkdirAll(filepath.Join("plugin", "admission"), 0700)
		admissionInitializerFileName := "initializer.go"
		path := filepath.Join(dir, "plugin", "admission", admissionInitializerFileName)
		created := util.WriteIfNotFound(path, "admission-initializer-template", admissionControllerInitializerTemplate, a)
		if !created {
			if !found {
				klog.Infof("admission initializer already exists.")
				// found = true
			}
		}

		// write the admission controller if it is missing
		os.MkdirAll(filepath.Join("plugin", "admission", strings.ToLower(kindName)), 0700)
		admissionControllerFileName := "admission.go"
		path = filepath.Join(dir, "plugin", "admission", strings.ToLower(kindName), admissionControllerFileName)
		created = util.WriteIfNotFound(path, "admission-controller-template", admissionControllerTemplate, a)
		if !created {
			if !found {
				klog.Infof("admission controller for kind %s test already exists.", kindName)
				found = true
			}
		}
	}

	if !skipGenerateController {
		// write controller-runtime scaffolding templates
		r := &resource.Resource{
			Namespaced: !nonNamespacedKind,
			Group:      groupName,
			Version:    versionName,
			Kind:       kindName,
			Resource:   resourceName,
		}

		err = (&scaffold.Scaffold{}).Execute(input.Options{
			BoilerplatePath: "boilerplate.go.txt",
		}, &Controller{
			Resource: r,
			Input: input.Input{
				IfExistsAction: input.Skip,
			},
		})
		if err != nil {
			klog.Warningf("failed generating %v controller: %v", kindName, err)
		}

		err = (&scaffold.Scaffold{}).Execute(input.Options{
			BoilerplatePath: "boilerplate.go.txt",
		},
			&manager.Controller{
				Input: input.Input{
					IfExistsAction: input.Skip,
				},
			},
			&controller.AddController{
				Resource: r,
				Input: input.Input{
					IfExistsAction: input.Skip,
				},
			},
			&SuiteTest{
				Resource: r,
				Input: input.Input{
					IfExistsAction: input.Skip,
				},
			},
			&Test{
				Resource: r,
				Input: input.Input{
					IfExistsAction: input.Skip,
				},
			},
		)
		if err != nil {
			klog.Warningf("failed generating controller basic packages: %v", err)
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
	ShortName         string
	Repo              string
	PluralizedKind    string
	NonNamespacedKind bool
}

var unversionedStrategyTemplate = `
{{.BoilerPlate}}

package {{.Group}}

import (
	"context"

	"k8s.io/klog"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// Validate checks that an instance of {{.Kind}} is well formed
func ({{.Kind}}Strategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	o := obj.(*{{.Kind}})
	klog.V(5).Infof("Validating fields for {{.Kind}} %s", o.Name)
	errors := field.ErrorList{}
	// perform validation here and add to errors using field.Invalid
	return errors
}

{{- if .NonNamespacedKind }}

func ({{.Kind}}Strategy) NamespaceScoped() bool { return false }

func ({{.Kind}}StatusStrategy) NamespaceScoped() bool { return false }
{{- end }}
`

var versionedResourceTemplate = `
{{.BoilerPlate}}

package {{.Version}}

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
{{- if .NonNamespacedKind }}
// +genclient:nonNamespaced
{{- end }}

// {{.Kind}}
// +k8s:openapi-gen=true
// +resource:path={{.Resource}},strategy={{.Kind}}Strategy{{ if .ShortName }},shortname={{.ShortName}}{{ end }}
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
`

var resourceSuiteTestTemplate = `
{{.BoilerPlate}}

package {{.Version}}_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/apiserver-builder-alpha/pkg/test"
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
	testenv = test.NewTestEnvironment(apis.GetAllApiBuilders(), openapi.GetOpenAPIDefinitions)
	config = testenv.Start()
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
	"context"

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
		client.Delete(context.TODO(), instance.Name, metav1.DeleteOptions{})
	})

	Describe("when sending a storage request", func() {
		Context("for a valid config", func() {
			It("should provide CRUD access to the object", func() {
				client = cs.{{ title .Group}}{{title .Version}}().{{plural .Kind}}({{ if not .NonNamespacedKind }}"{{lower .Kind}}-test-valid"{{ end }})

				By("returning success from the create request")
				actual, err := client.Create(context.TODO(), &instance, metav1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				By("defaulting the expected fields")
				Expect(actual.Spec).To(Equal(expected.Spec))

				By("returning the item for list requests")
				result, err := client.List(context.TODO(), metav1.ListOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Items).To(HaveLen(1))
				Expect(result.Items[0].Spec).To(Equal(expected.Spec))

				By("returning the item for get requests")
				actual, err = client.Get(context.TODO(), instance.Name, metav1.GetOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(actual.Spec).To(Equal(expected.Spec))

				By("deleting the item for delete requests")
				err = client.Delete(context.TODO(), instance.Name, metav1.DeleteOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				result, err = client.List(context.TODO(), metav1.ListOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Items).To(HaveLen(0))
			})
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

var admissionControllerTemplate = `
{{.BoilerPlate}}

package {{ lower .Kind }}admission

import (
	"context"
	aggregatedadmission "{{.Repo}}/plugin/admission"
	aggregatedinformerfactory "{{.Repo}}/pkg/client/informers_generated/externalversions"
	aggregatedclientset "{{.Repo}}/pkg/client/clientset_generated/clientset"
	genericadmissioninitializer "k8s.io/apiserver/pkg/admission/initializer"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/apiserver/pkg/admission"
)

var _ admission.Interface 											= &{{ lower .Kind }}Plugin{}
var _ admission.MutationInterface 									= &{{ lower .Kind }}Plugin{}
var _ admission.ValidationInterface 								= &{{ lower .Kind }}Plugin{}
var _ genericadmissioninitializer.WantsExternalKubeInformerFactory 	= &{{ lower .Kind }}Plugin{}
var _ genericadmissioninitializer.WantsExternalKubeClientSet 		= &{{ lower .Kind }}Plugin{}
var _ aggregatedadmission.WantsAggregatedResourceInformerFactory 	= &{{ lower .Kind }}Plugin{}
var _ aggregatedadmission.WantsAggregatedResourceClientSet 			= &{{ lower .Kind }}Plugin{}

func New{{ .Kind }}Plugin() *{{ lower .Kind }}Plugin {
	return &{{ lower .Kind }}Plugin{
		Handler: admission.NewHandler(admission.Create, admission.Update),
	}
}

type {{ lower .Kind }}Plugin struct {
	*admission.Handler
}

func (p *{{ lower .Kind }}Plugin) ValidateInitialization() error {
	return nil
}

func (p *{{ lower .Kind }}Plugin) Admit(ctx context.Context, a admission.Attributes, o admission.ObjectInterfaces) error {
	return nil
}

func (p *{{ lower .Kind }}Plugin) Validate(ctx context.Context, a admission.Attributes, o admission.ObjectInterfaces) error {
	return nil
}

func (p *{{ lower .Kind }}Plugin) SetAggregatedResourceInformerFactory(aggregatedinformerfactory.SharedInformerFactory) {}

func (p *{{ lower .Kind }}Plugin) SetAggregatedResourceClientSet(aggregatedclientset.Interface) {}

func (p *{{ lower .Kind }}Plugin) SetExternalKubeInformerFactory(informers.SharedInformerFactory) {}

func (p *{{ lower .Kind }}Plugin) SetExternalKubeClientSet(kubernetes.Interface) {}
`

var admissionControllerInitializerTemplate = `
{{.BoilerPlate}}

package admission

import (
	aggregatedclientset "{{.Repo}}/pkg/client/clientset_generated/clientset"
	aggregatedinformerfactory "{{.Repo}}/pkg/client/informers_generated/externalversions"
	"k8s.io/apiserver/pkg/admission"
)

// WantsAggregatedResourceClientSet defines a function which sets external ClientSet for admission plugins that need it
type WantsAggregatedResourceClientSet interface {
	SetAggregatedResourceClientSet(aggregatedclientset.Interface)
	admission.InitializationValidator
}

// WantsAggregatedResourceInformerFactory defines a function which sets InformerFactory for admission plugins that need it
type WantsAggregatedResourceInformerFactory interface {
	SetAggregatedResourceInformerFactory(aggregatedinformerfactory.SharedInformerFactory)
	admission.InitializationValidator
}

// New creates an instance of admission plugins initializer.
func New(
	clientset aggregatedclientset.Interface,
	informers aggregatedinformerfactory.SharedInformerFactory,
) pluginInitializer {
	return pluginInitializer{
		aggregatedResourceClient:    clientset,
		aggregatedResourceInformers: informers,
	}
}

type pluginInitializer struct {
	aggregatedResourceClient    aggregatedclientset.Interface
	aggregatedResourceInformers aggregatedinformerfactory.SharedInformerFactory
}

func (i pluginInitializer) Initialize(plugin admission.Interface) {
	if wants, ok := plugin.(WantsAggregatedResourceClientSet); ok {
		wants.SetAggregatedResourceClientSet(i.aggregatedResourceClient)
	}
	if wants, ok := plugin.(WantsAggregatedResourceInformerFactory); ok {
		wants.SetAggregatedResourceInformerFactory(i.aggregatedResourceInformers)
	}
}

`
