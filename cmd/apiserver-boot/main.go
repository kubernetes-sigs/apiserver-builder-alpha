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

package main

import (
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
	"sigs.k8s.io/apiserver-builder-alpha/v2/pkg/boot/build"
	"sigs.k8s.io/apiserver-builder-alpha/v2/pkg/boot/create"
	"sigs.k8s.io/apiserver-builder-alpha/v2/pkg/boot/init_repo"
	"sigs.k8s.io/apiserver-builder-alpha/v2/pkg/boot/run"
	"sigs.k8s.io/apiserver-builder-alpha/v2/pkg/boot/show"
	"sigs.k8s.io/apiserver-builder-alpha/v2/pkg/boot/version"
)

func main() {

	init_repo.AddInit(cmd)
	create.AddCreate(cmd)
	build.AddBuild(cmd)
	run.AddRun(cmd)
	version.AddVersion(cmd)
	show.AddShow(cmd)

	if err := cmd.Execute(); err != nil {
		klog.Fatal(err)
	}
}

var cmd = &cobra.Command{
	Use:   "apiserver-boot",
	Short: "apiserver-boot development kit for building Kubernetes extensions in go.",
	Long:  `apiserver-boot development kit for building Kubernetes extensions in go.`,
	Example: `# Initialize your repository with scaffolding directories and go files. Specify --module-name if the project works outside GOPATH.
apiserver-boot init repo --domain example.com

# Create new resource "Bee" in the "insect" group with version "v1beta1"
apiserver-boot create group version resource --group insect --version v1beta1 --kind Bee

# Build the generated code, apiserver and controller-manager so they be run locally.
apiserver-boot build executables

# Run the tests that were created for your resources
# Requires generated code was already built by "build executables" or "build generated"
go test ./pkg/...

# Run locally by starting a local etcd, apiserver and controller-manager
# Produces a kubeconfig to talk to the local server
apiserver-boot run local

# Check the api versions of the locally running server
kubectl --kubeconfig kubeconfig api-versions

# Build an image and run in a cluster in the default namespace
# Note: after running this you should clear the discovery service
# cache before running kubectl with "rm -rf ~/.kube/cache/discovery/"
apiserver-boot run in-cluster --name creatures --namespace default --image repo/name:tag`,
	Run: RunMain,
}

func RunMain(cmd *cobra.Command, args []string) {
	cmd.Help()
}
