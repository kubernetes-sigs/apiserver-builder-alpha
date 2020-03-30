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
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/klog"
	"sigs.k8s.io/apiserver-builder-alpha/cmd/apiserver-boot/boot/util"
)

var versionedAPIs []string
var unversionedAPIs []string
var codegenerators []string
var copyright string
var generators = sets.String{}
var vendorDir string

var generateCmd = &cobra.Command{
	Use:   "generated",
	Short: "Run code generators against repo.",
	Long:  `Automatically run by most build commands.  Writes generated source code for a repo.`,
	Example: `# Run code generators.
apiserver-boot build generated`,
	Run: RunGenerate,
}

var extraAPI = strings.Join([]string{
	"k8s.io/apimachinery/pkg/apis/meta/v1",
	"k8s.io/apimachinery/pkg/conversion",
	"k8s.io/apimachinery/pkg/runtime"}, ",")

func AddGenerate(cmd *cobra.Command) {
	cmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVar(&copyright, "copyright", "boilerplate.go.txt", "Location of copyright boilerplate file.")
	generateCmd.Flags().StringVar(&vendorDir, "vendor-dir", "", "Location of directory containing vendor files.")
	generateCmd.Flags().StringArrayVar(&versionedAPIs, "api-versions", []string{}, "API version to generate code for.  Can be specified multiple times.  e.g. --api-versions foo/v1beta1 --api-versions bar/v1  defaults to all versions found under directories pkg/apis/<group>/<version>")
	generateCmd.Flags().StringArrayVar(&codegenerators, "generator", []string{}, "list of generators to run.  e.g. --generator apiregister --generator conversion Valid values: [apiregister,conversion,client,deepcopy,defaulter,openapi,protobuf]")
	generateCmd.AddCommand(generateCleanCmd)

	generateCleanCmd.Flags().MarkDeprecated("gen-unversioned-client", "generate unversioned client is highly unrecommended, please use versioned client instead")
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

func doGen(g string) bool {
	switch g {
	case "go-to-protobuf":
		// disable protobuf generation by default
		return generators.Has("protobuf")
	default:
		g = strings.Replace(g, "-gen", "", -1)
	}
	return generators.Has(g) || generators.Len() == 0
}

func RunGenerate(cmd *cobra.Command, args []string) {
	initApis()

	for _, g := range codegenerators {
		generators.Insert(strings.Replace(g, "-gen", "", -1))
	}

	util.GetCopyright(copyright)

	root, err := os.Executable()
	if err != nil {
		klog.Fatalf("error: %v", err)
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
	}

	if doGen("apiregister-gen") {
		inputDirsArgs := []string{
			"--input-dirs", filepath.Join(util.Repo, "pkg", "apis", "..."),
		}
		controllerPkgs := filepath.Join(util.Repo, "pkg", "controller", "...")
		if _, err := os.Stat(filepath.Join(util.GoSrc, util.Repo, "pkg", "controller")); err == nil {
			inputDirsArgs = append(inputDirsArgs, "--input-dirs", controllerPkgs)
		} else {
			klog.Warningf("ignoring controller package code-generation due to %v", err)
		}
		inputDirsArgs = append(inputDirsArgs, "--go-header-file", copyright)

		c := exec.Command(filepath.Join(root, "apiregister-gen"), inputDirsArgs...)
		klog.Infof("%s", strings.Join(c.Args, " "))
		out, err := c.CombinedOutput()
		if err != nil {
			klog.Fatalf("failed to run apiregister-gen %s %v", out, err)
		}
	}

	if doGen("conversion-gen") {
		c := exec.Command(filepath.Join(root, "conversion-gen"),
			append(append(all, unversioned...),
				"-o", util.GoSrc,
				"--go-header-file", copyright,
				"-O", "zz_generated.conversion",
				"--extra-peer-dirs", extraAPI)...,
		)
		klog.Infof("%s", strings.Join(c.Args, " "))
		out, err := c.CombinedOutput()
		if err != nil {
			klog.Fatalf("failed to run conversion-gen %s %v", out, err)
		}
	}

	if doGen("deepcopy-gen") {
		c := exec.Command(filepath.Join(root, "deepcopy-gen"),
			append(append(all, unversioned...),
				"-o", util.GoSrc,
				"--go-header-file", copyright,
				"-O", "zz_generated.deepcopy")...,
		)
		klog.Infof("%s", strings.Join(c.Args, " "))
		out, err := c.CombinedOutput()
		if err != nil {
			klog.Fatalf("failed to run deepcopy-gen %s %v", out, err)
		}
	}

	if doGen("openapi-gen") {
		apis := []string{
			"k8s.io/apimachinery/pkg/apis/meta/v1",
			"k8s.io/apimachinery/pkg/api/resource",
			"k8s.io/apimachinery/pkg/version",
			"k8s.io/apimachinery/pkg/runtime",
			"k8s.io/apimachinery/pkg/util/intstr",
			"k8s.io/api/core/v1",
			"k8s.io/api/apps/v1",
		}

		// Add any vendored apis from core
		apis = append(apis, getVendorApis(filepath.Join("k8s.io", "api"))...)
		apis = append(apis, getVendorApis(filepath.Join("k8s.io", "client-go", "pkg", "apis"))...)

		// Special case 'k8s.io/client-go/pkg/api/v1' because it does not have a group
		if _, err := os.Stat(filepath.Join("vendor", "k8s.io", "client-go", "pkg", "api", "v1", "doc.go")); err == nil {
			apis = append(apis, filepath.Join("k8s.io", "client-go", "pkg", "api", "v1"))
		}

		if _, err := os.Stat(filepath.Join("vendor", "k8s.io", "api", "core", "v1", "doc.go")); err == nil {
			apis = append(apis, filepath.Join("k8s.io", "api", "core", "v1"))
		}

		c := exec.Command(filepath.Join(root, "openapi-gen"),
			append(all,
				"-o", util.GoSrc,
				"--go-header-file", copyright,
				"-i", strings.Join(apis, ","),
				"--report-filename", "violations.report",
				"--output-package", filepath.Join(util.Repo, "pkg", "openapi"))...,
		)

		// HACK: ensure GOROOT env var
		c.Env = os.Environ()
		if len(os.Getenv("GOROOT")) == 0 {
			if p, err := exec.Command("which", "go").CombinedOutput(); err == nil {
				// The returned string will have some/path/bin/go, so remove the last two elements.
				c.Env = append(c.Env,
					fmt.Sprintf("GOROOT=%s", filepath.Dir(filepath.Dir(strings.Trim(string(p), "\n")))))
			} else {
				klog.Warningf("Warning: $GOROOT not set, and unable to run `which go` to find it: %v", err)
			}
		}

		klog.Infof("%s", strings.Join(c.Args, " "))
		out, err := c.CombinedOutput()
		if err != nil {
			klog.Fatalf("failed to run openapi-gen %s %v", out, err)
		}
	}

	if doGen("defaulter-gen") {
		c := exec.Command(filepath.Join(root, "defaulter-gen"),
			append(append(all, unversioned...),
				"-o", util.GoSrc,
				"--go-header-file", copyright,
				"-O", "zz_generated.defaults",
				"--extra-peer-dirs=", extraAPI)...,
		)
		klog.Infof("%s", strings.Join(c.Args, " "))
		out, err := c.CombinedOutput()
		if err != nil {
			klog.Fatalf("failed to run defaulter-gen %s %v", out, err)
		}
	}

	if doGen("client-gen") {
		// Builder the versioned apis client
		clientPkg := filepath.Join(util.Repo, "pkg", "client")
		clientset := filepath.Join(clientPkg, "clientset_generated")
		c := exec.Command(filepath.Join(root, "client-gen"),
			"-o", util.GoSrc,
			"--go-header-file", copyright,
			"--input-base", filepath.Join(util.Repo, "pkg", "apis"),
			"--input", strings.Join(versionedAPIs, ","),
			"--clientset-path", clientset,
			"--clientset-name", "clientset",
		)
		klog.Infof("%s", strings.Join(c.Args, " "))
		out, err := c.CombinedOutput()
		if err != nil {
			klog.Fatalf("failed to run client-gen %s %v", out, err)
		}

		toGen := versioned
		if false /* generating unversioned client is deprecated */ {
			toGen = all
			c = exec.Command(filepath.Join(root, "client-gen"),
				"-o", util.GoSrc,
				"--go-header-file", copyright,
				"--input-base", filepath.Join(util.Repo, "pkg", "apis"),
				"--input", strings.Join(unversionedAPIs, ","),
				"--clientset-path", clientset,
				"--clientset-name", "internalclientset")
			klog.Infof("%s", strings.Join(c.Args, " "))
			out, err = c.CombinedOutput()
			if err != nil {
				klog.Fatalf("failed to run client-gen for unversioned APIs %s %v", out, err)
			}
		}

		listerPkg := filepath.Join(clientPkg, "listers_generated")
		c = exec.Command(filepath.Join(root, "lister-gen"),
			append(toGen,
				"-o", util.GoSrc,
				"--go-header-file", copyright,
				"--output-package", listerPkg)...,
		)
		klog.Infof("%s", strings.Join(c.Args, " "))
		out, err = c.CombinedOutput()
		if err != nil {
			klog.Fatalf("failed to run lister-gen %s %v", out, err)
		}

		informerPkg := filepath.Join(clientPkg, "informers_generated")
		c = exec.Command(filepath.Join(root, "informer-gen"),
			append(toGen,
				"-o", util.GoSrc,
				"--go-header-file", copyright,
				"--output-package", informerPkg,
				"--listers-package", listerPkg,
				"--versioned-clientset-package", filepath.Join(clientset, "clientset"))...,
		)
		klog.Infof("%s", strings.Join(c.Args, " "))
		out, err = c.CombinedOutput()
		if err != nil {
			klog.Fatalf("failed to run informer-gen %s %v", out, err)
		}
	}

	if doGen("go-to-protobuf") {
		versionedAPIPackages := sets.NewString()
		for _, versionedAPI := range versionedAPIs {
			versionedAPIPackages.Insert(filepath.Join(util.Repo, "pkg", "apis", versionedAPI))
		}
		c := exec.Command(filepath.Join(root, "go-to-protobuf"),
			"--packages", strings.Join(versionedAPIPackages.List(), ","),
			"--apimachinery-packages", strings.Join([]string{
				"-k8s.io/apimachinery/pkg/util/intstr",
				"-k8s.io/apimachinery/pkg/api/resource",
				"-k8s.io/apimachinery/pkg/runtime/schema",
				"-k8s.io/apimachinery/pkg/runtime",
				"-k8s.io/apimachinery/pkg/apis/meta/v1",
				"-sigs.k8s.io/apiserver-builder-alpha/pkg/builders",
			}, ","),
			"--drop-embedded-fields", strings.Join([]string{
				"k8s.io/apimachinery/pkg/apis/meta/v1.TypeMeta",
				"k8s.io/apimachinery/pkg/runtime.Serializer",
			}, ","),
			"--proto-import=./vendor",
			"--vendor-output-base=./vendor/",
		)
		klog.Infof("%s", strings.Join(c.Args, " "))
		out, err := c.CombinedOutput()
		if err != nil {
			klog.Fatalf("failed to run go-to-protobuf %s %v", out, err)
		}
	}
}

func getVendorApis(pkg string) []string {
	dir := filepath.Join("vendor", pkg)
	if len(vendorDir) >= 0 {
		dir = filepath.Join(vendorDir, dir)
	}
	apis := []string{}
	if groups, err := ioutil.ReadDir(dir); err == nil {
		for _, g := range groups {
			p := filepath.Join(dir, g.Name())
			if g.IsDir() {
				if versions, err := ioutil.ReadDir(p); err == nil {
					for _, v := range versions {
						versionMatch := regexp.MustCompile("^v\\d+(alpha\\d+|beta\\d+)*$")
						if v.IsDir() && versionMatch.MatchString(v.Name()) {
							apis = append(apis, filepath.Join(pkg, g.Name(), v.Name()))
						}
					}
				}
			}
		}
	}
	return apis
}

func initApis() {
	if len(versionedAPIs) == 0 {
		groups, err := ioutil.ReadDir(filepath.Join("pkg", "apis"))
		if err != nil {
			klog.Fatalf("could not read pkg/apis directory to find api Versions")
		}
		for _, g := range groups {
			if g.IsDir() {
				versionFiles, err := ioutil.ReadDir(filepath.Join("pkg", "apis", g.Name()))
				if err != nil {
					klog.Fatalf("could not read pkg/apis/%s directory to find api Versions", g.Name())
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
