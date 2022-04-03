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

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
	"sigs.k8s.io/apiserver-builder-alpha/v2/pkg/boot/util"
	config "sigs.k8s.io/kubebuilder/v3/pkg/config/v3"
	"sigs.k8s.io/kubebuilder/v3/pkg/machinery"
	"sigs.k8s.io/kubebuilder/v3/pkg/plugins/golang/v3/scaffolds"
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
var moduleName string

func AddInitRepo(cmd *cobra.Command) {
	cmd.AddCommand(repoCmd)
	repoCmd.Flags().StringVar(&domain, "domain", "", "domain the api groups live under")

	// Hide this flag by default
	repoCmd.Flags().StringVar(&copyright, "copyright", filepath.Join("hack", "boilerplate.go.txt"), "Location of copyright boilerplate file.")
	repoCmd.Flags().StringVar(&moduleName, "module-name", "",
		"the module name of the go mod project, required if the project uses go module outside GOPATH")
}

func RunInitRepo(cmd *cobra.Command, args []string) {
	if len(domain) == 0 {
		klog.Fatal("Must specify --domain")
	}

	if len(moduleName) == 0 {
		if err := util.LoadRepoFromGoPath(); err != nil {
			klog.Fatal(err)
		}
	} else {
		util.SetRepo(moduleName)
	}
	createControllerManager()
	os.RemoveAll(filepath.Join("config")) // removes kubebuilder config scaffolding

	cr := util.GetCopyright(copyright)
	createGoMod()
	createApiserver(cr)
	createAPIs(cr)

	//createPackage(cr, filepath.Join("pkg"), "")
	//createPackage(cr, filepath.Join("pkg", "controller"), "")
	//createPackage(cr, filepath.Join("pkg", "openapi"), "//go:generate "+
	//	"openapi-gen"+
	//	"-o . "+
	//	"--output-package ../../pkg/openapi "+
	//	"--report-filename violations.report "+
	//	"-i ../../pkg/apis/...,../../vendor/k8s.io/api/core/v1,../../vendor/k8s.io/apimachinery/pkg/apis/meta/v1 "+
	//	"-h ../../boilerplate.go.txt")

	os.MkdirAll("bin", 0700)

}

func createControllerManager() {
	cfg := config.New()
	cfg.SetMultiGroup()
	cfg.SetRepository(util.GetRepo())
	cfg.SetDomain(util.Domain)
	scaffolder := scaffolds.NewInitScaffolder(
		cfg,
		"",
		"",
	)

	scaffolder.InjectFS(machinery.Filesystem{FS: afero.NewOsFs()})
	if err := scaffolder.Scaffold(); err != nil {
		klog.Fatal(err)
	}

	os.MkdirAll(filepath.Join("cmd", "manager"), 0700)
	os.Symlink(filepath.Join("..", "..", "main.go"), filepath.Join("cmd", "manager", "main.go"))
}

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
			util.GetRepo(),
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

func createGoMod() {
	dir, err := os.Getwd()
	if err != nil {
		klog.Fatal(err)
	}
	path := filepath.Join(dir, "go.mod")
	util.Overwrite(path, "gomod-template", goModTemplate,
		goModTemplateArguments{
			util.GetRepo(),
		})
}

type apiserverTemplateArguments struct {
	Domain      string
	BoilerPlate string
	Repo        string
}

type apisDocTemplateArguments struct {
	BoilerPlate string
	Domain      string
}

type goModTemplateArguments struct {
	Repo string
}
