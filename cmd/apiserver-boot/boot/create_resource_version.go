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
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var createVersionCmd = &cobra.Command{
	Use:   "create-version",
	Short: "Creates an API version",
	Long:  `Creates an API version.  Automatically run with create-resource.`,
	Run:   RunCreateVersion,
}

func AddCreateVersion(cmd *cobra.Command) {
	createVersionCmd.Flags().StringVar(&groupName, "group", "", "name of the API group")
	createVersionCmd.Flags().StringVar(&versionName, "version", "", "name of the API version to create")
	createVersionCmd.Flags().StringVar(&copyright, "copyright", "", "path to copyright file. defaults to boilerplate.go.txt")
	createVersionCmd.Flags().StringVar(&domain, "domain", "", "domain the api group lives under")
	cmd.AddCommand(createVersionCmd)
}

func RunCreateVersion(cmd *cobra.Command, args []string) {
	if _, err := os.Stat("pkg"); err != nil {
		log.Fatalf("could not find 'pkg' directory.  must run apiserver-boot init before creating resources")
	}

	if len(domain) == 0 {
		log.Fatalf("apiserver-boot create-version requires the --domain flag")
	}
	if len(groupName) == 0 {
		log.Fatalf("apiserver-boot create-version requires the --group flag")
	}
	if len(versionName) == 0 {
		log.Fatalf("apiserver-boot create-version requires the --version flag")
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

	cr := getCopyright()

	ignoreExists = true
	createGroup(cr)

	ignoreExists = false
	createVersion(cr)
}

func createVersion(boilerplate string) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("%v\n", err)
		os.Exit(-1)
	}
	path := filepath.Join(dir, "pkg", "apis", groupName, versionName, "doc.go")
	created := writeIfNotFound(path, "version-template", versionTemplate, versionTemplateArgs{
		boilerplate,
		domain,
		groupName,
		versionName,
		Repo,
	})
	if !created && !ignoreExists {
		log.Fatalf("API group version %s/%s already exists.", groupName, versionName)
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
