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

package build

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/kubernetes-incubator/apiserver-builder/cmd/apiserver-boot/boot/util"
	"github.com/spf13/cobra"
	"regexp"
)

var versionedAPIs []string
var unversionedAPIs []string
var copyright string = "boilerplate.go.txt"

var generateCmd = &cobra.Command{
	Use:   "generated",
	Short: "Run code generators against repo.",
	Long:  `Automatically run by most build commands.  Writes generated source code for a repo.`,
	Example: `# Run code generators.
apiserver-boot build generated`,
	Run: RunGenerate,
}

var genericAPI = strings.Join([]string{
	"k8s.io/client-go/pkg/api/v1",
	"k8s.io/client-go/pkg/apis/apps/v1beta1",
	"k8s.io/client-go/pkg/apis/authentication/v1",
	"k8s.io/client-go/pkg/apis/authentication/v1beta1",
	"k8s.io/client-go/pkg/apis/authorization/v1",
	"k8s.io/client-go/pkg/apis/authorization/v1beta1",
	"k8s.io/client-go/pkg/apis/autoscaling/v1",
	"k8s.io/client-go/pkg/apis/autoscaling/v2alpha1",
	"k8s.io/client-go/pkg/apis/batch/v1",
	"k8s.io/client-go/pkg/apis/batch/v2alpha1",
	"k8s.io/client-go/pkg/apis/certificates/v1beta1",
	"k8s.io/client-go/pkg/apis/extensions/v1beta1",
	"k8s.io/client-go/pkg/apis/policy/v1beta1",
	"k8s.io/client-go/pkg/apis/rbac/v1alpha1",
	"k8s.io/client-go/pkg/apis/rbac/v1beta1",
	"k8s.io/client-go/pkg/apis/settings/v1alpha1",
	"k8s.io/client-go/pkg/apis/storage/v1",
	"k8s.io/client-go/pkg/apis/storage/v1beta1",
	"k8s.io/apimachinery/pkg/apis/meta/v1",
	"k8s.io/apimachinery/pkg/api/resource",
	"k8s.io/apimachinery/pkg/version",
	"k8s.io/apimachinery/pkg/runtime",
	"k8s.io/apimachinery/pkg/util/intstr"}, ",")

var extraAPI = strings.Join([]string{
	"k8s.io/apimachinery/pkg/apis/meta/v1",
	"k8s.io/apimachinery/pkg/conversion",
	"k8s.io/apimachinery/pkg/runtime"}, ",")

func AddGenerate(cmd *cobra.Command) {
	cmd.AddCommand(generateCmd)
	generateCmd.Flags().StringArrayVar(&versionedAPIs, "api-versions", []string{}, "comma separated list of APIs Versions.  e.g. foo/v1beta1,bar/v1  defaults to all directories under pkd/apis/group/version")
	generateCmd.AddCommand(generateCleanCmd)
}

var generateCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Removes generated source code",
	Long:  `Removes generated source code`,
	Run:   RunCleanGenerate,
}

func RunCleanGenerate(cmd *cobra.Command, args []string) {
	os.RemoveAll(filepath.Join("pkg", "client", "clientset_generated"))
	os.RemoveAll(filepath.Join("pkg", "client", "informers_generated"))
	os.RemoveAll(filepath.Join("pkg", "client", "listers_generated"))
	os.Remove(filepath.Join("pkg", "openapi", "openapi_generated.go"))

	filepath.Walk("pkg", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasPrefix(info.Name(), "zz_generated.") {
			return os.Remove(path)
		}
		return nil
	})
}

