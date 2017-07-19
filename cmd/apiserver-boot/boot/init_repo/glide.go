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
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/kubernetes-incubator/apiserver-builder/cmd/apiserver-boot/boot/util"
	"github.com/spf13/cobra"
)

var glideInstallCmd = &cobra.Command{
	Use:   "glide",
	Short: "Install glide.yaml, glide.lock and vendor/.",
	Long:  `Install glide.yaml, glide.lock and vendor/.`,
	Example: `# Bootstrap vendor/ from the src packaged with apiserver-boot
apiserver-boot init glide

# Install vendor/ from using "glide install --strip-vendor"
apiserver-boot init glide --fetch
`,
	Run: RunGlideInstall,
}

var fetch bool

func AddGlideInstallCmd(cmd *cobra.Command) {
	glideInstallCmd.Flags().BoolVar(&fetch, "fetch", true, "if true, fetch new glide deps instead of copying the ones packaged with the tools")
	cmd.AddCommand(glideInstallCmd)
}

func fetchGlide() {
	o, err := exec.Command("glide", "-v").CombinedOutput()
	if err != nil {
		log.Fatal("must install glide v0.12 or later")
	}
	if !strings.HasPrefix(string(o), "glide version v0.12") &&
		!strings.HasPrefix(string(o), "glide version v0.13") {
		log.Fatalf("must install glide  or later, was %s", o)
	}

	c := exec.Command("glide", "install", "--strip-vendor")
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	err = c.Run()
	if err != nil {
		log.Fatalf("failed to run glide install\n%v\n", err)
	}
}

func copyGlide() {
	// Move up two directories from the location of the `apiserver-boot`
	// executable to find the `vendor` directory we package with our
	// releases. TODO(campbellalex@google.com): this doesn't work for people
	// who used `go install` to put `apiserver-boot` in their $GOPATH/bin.
	e, err := os.Executable()
	if err != nil {
		log.Fatal("unable to get directory of apiserver-builder tools")
	}

	e = filepath.Dir(filepath.Dir(e))

	doCmd := func(cmd string, args ...string) {
		c := exec.Command(cmd, args...)
		c.Stderr = os.Stderr
		c.Stdout = os.Stdout
		err = c.Run()
		if err != nil {
			log.Fatalf("failed to copy go dependencies %v", err)
		}
	}

	doCmd("tar", "-xzvf", filepath.Join(e, "bin", "glide.tar.gz"))
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
- package: k8s.io/api
  version: c9fffff41e45e3c00186ac6b00d2cb585734d43e
- package: k8s.io/apimachinery
  version: 7da60ba7ddca684051555f2c558eef2dfebc70d5
- package: k8s.io/apiserver
  version: e24df9a2e58151a85874948908a454d511066460
- package: k8s.io/client-go
  version: 1be407b92aa39a2f63ddbb3d46104a1fd425fda0
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
		log.Fatal(err)
	}
	path := filepath.Join(dir, "glide.yaml")
	util.WriteIfNotFound(path, "glide-template", glideTemplate,
		glideTemplateArguments{
			util.Repo,
		})
}
