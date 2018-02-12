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
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "build the binaries",
	Long:  `build the binaries`,
	Run:   RunBuild,
}

func RunBuild(cmd *cobra.Command, args []string) {
	if len(version) == 0 {
		log.Fatal("must specify the --version flag")
	}
	if len(targets) == 0 {
		log.Fatal("must provide at least one --targets flag when building tools")
	}

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dir = filepath.Join(dir, "release", version)
	vendor := filepath.Join(dir, "src")

	if _, err := os.Stat(vendor); os.IsNotExist(err) {
		log.Fatalf("must first run `kubebuilder-release vendor`.  could not find %s", vendor)
	}

	if !useBazel {
		// Build binaries for the targeted platforms in then tar
		for _, target := range targets {
			// Build binaries for this os:arch
			parts := strings.Split(target, ":")
			if len(parts) != 2 {
				log.Fatalf("--targets flags must be GOOS:GOARCH pairs [%s]", target)
			}
			goos := parts[0]
			goarch := parts[1]
			// Cleanup old binaries
			if !useCached {
				os.RemoveAll(filepath.Join(dir, "bin"))
			}
			os.Mkdir(filepath.Join(dir, "bin"), 0700)

			BuildVendorTar(dir)
			BuildKubernetes(dir, goos, goarch)

			for _, pkg := range VendoredBuildPackages {
				if _, err := os.Stat(filepath.Join(dir, "bin", filepath.Base(pkg))); !useCached || err != nil {
					Build(filepath.Join("cmd", "vendor", pkg, "main.go"),
						filepath.Join(dir, "bin", filepath.Base(pkg)),
						goos, goarch,
					)
				}
			}
			for _, pkg := range OwnedBuildPackages {
				Build(filepath.Join(pkg, "main.go"),
					filepath.Join(dir, "bin", filepath.Base(pkg)),
					goos, goarch,
				)
			}
			PackageTar(goos, goarch, dir, vendor)
		}
	} else {
		os.MkdirAll(filepath.Join(dir, "bin"), 0700)
		BuildVendorTar(dir)
		BazelBuildCopy(dir, []string{
			"//cmd/kubebuilder-gen",
			"//cmd/kubebuilder",
			"//cmd/vendor/github.com/kubernetes-incubator/reference-docs/gen-apidocs",
			"//cmd/vendor/k8s.io/code-generator/cmd/client-gen",
			"//cmd/vendor/k8s.io/code-generator/cmd/conversion-gen",
			"//cmd/vendor/k8s.io/code-generator/cmd/deepcopy-gen",
			"//cmd/vendor/k8s.io/code-generator/cmd/defaulter-gen",
			"//cmd/vendor/k8s.io/code-generator/cmd/informer-gen",
			"//cmd/vendor/k8s.io/code-generator/cmd/lister-gen",
			"//cmd/vendor/k8s.io/code-generator/cmd/openapi-gen",
		}...)
		PackageTar("", "", dir, vendor)
	}
}

