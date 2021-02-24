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

var _ resource.Object = &{{.Kind}}{}
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
	// TODO(user): Modify it, adding your API validation here.
	return nil
}

var _ resource.ObjectList = &{{.Kind}}List{}

func (in *{{.Kind}}List) GetListMeta() *metav1.ListMeta {
	return &in.ListMeta
}

{{- if .WithStatusSubResource }}
// {{.Kind}}Status defines the observed state of {{.Kind}}
type {{.Kind}}Status struct {
}

func (in {{.Kind}}Status) SubResourceName() string {
	return "status"
}

// {{.Kind}} implements ObjectWithStatusSubResource interface.
var _ resource.ObjectWithStatusSubResource = &{{.Kind}}{}

func (in *{{.Kind}}) GetStatus() resource.StatusSubResource {
	return in.Status
}

// {{.Kind}}Status{} implements StatusSubResource interface.
var _ resource.StatusSubResource = &{{.Kind}}Status{}

func (in {{.Kind}}Status) CopyTo(parent resource.ObjectWithStatusSubResource) {
	parent.(*{{.Kind}}).Status = in
}
{{- end }}
`
