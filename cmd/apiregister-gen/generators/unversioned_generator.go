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

package generators

import (
	"io"
	"text/template"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/gengo/generator"
)

type unversionedGenerator struct {
	generator.DefaultGen
	apigroup *APIGroup
}

var _ generator.Generator = &unversionedGenerator{}

func CreateUnversionedGenerator(apigroup *APIGroup, filename string) generator.Generator {
	return &unversionedGenerator{
		generator.DefaultGen{OptionalName: filename},
		apigroup,
	}
}

func (d *unversionedGenerator) Imports(c *generator.Context) []string {
	imports := sets.NewString(
		"fmt",
		"github.com/kubernetes-incubator/apiserver-builder/pkg/builders",
		"k8s.io/apimachinery/pkg/apis/meta/internalversion",
		"k8s.io/apimachinery/pkg/runtime",
		"k8s.io/apimachinery/pkg/runtime/schema",
		"k8s.io/apiserver/pkg/endpoints/request",
		"k8s.io/apiserver/pkg/registry/rest",
		"k8s.io/client-go/pkg/api")

	// Get imports for all fields
	for _, s := range d.apigroup.Structs {
		for _, f := range s.Fields {
			if len(f.UnversionedImport) > 0 {
				imports.Insert(f.UnversionedImport)
			}
		}
	}

	return imports.List()
}

func (d *unversionedGenerator) Finalize(context *generator.Context, w io.Writer) error {
	temp := template.Must(template.New("unversioned-wiring-template").Parse(UnversionedAPITemplate))
	err := temp.Execute(w, d.apigroup)
	if err != nil {
		return err
	}
	return err
}

