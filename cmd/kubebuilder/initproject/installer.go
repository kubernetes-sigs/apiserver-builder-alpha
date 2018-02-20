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

package initproject

import (
    "fmt"
	"path/filepath"

	"github.com/kubernetes-sigs/kubebuilder/cmd/kubebuilder/util"
)

func createInstaller(boilerplate string) {
	fmt.Printf("\t%s/\n", filepath.Join("cmd", "installer"))
	execute(
		filepath.Join("cmd", "installer", "main.go"),
		"installer-template",
		installerTemplate,
		installerTemplateArguments{boilerplate, util.Repo},
	)
}

type installerTemplateArguments struct {
	BoilerPlate string
	Repo        string
}

var installerTemplate = `
{{.BoilerPlate}}

package main

import (
    "flag"
    "log"

    "github.com/kubernetes-sigs/kubebuilder/pkg/install"
    controllerlib "github.com/kubernetes-sigs/kubebuilder/pkg/controller"

	"{{ .Repo }}/pkg/apis"
)

var kubeconfig = flag.String("kubeconfig", "", "path to kubeconfig")
var controllerImage = flag.String("controller-image", "", "name of container image containing the controller binary")
var docsImage = flag.String("docs-image", "", "name of container image the reference docs")
var apiserverImage = flag.String("apiserver-image", "", "name of apiserver image")
var name = flag.String("name", "", "name of the installation")
var apiserverAggregation = flag.Bool("apiserver-aggregation", false, "use apiserver aggregation instead of CRDs")
var uninstall = flag.Bool("uninstall", false, "uninstall the API")

func main() {
    flag.Parse()
    config, err := controllerlib.GetConfig(*kubeconfig)
    if err != nil {
        log.Fatalf("Could not create Config for talking to the apiserver: %v", err)
    }

    // Install the API components into the cluster
    var strategy install.InstallStrategy
        if *apiserverAggregation {
            // Do not use - doesn't work yet
            strategy = &install.ApiserverInstallStrategy{
                Name: *name,
                APIMeta: apis.APIMeta,
                ApiserverImage: *apiserverImage,
                ControllerManagerImage: *controllerImage,
                DocsImage: *docsImage,
            }
        } else {
            strategy = &install.CRDInstallStrategy{
                Name: *name,
                APIMeta: apis.APIMeta,
                ControllerManagerImage: *controllerImage,
                DocsImage: *docsImage,
            }
        }

    if !*uninstall {
        err = install.NewInstaller(config).Install(strategy)
        if err != nil {
            log.Fatalf("Failed to install API: %v", err)
        }
    } else {
        err = install.NewUninstaller(config).Uninstall(strategy)
        if err != nil {
            log.Fatalf("Failed to uninstall API: %v", err)
        }
    }
}
`
