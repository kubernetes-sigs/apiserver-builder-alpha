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
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/klog"
	"sigs.k8s.io/apiserver-builder-alpha/cmd/apiserver-boot/boot/util"
)

var versionName string
var ignoreVersionExists = false

var createVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Creates an API group and version",
	Long:  `Creates an API group and version.  Will not recreate group if already exists.`,
	Run:   RunCreateVersion,
}

func AddCreateVersion(cmd *cobra.Command) {
	createVersionCmd.Flags().StringVar(&groupName, "group", "", "name of the API group to create")
	createVersionCmd.Flags().StringVar(&versionName, "version", "", "name of the API version to create")

	cmd.AddCommand(createVersionCmd)
	createVersionCmd.AddCommand(createResourceCmd)
}

func RunCreateVersion(cmd *cobra.Command, args []string) {
	if _, err := os.Stat("pkg"); err != nil {
		klog.Fatalf("could not find 'pkg' directory.  must run apiserver-boot init before creating resources")
	}

	util.GetDomain()
	if len(groupName) == 0 {
		klog.Fatalf("Must specify --group")
	}
	if len(versionName) == 0 {
		klog.Fatalf("Must specify --version")
	}

	if strings.ToLower(groupName) != groupName {
		klog.Fatalf("--group must be lowercase was (%s)", groupName)
	}
	versionMatch := regexp.MustCompile("^v\\d+(alpha\\d+|beta\\d+)*$")
	if !versionMatch.MatchString(versionName) {
		klog.Fatalf(
			"--version has bad format. must match ^v\\d+(alpha\\d+|beta\\d+)*$.  "+
				"e.g. v1alpha1,v1beta1,v1 was(%s)", versionName)
	}

	cr := util.GetCopyright(copyright)

	ignoreGroupExists = true
	createGroup(cr)
	createVersion(cr)
}

func createVersion(boilerplate string) {
	dir, err := os.Getwd()
	if err != nil {
		klog.Fatalf("%v", err)
		os.Exit(-1)
	}
	path := filepath.Join(dir, "pkg", "apis", groupName, versionName, "doc.go")
	created := util.WriteIfNotFound(path, "version-template", versionTemplate, versionTemplateArgs{
		boilerplate,
		util.Domain,
		groupName,
		versionName,
		util.GetRepo(),
	})

	path = filepath.Join(dir, "pkg", "apis", groupName, versionName, "register.go")
	created = util.WriteIfNotFound(path, "register-template", registerTemplate, registerTemplateArgs{
		boilerplate,
		util.Domain,
		groupName,
		versionName,
	})

	if !created && !ignoreVersionExists {
		klog.Fatalf("API group version %s/%s already exists.", groupName, versionName)
	}
}

type versionTemplateArgs struct {
	BoilerPlate string
	Domain      string
	Group       string
	Version     string
	Repo        string
}

var versionTemplate = `
{{.BoilerPlate}}

// Api versions allow the api contract for a resource to be changed while keeping
// backward compatibility by support multiple concurrent versions
// of the same resource

// +k8s:openapi-gen=true
// +k8s:deepcopy-gen=package,register
// +k8s:conversion-gen={{.Repo}}/pkg/apis/{{.Group}}
// +k8s:defaulter-gen=TypeMeta
// +groupName={{.Group}}.{{.Domain}}
package {{.Version}} // import "{{.Repo}}/pkg/apis/{{.Group}}/{{.Version}}"

`

type registerTemplateArgs struct {
	BoilerPlate string
	Domain      string
	Group       string
	Version     string
}

var registerTemplate = `
{{.BoilerPlate}}

package {{.Version}}

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)


var AddToScheme = func(scheme *runtime.Scheme) error {
	metav1.AddToGroupVersion(scheme, schema.GroupVersion{
		Group:   "{{.Group}}.{{.Domain}}",
		Version: "{{.Version}}",
	})
	// +kubebuilder:scaffold:install
	return nil
}
`
