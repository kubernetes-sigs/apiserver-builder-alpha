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
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/parser"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/klog"
	"sigs.k8s.io/apiserver-builder-alpha/cmd/apiserver-boot/boot/util"

	gengoargs "k8s.io/gengo/args"

	defaultergenargs "k8s.io/code-generator/cmd/defaulter-gen/args"
	defaultergen "k8s.io/gengo/examples/defaulter-gen/generators"

	conversiongenargs "k8s.io/code-generator/cmd/conversion-gen/args"
	conversiongen "k8s.io/code-generator/cmd/conversion-gen/generators"

	deepcopygenargs "k8s.io/code-generator/cmd/deepcopy-gen/args"
	deepcopygen "k8s.io/gengo/examples/deepcopy-gen/generators"

	openapigenargs "k8s.io/kube-openapi/cmd/openapi-gen/args"
	openapigen "k8s.io/kube-openapi/pkg/generators"

	apiregistergen "sigs.k8s.io/apiserver-builder-alpha/cmd/apiregister-gen/generators"
)

var versionedAPIs []string
var unversionedAPIs []string
var codegenerators []string
var copyright string
var generators = sets.String{}
var vendorDir string

var parserBuilder *parser.Builder

var generateCmd = &cobra.Command{
	Use:   "generated",
	Short: "Run code generators against repo.",
	Long:  `Automatically run by most build commands.  Writes generated source code for a repo.`,
	Example: `# Run code generators.
apiserver-boot build generated`,
	Run: RunGenerate,
}

