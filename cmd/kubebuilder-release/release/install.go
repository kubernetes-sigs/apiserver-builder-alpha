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
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install release locally",
	Long:  `install release locally`,
	Run:   RunInstall,
}

func RunInstall(cmd *cobra.Command, args []string) {
	if len(version) == 0 {
		log.Fatal("must specify the --version flag")
	}

	// Untar to to /usr/local/apiserver-build/
	os.Mkdir(filepath.Join("/", "usr", "local", "kubebuilder"), 0700)
	c := exec.Command("tar", "-xzvf", fmt.Sprintf("%s-%s-%s-%s.tar.gz", output, version, "", ""),
		"-C", filepath.Join("/", "usr", "local", "kubebuilder"),
	)
	RunCmd(c, "")
}
