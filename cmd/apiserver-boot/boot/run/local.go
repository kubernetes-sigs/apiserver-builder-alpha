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

package run

import (
	"bytes"
	"context"
	"fmt"
	"k8s.io/klog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/apiserver-builder-alpha/cmd/apiserver-boot/boot/build"
	"sigs.k8s.io/apiserver-builder-alpha/cmd/apiserver-boot/boot/util"
)

var localCmd = &cobra.Command{
	Use:   "local",
	Short: "run the etcd, apiserver and controller",
	Long: `run the etcd, apiserver and controller, Note that the aggregated apiserver in the local mode 
will not be attempting to delegate any requests to an acutal kube-apiserver, hence neither authentication 
nor authorization will be performed.`,
	Example: `
# Regenerate code and build binaries then run them. 

apiserver-boot run local

# Check the api versions of the locally running server
kubectl --kubeconfig kubeconfig api-versions

# Run locally without rebuilding
apiserver-boot run local --build=false

# Create an instance and fetch it
nano -w samples/<type>.yaml
kubectl --kubeconfig kubeconfig apply -f samples/<type>.yaml
kubectl --kubeconfig kubeconfig get <type>`,
	Run: RunLocal,
}

var etcd string
var config string
var printapiserver bool
var printcontrollermanager bool
var printetcd bool
var buildBin bool

var server string
var controllermanager string
var toRun []string
var disableMTLS bool
var certDir string
var securePort int32

func AddLocal(cmd *cobra.Command) {
	localCmd.Flags().StringSliceVar(&toRun, "run", []string{"etcd", "apiserver", "controller"}, "path to apiserver binary to run")
	localCmd.Flags().BoolVar(&disableMTLS, "disable-mtls", true,
		`If true, disable mTLS serving in the apiserver with --standalone-debug-mode. `+
			`This optional requries the apiserver to build with "WithLocalDebugExtension" from apiserver-runtime.`)

	localCmd.Flags().StringVar(&server, "apiserver", "", "path to apiserver binary to run")
	localCmd.Flags().StringVar(&controllermanager, "controller-manager", "", "path to controller-manager binary to run")
	localCmd.Flags().StringVar(&etcd, "etcd", "", "if non-empty, use this etcd instead of starting a new one")

	localCmd.Flags().StringVar(&config, "config", "kubeconfig", "path to the kubeconfig to write for using kubectl")

	localCmd.Flags().BoolVar(&printapiserver, "print-apiserver", true, "if true, pipe the apiserver stdout and stderr")
	localCmd.Flags().BoolVar(&printcontrollermanager, "print-controller-manager", true, "if true, pipe the controller-manager stdout and stderr")
	localCmd.Flags().BoolVar(&printetcd, "printetcd", false, "if true, pipe the etcd stdout and stderr")
	localCmd.Flags().BoolVar(&buildBin, "build", true, "if true, build the binaries before running")

	localCmd.Flags().Int32Var(&securePort, "secure-port", 9443, "Secure port from apiserver to serve requests")
	localCmd.Flags().StringVar(&certDir, "cert-dir", filepath.Join("config", "certificates"), "directory containing apiserver certificates")

	cmd.AddCommand(localCmd)
}

func RunLocal(cmd *cobra.Command, args []string) {
	if buildBin {
		build.BuildTargets = toRun
		build.RunBuildExecutables(cmd, args)
	}

	WriteKubeConfig()

	// parent context to indicate whether cmds quit
	ctx, cancel := context.WithCancel(context.Background())
	ctx = util.CancelWhenSignaled(ctx)

	r := map[string]interface{}{}
	for _, s := range toRun {
		r[s] = nil
	}

	startedCommands := map[string]*exec.Cmd{}
	defer func() {
		klog.Info("Cleaning up processes")
		for _, cmd := range startedCommands {
			WaitUntilCommandCompleted(cmd)
		}
	}()
	// Start etcd
	if _, f := r["etcd"]; f {
		etcd = "http://localhost:2379"
		startedCommands["etcd"] = RunEtcd(ctx, cancel)
		time.Sleep(time.Second * 2)
	}

	// Start apiserver
	if _, f := r["apiserver"]; f {
		startedCommands["apiserver"] = RunApiserver(ctx, cancel)
		time.Sleep(time.Second * 2)
		klog.Info("Aggregated apiserver successfully started")
	}

	// Start controller manager
	if _, f := r["controller"]; f {
		startedCommands["controller"] = RunControllerManager(ctx, cancel)
		klog.Info("Controller manager successfully started")
	}

	klog.Infof(`
==================================================
| Now you're all set!
| To test the server, try the following commands:
|
| >> "kubectl --kubeconfig %s api-versions" or "KUBECONFIG=%s kubectl api-resources"
|
==================================================`,
		config, config)
	<-ctx.Done() // wait forever
}

