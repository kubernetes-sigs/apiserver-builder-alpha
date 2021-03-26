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
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/klog"
	"sigs.k8s.io/apiserver-builder-alpha/cmd/apiserver-boot/boot/util"
)

var subresourceName string
var targetSubresourceType string

type subresourceType string

var (
	subresourceTypeArbitrary subresourceType = "arbitrary"
	subresourceTypeScale     subresourceType = "scale"
	subresourceTypeConnector subresourceType = "connector"
)

var (
	supportedSubresourceTypes = []string{
		string(subresourceTypeArbitrary),
		string(subresourceTypeScale),
		string(subresourceTypeConnector),
	}
)

var createSubresourceCmd = &cobra.Command{
	Use:   "subresource",
	Short: "Creates a subresource",
	Long:  `Creates a subresource.  Creates file pkg/apis/<group>/<version>/<subresourceName>_<kind>_types.go and updates pkg/apis/<group>/<version>/<kind>_types.go with the subresource comment directive.`,
	Example: `# Create new subresource "pollinate" of resource "Bee" in the "insect" group with version "v1beta1"
apiserver-boot create subresource --subresource pollinate --group insect --version v1beta1 --kind Bee`,
	Run: RunCreateSubresource,
}

func AddCreateSubresource(cmd *cobra.Command) {
	RegisterResourceFlags(createSubresourceCmd)

	createSubresourceCmd.Flags().StringVar(&subresourceName, "subresource", "", "name of the subresource, must be singular lowercase")
	createSubresourceCmd.Flags().StringVar(&targetSubresourceType, "type", string(subresourceTypeArbitrary),
		fmt.Sprintf("type of the subresource, supported values: %v", supportedSubresourceTypes))

	cmd.AddCommand(createSubresourceCmd)
}

func RunCreateSubresource(cmd *cobra.Command, args []string) {
	ValidateSubresourceFlags()
	ValidateResourceFlags()

	if _, err := os.Stat("pkg"); err != nil {
		klog.Fatalf("could not find 'pkg' directory.  must run apiserver-boot init before creating resources")
	}

	cr := util.GetCopyright(copyright)
	createSubresource(cr)
}

func createSubresource(boilerplate string) {
	dir, err := os.Getwd()
	if err != nil {
		klog.Fatal(err)
	}
	a := subresourceTemplateArgs{
		boilerplate,
		subresourceName,
		strings.Title(kindName) + strings.Title(subresourceName),
		util.GetRepo(),
		groupName,
		versionName,
		kindName,
		resourceName,
	}

	created := false

	subResourceFileName := fmt.Sprintf("%s_%s.go", strings.ToLower(kindName), strings.ToLower(subresourceName))
	switch targetSubresourceType {
	case string(subresourceTypeArbitrary):
		path := filepath.Join(dir, "pkg", "apis", groupName, versionName, subResourceFileName)
		created = util.WriteIfNotFound(
			path,
			"subresource-arbitrary-template",
			subresourceArbitraryTemplate, a)
	case string(subresourceTypeScale):
		path := filepath.Join(dir, "pkg", "apis", groupName, versionName, subResourceFileName)
		created = util.WriteIfNotFound(
			path,
			"subresource-scale-template",
			subresourceScaleTemplate, a)
	case string(subresourceTypeConnector):
		path := filepath.Join(dir, "pkg", "apis", groupName, versionName, subResourceFileName)
		created = util.WriteIfNotFound(
			path,
			"subresource-connector-template",
			subresourceConnectorTemplate, a)
	}

	if !created {
		klog.Warningf("File %v already exists", subResourceFileName)
		os.Exit(-1)
	}

	func() {
		const (
			scaffoldSubresource = "// +kubebuilder:scaffold:subresource"
		)
		typeFile := filepath.Join("pkg", "apis", groupName, versionName, strings.ToLower(kindName)+"_types.go")
		typeFileData, err := ioutil.ReadFile(typeFile)
		if err != nil {
			klog.Fatalf("Failed reading file %v: %v", typeFile, err)
		}
		if !strings.Contains(string(typeFileData), scaffoldSubresource) {
			subresourceAppending := fmt.Sprintf(`
var _ resource.ObjectWithArbitrarySubResource = &%s{}

func (in *%s) GetArbitrarySubResources() []resource.ArbitrarySubResource {
	return []resource.ArbitrarySubResource{
		%s
	}
}
`, kindName, kindName, scaffoldSubresource)
			appendedTypeFileData := string(typeFileData) + subresourceAppending
			if err := ioutil.WriteFile(typeFile, []byte(appendedTypeFileData), 0644); err != nil {
				klog.Fatalf("Failed writing file %v: %v", typeFile, err)
			}
		}
		newRegister := fmt.Sprintf(`&%s{},`, strings.Title(kindName)+strings.Title(subresourceName))
		if err := appendMixin(typeFile, scaffoldSubresource, newRegister); err != nil {
			klog.Fatal(err)
		}
		format(typeFile)
	}()
}

