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
	"sigs.k8s.io/kubebuilder/pkg/model/config"
	"sigs.k8s.io/kubebuilder/pkg/model/resource"
	"sigs.k8s.io/kubebuilder/pkg/plugin/v3/scaffolds"
)

var kindName string
var resourceName string
var shortName string
var nonNamespacedKind bool
var skipGenerateAdmissionController bool
var skipGenerateResource bool
var skipGenerateController bool
var withStatusSubresource bool

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
	createResourceCmd.Flags().MarkDeprecated("skip-admission-controller", "")
	createResourceCmd.Flags().BoolVar(&withStatusSubresource, "with-status-subresource", true, "if set, the status sub-resource will be generated")

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

	// TODO: admission controller scaffolding
	//if !cmd.Flag("skip-admission-controller").Changed {
	//	fmt.Println("Create Admission Controller [y/n]")
	//	skipGenerateAdmissionController = !Yesno(reader)
	//}

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
		util.GetRepo(),
		inflect.NewDefaultRuleset().Pluralize(kindName),
		nonNamespacedKind,
		withStatusSubresource,
	}

	found := false

	if !skipGenerateResource {

		func() {
			// creates resource source file
			typesFileName := fmt.Sprintf("%s_types.go", strings.ToLower(kindName))
			path := filepath.Join(dir, "pkg", "apis", groupName, versionName, typesFileName)
			created := util.WriteIfNotFound(path, "versioned-resource-template", versionedResourceTemplate, a)
			if !created {
				if !found {
					klog.Infof("API group version kind %s/%s/%s already exists.",
						groupName, versionName, kindName)
					found = true
				}
			}
		}()

		func() {
			// re-render cmd/apiserver/main.go
			const (
				scaffoldImports  = "// +kubebuilder:scaffold:resource-imports"
				scaffoldRegister = "// +kubebuilder:scaffold:resource-register"
			)
			mainFile := filepath.Join("cmd", "apiserver", "main.go")
			newImport := fmt.Sprintf(`%s%s "%s/pkg/apis/%s/%s"`, groupName, versionName, util.GetRepo(), groupName, versionName)
			if err := appendMixin(mainFile, scaffoldImports, newImport); err != nil {
				klog.Fatal(err)
			}

			newRegister := fmt.Sprintf("WithResource(&%s%s.%s{}).", groupName, versionName, kindName)
			if err := appendMixin(mainFile, scaffoldRegister, newRegister); err != nil {
				klog.Fatal(err)
			}
			format(mainFile)
		}()

		func() {
			// re-render register.go
			const (
				scaffoldInstall = "// +kubebuilder:scaffold:install"
			)
			registerFile := filepath.Join("pkg", "apis", groupName, versionName, "register.go")
			fullGroupName := groupName + "." + util.Domain
			newRegister := fmt.Sprintf(`
	scheme.AddKnownTypes(schema.GroupVersion{
		Group:   "%s",
		Version: "%s",
	}, &%s{}, &%sList{})`,
				fullGroupName, versionName, kindName, kindName)
			if err := appendMixin(registerFile, scaffoldInstall, newRegister); err != nil {
				klog.Fatal(err)
			}
			format(registerFile)
		}()
	}

	if false && !skipGenerateAdmissionController {
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
			Namespaced:  !nonNamespacedKind,
			Group:       groupName,
			Version:     versionName,
			Kind:        kindName,
			Plural:      resourceName,
			Package:     filepath.Join(util.GetRepo(), "pkg", "apis", groupName, versionName),
			ImportAlias: resourceName,
		}
		scaffolder := scaffolds.NewAPIScaffolder(
			&config.Config{
				MultiGroup: true,
				Domain:     util.Domain,
				Repo:       util.GetRepo(),
				Version:    config.Version3Alpha,
			},
			boilerplate, // TODO
			r,
			false,
			true,
			nil,
		)
		err := scaffolder.Scaffold()
		if err != nil {
			klog.Warningf("failed generating controller basic packages: %v", err)
		}
		os.Remove(filepath.Join("controllers", groupName, "suite_test.go"))
	}

	if found {
		os.Exit(-1)
	}
}

type resourceTemplateArgs struct {
	BoilerPlate           string
	Domain                string
	Group                 string
	Version               string
	Kind                  string
	Resource              string
	ShortName             string
	Repo                  string
	PluralizedKind        string
	NonNamespacedKind     bool
	WithStatusSubResource bool
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
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
 	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource/resourcestrategy"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
{{- if .NonNamespacedKind }}
// +genclient:nonNamespaced
{{- end }}

// {{.Kind}}
// +k8s:openapi-gen=true
type {{.Kind}} struct {
	metav1.TypeMeta   ` + "`" + `json:",inline"` + "`" + `
	metav1.ObjectMeta ` + "`" + `json:"metadata,omitempty"` + "`" + `

	Spec   {{.Kind}}Spec   ` + "`" + `json:"spec,omitempty"` + "`" + `
{{- if .WithStatusSubResource }}
	Status {{.Kind}}Status ` + "`" + `json:"status,omitempty"` + "`" + `
{{- end }}
}

// {{.Kind}}List
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type {{.Kind}}List struct {
	metav1.TypeMeta   ` + "`" + `json:",inline"` + "`" + `
	metav1.ListMeta ` + "`" + `json:"metadata,omitempty"` + "`" + `

	Items []{{.Kind}} ` + "`" + `json:"items"` + "`" + `
}

// {{.Kind}}Spec defines the desired state of {{.Kind}}
type {{.Kind}}Spec struct {
}

{{- if .WithStatusSubResource }}
// {{.Kind}}Status defines the observed state of {{.Kind}}
type {{.Kind}}Status struct {
}
{{- end }}

var _ resource.Object = &{{.Kind}}{}
var _ resource.ObjectList = &{{.Kind}}List{}
var _ resourcestrategy.Validater = &{{.Kind}}{}


func (in *{{.Kind}}) GetObjectMeta() *metav1.ObjectMeta {
	return &in.ObjectMeta
}

func (in *{{.Kind}}) NamespaceScoped() bool {
	return false
}

func (in *{{.Kind}}) New() runtime.Object {
	return &{{.Kind}}{}
}

func (in *{{.Kind}}) NewList() runtime.Object {
	return &{{.Kind}}List{}
}

func (in *{{.Kind}}) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "{{.Group}}.{{.Domain}}",
		Version:  "{{.Version}}",
		Resource: "{{.Resource}}",
	}
}

func (in *{{.Kind}}) IsStorageVersion() bool {
	return true
}

func (in *{{.Kind}}) Validate(ctx context.Context) field.ErrorList {
	return nil
}

func (in *{{.Kind}}List) GetListMeta() *metav1.ListMeta {
	return &in.ListMeta
}

{{- if .WithStatusSubResource }}
var _ resource.ObjectWithStatusSubResource = &{{.Kind}}{}

func (in *{{.Kind}}) GetStatus() resource.StatusSubResource {
	return in.Status
}

var _ resource.StatusSubResource = &{{.Kind}}Status{}

func (in {{.Kind}}Status) CopyTo(parent resource.ObjectWithStatusSubResource) {
	parent.(*{{.Kind}}).Status = in
}
{{- end }}
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
