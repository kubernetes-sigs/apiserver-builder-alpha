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
	"github.com/spf13/cobra"
	"log"
)

var targets []string
var output string
var dovendor bool
var test bool
var version string
var kubernetesVersion string
var commit string
var useBazel bool
var useCached bool

var cachevendordir string

var DefaultTargets = []string{"linux:amd64", "darwin:amd64"}

func Run() {
	buildCmd.Flags().StringSliceVar(&targets, "targets",
		DefaultTargets, "GOOS:GOARCH pair.  maybe specified multiple times.")
	buildCmd.Flags().StringVar(&cachevendordir, "vendordir", "",
		"if specified, use this directory for setting up vendor instead of creating a tmp directory.")
	buildCmd.Flags().StringVar(&output, "output", "kubebuilder",
		"value name of the tar file to build")
	buildCmd.Flags().StringVar(&version, "version", "", "version name")
	buildCmd.Flags().BoolVar(&useBazel, "bazel", false, "use bazel to compile (faster, but no X-compile)")
	buildCmd.Flags().BoolVar(&useCached, "cached", false, "use cached binaries")

	buildCmd.Flags().BoolVar(&dovendor, "vendor", true, "if true, fetch packages to vendor")
	buildCmd.Flags().BoolVar(&test, "test", true, "if true, run tests")
	cmd.AddCommand(buildCmd)

	vendorCmd.Flags().StringVar(&commit, "commit", "", "kubebuilder commit")
	vendorCmd.Flags().StringVar(&version, "version", "", "version name")
	vendorCmd.Flags().StringVar(&kubernetesVersion, "kubernetesVersion", "1.9", "version of kubernetes libs")
	vendorCmd.Flags().StringVar(&cachevendordir, "vendordir", "",
		"if specified, use this directory for setting up vendor instead of creating a tmp directory.")
	cmd.AddCommand(vendorCmd)

	installCmd.Flags().StringVar(&version, "version", "", "version name")
	cmd.AddCommand(installCmd)

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

var cmd = &cobra.Command{
	Use:   "kubebuilder-release",
	Short: "kubebuilder-release builds a .tar.gz release package",
	Long:  `kubebuilder-release builds a .tar.gz release package`,
	Run:   RunMain,
}

func RunMain(cmd *cobra.Command, args []string) {
	cmd.Help()
}