func ValidateSubresourceFlags() {
	switch targetSubresourceType {
	case string(subresourceTypeArbitrary):
	case string(subresourceTypeScale):
		if subresourceName != "scale" {
			klog.Infof(`Overriding subresource name to "scale" because the type is set to "scale"`)
			subresourceName = "scale"
		}
	case string(subresourceTypeConnector):
	}
	if len(subresourceName) == 0 {
		klog.Fatalf("Must specify --subresource")
	} else {
		if strings.ToLower(subresourceName) != subresourceName {
			klog.Fatalf("Subresource name %v must be lowercased", subresourceName)
		}
	}
	if !sets.NewString(supportedSubresourceTypes...).Has(targetSubresourceType) {
		klog.Fatalf("Subresource type %v not supported", targetSubresourceType)
	} else {
		subresourceMatch := regexp.MustCompile("^[a-z]+$")
		if !subresourceMatch.MatchString(subresourceName) {
			klog.Fatalf("--subresource must match regex ^[a-z]+$ but was (%s)", subresourceName)
		}
	}
}

type subresourceTemplateArgs struct {
	BoilerPlate     string
	Subresource     string
	SubresourceKind string
	Repo            string
	Group           string
	Version         string
	Kind            string
	Resource        string
}

var subresourceScaleTemplate = `
{{.BoilerPlate}}

package {{.Version}}

import (
	v1 "k8s.io/api/autoscaling/v1"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource"
)

var _ resource.ObjectWithScaleSubResource = &{{.Kind}}{}

func (in *{{.Kind}}) SetScale(scaleSubResource *v1.Scale) {
	// EDIT IT
}

func (in *{{.Kind}}) GetScale() (scaleSubResource *v1.Scale) {
	// EDIT IT
	return &v1.Scale{}
}
`

var subresourceArbitraryTemplate = `
{{.BoilerPlate}}

package {{.Version}}

import (
	"context"
	"fmt"

	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource/resourcerest"
	contextutil "sigs.k8s.io/apiserver-runtime/pkg/util/context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/rest"
)

var _ resource.SubResource = &{{.SubresourceKind}}{}
var _ resourcerest.Getter = &{{.SubresourceKind}}{}
var _ resourcerest.Updater = &{{.SubresourceKind}}{}

// {{.Kind}}{{.SubresourceKind}}
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type {{.SubresourceKind}} struct {
	metav1.TypeMeta ` + "`" + `json:",inline" ` + "`" + `
}

func (c *{{.SubresourceKind}}) SubResourceName() string {
	return "{{.Subresource}}"
}

func (c *{{.SubresourceKind}}) New() runtime.Object {
	return &{{.SubresourceKind}}{}
}

func (c *{{.SubresourceKind}}) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
    // EDIT IT
	parentStorage, ok := contextutil.GetParentStorage(ctx)
	if !ok {
		return nil, fmt.Errorf("no parent storage found in the context")
	}
	return parentStorage.Get(ctx, name, options)
}

func (c *{{.SubresourceKind}}) Update(
	ctx context.Context,
	name string,
	objInfo rest.UpdatedObjectInfo,
	createValidation rest.ValidateObjectFunc,
	updateValidation rest.ValidateObjectUpdateFunc,
	forceAllowCreate bool,
	options *metav1.UpdateOptions) (runtime.Object, bool, error) {
    // EDIT IT
	parentStorage, ok := contextutil.GetParentStorage(ctx)
	if !ok {
		return nil, false, fmt.Errorf("no parent storage found in the context")
	}
	return parentStorage.Update(ctx, name, objInfo, createValidation, updateValidation, forceAllowCreate, options)
}
`

var subresourceConnectorTemplate = `
{{.BoilerPlate}}

package {{.Version}}

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/rest"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource/resourcerest"
	contextutil "sigs.k8s.io/apiserver-runtime/pkg/util/context"
)

var _ resource.SubResource = &{{.SubresourceKind}}{}
var _ rest.Storage = &{{.SubresourceKind}}{}
var _ resourcerest.Connecter = &{{.SubresourceKind}}{}

var {{.Subresource}}ProxyMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

// {{.SubresourceKind}}
type {{.SubresourceKind}} struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type {{.SubresourceKind}}Options struct {
	metav1.TypeMeta

	// Path is the target api path of the proxy request.
	Path string ` + "`" + `json:"path"` + "`" + `
}

func (c *{{.SubresourceKind}}) SubResourceName() string {
	return "proxy"
}

func (c *{{.SubresourceKind}}) New() runtime.Object {
	return &{{.SubresourceKind}}Options{}
}

func (c *{{.SubresourceKind}}) Connect(ctx context.Context, id string, options runtime.Object, r rest.Responder) (http.Handler, error) {
	// EDIT IT
	parentStorage, ok := contextutil.GetParentStorage(ctx)
	if !ok {
		return nil, fmt.Errorf("no parent storage found")
	}
	_, err := parentStorage.Get(ctx, id, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return http.Handler(nil), nil
}

func (c *{{.SubresourceKind}}) NewConnectOptions() (runtime.Object, bool, string) {
	return &{{.SubresourceKind}}Options{}, false, "path"
}

func (c *{{.SubresourceKind}}) ConnectMethods() []string {
	return {{.Subresource}}ProxyMethods
}

var _ resource.QueryParameterObject = &{{.SubresourceKind}}Options{}

func (in *{{.SubresourceKind}}Options) ConvertFromUrlValues(values *url.Values) error {
	in.Path = values.Get("path")
	return nil
}
`
