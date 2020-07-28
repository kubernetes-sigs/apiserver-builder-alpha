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
	"github.com/spf13/cobra"
	"k8s.io/klog"
	"os"
	"path/filepath"
	"sigs.k8s.io/apiserver-builder-alpha/cmd/apiserver-boot/boot/util"
	"sigs.k8s.io/kubebuilder/pkg/scaffold"
	"sigs.k8s.io/kubebuilder/pkg/scaffold/input"
	"sigs.k8s.io/kubebuilder/pkg/scaffold/manager"
)

var repoCmd = &cobra.Command{
	Use:     "repo",
	Short:   "Initialize a repo with the apiserver scaffolding",
	Long:    `Initialize a repo with the apiserver scaffolding`,
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
		"openapi-gen"+
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
	path = filepath.Join(dir, "go.mod")
	util.WriteIfNotFound(path, "go-mod-template",
		goModTemplate, struct {
			Repo string
		}{util.Repo})
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


//go:generate apiregister-gen --input-dirs ./... -h ../../boilerplate.go.txt

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

var goModTemplate = `
module {{.Repo}}

go 1.13

require (
	github.com/go-logr/zapr v0.1.1 // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/protobuf v1.3.4 // indirect
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/gorilla/websocket v1.4.1 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.1-0.20190723091251-e0797f438f94 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/kubernetes-incubator/reference-docs v0.0.0 // indirect
	github.com/markbates/inflect v0.0.0-00010101000000-000000000000
	github.com/onsi/ginkgo v1.11.0
	github.com/onsi/gomega v1.8.1
	github.com/pkg/errors v0.8.1
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	github.com/tmc/grpc-websocket-proxy v0.0.0-20190109142713-0ad062ec5ee5 // indirect
	golang.org/x/xerrors v0.0.0-20191011141410-1b5146add898 // indirect
	k8s.io/api v0.18.4
	k8s.io/apimachinery v0.18.4
	k8s.io/apiserver v0.18.4
	k8s.io/client-go v0.18.4
	k8s.io/gengo v0.0.0-20190822140433-26a664648505
	k8s.io/klog v1.0.0
	k8s.io/kube-aggregator v0.18.4
    k8s.io/code-generator v0.18.4
	k8s.io/kube-openapi v0.0.0-20191107075043-30be4d16710a
	k8s.io/utils v0.0.0-20191114184206-e782cd3c129f
	sigs.k8s.io/controller-runtime v0.6.0
	sigs.k8s.io/controller-tools v0.1.12 // indirect
	sigs.k8s.io/kubebuilder v1.0.8
	sigs.k8s.io/testing_frameworks v0.1.1
	sigs.k8s.io/apiserver-builder-alpha v1.18.0
)

replace sigs.k8s.io/controller-tools => sigs.k8s.io/controller-tools v0.1.12

replace sigs.k8s.io/kubebuilder => sigs.k8s.io/kubebuilder v1.0.8

replace github.com/markbates/inflect => github.com/markbates/inflect v1.0.4

replace github.com/kubernetes-incubator/reference-docs => github.com/kubernetes-sigs/reference-docs v0.0.0-20170929004150-fcf65347b256

replace sigs.k8s.io/structured-merge-diff => sigs.k8s.io/structured-merge-diff v1.0.1-0.20191108220359-b1b620dd3f06
`
