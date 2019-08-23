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
	"sigs.k8s.io/apiserver-builder-alpha/pkg/builders"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
)

// LogREST implements GetterWithOptions
var _ = rest.GetterWithOptions(&StudentBookREST{})

// +k8s:deepcopy-gen=false
type StudentBookREST struct {
	StreamingHandler func(book *StudentBook) (streamer runtime.Object)
}

// New creates a new Pod log options object
func (r *StudentBookREST) New() runtime.Object {
	// TODO - return a resource that represents a log
	return &StudentBook{}
}

// LogREST implements StorageMetadata
func (r *StudentBookREST) ProducesMIMETypes(verb string) []string {
	// Since the default list does not include "plain/text", we need to
	// explicitly override ProducesMIMETypes, so that it gets added to
	// the "produces" section for pods/{name}/log
	return []string{
		"text/plain",
	}
}

// LogREST implements StorageMetadata, return string as the generating object
func (r *StudentBookREST) ProducesObject(verb string) interface{} {
	return ""
}

func (r *StudentBookREST) Get(ctx context.Context, name string, opts runtime.Object) (runtime.Object, error) {
	return r.StreamingHandler(opts.(*StudentBook)), nil
}

// NewGetOptions creates a new options object
func (r *StudentBookREST) NewGetOptions() (runtime.Object, bool, string) {
	return &StudentBook{}, false, ""
}

// OverrideMetricsVerb override the GET verb to CONNECT for pod log resource
func (r *StudentBookREST) OverrideMetricsVerb(oldVerb string) (newVerb string) {
	newVerb = oldVerb

	if oldVerb == "GET" {
		newVerb = "CONNECT"
	}

	return
}

// Custom REST storage that delegates to the generated standard Registry
func NewStudentBookREST(getter generic.RESTOptionsGetter) rest.Storage {
	builders.ParameterScheme.AddKnownTypes(SchemeGroupVersion, &StudentBook{})
	return &StudentBookREST{
		StreamingHandler: func(book *StudentBook) runtime.Object {
			return &Student{
				Spec: StudentSpec{
					ID: 3,
				},
			}
		},
	}
}