func BuildKubernetes(output, goos, goarch string) {
	output = filepath.Join(output, "bin")
	input := filepath.Join("_output", "local", "bin", goos, goarch)
	dir := filepath.Join("cmd", "vendor", "k8s.io", "kubernetes")

	apiserverbin := "kube-apiserver"
	if goos == "windows" {
		apiserverbin += ".exe"
	}

	if _, err := os.Stat(filepath.Join(output, apiserverbin)); !useCached || err != nil {
		fmt.Printf("dir\n%v %v %v\n", filepath.Join(output, apiserverbin), useCached, err)

		cmd := exec.Command("bash", "-c", "echo $PATH; make")
		cmd.Env = append(cmd.Env,
			"WHAT=cmd/kube-apiserver",
			fmt.Sprintf("PATH=/usr/gnu/bin:/usr/local/bin:/bin:/usr/bin:/usr/local/go/bin"))
		cmd.Env = append(cmd.Env, fmt.Sprintf("KUBE_BUILD_PLATFORMS=%s/%s", goos, goarch))
		cmd.Dir = dir
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			log.Fatal("%v", err)
		}

		cmd = exec.Command("cp", filepath.Join(input, apiserverbin), filepath.Join(output, apiserverbin))
		cmd.Dir = dir
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		err = cmd.Run()
		if err != nil {
			log.Fatal("%v", err)
		}
	}

	kubectlbin := "kubectl"
	if goos == "windows" {
		kubectlbin += ".exe"
	}
	if _, err := os.Stat(filepath.Join(output, kubectlbin)); !useCached || err != nil {
		cmd := exec.Command("make")
		cmd.Env = append(cmd.Env, "WHAT=cmd/kubectl",
			fmt.Sprintf("PATH=/usr/gnu/bin:/usr/local/bin:/bin:/usr/bin:/usr/local/go/bin"))
		cmd.Env = append(cmd.Env, fmt.Sprintf("KUBE_BUILD_PLATFORMS=%s/%s", goos, goarch))
		cmd.Dir = dir
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		err = cmd.Run()
		if err != nil {
			log.Fatal("%v", err)
		}

		cmd = exec.Command("cp", filepath.Join(input, kubectlbin), filepath.Join(output, kubectlbin))
		cmd.Dir = dir
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		err = cmd.Run()
		if err != nil {
			log.Fatal("%v", err)
		}
	} else {
	}
}

func Build(input, output, goos, goarch string) {
	var cmd *exec.Cmd
	if strings.HasSuffix(output, "kubebuilder") {
		commit, err := exec.Command("git", "rev-parse", "HEAD").CombinedOutput()
		if err != nil {
			log.Fatalf("%v", err)
		}

		t := time.Now().Local()
		p := "github.com/kubernetes-sigs/kubebuilder/cmd/kubebuilder/version"
		ldflags := []string{
			fmt.Sprintf("-X %s.kubeBuilderVersion=%s", p, version),
			fmt.Sprintf("-X %s.kubernetesVendorVersion=%s", p, kubernetesVersion),
			fmt.Sprintf("-X %s.goos=%s", p, goos),
			fmt.Sprintf("-X %s.goarch=%s", p, goarch),
			fmt.Sprintf("-X %s.gitCommit=%s", p, commit),
			fmt.Sprintf("-X %s.buildDate=%s", p, t.Format("2006-01-02-15:04:05")),
		}
		cmd = exec.Command("go", "build",
			"-ldflags", strings.Join(ldflags, " "),
			"-o", output, input)
	} else {
		cmd = exec.Command("go", "build", "-o", output, input)
	}

	// CGO_ENABLED=0 for statically compile binaries
	cmd.Env = []string{"CGO_ENABLED=0"}
	if len(goos) > 0 {
		cmd.Env = append(cmd.Env, fmt.Sprintf("GOOS=%s", goos))
	}
	if len(goarch) > 0 {
		cmd.Env = append(cmd.Env, fmt.Sprintf("GOARCH=%s", goarch))
	}
	for _, v := range os.Environ() {
		if strings.HasPrefix(v, "CGO_ENABLED=") {
			continue
		}
		if strings.HasPrefix(v, "GOOS=") && len(goos) > 0 {
			continue
		}
		if strings.HasPrefix(v, "GOARCH=") && len(goarch) > 0 {
			continue
		}
		cmd.Env = append(cmd.Env, v)
	}
	RunCmd(cmd, "")
}

func BazelBuildCopy(dest string, targets ...string) {
	args := append([]string{"build"}, targets...)
	c := exec.Command("bazel", args...)
	RunCmd(c, "")

	// Copy the binaries
	for _, t := range targets {
		name := filepath.Base(t)
		c := exec.Command("cp", filepath.Join("bazel-bin", t, name), filepath.Join(dest, "bin", name))
		RunCmd(c, "")
	}
}
