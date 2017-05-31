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

package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kubernetes-incubator/apiserver-builder/cmd/apiserver-boot/boot"
	"github.com/spf13/cobra"
)

var gopath string
var wd string

func main() {
	gopath = os.Getenv("GOPATH")
	if len(gopath) == 0 {
		log.Fatal("GOPATH not defined")
	}
	boot.GoSrc = filepath.Join(gopath, "src")

	var err error
	wd, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	if !strings.HasPrefix(filepath.Dir(wd), boot.GoSrc) {
		log.Fatalf("apiserver-boot must be run from the directory containing the go package to "+
			"bootstrap. This must be under $GOPATH/src/<package>. "+
			"\nCurrent GOPATH=%s.  \nCurrent directory=%s", gopath, wd)
	}
	boot.Repo = strings.Replace(wd, boot.GoSrc+string(filepath.Separator), "", 1)
	boot.AddCreateGroup(cmd)
	boot.AddCreateResource(cmd)
	boot.AddCreateVersion(cmd)
	boot.AddDocs(cmd)
	boot.AddGenerate(cmd)
	boot.AddGlideInstallCmd(cmd)
	boot.AddInit(cmd)
	boot.AddRunCmd(cmd)
	boot.AddBuild(cmd)

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

var cmd = &cobra.Command{
	Use:   "apiserver-boot",
	Short: "apiserver-boot bootstraps building Kubernetes extensions",
	Long:  `apiserver-boot bootstraps building Kubernetes extensions`,
	Run:   RunMain,
}

func RunMain(cmd *cobra.Command, args []string) {
	cmd.Help()
}