func RunGenerate(cmd *cobra.Command, args []string) {
	initApis()

	util.GetCopyright(copyright)

	root, err := os.Executable()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	root = filepath.Dir(root)

	all := []string{}
	versioned := []string{}
	for _, v := range versionedAPIs {
		v = filepath.Join(util.Repo, "pkg", "apis", v)
		versioned = append(versioned, "--input-dirs", v)
		all = append(all, "--input-dirs", v)
	}
	unversioned := []string{}
	for _, u := range unversionedAPIs {
		u = filepath.Join(util.Repo, "pkg", "apis", u)
		unversioned = append(unversioned, "--input-dirs", u)
		all = append(all, "--input-dirs", u)
	}

	c := exec.Command(filepath.Join(root, "apiregister-gen"),
		"--input-dirs", filepath.Join(util.Repo, "pkg", "apis", "..."),
		"--input-dirs", filepath.Join(util.Repo, "pkg", "controller", "..."),
	)
	fmt.Printf("%s\n", strings.Join(c.Args, " "))
	out, err := c.CombinedOutput()
	if err != nil {
		log.Fatalf("failed to run apiregister-gen %s %v", out, err)
	}

	c = exec.Command(filepath.Join(root, "conversion-gen"),
		append(all,
			"-o", util.GoSrc,
			"--go-header-file", copyright,
			"-O", "zz_generated.conversion",
			"--extra-peer-dirs", extraAPI)...,
	)
	fmt.Printf("%s\n", strings.Join(c.Args, " "))
	out, err = c.CombinedOutput()
	if err != nil {
		log.Fatalf("failed to run conversion-gen %s %v", out, err)
	}

	c = exec.Command(filepath.Join(root, "deepcopy-gen"),
		append(all,
			"-o", util.GoSrc,
			"--go-header-file", copyright,
			"-O", "zz_generated.deepcopy")...,
	)
	fmt.Printf("%s\n", strings.Join(c.Args, " "))
	out, err = c.CombinedOutput()
	if err != nil {
		log.Fatalf("failed to run deepcopy-gen %s %v", out, err)
	}

	c = exec.Command(filepath.Join(root, "openapi-gen"),
		append(all,
			"-o", util.GoSrc,
			"--go-header-file", copyright,
			"-i", genericAPI,
			"--output-package", filepath.Join(util.Repo, "pkg", "openapi"))...,
	)
	fmt.Printf("%s\n", strings.Join(c.Args, " "))
	out, err = c.CombinedOutput()
	if err != nil {
		log.Fatalf("failed to run openapi-gen %s %v", out, err)
	}

	c = exec.Command(filepath.Join(root, "defaulter-gen"),
		append(all,
			"-o", util.GoSrc,
			"--go-header-file", copyright,
			"-O", "zz_generated.defaults",
			"--extra-peer-dirs=", extraAPI)...,
	)
	fmt.Printf("%s\n", strings.Join(c.Args, " "))
	out, err = c.CombinedOutput()
	if err != nil {
		log.Fatalf("failed to run defaulter-gen %s %v", out, err)
	}

	// Builder the versioned apis client
	clientPkg := filepath.Join(util.Repo, "pkg", "client")
	clientset := filepath.Join(clientPkg, "clientset_generated")
	c = exec.Command(filepath.Join(root, "client-gen"),
		"-o", util.GoSrc,
		"--go-header-file", copyright,
		"--input-base", filepath.Join(util.Repo, "pkg", "apis"),
		"--input", strings.Join(versionedAPIs, ","),
		"--clientset-path", clientset,
		"--clientset-name", "clientset",
	)
	fmt.Printf("%s\n", strings.Join(c.Args, " "))
	out, err = c.CombinedOutput()
	if err != nil {
		log.Fatalf("failed to run client-gen %s %v", out, err)
	}

	c = exec.Command(filepath.Join(root, "client-gen"),
		"-o", util.GoSrc,
		"--go-header-file", copyright,
		"--input-base", filepath.Join(util.Repo, "pkg", "apis"),
		"--input", strings.Join(unversionedAPIs, ","),
		"--clientset-path", clientset,
		"--clientset-name", "internalclientset")
	fmt.Printf("%s\n", strings.Join(c.Args, " "))
	out, err = c.CombinedOutput()
	if err != nil {
		log.Fatalf("failed to run client-gen for unversioned APIs %s %v", out, err)
	}

	listerPkg := filepath.Join(clientPkg, "listers_generated")
	c = exec.Command(filepath.Join(root, "lister-gen"),
		append(all,
			"-o", util.GoSrc,
			"--go-header-file", copyright,
			"--output-package", listerPkg)...,
	)
	fmt.Printf("%s\n", strings.Join(c.Args, " "))
	out, err = c.CombinedOutput()
	if err != nil {
		log.Fatalf("failed to run lister-gen %s %v", out, err)
	}

	informerPkg := filepath.Join(clientPkg, "informers_generated")
	c = exec.Command(filepath.Join(root, "informer-gen"),
		append(all,
			"-o", util.GoSrc,
			"--go-header-file", copyright,
			"--output-package", informerPkg,
			"--listers-package", listerPkg,
			"--versioned-clientset-package", filepath.Join(clientset, "clientset"),
			"--internal-clientset-package", filepath.Join(clientset, "internalclientset"))...,
	)
	fmt.Printf("%s\n", strings.Join(c.Args, " "))
	out, err = c.CombinedOutput()
	if err != nil {
		log.Fatalf("failed to run informer-gen %s %v", out, err)
	}
}

func initApis() {
	if len(versionedAPIs) == 0 {
		groups, err := ioutil.ReadDir(filepath.Join("pkg", "apis"))
		if err != nil {
			log.Fatalf("could not read pkg/apis directory to find api Versions")
		}
		for _, g := range groups {
			if g.IsDir() {
				versionFiles, err := ioutil.ReadDir(filepath.Join("pkg", "apis", g.Name()))
				if err != nil {
					log.Fatalf("could not read pkg/apis/%s directory to find api Versions", g.Name())
				}
				versionMatch := regexp.MustCompile("^v\\d+(alpha\\d+|beta\\d+)*$")
				for _, v := range versionFiles {
					if v.IsDir() && versionMatch.MatchString(v.Name()) {
						versionedAPIs = append(versionedAPIs, filepath.Join(g.Name(), v.Name()))
					}
				}
			}
		}
	}
	u := map[string]bool{}
	for _, a := range versionedAPIs {
		u[path.Dir(a)] = true
	}
	for a, _ := range u {
		unversionedAPIs = append(unversionedAPIs, a)
	}
}