var UnversionedAPITemplate = `
var (
	{{ range $api := .UnversionedResources -}}
	Internal{{ $api.Kind }} = builders.NewInternalResource(
		"{{ $api.Resource }}",
		func() runtime.Object { return &{{ $api.Kind }}{} },
		func() runtime.Object { return &{{ $api.Kind }}List{} },
	)
	Internal{{ $api.Kind }}Status = builders.NewInternalResourceStatus(
		"{{ $api.Resource }}",
		func() runtime.Object { return &{{ $api.Kind }}{} },
		func() runtime.Object { return &{{ $api.Kind }}List{} },
	)
	{{ range $subresource := .Subresources -}}
	Internal{{$subresource.REST}} = builders.NewInternalSubresource(
		"{{$subresource.Resource}}", "{{$subresource.Path}}",
		func() runtime.Object { return &{{$subresource.Request}}{} },
	)
	{{ end -}}
	{{ end -}}

	// Registered resources and subresources
	ApiVersion = builders.NewApiGroup("{{.Group}}.{{.Domain}}").WithKinds(
		{{ range $api := .UnversionedResources -}}
		Internal{{$api.Kind}},
		Internal{{$api.Kind}}Status,
		{{ range $subresource := $api.Subresources -}}
		Internal{{$subresource.REST}},
		{{ end -}}
		{{ end -}}
	)

	// Required by code generated by go2idl
	AddToScheme = ApiVersion.SchemaBuilder.AddToScheme
	SchemeBuilder = ApiVersion.SchemaBuilder
	SchemeGroupVersion = ApiVersion.GroupVersion
)

// Required by code generated by go2idl
// Kind takes an unqualified kind and returns a Group qualified GroupKind
func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

// Required by code generated by go2idl
// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

{{ range $s := .Structs -}}
// +genclient=true

type {{ $s.Name }} struct {
{{ range $f := $s.Fields -}}
    {{ $f.Name }} {{ $f.UnversionedType }}
{{ end -}}
}

{{ end -}}

{{ range $api := .UnversionedResources -}}
//
// {{.Kind}} Functions and Structs
//
type {{.Kind}}Strategy struct {
	builders.DefaultStorageStrategy
}

{{ if .NonNamespaced -}}
var {{.Kind}}StrategySingleton = {{.Kind}}Strategy{builders.StorageStrategySingleton}

func ({{.Kind}}Strategy) NamespaceScoped() bool { return false }
{{ end -}}

type {{$api.Kind}}StatusStrategy struct {
	builders.DefaultStatusStorageStrategy
}

{{ if .NonNamespaced -}}
var {{.Kind}}StatusStrategySingleton = {{.Kind}}StatusStrategy{builders.StatusStorageStrategySingleton}

func ({{.Kind}}StatusStrategy) NamespaceScoped() bool { return false }
{{ end -}}

type {{$api.Kind}}List struct {
	metav1.TypeMeta
	metav1.ListMeta
	Items []{{$api.Kind}}
}

{{ range $subresource := $api.Subresources -}}
type {{$subresource.Request}}List struct {
	metav1.TypeMeta
	metav1.ListMeta
	Items []{{$subresource.Request}}
}
{{ end -}}

func ({{$api.Kind}}) NewStatus() interface{} {
	return {{$api.Kind}}Status{}
}

func (pc *{{$api.Kind}}) GetStatus() interface{} {
	return pc.Status
}

func (pc *{{$api.Kind}}) SetStatus(s interface{}) {
	pc.Status = s.({{$api.Kind}}Status)
}

func (pc *{{$api.Kind}}) GetSpec() interface{} {
	return pc.Status
}

func (pc *{{$api.Kind}}) SetSpec(s interface{}) {
	pc.Spec = s.({{$api.Kind}}Spec)
}

func (pc *{{$api.Kind}}) GetObjectMeta() *metav1.ObjectMeta {
	return &pc.ObjectMeta
}

func (pc *{{$api.Kind}}) SetGeneration(generation int64) {
	pc.ObjectMeta.Generation = generation
}

func (pc {{$api.Kind}}) GetGeneration() int64 {
	return pc.ObjectMeta.Generation
}

// Registry is an interface for things that know how to store {{.Kind}}.
type {{.Kind}}Registry interface {
	List{{.Kind}}s(ctx request.Context, options *internalversion.ListOptions) (*{{.Kind}}List, error)
	Get{{.Kind}}(ctx request.Context, id string, options *metav1.GetOptions) (*{{.Kind}}, error)
	Create{{.Kind}}(ctx request.Context, id *{{.Kind}}) (*{{.Kind}}, error)
	Update{{.Kind}}(ctx request.Context, id *{{.Kind}}) (*{{.Kind}}, error)
	Delete{{.Kind}}(ctx request.Context, id string) (bool, error)
}

// NewRegistry returns a new Registry interface for the given Storage. Any mismatched types will panic.
func New{{.Kind}}Registry(sp builders.StandardStorageProvider) {{.Kind}}Registry {
	return &storage{{.Kind}}{sp}
}

// Implement Registry
// storage puts strong typing around storage calls
type storage{{.Kind}} struct {
	builders.StandardStorageProvider
}

func (s *storage{{.Kind}}) List{{.Kind}}s(ctx request.Context, options *internalversion.ListOptions) (*{{.Kind}}List, error) {
	if options != nil && options.FieldSelector != nil && !options.FieldSelector.Empty() {
		return nil, fmt.Errorf("field selector not supported yet")
	}
	st := s.GetStandardStorage()
	obj, err := st.List(ctx, options)
	if err != nil {
		return nil, err
	}
	return obj.(*{{.Kind}}List), err
}

func (s *storage{{.Kind}}) Get{{.Kind}}(ctx request.Context, id string, options *metav1.GetOptions) (*{{.Kind}}, error) {
	st := s.GetStandardStorage()
	obj, err := st.Get(ctx, id, options)
	if err != nil {
		return nil, err
	}
	return obj.(*{{.Kind}}), nil
}

func (s *storage{{.Kind}}) Create{{.Kind}}(ctx request.Context, object *{{.Kind}}) (*{{.Kind}}, error) {
	st := s.GetStandardStorage()
	obj, err := st.Create(ctx, object)
	if err != nil {
		return nil, err
	}
	return obj.(*{{.Kind}}), nil
}

func (s *storage{{.Kind}}) Update{{.Kind}}(ctx request.Context, object *{{.Kind}}) (*{{.Kind}}, error) {
	st := s.GetStandardStorage()
	obj, _, err := st.Update(ctx, object.Name, rest.DefaultUpdatedObjectInfo(object, api.Scheme))
	if err != nil {
		return nil, err
	}
	return obj.(*{{.Kind}}), nil
}

func (s *storage{{.Kind}}) Delete{{.Kind}}(ctx request.Context, id string) (bool, error) {
	st := s.GetStandardStorage()
	_, sync, err := st.Delete(ctx, id, nil)
	return sync, err
}

{{ end -}}
`
