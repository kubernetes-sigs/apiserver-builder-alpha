
/*
Copyright 2022.

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

package v1

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/rest"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource/resourcerest"
	contextutil "sigs.k8s.io/apiserver-runtime/pkg/util/context"
)

var _ resource.SubResource = &SnapshotBar{}
var _ rest.Storage = &SnapshotBar{}
var _ resourcerest.Connecter = &SnapshotBar{}

var barProxyMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

// SnapshotBar
type SnapshotBar struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SnapshotBarOptions struct {
	metav1.TypeMeta

	// Path is the target api path of the proxy request.
	Path string `json:"path"`
}

func (c *SnapshotBar) SubResourceName() string {
	return "proxy"
}

func (c *SnapshotBar) New() runtime.Object {
	return &SnapshotBarOptions{}
}

func (c *SnapshotBar) Connect(ctx context.Context, id string, options runtime.Object, r rest.Responder) (http.Handler, error) {
	// EDIT IT
	parentStorage, ok := contextutil.GetParentStorage(ctx)
	if !ok {
		return nil, fmt.Errorf("no parent storage found")
	}
	_, err := parentStorage.Get(ctx, id, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return http.Handler(nil), nil
}

func (c *SnapshotBar) NewConnectOptions() (runtime.Object, bool, string) {
	return &SnapshotBarOptions{}, false, "path"
}

func (c *SnapshotBar) ConnectMethods() []string {
	return barProxyMethods
}

var _ resource.QueryParameterObject = &SnapshotBarOptions{}

func (in *SnapshotBarOptions) ConvertFromUrlValues(values *url.Values) error {
	in.Path = values.Get("path")
	return nil
}
