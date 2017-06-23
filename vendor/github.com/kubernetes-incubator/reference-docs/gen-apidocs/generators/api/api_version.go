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

package api

type ApiVersion string

func (this ApiVersion) LessThan(that ApiVersion) bool {
	v1 := 100
	switch this {
	case "v1":
		v1 = 0
	case "v1beta1":
		v1 = 1
	case "v1alpha1":
		v1 = 2
	}
	v2 := 100
	switch that {
	case "v1":
		v2 = 0
	case "v1beta1":
		v2 = 1
	case "v1alpha1":
		v2 = 2
	}
	return v1 < v2
}
func (a ApiVersion) String() string {
	return string(a)
}
