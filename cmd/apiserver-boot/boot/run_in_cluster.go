/*
Copyright 2016 The Kubernetes Authors.

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
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var runInClusterCmd = &cobra.Command{
	Use:   "run-in-cluster",
	Short: "run the etcd, apiserver and the controller-manager as an aggegrated apiserver in a cluster",
	Long:  `run the etcd, apiserver and the controller-manager as an aggegrated apiserver in a cluster`,
	Run:   RunRunInCluster,
}

var buildImage bool

func AddRunInCluster(cmd *cobra.Command) {
	runInClusterCmd.Flags().StringVar(&name, "name", "", "")
	runInClusterCmd.Flags().StringVar(&namespace, "namespace", "", "")
	runInClusterCmd.Flags().StringVar(&image, "image", "", "name of the image to use")
	runInClusterCmd.Flags().StringVar(&resourceConfigDir, "output", "config", "directory to output resourceconfig")

	runInClusterCmd.Flags().BoolVar(&buildImage, "build-image", true, "if true, build and push the image")
	runInClusterCmd.Flags().BoolVar(&generateForBuild, "generate", true, "if true, generate code before building image")

	cmd.AddCommand(runInClusterCmd)
}

func RunRunInCluster(cmd *cobra.Command, args []string) {
	if buildImage {
		// Build the container first
		RunBuildContainer(cmd, args)

		// Push the image
		doCmd("docker", "push", image)
	}

	// Build the resource config
	os.Remove(filepath.Join(resourceConfigDir, "apiserver.yaml"))
	RunBuildResourceConfig(cmd, args)

	// Apply the new config
	doCmd("kubectl", "apply", "-f", filepath.Join(resourceConfigDir, "apiserver.yaml"))
}
