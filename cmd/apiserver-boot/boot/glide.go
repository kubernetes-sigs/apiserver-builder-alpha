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
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var glideInstallCmd = &cobra.Command{
	Use:   "glide-install",
	Short: "Runs glide install and flatten vendored directories",
	Long:  `Runs glide install and flatten vendored directories`,
	Run:   RunGlideInstall,
}

func AddGlideInstall(cmd *cobra.Command) {
	cmd.AddCommand(glideInstallCmd)
}

func RunGlideInstall(cmd *cobra.Command, args []string) {
	createGlide()

	c := exec.Command("glide", "install", "--strip-vendor")
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	err := c.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to run glide install\n%v\n", err)
		os.Exit(-1)
	}
}

type glideTemplateArguments struct {
	Repo string
}

var glideTemplate = `
package: {{.Repo}}
import:
- package: github.com/go-openapi/spec
- package: github.com/go-openapi/loads
- package: github.com/golang/glog
- package: github.com/pkg/errors
- package: github.com/spf13/cobra
- package: github.com/spf13/pflag
  version: d90f37a48761fe767528f31db1955e4f795d652f
- package: k8s.io/apimachinery
- package: k8s.io/apiserver
- package: k8s.io/client-go
- package: k8s.io/gengo
- package: k8s.io/kubernetes
  subpackages:
  - pkg/api
- package: k8s.io/apimachinery
  subpackages:
  - pkg/apis/meta/v1
  - pkg/apis/meta
ignore:
- {{.Repo}}
`

func createGlide() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(-1)
	}
	path := filepath.Join(dir, "glide.yaml")
	writeIfNotFound(path, "glide-template", glideTemplate, glideTemplateArguments{Repo})
}
