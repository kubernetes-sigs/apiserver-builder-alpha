/*
Copyright YEAR The Kubernetes Authors.

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

package miskatonic

import (
	"context"
	"fmt"
	"net/http"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
)

var _ rest.Connecter = &StudentPencilREST{}

// +k8s:deepcopy-gen=false
type StudentPencilREST struct {
	HTTPHandlerGetter func(pencil *StudentPencil) http.Handler
}

// New creates a new podAttachOptions object.
func (r *StudentPencilREST) New() runtime.Object {
	return &StudentPencil{}
}

// Connect returns a handler for the pod exec proxy
func (r *StudentPencilREST) Connect(ctx context.Context, name string, opts runtime.Object, responder rest.Responder) (http.Handler, error) {
	pencil, ok := opts.(*StudentPencil)
	if !ok {
		return nil, fmt.Errorf("Invalid options object: %#v", opts)
	}
	return r.HTTPHandlerGetter(pencil), nil
}

// NewConnectOptions returns the versioned object that represents exec parameters
func (r *StudentPencilREST) NewConnectOptions() (runtime.Object, bool, string) {
	return &StudentPencil{}, false, ""
}

// ConnectMethods returns the methods supported by exec
func (r *StudentPencilREST) ConnectMethods() []string {
	return []string{"GET", "POST"}
}

// Custom REST storage that delegates to the generated standard Registry
func NewStudentPencilREST(getter generic.RESTOptionsGetter) rest.Storage {
	return &StudentPencilREST{
		HTTPHandlerGetter: func(pencil *StudentPencil) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`all is well`))
			})
		},
	}
}