func AddGenerate(cmd *cobra.Command) {
	cmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVar(&copyright, "copyright", "boilerplate.go.txt", "Location of copyright boilerplate file.")
	generateCmd.Flags().StringVar(&vendorDir, "vendor-dir", "", "Location of directory containing vendor files.")
	generateCmd.Flags().StringArrayVar(&versionedAPIs, "api-versions", []string{}, "API version to generate code for.  Can be specified multiple times.  e.g. --api-versions foo/v1beta1 --api-versions bar/v1  defaults to all versions found under directories pkg/apis/<group>/<version>")
	generateCmd.Flags().StringArrayVar(&codegenerators, "generator", []string{}, "list of generators to run.  e.g. --generator apiregister --generator conversion Valid values: [apiregister,conversion,client,deepcopy,defaulter,openapi,protobuf]")
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
	versionedAPIPkgs, unversionedAPIPkgs := make([]string, 0), make([]string, 0)
	for _, api := range versionedAPIs {
		versionedAPIPkgs = append(versionedAPIPkgs, filepath.Join(util.Repo, "pkg", "apis", api))
	}
	for _, api := range unversionedAPIs {
		unversionedAPIPkgs = append(unversionedAPIPkgs, filepath.Join(util.Repo, "pkg", "apis", api))
	}

	if doGen("apiregister-gen") {
		apisPkg := filepath.Join(util.Repo, "pkg", "apis")
		runApiRegisterGen(apisPkg, versionedAPIPkgs, unversionedAPIPkgs, []string{
			"k8s.io/apimachinery/pkg/apis/meta/v1",
			"k8s.io/apimachinery/pkg/api/resource",
			"k8s.io/apimachinery/pkg/version",
			"k8s.io/apimachinery/pkg/runtime",
			"k8s.io/apimachinery/pkg/util/intstr",
			"k8s.io/api/core/v1",
			"k8s.io/api/apps/v1",
		})
	}

	if doGen("defaulter-gen") {
		runDefaulterGen(versionedAPIPkgs, unversionedAPIPkgs)
	}

	if doGen("conversion-gen") {
		runConversionGen(versionedAPIPkgs, unversionedAPIPkgs)
	}

	if doGen("deepcopy-gen") {
		runDeepcopyGen(versionedAPIPkgs, unversionedAPIPkgs)
	}

	if doGen("openapi-gen") {
		runOpenApiGen()
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

func runApiRegisterGen(apisPkg string, versionedAPIPkgs, unversionedAPIPkgs, additionalPkgs []string) {
	klog.Info("running apiregister-gen..")
	apiregisterGenArgs := gengoargs.Default()
	apiregisterGenArgs.GoHeaderFilePath = filepath.Join(gengoargs.DefaultSourceTree(), util.Repo, copyright)
	apiregisterGenArgs.InputDirs = append(versionedAPIPkgs,
		append(unversionedAPIPkgs,
			append(additionalPkgs, apisPkg)...)...)
	apiregisterGenArgs.OutputFileBaseName = "zz_generated.api.register"
	apiregisterGenArgs.CustomArgs = apiregistergen.CustomArgs{}
	gen := &apiregistergen.Gen{}
	if err := execute(
		apiregisterGenArgs,
		gen.NameSystems(),
		gen.DefaultNameSystem(),
		gen.Packages); err != nil {
		klog.Fatalf("Error: %v", err)
	}
}

func runConversionGen(versionedAPIPkgs, unversionedAPIPkgs []string) {
	klog.Info("running conversion-gen..")
	conversionGenArgs, customArgs := conversiongenargs.NewDefaults()
	conversionGenArgs.OutputBase = util.GoSrc
	conversionGenArgs.GoHeaderFilePath = filepath.Join(gengoargs.DefaultSourceTree(), util.Repo, copyright)
	conversionGenArgs.InputDirs = append(
		append(versionedAPIPkgs, unversionedAPIPkgs...),
		"k8s.io/apimachinery/pkg/runtime")
	conversionGenArgs.OutputFileBaseName = "zz_generated.conversion"
	customArgs.ExtraPeerDirs = []string{
		"k8s.io/apimachinery/pkg/apis/meta/v1",
		"k8s.io/apimachinery/pkg/conversion",
		"k8s.io/apimachinery/pkg/runtime",
	}
	if err := execute(
		conversionGenArgs,
		conversiongen.NameSystems(),
		conversiongen.DefaultNameSystem(),
		conversiongen.Packages); err != nil {
		klog.Fatalf("failed to run defaulter-gen: %v", err)
	}
}

func runDefaulterGen(versionedAPIPkgs, unversionedAPIPkgs []string) {
	klog.Info("running defaulter-gen..")
	defaulterGenArgs, customArgs := defaultergenargs.NewDefaults()
	defaulterGenArgs.InputDirs = append(versionedAPIPkgs, unversionedAPIPkgs...)
	defaulterGenArgs.GoHeaderFilePath = filepath.Join(gengoargs.DefaultSourceTree(), util.Repo, copyright)
	defaulterGenArgs.OutputBase = util.GoSrc
	defaulterGenArgs.OutputFileBaseName = "zz_generated.defaults"
	customArgs.ExtraPeerDirs = []string{
		"k8s.io/apimachinery/pkg/apis/meta/v1",
		"k8s.io/apimachinery/pkg/conversion",
		"k8s.io/apimachinery/pkg/runtime",
	}

	if err := execute(
		defaulterGenArgs,
		defaultergen.NameSystems(),
		defaultergen.DefaultNameSystem(),
		defaultergen.Packages); err != nil {
		klog.Fatalf("failed to run defaulter-gen: %v", err)
	}
}

func runDeepcopyGen(versionedAPIPkgs, unversionedAPIPkgs []string) {
	klog.Info("running deepcopy-gen..")
	deepcopyGenArgs, _ := deepcopygenargs.NewDefaults()
	deepcopyGenArgs.InputDirs = append(versionedAPIPkgs, unversionedAPIPkgs...)
	deepcopyGenArgs.GoHeaderFilePath = filepath.Join(gengoargs.DefaultSourceTree(), util.Repo, copyright)
	deepcopyGenArgs.OutputBase = util.GoSrc
	deepcopyGenArgs.OutputFileBaseName = "zz_generated.deepcopy"
	if err := execute(
		deepcopyGenArgs,
		deepcopygen.NameSystems(),
		deepcopygen.DefaultNameSystem(),
		deepcopygen.Packages); err != nil {
		klog.Fatalf("Error: %v", err)
	}
}

func runOpenApiGen() {
	klog.Info("running openapi-gen..")
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

	openapiGenArgs, customArgs := openapigenargs.NewDefaults()
	openapiGenArgs.OutputBase = util.GoSrc
	openapiGenArgs.GoHeaderFilePath = filepath.Join(gengoargs.DefaultSourceTree(), util.Repo, copyright)
	openapiGenArgs.InputDirs = apis
	openapiGenArgs.OutputPackagePath = filepath.Join(util.Repo, "pkg", "openapi")
	openapiGenArgs.OutputFileBaseName = "zz_generated.openapi"
	customArgs.ReportFilename = "violations.report"

	// Generates the code for the OpenAPIDefinitions.
	if err := execute(
		openapiGenArgs,
		openapigen.NameSystems(),
		openapigen.DefaultNameSystem(),
		openapigen.Packages); err != nil {
		klog.Fatalf("OpenAPI code generation error: %v", err)
	}
}

func execute(
	g *gengoargs.GeneratorArgs,
	nameSystems namer.NameSystems,
	defaultSystem string,
	pkgs func(*generator.Context, *gengoargs.GeneratorArgs) generator.Packages) error {

	if parserBuilder == nil {
		klog.Info("scanning packages (can take a few minutes)...")
		b, err := g.NewBuilder()
		if err != nil {
			return fmt.Errorf("Failed making a parser: %v", err)
		}
		parserBuilder = b
	}
	b := parserBuilder

	// pass through the flag on whether to include *_test.go files
	b.IncludeTestFiles = g.IncludeTestFiles

	c, err := generator.NewContext(b, nameSystems, defaultSystem)
	if err != nil {
		return fmt.Errorf("Failed making a context: %v", err)
	}

	c.Verify = g.VerifyOnly
	packages := pkgs(c, g)
	if err := c.ExecutePackages(g.OutputBase, packages); err != nil {
		return fmt.Errorf("Failed executing generator: %v", err)
	}

	return nil

}
