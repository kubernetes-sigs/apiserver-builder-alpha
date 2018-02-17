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

package release

import (
	"archive/tar"
	"compress/gzip"
//	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var repo = filepath.Join("github.com", "kubernetes-sigs", "kubebuilder")

var vendorCmd = &cobra.Command{
	Use:   "vendor",
	Short: "create vendored libraries for release",
	Long:  `create vendored libraries for release`,
	Run:   RunVendor,
}

func RunVendor(cmd *cobra.Command, args []string) {
	if len(version) == 0 {
		log.Fatal("must specify the --version flag")
	}

	// Create the release directory
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dir = filepath.Join(dir, "release", version)
	os.MkdirAll(dir, 0700)

	BuildLocalVendor(dir)
}

func BuildLocalVendor(tooldir string) {
	os.MkdirAll(filepath.Join(tooldir, "src"), 0700)
	c := exec.Command("cp", "-R", "-H",
		filepath.Join("vendor"),
		filepath.Join(tooldir, "src"))
	RunCmd(c, "")
	os.MkdirAll(filepath.Join(tooldir, "src", "vendor", repo), 0700)
	c = exec.Command("cp", "-R", "-H",
		filepath.Join("pkg"),
		filepath.Join(tooldir, "src", "vendor", repo, "pkg"))
	RunCmd(c, "")

/*	c = exec.Command("bash", "-c",
		fmt.Sprintf("find %s -name BUILD.bazel| xargs sed -i s'|//pkg|//vendor/github.com/kubernetes-sigs/kubebuilder/pkg|g'",
			filepath.Join(tooldir, "src", "vendor", repo, "pkg"),
		))
	RunCmd(c, "")
*/
	c = exec.Command("cp", "-R", "-H",
		filepath.Join("Gopkg.toml"),
		filepath.Join(tooldir, "src", "Gopkg.toml"))
	RunCmd(c, "")
	c = exec.Command("cp", "-R", "-H",
		filepath.Join("Gopkg.lock"),
		filepath.Join(tooldir, "src", "Gopkg.lock"))
	RunCmd(c, "")

}

var VendoredBuildPackages = []string{
	"github.com/coreos/etcd",
	"github.com/kubernetes-incubator/reference-docs/gen-apidocs",
	"k8s.io/code-generator/cmd/client-gen",
	"k8s.io/code-generator/cmd/conversion-gen",
	"k8s.io/code-generator/cmd/deepcopy-gen",
	"k8s.io/code-generator/cmd/defaulter-gen",
	//"k8s.io/code-generator/cmd/go-to-protobuf",
	//"k8s.io/code-generator/cmd/import-boss",
	"k8s.io/code-generator/cmd/informer-gen",
	"k8s.io/code-generator/cmd/lister-gen",
	"k8s.io/code-generator/cmd/openapi-gen",
	//"k8s.io/code-generator/cmd/set-gen",
}

var OwnedBuildPackages = []string{
	"cmd/kubebuilder-gen",
	"cmd/kubebuilder",
}

func BuildVendorTar(dir string) {
	// create the new file
	f := filepath.Join(dir, "bin", "vendor.tar.gz")
	fw, err := os.Create(f)
	if err != nil {
		log.Fatalf("failed to create vendor tar file %s %v", f, err)
	}
	defer fw.Close()

	// setup gzip of tar
	gw := gzip.NewWriter(fw)
	defer gw.Close()

	// setup tar writer
	tw := tar.NewWriter(gw)
	defer tw.Close()

	srcdir := filepath.Join(dir)
	filepath.Walk(srcdir, TarFile{
		tw,
		0644,
		filepath.Join(srcdir, "src"),
		"",
	}.Do)
}
