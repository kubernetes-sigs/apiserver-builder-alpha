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

package v1beta1

import (
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource"
)

var _ resource.ObjectWithArbitrarySubResource = &University{}

type UniversityCampus struct {
	Faculty int `json:"faculty,omitempty"`
}

func (in *University) SubResourceNames() []string {
	return []string{"campus"}
}

func (in *University) SetSubResource(subResourceName string, subResource interface{}) {
	switch subResourceName {
	case "campus":
		campus := subResource.(UniversityCampus)
		in.Spec.FacultySize = campus.Faculty
	}
	panic("unknown subresource " + subResourceName)
}

func (in *University) GetSubResource(subResourceName string) (subResource interface{}) {
	switch subResourceName {
	case "campus":
		return UniversityCampus{
			Faculty: in.Spec.FacultySize,
		}
	}
	panic("unknown subresource " + subResourceName)
}
