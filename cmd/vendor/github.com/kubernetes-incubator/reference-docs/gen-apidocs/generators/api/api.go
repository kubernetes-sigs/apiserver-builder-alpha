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

import (
	"regexp"
	"strings"
)

func (a ApiGroup) String() string {
	return string(a)
}

func (a ApiGroups) Len() int      { return len(a) }
func (a ApiGroups) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ApiGroups) Less(i, j int) bool {
	// "apps" group APIs are newer than "extensions" group APIs
	if a[i].String() == "apps" && a[j].String() == "extensions" {
		return false
	}
	if a[j].String() == "apps" && a[i].String() == "extensions" {
		return true
	}
	return strings.Compare(a[i].String(), a[j].String()) < 0
}

func (a ApiVersions) Len() int      { return len(a) }
func (a ApiVersions) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ApiVersions) Less(i, j int) bool {
	return a[i].LessThan(a[j])
}

func (k ApiKind) String() string {
	return string(k)
}

func (this ApiVersion) LessThan(that ApiVersion) bool {
	re := regexp.MustCompile("(v\\d+)(alpha|beta|)(\\d*)")
	thisMatches := re.FindStringSubmatch(string(this))
	thatMatches := re.FindStringSubmatch(string(that))

	a := []string{thisMatches[1]}
	if len(thisMatches) > 2 {
		a = []string{thisMatches[2], thisMatches[1], thisMatches[0]}
	}

	b := []string{thatMatches[1]}
	if len(thatMatches) > 2 {
		b = []string{thatMatches[2], thatMatches[1], thatMatches[0]}
	}

	for i := 0; i < len(a) && i < len(b); i++ {
		v1 := ""
		v2 := ""
		if i < len(a) {
			v1 = a[i]
		}
		if i < len(b) {
			v2 = b[i]
		}
		// If the "beta" or "alpha" is missing, then it is ga (empty string comes before non-empty string)
		if len(v1) == 0 || len(v2) == 0 {
			return v1 < v2
		}
		// The string with the higher number comes first (or in the case of alpha/beta, beta comes first)
		if v1 != v2 {
			return v1 > v2
		}
	}

	return false
}
