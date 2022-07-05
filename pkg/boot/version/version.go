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

package version

import (
	"fmt"
	"k8s.io/klog"
	"log"
	"runtime"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

var (
	apiserverBuilderVersion = "unknown"
	kubernetesVendorVersion = "unknown"
	goos                    = runtime.GOOS
	goarch                  = runtime.GOARCH
	gitCommit               = "$Format:%H$" // sha1 from git, output of $(git rev-parse HEAD)

	buildDate = "1970-01-01T00:00:00Z" // build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
)

type Version struct {
	ApiserverBuilderVersion string `json:"apiserverBuilderVersion"`
	KubernetesVendor        string `json:"kubernetesVendor"`
	GitCommit               string `json:"gitCommit"`
	BuildDate               string `json:"buildDate"`
	GoOs                    string `json:"goOs"`
	GoArch                  string `json:"goArch"`
}

// GetVersion returns the version
func GetVersion() string {
	return fmt.Sprintf("Version: %#v", Version{
		apiserverBuilderVersion,
		kubernetesVendorVersion,
		gitCommit,
		buildDate,
		goos,
		goarch,
	})
}

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Print the apisever-builder version.",
	Long:    `Print the apisever-builder version.`,
	Example: `apiserver-boot version`,
	Run:     RunVersion,
}

func AddVersion(cmd *cobra.Command) {
	cmd.AddCommand(versionCmd)
}

func RunVersion(cmd *cobra.Command, args []string) {
	version := GetVersion()
	log.Printf(version)
}
