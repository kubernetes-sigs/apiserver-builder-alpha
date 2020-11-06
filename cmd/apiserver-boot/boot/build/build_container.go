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
	"k8s.io/klog"
	"os"
	"path/filepath"

	"io/ioutil"

	"github.com/spf13/cobra"
	"sigs.k8s.io/apiserver-builder-alpha/cmd/apiserver-boot/boot/util"
)

var Image string

var createBuildContainerCmd = &cobra.Command{
	Use:   "container",
	Short: "Builds a container with the apiserver and controller-manager binaries",
	Long:  `Builds a container with the apiserver and controller-manager binaries`,
	Example: `# Build an image containing the apiserver
# and controller-manager binaries (built for linux:amd64)
apiserver-boot build container --image gcr.io/myrepo/myimage:mytag

# Push the newly built image to the image repo
docker push gcr.io/myrepo/myimage:mytag`,
	Run: RunBuildContainer,
}

func AddBuildContainer(cmd *cobra.Command) {
	cmd.AddCommand(createBuildContainerCmd)
	AddBuildContainerFlags(createBuildContainerCmd)
}

func AddBuildContainerFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Image, "image", "", "name of the image with tag")
	cmd.Flags().StringArrayVar(&BuildTargets, "targets", []string{apiserverTarget, controllerTarget}, "The target binaries to build")
}

func RunBuildContainer(cmd *cobra.Command, args []string) {
	if len(Image) == 0 {
		klog.Fatalf("Must specify --image")
	}

	dir, err := ioutil.TempDir(os.TempDir(), "apiserver-boot-build-container")
	if err != nil {
		klog.Fatalf("failed to create temp directory %s %v", dir, err)
	}
	klog.Infof("Will build docker Image from directory %s", dir)

	klog.Infof("Writing the Dockerfile.")

	path := filepath.Join(dir, "Dockerfile")
	util.WriteIfNotFound(path, "dockerfile-template", dockerfileTemplate, dockerfileTemplateArguments{
		BuildApiserver:  buildApiserver(),
		BuildController: buildController(),
	})

	klog.Infof("Building binaries for linux amd64.")

	// Set the goos and goarch
	goos = "linux"
	goarch = "amd64"
	outputdir = dir
	RunBuildExecutables(cmd, args)

	klog.Infof("Building the docker Image using %s.", path)

	util.DoCmd("docker", "build", "-t", Image, dir)
}

type dockerfileTemplateArguments struct {
	BuildApiserver  bool
	BuildController bool
}

var dockerfileTemplate = `
FROM ubuntu:14.04

RUN apt-get update
RUN apt-get install -y ca-certificates

{{ if .BuildApiserver }}
ADD apiserver .
{{ end }}
{{ if .BuildController }}
ADD controller-manager .
{{ end }}
`
