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
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a directory",
	Long:  `Initialize a directory`,
	Run:   RunInit,
}

var installDeps bool

func AddInit(cmd *cobra.Command) {
	cmd.AddCommand(initCmd)
	initCmd.Flags().StringVar(&domain, "domain", "", "domain the api groups live under")
	initCmd.Flags().StringVar(&copyright, "copyright", "", "path to copyright file.  defaults to boilerplate.go.txt")
	initCmd.Flags().BoolVar(&installDeps, "install-deps", true, "if true, install the vendored deps")
}

func RunInit(cmd *cobra.Command, args []string) {
	if len(domain) == 0 {
		log.Fatal("apiserver-boot init requires the --domain flag")
	}
	cr := getCopyright()

	createApiserver(cr)
	createControllerManager(cr)
	createAPIs(cr)
	createDocs()

	createPackage(cr, filepath.Join("pkg"))
	createPackage(cr, filepath.Join("pkg", "controller"))
	createPackage(cr, filepath.Join("pkg", "controller", "sharedinformers"))
	createPackage(cr, filepath.Join("pkg", "openapi"))

	exec.Command("mkdir", "-p", filepath.Join("bin")).CombinedOutput()

	if installDeps {
		log.Printf("installing godeps.  To disable this, run with --install-deps=false.")
		copyGlide()
	}
}

type controllerManagerTemplateArguments struct {
	BoilerPlate string
	Repo        string
}

var controllerManagerTemplate = `
{{.BoilerPlate}}

package main

import (
	"flag"
	"log"

	controllerlib "github.com/kubernetes-incubator/apiserver-builder/pkg/controller"

	"{{ .Repo }}/pkg/controller"
)

var kubeconfig = flag.String("kubeconfig", "", "path to kubeconfig")

func main() {
	flag.Parse()
	config, err := controllerlib.GetConfig(*kubeconfig)
	if err != nil {
		log.Fatalf("Could not create Config for talking to the apiserver: %v", err)
	}

	controllers, _ := controller.GetAllControllers(config)
	controllerlib.StartControllerManager(controllers...)

	// Blockforever
	select {}
}
`

func createControllerManager(boilerplate string) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	path := filepath.Join(dir, "cmd", "controller", "main.go")
	writeIfNotFound(path, "main-template", controllerManagerTemplate, controllerManagerTemplateArguments{boilerplate, Repo})

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
	// Make sure glide gets these dependencies
	_ "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "github.com/go-openapi/loads"

	"github.com/kubernetes-incubator/apiserver-builder/pkg/cmd/server"
	_ "k8s.io/client-go/plugin/pkg/client/auth" // Enable cloud provider auth

	"{{.Repo}}/pkg/apis"
	"{{.Repo}}/pkg/openapi"
)

func main() {
	server.StartApiServer("/registry/{{ .Domain }}", apis.GetAllApiBuilders(), openapi.GetOpenAPIDefinitions)
}
`

func createApiserver(boilerplate string) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	path := filepath.Join(dir, "cmd", "apiserver", "main.go")
	writeIfNotFound(path, "apiserver-template", apiserverTemplate, apiserverTemplateArguments{domain, boilerplate, Repo})

}

func createPackage(boilerplate, path string) {
	pkg := filepath.Base(path)
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	path = filepath.Join(dir, path, "doc.go")
	writeIfNotFound(path, "pkg-template", packageDocTemplate, packageDocTemplateArguments{boilerplate, pkg})
}

type packageDocTemplateArguments struct {
	BoilerPlate string
	Package     string
}

var packageDocTemplate = `
{{.BoilerPlate}}


package {{.Package}}

`

func createAPIs(boilerplate string) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	path := filepath.Join(dir, "pkg", "apis", "doc.go")
	writeIfNotFound(path, "apis-template", apisDocTemplate, apisDocTemplateArguments{boilerplate, domain})
}

type apisDocTemplateArguments struct {
	BoilerPlate string
	Domain      string
}

var apisDocTemplate = `
{{.BoilerPlate}}


//
// +domain={{.Domain}}

package apis

`

func createDocs() {
	exec.Command("mkdir", "-p", filepath.Join("docs", "openapi-spec")).CombinedOutput()
	exec.Command("mkdir", "-p", filepath.Join("docs", "static_includes")).CombinedOutput()
	exec.Command("mkdir", "-p", filepath.Join("docs", "examples")).CombinedOutput()
}
