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
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var versionedAPIs []string
var unversionedAPIs []string

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates source code for a repo",
	Long:  `Generates source code for a repo`,
	Run:   RunGenerate,
}

const genericAPI = "k8s.io/apimachinery/pkg/apis/meta/v1,k8s.io/apimachinery/pkg/api/resource/,k8s.io/apimachinery/pkg/version/,k8s.io/apimachinery/pkg/runtime/,k8s.io/apimachinery/pkg//util/intstr/"
const extraAPI = "k8s.io/apimachinery/pkg/apis/meta/v1,k8s.io/apimachinery/pkg/conversion,k8s.io/apimachinery/pkg/runtime"

func AddGenerate(cmd *cobra.Command) {
	cmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVar(&copyright, "copyright", "boilerplate.go.txt", "path to copyright file.  defaults to boilerplate.go.txt")
	generateCmd.Flags().StringArrayVar(&versionedAPIs, "api-versions", []string{}, "comma separated list of APIs versions.  e.g. foo/v1beta1,bar/v1  defaults to all directories under pkd/apis/group/version")
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

	getCopyright()

	root, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(-1)
	}
	root = filepath.Dir(root)

	all := []string{}
	versioned := []string{}
	for _, v := range versionedAPIs {
		v = filepath.Join(Repo, "pkg", "apis", v)
		versioned = append(versioned, "--input-dirs", v)
		all = append(all, "--input-dirs", v)
	}
	unversioned := []string{}
	for _, u := range unversionedAPIs {
		u = filepath.Join(Repo, "pkg", "apis", u)
		unversioned = append(unversioned, "--input-dirs", u)
		all = append(all, "--input-dirs", u)
	}

	c := exec.Command(filepath.Join(root, "apiregister-gen"),
		"--input-dirs", filepath.Join(Repo, "pkg", "apis", "..."),
		"--input-dirs", filepath.Join(Repo, "pkg", "controller", "..."),
	)
	fmt.Printf("%s\n", strings.Join(c.Args, " "))
	out, err := c.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to run apiregister-gen %s %v\n", out, err)
		os.Exit(-1)
	}

	c = exec.Command(filepath.Join(root, "conversion-gen"),
		append(all,
			"-o", GoSrc,
			"--go-header-file", copyright,
			"-O", "zz_generated.conversion",
			"--extra-peer-dirs", extraAPI)...,
	)
	fmt.Printf("%s\n", strings.Join(c.Args, " "))
	out, err = c.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to run conversion-gen %s %v\n", out, err)
		os.Exit(-1)
	}

	c = exec.Command(filepath.Join(root, "deepcopy-gen"),
		append(all,
			"-o", GoSrc,
			"--go-header-file", copyright,
			"-O", "zz_generated.deepcopy")...,
	)
	fmt.Printf("%s\n", strings.Join(c.Args, " "))
	out, err = c.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to run deepcopy-gen %s %v\n", out, err)
		os.Exit(-1)
	}

	c = exec.Command(filepath.Join(root, "openapi-gen"),
		append(all,
			"-o", GoSrc,
			"--go-header-file", copyright,
			"-i", genericAPI,
			"--output-package", filepath.Join(Repo, "pkg", "openapi"))...,
	)
	fmt.Printf("%s\n", strings.Join(c.Args, " "))
	out, err = c.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to run openapi-gen %s %v\n", out, err)
		os.Exit(-1)
	}

	c = exec.Command(filepath.Join(root, "defaulter-gen"),
		append(all,
			"-o", GoSrc,
			"--go-header-file", copyright,
			"-O", "zz_generated.defaults",
			"--extra-peer-dirs=", extraAPI)...,
	)
	fmt.Printf("%s\n", strings.Join(c.Args, " "))
	out, err = c.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to run defaulter-gen %s %v\n", out, err)
		os.Exit(-1)
	}

	// Builder the versioned apis client
	clientPkg := filepath.Join(Repo, "pkg", "client")
	clientset := filepath.Join(clientPkg, "clientset_generated")
	c = exec.Command(filepath.Join(root, "client-gen"),
		"-o", GoSrc,
		"--go-header-file", copyright,
		"--input-base", filepath.Join(Repo, "pkg", "apis"),
		"--input", strings.Join(versionedAPIs, ","),
		"--clientset-path", clientset,
		"--clientset-name", "clientset",
	)
	fmt.Printf("%s\n", strings.Join(c.Args, " "))
	out, err = c.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to run client-gen %s %v\n", out, err)
		os.Exit(-1)
	}

	c = exec.Command(filepath.Join(root, "client-gen"),
		"-o", GoSrc,
		"--go-header-file", copyright,
		"--input-base", filepath.Join(Repo, "pkg", "apis"),
		"--input", strings.Join(unversionedAPIs, ","),
		"--clientset-path", clientset,
		"--clientset-name", "internalclientset")
	fmt.Printf("%s\n", strings.Join(c.Args, " "))
	out, err = c.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to run client-gen for unversioned APIs %s %v\n", out, err)
		os.Exit(-1)
	}

	listerPkg := filepath.Join(clientPkg, "listers_generated")
	c = exec.Command(filepath.Join(root, "lister-gen"),
		append(all,
			"-o", GoSrc,
			"--go-header-file", copyright,
			"--output-package", listerPkg)...,
	)
	fmt.Printf("%s\n", strings.Join(c.Args, " "))
	out, err = c.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to run lister-gen %s %v\n", out, err)
		os.Exit(-1)
	}

	informerPkg := filepath.Join(clientPkg, "informers_generated")
	c = exec.Command(filepath.Join(root, "informer-gen"),
		append(all,
			"-o", GoSrc,
			"--go-header-file", copyright,
			"--output-package", informerPkg,
			"--listers-package", listerPkg,
			"--versioned-clientset-package", filepath.Join(clientset, "clientset"),
			"--internal-clientset-package", filepath.Join(clientset, "internalclientset"))...,
	)
	fmt.Printf("%s\n", strings.Join(c.Args, " "))
	out, err = c.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to run informer-gen %s %v\n", out, err)
		os.Exit(-1)
	}
}

func initApis() {
	if len(versionedAPIs) == 0 {
		groups, err := ioutil.ReadDir(filepath.Join("pkg", "apis"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not read pkg/apis directory to find api versions\n")
			os.Exit(-1)
		}
		for _, g := range groups {
			if g.IsDir() {
				versionFiles, err := ioutil.ReadDir(filepath.Join("pkg", "apis", g.Name()))
				if err != nil {
					fmt.Fprintf(os.Stderr, "could not read pkg/apis/%s directory to find api versions\n", g.Name())
					os.Exit(-1)
				}
				for _, v := range versionFiles {
					if v.IsDir() {
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
