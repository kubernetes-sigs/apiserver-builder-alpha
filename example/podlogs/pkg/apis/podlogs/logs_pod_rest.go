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

package podlogs

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/generic"
	genericrest "k8s.io/apiserver/pkg/registry/generic/rest"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/client-go/kubernetes"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	restclient "k8s.io/client-go/rest"
	"sigs.k8s.io/apiserver-builder-alpha/example/podlogs/pkg/kubelet"
	"sigs.k8s.io/apiserver-builder-alpha/pkg/builders"
)

var _ rest.GetterWithOptions = &PodLogsREST{}

func NewPodLogsREST(getter generic.RESTOptionsGetter) rest.Storage {

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

	return &PodLogsREST{
		podClient:   client.CoreV1(),
		KubeletConn: nodeConnGetter,
	}
}

// +k8s:deepcopy-gen=false
type PodLogsREST struct {
	podClient   corev1client.CoreV1Interface
	KubeletConn kubelet.ConnectionInfoGetter
}

func (r *PodLogsREST) NewGetOptions() (runtime.Object, bool, string) {
	builders.ParameterScheme.AddKnownTypes(SchemeGroupVersion, &PodLogs{})
	return &PodLogs{}, false, ""
}

// Connect returns a handler for the pod exec proxy
func (r *PodLogsREST) Get(ctx context.Context, name string, opts runtime.Object) (runtime.Object, error) {
	logOpts, ok := opts.(*PodLogs)
	if !ok {
		return nil, fmt.Errorf("invalid options object: %#v", opts)
	}
	// TODO(developers): apply validations here
	//if errs := validation.ValidatePodLogOptions(logOpts); len(errs) > 0 {
	//	return nil, errors.NewInvalid(api.Kind("PodLogs"), name, errs)
	//}
	location, transport, err := LogLocation(func(name string) (*corev1.Pod, error) {
		ns, _ := request.NamespaceFrom(ctx)
		return r.podClient.Pods(ns).Get(name, metav1.GetOptions{})
	}, r.KubeletConn, ctx, name, logOpts)
	if err != nil {
		return nil, err
	}
	return &genericrest.LocationStreamer{
		Location:    location,
		Transport:   transport,
		ContentType: "text/plain",
		Flush:       logOpts.Follow,
		ResponseChecker: genericrest.NewGenericHttpResponseChecker(
			schema.GroupResource{Group: "podlogs", Resource: "pods/log"},
			name),
		RedirectChecker: genericrest.PreventRedirects,
	}, nil

}

// NewConnectOptions returns the versioned object that represents exec parameters
func (r *PodLogsREST) NewConnectOptions() (runtime.Object, bool, string) {
	return &PodLogs{}, false, ""
}

// ConnectMethods returns the methods supported by exec
func (r *PodLogsREST) ConnectMethods() []string {
	return []string{"GET", "POST"}
}

func (r *PodLogsREST) New() runtime.Object {
	builders.ParameterScheme.AddKnownTypes(SchemeGroupVersion, &PodLogs{})
	return &PodLogs{}
}

type PodGetter func(name string) (*corev1.Pod, error)

// LogLocation returns the log URL for a pod container. If opts.Container is blank
// and only one container is present in the pod, that container is used.
func LogLocation(
	getter PodGetter,
	connInfo kubelet.ConnectionInfoGetter,
	ctx context.Context,
	name string,
	opts *PodLogs,
) (*url.URL, http.RoundTripper, error) {
	pod, err := getter(name)
	if err != nil {
		return nil, nil, err
	}

	// Try to figure out a container
	// If a container was provided, it must be valid
	container := opts.Container
	if len(container) == 0 {
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
		return nil, nil, nil
	}
	nodeInfo, err := connInfo.GetConnectionInfo(ctx, nodeName)
	if err != nil {
		return nil, nil, err
	}
	params := url.Values{}
	if opts.Follow {
		params.Add("follow", "true")
	}
	if opts.Previous {
		params.Add("previous", "true")
	}
	if opts.Timestamps {
		params.Add("timestamps", "true")
	}
	if opts.SinceSeconds != nil {
		params.Add("sinceSeconds", strconv.FormatInt(*opts.SinceSeconds, 10))
	}
	if opts.SinceTime != nil {
		params.Add("sinceTime", opts.SinceTime.Format(time.RFC3339))
	}
	if opts.TailLines != nil {
		params.Add("tailLines", strconv.FormatInt(*opts.TailLines, 10))
	}
	if opts.LimitBytes != nil {
		params.Add("limitBytes", strconv.FormatInt(*opts.LimitBytes, 10))
	}
	loc := &url.URL{
		Scheme:   nodeInfo.Scheme,
		Host:     net.JoinHostPort(nodeInfo.Hostname, nodeInfo.Port),
		Path:     fmt.Sprintf("/containerLogs/%s/%s/%s", pod.Namespace, pod.Name, container),
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
