/*
Copyright 2014 The Kubernetes Authors.

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

package generators

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	ccmapp "k8s.io/kubernetes/cmd/cloud-controller-manager/app"
	"k8s.io/kubernetes/cmd/genutils"
	apiservapp "k8s.io/kubernetes/cmd/kube-apiserver/app"
	cmapp "k8s.io/kubernetes/cmd/kube-controller-manager/app"
	proxyapp "k8s.io/kubernetes/cmd/kube-proxy/app"
	schapp "k8s.io/kubernetes/cmd/kube-scheduler/app"
	kubeadmapp "k8s.io/kubernetes/cmd/kubeadm/app/cmd"
	kubeletapp "k8s.io/kubernetes/cmd/kubelet/app"
	kubectlcmd "k8s.io/kubernetes/pkg/kubectl/cmd"
)

func GenerateFiles(path, module string) {

	outDir, err := genutils.OutDir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get output directory: %v\n", err)
		os.Exit(1)
	}

	stop := make(chan struct{})
	switch module {
	case "kube-apiserver":
		apiserver := apiservapp.NewAPIServerCommand(stop)
		GenMarkdownTree(apiserver, outDir, true)

	case "kube-controller-manager":
		controllermanager := cmapp.NewControllerManagerCommand()
		GenMarkdownTree(controllermanager, outDir, true)

	case "cloud-controller-manager":
		cloudcontrollermanager := ccmapp.NewCloudControllerManagerCommand()
		GenMarkdownTree(cloudcontrollermanager, outDir, true)

	case "kube-proxy":
		proxy := proxyapp.NewProxyCommand()
		GenMarkdownTree(proxy, outDir, true)

	case "kube-scheduler":
		scheduler := schapp.NewSchedulerCommand()
		GenMarkdownTree(scheduler, outDir, true)

	case "kubelet":
		kubelet := kubeletapp.NewKubeletCommand(stop)
		GenMarkdownTree(kubelet, outDir, true)

	case "kubeadm":
		// resets global flags created by kubelet or other commands e.g.
		// --azure-container-registry-config from pkg/credentialprovider/azure
		// --google-json-key from pkg/credentialprovider/gcp
		// --version pkg/version/verflag
		pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

		// generate docs for kubeadm
		kubeadm := kubeadmapp.NewKubeadmCommand(os.Stdin, os.Stdout, os.Stderr)
		GenMarkdownTree(kubeadm, outDir, false)

		// cleanup generated code for usage as include in the website
		MarkdownPostProcessing(kubeadm, outDir, cleanupForInclude)

	case "kubectl":
		kubectl := kubectlcmd.NewDefaultKubectlCommand()
		GenMarkdownTree(kubectl, outDir, true)

	default:
		fmt.Fprintf(os.Stderr, "Module %s is not supported", module)
		os.Exit(1)
	}
}
