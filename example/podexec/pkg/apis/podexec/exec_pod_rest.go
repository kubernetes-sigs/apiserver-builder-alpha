/*
Copyright 2019 The Kubernetes Authors.

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

package podexec

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/proxy"
	"k8s.io/apiserver/pkg/endpoints/request"
	genericfeatures "k8s.io/apiserver/pkg/features"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/client-go/kubernetes"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	restclient "k8s.io/client-go/rest"
	"sigs.k8s.io/apiserver-builder-alpha/example/podexec/pkg/kubelet"
)

var _ rest.Connecter = &PodExecREST{}

func NewPodExecREST(getter generic.RESTOptionsGetter) rest.Storage {

	inClusterClientConfig, err := restclient.InClusterConfig()
	if err != nil {
		panic(err)
	}

	client := kubernetes.NewForConfigOrDie(inClusterClientConfig)
	nodeConnGetter, err := kubelet.NewNodeConnectionInfoGetter(
		func(name string) (*corev1.Node, error) {
			return client.CoreV1().Nodes().Get(name, metav1.GetOptions{})
		},
		kubelet.KubeletClientConfig{
			Port:         10250,
			ReadOnlyPort: 10255,
			PreferredAddressTypes: []string{
				// --override-hostname
				string(corev1.NodeHostName),

				// internal, preferring DNS if reported
				string(corev1.NodeInternalDNS),
				string(corev1.NodeInternalIP),

				// external, preferring DNS if reported
				string(corev1.NodeExternalDNS),
				string(corev1.NodeExternalIP),
			},
			EnableHTTPS: true,
			HTTPTimeout: time.Duration(5) * time.Second,
		})

	return &PodExecREST{
		podClient:   client.CoreV1(),
		KubeletConn: nodeConnGetter,
	}
}

// +k8s:deepcopy-gen=false
type PodExecREST struct {
	podClient   corev1client.CoreV1Interface
	KubeletConn kubelet.ConnectionInfoGetter
}

func (r *PodExecREST) Connect(ctx context.Context, name string, opts runtime.Object, responder rest.Responder) (http.Handler, error) {
	execOpts, ok := opts.(*PodExec)
	if !ok {
		return nil, fmt.Errorf("invalid options object: %#v", opts)
	}
	location, transport, err := ExecLocation(r.podClient, r.KubeletConn, ctx, name, execOpts)
	if err != nil {
		return nil, err
	}
	return newThrottledUpgradeAwareProxyHandler(location, transport, false, true, true, responder), nil
}

func (r *PodExecREST) NewConnectOptions() (runtime.Object, bool, string) {
	return &PodExec{}, false, ""
}

func (r *PodExecREST) ConnectMethods() []string {
	return []string{"GET", "POST"}
}

func (r *PodExecREST) New() runtime.Object {
	return &PodExec{}
}

func newThrottledUpgradeAwareProxyHandler(location *url.URL, transport http.RoundTripper, wrapTransport, upgradeRequired, interceptRedirects bool, responder rest.Responder) *proxy.UpgradeAwareHandler {
	handler := proxy.NewUpgradeAwareHandler(location, transport, wrapTransport, upgradeRequired, proxy.NewErrorResponder(responder))
	handler.InterceptRedirects = interceptRedirects && utilfeature.DefaultFeatureGate.Enabled(genericfeatures.StreamingProxyRedirects)
	handler.RequireSameHostRedirects = utilfeature.DefaultFeatureGate.Enabled(genericfeatures.ValidateProxyRedirects)
	handler.MaxBytesPerSec = 0
	return handler
}

// ExecLocation returns the exec URL for a pod container. If opts.Container is blank
// and only one container is present in the pod, that container is used.
func ExecLocation(
	getter corev1client.CoreV1Interface,
	connInfo kubelet.ConnectionInfoGetter,
	ctx context.Context,
	name string,
	opts *PodExec,
) (*url.URL, http.RoundTripper, error) {
	return streamLocation(getter, connInfo, ctx, name, opts, opts.Container, "exec")
}

func streamLocation(
	getter corev1client.CoreV1Interface,
	connInfo kubelet.ConnectionInfoGetter,
	ctx context.Context,
	name string,
	opts runtime.Object,
	container,
	path string,
) (*url.URL, http.RoundTripper, error) {
	ns, _ := request.NamespaceFrom(ctx)
	pod, err := getter.Pods(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, nil, err
	}

	// Try to figure out a container
	// If a container was provided, it must be valid
	if container == "" {
		switch len(pod.Spec.Containers) {
		case 1:
			container = pod.Spec.Containers[0].Name
		case 0:
			return nil, nil, errors.NewBadRequest(fmt.Sprintf("a container name must be specified for pod %s", name))
		default:
			containerNames := getContainerNames(pod.Spec.Containers)
			initContainerNames := getContainerNames(pod.Spec.InitContainers)
			err := fmt.Sprintf("a container name must be specified for pod %s, choose one of: [%s]", name, containerNames)
			if len(initContainerNames) > 0 {
				err += fmt.Sprintf(" or one of the init containers: [%s]", initContainerNames)
			}
			return nil, nil, errors.NewBadRequest(err)
		}
	} else {
		if !podHasContainerWithName(pod, container) {
			return nil, nil, errors.NewBadRequest(fmt.Sprintf("container %s is not valid for pod %s", container, name))
		}
	}
	nodeName := types.NodeName(pod.Spec.NodeName)
	if len(nodeName) == 0 {
		// If pod has not been assigned a host, return an empty location
		return nil, nil, errors.NewBadRequest(fmt.Sprintf("pod %s does not have a host assigned", name))
	}
	nodeInfo, err := connInfo.GetConnectionInfo(ctx, nodeName)
	if err != nil {
		return nil, nil, err
	}
	params := url.Values{}
	if err := streamParams(params, opts); err != nil {
		return nil, nil, err
	}
	loc := &url.URL{
		Scheme:   nodeInfo.Scheme,
		Host:     net.JoinHostPort(nodeInfo.Hostname, nodeInfo.Port),
		Path:     fmt.Sprintf("/%s/%s/%s/%s", path, pod.Namespace, pod.Name, container),
		RawQuery: params.Encode(),
	}
	return loc, nodeInfo.Transport, nil
}

// getContainerNames returns a formatted string containing the container names
func getContainerNames(containers []corev1.Container) string {
	names := []string{}
	for _, c := range containers {
		names = append(names, c.Name)
	}
	return strings.Join(names, " ")
}

func podHasContainerWithName(pod *corev1.Pod, containerName string) bool {
	var hasContainer bool = false
	for _, c := range pod.Spec.Containers {
		if c.Name == containerName {
			hasContainer = true
		}
	}
	return hasContainer
}

func streamParams(params url.Values, opts runtime.Object) error {
	switch opts := opts.(type) {
	case *PodExec:
		if opts.Stdin {
			params.Add(corev1.ExecStdinParam, "1")
		}
		if opts.Stdout {
			params.Add(corev1.ExecStdoutParam, "1")
		}
		if opts.Stderr {
			params.Add(corev1.ExecStderrParam, "1")
		}
		if opts.TTY {
			params.Add(corev1.ExecTTYParam, "1")
		}
		for _, c := range opts.Command {
			params.Add("command", c)
		}
	default:
		return fmt.Errorf("Unknown object for streaming: %v", opts)
	}
	return nil
}
