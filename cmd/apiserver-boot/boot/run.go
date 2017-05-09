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
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"os/exec"
	"path/filepath"
	"strings"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run an apiserver",
	Long:  `run an apiserver`,
	Run:   RunRun,
}

var etcd string
var config string
var printapiserver bool
var printetcd bool

func AddRunCmd(cmd *cobra.Command) {
	runCmd.Flags().StringVar(&server, "server", "", "path to apiserver binary to run")
	runCmd.Flags().StringVar(&etcd, "etcd", "", "if non-empty, use this etcd instead of starting a new one")
	runCmd.Flags().StringVar(&config, "config", "kubeconfig", "path to the kubeconfig to write for using kubectl")
	runCmd.Flags().BoolVar(&printapiserver, "printapiserver", false, "if true, pipe the apiserver stdout and stderr")
	runCmd.Flags().BoolVar(&printetcd, "printetcd", false, "if true, pipe the etcd stdout and stderr")
	cmd.AddCommand(runCmd)
}

func RunRun(cmd *cobra.Command, args []string) {
	if len(server) == 0 {
		fmt.Fprintf(os.Stderr, "apiserver-boot run requires the --server flag\n")
		os.Exit(-1)
	}

	WriteKubeConfig()

	// Start etcd
	if len(etcd) == 0 {
		etcd = "http://localhost:2379"
		etcdCmd := RunEtcd()
		defer etcdCmd.Process.Kill()
	}

	// Start apiserver
	RunApiserver()
}

func RunEtcd() *exec.Cmd {
	etcdCmd := exec.Command("etcd")
	if printetcd {
		etcdCmd.Stderr = os.Stderr
		etcdCmd.Stdout = os.Stdout
	}

	fmt.Printf("%s\n", strings.Join(etcdCmd.Args, " "))
	go func() {
		err := etcdCmd.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to run etcd %v\n", err)
			os.Exit(-1)
		}
	}()
	return etcdCmd
}

func RunApiserver() *exec.Cmd {
	apiserverCmd := exec.Command(server,
		"--delegated-auth=false",
		fmt.Sprintf("--etcd-servers=%s", etcd),
		"--secure-port=9443",
		"--print-bearer-token",
	)
	fmt.Printf("%s\n", strings.Join(apiserverCmd.Args, " "))
	if printapiserver {
		apiserverCmd.Stderr = os.Stderr
		apiserverCmd.Stdout = os.Stdout
	}

	fmt.Printf("to test the server run `kubectl --kubeconfig %s version`\n", config)

	err := apiserverCmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run apiserver %v\n", err)
		os.Exit(-1)
	}

	return apiserverCmd
}

func WriteKubeConfig() {
	// Write a kubeconfig
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot get working directory %v\n", err)
		os.Exit(-1)
	}
	path := filepath.Join(dir, "apiserver.local.config", "certificates", "apiserver")
	writeIfNotFound(config, "kubeconfig-template", configTemplate, path)
}

var configTemplate = `
apiVersion: v1
clusters:
- cluster:
    certificate-authority: {{ . }}.crt
    server: https://localhost:9443
  name: apiserver
contexts:
- context:
    cluster: apiserver
    user: apiserver
  name: apiserver
current-context: apiserver
kind: Config
preferences: {}
users:
- name: apiserver
  user:
    client-certificate: {{ . }}.crt
    client-key: {{ . }}.key
`
