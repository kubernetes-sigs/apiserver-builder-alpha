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
	"strings"

	"github.com/spf13/cobra"
)

var glideInstallCmd = &cobra.Command{
	Use:   "glide-install",
	Short: "Runs glide install and flatten vendored directories",
	Long:  `Runs glide install and flatten vendored directories`,
	Run:   RunGlideInstall,
}

var fetch bool

func AddGlideInstallCmd(cmd *cobra.Command) {
	glideInstallCmd.Flags().BoolVar(&fetch, "fetch", false, "if true, fetch new glide deps instead of copying the ones packaged with the tools")
	cmd.AddCommand(glideInstallCmd)
}

func fetchGlide() {
	o, err := exec.Command("glide", "-v").CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "must install glide v0.12 or later\n")
		os.Exit(-1)
	}
	if !strings.HasPrefix(string(o), "glide version v0.12") &&
		!strings.HasPrefix(string(o), "glide version v0.13") {
		fmt.Fprintf(os.Stderr, "must install glide  or later, was %s\n", o)
		os.Exit(-1)
	}

	c := exec.Command("glide", "install", "--strip-vendor")
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	err = c.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to run glide install\n%v\n", err)
		os.Exit(-1)
	}
}

func copyGlide() {
	// copy the files
	e, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to get directory of apiserver-builder tools")
	}
	e = filepath.Dir(filepath.Dir(e))

	doCmd := func(cmd string, args ...string) {
		c := exec.Command(cmd, args...)
		c.Stderr = os.Stderr
		c.Stdout = os.Stdout
		err = c.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to copy go dependencies\n%v\n", err)
			os.Exit(-1)
		}
	}

	doCmd("cp", "-r", filepath.Join(e, "src", "vendor"), "vendor")
	doCmd("cp", filepath.Join(e, "src", "glide.yaml"), "glide.yaml")
	doCmd("cp", filepath.Join(e, "src", "glide.lock"), "glide.lock")
}

func RunGlideInstall(cmd *cobra.Command, args []string) {
	createGlide()
	if fetch {
		fetchGlide()
	} else {
		copyGlide()
	}
}

type glideTemplateArguments struct {
	Repo string
}

var glideTemplate = `
package: {{.Repo}}
import:
- package: k8s.io/apimachinery
  version: 565bae4589e797e6474096f31f9e70a47132d5e5
- package: k8s.io/apiserver
  version: 8f71532ed814093d43a03eb58f55503f77816992
- package: k8s.io/client-go
  version: 1b8c2a3e22db89d0749437fb75717be7845a5880
- package: github.com/go-openapi/analysis
  version: b44dc874b601d9e4e2f6e19140e794ba24bead3b
- package: github.com/go-openapi/jsonpointer
  version: 46af16f9f7b149af66e5d1bd010e3574dc06de98
- package: github.com/go-openapi/jsonreference
  version: 13c6e3589ad90f49bd3e3bbe2c2cb3d7a4142272
- package: github.com/go-openapi/loads
  version: 18441dfa706d924a39a030ee2c3b1d8d81917b38
- package: github.com/go-openapi/spec
  version: 6aced65f8501fe1217321abf0749d354824ba2ff
- package: github.com/go-openapi/swag
  version: 1d0bd113de87027671077d3c71eb3ac5d7dbba72
- package: github.com/golang/glog
  version: 44145f04b68cf362d9c4df2182967c2275eaefed
- package: github.com/pkg/errors
  version: a22138067af1c4942683050411a841ade67fe1eb
- package: github.com/spf13/cobra
  version: 7b1b6e8dc027253d45fc029bc269d1c019f83a34
- package: github.com/spf13/pflag
  version: d90f37a48761fe767528f31db1955e4f795d652f
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