func RunEtcd(ctx context.Context, cancel context.CancelFunc) *exec.Cmd {
	etcdCmd := exec.Command("etcd")
	if printetcd {
		etcdCmd.Stderr = os.Stderr
		etcdCmd.Stdout = os.Stdout
	}

	go runCommon(etcdCmd, ctx, cancel)

	return etcdCmd
}

func RunApiserver(ctx context.Context, cancel context.CancelFunc) *exec.Cmd {
	if len(server) == 0 {
		server = "bin/apiserver"
	}

	// checking if apiserver supports local running
	apiserverTestLocalCmd := exec.Command(server, "-h")
	buf := &bytes.Buffer{}
	apiserverTestLocalCmd.Stderr = buf
	runCommon(apiserverTestLocalCmd, ctx, cancel)
	if !strings.Contains(string(buf.Bytes()), "--standalone-debug-mode") {
		klog.Fatalf(`
The apiserver binary doesn't seem to support --standalone-debug-mode, 
did you have WithLocalDebugExtension() in your apiserver? (if you're using kuberentes-sigs/apiserver-runtime')`)
	}
	klog.Info("The apiserver binary seems to support local-running, proceeding..")

	// starting apiserver process
	flags := []string{
		fmt.Sprintf("--etcd-servers=%s", etcd),
		fmt.Sprintf("--secure-port=%v", securePort),
		fmt.Sprintf("--feature-gates=APIPriorityAndFairness=false"), // TODO: remove this line after https://github.com/kubernetes/kubernetes/pull/97957 merged
	}

	if disableMTLS {
		flags = append(flags, "--standalone-debug-mode")
		flags = append(flags, "--bind-address=127.0.0.1")
	} else {
		flags = append(flags,
			fmt.Sprintf("--cert-dir=%s", certDir),
			fmt.Sprintf("--client-ca-file=%s/apiserver_ca.crt", certDir),
		)
	}

	apiserverCmd := exec.Command(server,
		flags...,
	)
	if printapiserver {
		apiserverCmd.Stderr = os.Stderr
		apiserverCmd.Stdout = os.Stdout
	}

	go runCommon(apiserverCmd, ctx, cancel)

	return apiserverCmd
}

func RunControllerManager(ctx context.Context, cancel context.CancelFunc) *exec.Cmd {
	if len(controllermanager) == 0 {
		controllermanager = "bin/controller-manager"
	}

	controllerManagerCmd := exec.Command(controllermanager,
		fmt.Sprintf("--kubeconfig=%s", config),
	)
	if printcontrollermanager {
		controllerManagerCmd.Stderr = os.Stderr
		controllerManagerCmd.Stdout = os.Stdout
	}

	go runCommon(controllerManagerCmd, ctx, cancel)

	return controllerManagerCmd
}

// run a command via goroutine
func runCommon(cmd *exec.Cmd, ctx context.Context, cancel context.CancelFunc) {
	stopCh := make(chan error)
	cmdName := cmd.Args[0]

	klog.Infof("Starting local component: %s", strings.Join(cmd.Args, " "))
	go func() {
		err := cmd.Run()
		if err != nil {
			klog.Infof("Failed to run %s, error: %v", cmdName, err)
		} else {
			klog.Infof("Command %s quitted normally", cmdName)
		}
		stopCh <- err
	}()

	select {
	case <-stopCh:
		// my command quited
		cancel()
	case <-ctx.Done():
		// other commands quited
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	}
}

func WriteKubeConfig() {
	klog.Infof("Writing kubeconfig to %s", config)
	// Write a kubeconfig
	dir, err := os.Getwd()
	if err != nil {
		klog.Fatalf("Cannot get working directory %v", err)
		os.Exit(-1)
	}
	path := filepath.Join(dir, certDir)
	util.WriteIfNotFound(config, "kubeconfig-template", configTemplate,
		ConfigArgs{
			DisabltMTLS: disableMTLS,
			Path:        path,
			Port:        fmt.Sprintf("%v", securePort),
		})
}

func WaitUntilCommandCompleted(cmd *exec.Cmd) {
	cmdName := cmd.Args[0]
	wait.PollImmediate(time.Millisecond*100, time.Second, func() (bool, error) {
		if cmd.ProcessState != nil {
			klog.Infof("Waiting for process of %s (pid=%v) to be completed", cmdName, cmd.ProcessState.Pid())
			return cmd.ProcessState.Exited(), nil
		}
		return true, nil
	})
	klog.Infof("Completed %s", cmdName)
}

type ConfigArgs struct {
	DisabltMTLS bool
	Path        string
	Port        string
}

var configTemplate = `
apiVersion: v1
clusters:
- cluster:
{{- if .DisabltMTLS }}
    insecure-skip-tls-verify: true
{{- else }}
    certificate-authority: {{ .Path }}/apiserver_ca.crt
{{- end }}
    server: https://localhost:{{ .Port }}
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
{{- if not .DisabltMTLS }}
    client-certificate: {{ .Path }}/apiserver.crt
    client-key: {{ .Path }}/apiserver.key
{{- else }}
    username: apiserver
{{- end }}
`
