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

package init_repo

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"k8s.io/klog"
	"sigs.k8s.io/apiserver-builder-alpha/cmd/apiserver-boot/boot/util"
	"sigs.k8s.io/kubebuilder/pkg/scaffold"
	"sigs.k8s.io/kubebuilder/pkg/scaffold/input"
	"sigs.k8s.io/kubebuilder/pkg/scaffold/manager"
)

var repoCmd = &cobra.Command{
	Use:     "repo",
	Short:   "Initialize a repo with the apiserver scaffolding and vendor/ dependencies",
	Long:    `Initialize a repo with the apiserver scaffolding and vendor/ dependencies`,
	Example: `apiserver-boot init repo --domain mydomain`,
	Run:     RunInitRepo,
}

var domain string
var copyright string

func AddInitRepo(cmd *cobra.Command) {
	cmd.AddCommand(repoCmd)
	repoCmd.Flags().StringVar(&domain, "domain", "", "domain the api groups live under")

	// Hide this flag by default
	repoCmd.Flags().StringVar(&copyright, "copyright", "boilerplate.go.txt", "Location of copyright boilerplate file.")
}

func RunInitRepo(cmd *cobra.Command, args []string) {
	if len(domain) == 0 {
		klog.Fatal("Must specify --domain")
	}
	cr := util.GetCopyright(copyright)

	createKubeBuilderProjectFile()
	createBazelWorkspace()
	createApiserver(cr)
	createControllerManager(cr)
	createAPIs(cr)

	createPackage(cr, filepath.Join("pkg"), "")
	createPackage(cr, filepath.Join("pkg", "controller"), "")
	createPackage(cr, filepath.Join("pkg", "openapi"), "//go:generate "+
		"go run ../../vendor/k8s.io/kube-openapi/cmd/openapi-gen/openapi-gen.go "+
		"-o . "+
		"--output-package ../../pkg/openapi "+
		"--report-filename violations.report "+
		"-i ../../pkg/apis/...,../../vendor/k8s.io/api/core/v1,../../vendor/k8s.io/apimachinery/pkg/apis/meta/v1 "+
		"-h ../../boilerplate.go.txt")

	os.MkdirAll("bin", 0700)

}

func createKubeBuilderProjectFile() {
	dir, err := os.Getwd()
	if err != nil {
		klog.Fatal(err)
	}
	path := filepath.Join(dir, "PROJECT")
	util.WriteIfNotFound(path, "project-template", projectFileTemplate,
		buildTemplateArguments{domain, util.Repo})
}

var projectFileTemplate = `
version: "1"
domain: {{.Domain}}
repo: {{.Repo}}
`

func createBazelWorkspace() {
	dir, err := os.Getwd()
	if err != nil {
		klog.Fatal(err)
	}
	path := filepath.Join(dir, "WORKSPACE")
	util.WriteIfNotFound(path, "bazel-workspace-template", workspaceTemplate, nil)
	path = filepath.Join(dir, "BUILD.bazel")
	util.WriteIfNotFound(path, "bazel-build-template",
		buildTemplate, buildTemplateArguments{domain, util.Repo})
}

func createControllerManager(boilerplate string) {
	err := (&scaffold.Scaffold{}).Execute(input.Options{
		BoilerplatePath: "boilerplate.go.txt",
	},
		&manager.Cmd{
			Input: input.Input{
				Boilerplate: boilerplate,
			},
		},
		&manager.Webhook{
			Input: input.Input{
				Boilerplate: boilerplate,
			},
		})
	if err != nil {
		klog.Fatal(err)
	}
}

type apiserverTemplateArguments struct {
	Domain      string
	BoilerPlate string
	Repo        string
}

var apiserverTemplate = `
{{.BoilerPlate}}

package main

import (
	// Make sure dep tools picks up these dependencies
	_ "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "github.com/go-openapi/loads"

	"sigs.k8s.io/apiserver-builder-alpha/pkg/cmd/server"
	_ "k8s.io/client-go/plugin/pkg/client/auth" // Enable cloud provider auth

	"{{.Repo}}/pkg/apis"
	"{{.Repo}}/pkg/openapi"
	_ "{{.Repo}}/plugin/admission/install"
)

func main() {
	version := "v0"

	err := server.StartApiServerWithOptions(&server.StartOptions{
		EtcdPath:         "/registry/{{ .Domain }}",
		Apis:             apis.GetAllApiBuilders(),
		Openapidefs:      openapi.GetOpenAPIDefinitions,
		Title:            "Api",
		Version:          version,

		// TweakConfigFuncs []func(apiServer *apiserver.Config) error
		// FlagConfigFuncs []func(*cobra.Command) error
	})
	if err != nil {
		panic(err)
	}
}
`

func createApiserver(boilerplate string) {
	dir, err := os.Getwd()
	if err != nil {
		klog.Fatal(err)
	}
	path := filepath.Join(dir, "cmd", "apiserver", "main.go")
	util.WriteIfNotFound(path, "apiserver-template", apiserverTemplate,
		apiserverTemplateArguments{
			domain,
			boilerplate,
			util.Repo,
		})

}

func createPackage(boilerplate, path, goGenerateCommand string) {
	pkg := filepath.Base(path)
	dir, err := os.Getwd()
	if err != nil {
		klog.Fatal(err)
	}
	path = filepath.Join(dir, path, "doc.go")
	util.WriteIfNotFound(path, "pkg-template", packageDocTemplate,
		packageDocTemplateArguments{
			boilerplate,
			pkg,
			goGenerateCommand,
		})
}

type packageDocTemplateArguments struct {
	BoilerPlate       string
	Package           string
	GoGenerateCommand string
}

var packageDocTemplate = `
{{.BoilerPlate}}

{{.GoGenerateCommand}}
package {{.Package}}

`

func createAPIs(boilerplate string) {
	dir, err := os.Getwd()
	if err != nil {
		klog.Fatal(err)
	}
	path := filepath.Join(dir, "pkg", "apis", "doc.go")
	util.WriteIfNotFound(path, "apis-template", apisDocTemplate,
		apisDocTemplateArguments{
			boilerplate,
			domain,
		})
}

type apisDocTemplateArguments struct {
	BoilerPlate string
	Domain      string
}

var apisDocTemplate = `
{{.BoilerPlate}}


//go:generate go run ../../vendor/sigs.k8s.io/apiserver-builder-alpha/cmd/apiregister-gen/main.go --input-dirs ./... -h ../../boilerplate.go.txt

//
// +domain={{.Domain}}

package apis

`

var workspaceTemplate = `
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_k8s_repo_infra",
    sha256 = "5ee2a8e306af0aaf2844b5e2c79b5f3f53fc9ce3532233f0615b8d0265902b2a",
    strip_prefix = "repo-infra-0.0.1-alpha.1",
    urls = [
        "https://github.com/kubernetes/repo-infra/archive/v0.0.1-alpha.1.tar.gz",
    ],
)

load("@io_k8s_repo_infra//:load.bzl", _repo_infra_repos = "repositories")

_repo_infra_repos()

load("@io_k8s_repo_infra//:repos.bzl", "configure")

# use k8s.io/repo-infra to configure go and bazel
# default minimum_bazel_version is 0.29.1
configure(
    go_version = "1.13",
    rbe_name = None,
)
`

type buildTemplateArguments struct {
	Domain string
	Repo   string
}

var buildTemplate = `
load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:proto disable_global
# gazelle:prefix {{.Repo}}
gazelle(
    name = "gazelle",
    command = "fix",
)
`
