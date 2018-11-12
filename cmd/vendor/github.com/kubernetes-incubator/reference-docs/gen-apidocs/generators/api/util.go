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
	"fmt"
	"strings"

	"errors"
	"github.com/go-openapi/spec"
)

var INLINE_DEFINITIONS = []InlineDefinition{
	{Name: "Spec", Match: "${resource}Spec"},
	{Name: "Status", Match: "${resource}Status"},
	{Name: "List", Match: "${resource}List"},
	{Name: "Strategy", Match: "${resource}Strategy"},
	{Name: "Rollback", Match: "${resource}Rollback"},
	{Name: "RollingUpdate", Match: "RollingUpdate${resource}"},
	{Name: "EventSource", Match: "${resource}EventSource"},
}


// GetDefinitionVersionKind returns the api version and kind for the spec.  This is the primary key of a Definition.
func GetDefinitionVersionKind(s spec.Schema) (string, string, string) {
	// Get the reference for complex types
	if IsDefinition(s) {
		s := fmt.Sprintf("%s", s.SchemaProps.Ref.GetPointer())
		s = strings.Replace(s, "/definitions/", "", -1)
		name := strings.Split(s, ".")

		var group, version, kind string
		if name[len(name)-3] == "api" {
			// e.g. "io.k8s.apimachinery.pkg.api.resource.Quantity"
			group = "core"
			version = name[len(name)-2]
			kind = name[len(name)-1]
		} else if name[len(name)-4] == "api" {
			// e.g. "io.k8s.api.core.v1.Pod"
			group = name[len(name)-3]
			version = name[len(name)-2]
			kind = name[len(name)-1]
		} else if name[len(name)-4] == "apis" {
			// e.g. "io.k8s.apimachinery.pkg.apis.meta.v1.Status"
			group = name[len(name)-3]
			version = name[len(name)-2]
			kind = name[len(name)-1]
		} else if name[len(name)-3] == "util" || name[len(name)-3] == "pkg" {
			// e.g. io.k8s.apimachinery.pkg.util.intstr.IntOrString
			// e.g. io.k8s.apimachinery.pkg.runtime.RawExtension
			return "", "", ""
		} else {
			panic(errors.New(fmt.Sprintf("Could not locate group for %s", name)))
		}
		return group, version, kind
	}
	// Recurse if type is array
	if IsArray(s) {
		return GetDefinitionVersionKind(*s.Items.Schema)
	}
	return "", "", ""
}

// GetTypeName returns the display name of a Schema.  This is the api kind for definitions and the type for
// primitive types.  Arrays of objects have "array" appended.
func GetTypeName(s spec.Schema) string {
	// Get the reference for complex types
	if IsDefinition(s) {
		_, _, name := GetDefinitionVersionKind(s)
		return name
	}
	// Recurse if type is array
	if IsArray(s) {
		return fmt.Sprintf("%s array", GetTypeName(*s.Items.Schema))
	}
	// Get the value for primitive types
	if len(s.Type) > 0 {
		return fmt.Sprintf("%s", s.Type[0])
	}
	panic(fmt.Errorf("No type found for object %v", s))
}

// IsArray returns true if the type is an array type.
func IsArray(s spec.Schema) bool {
	return len(s.Type) > 0 && s.Type[0] == "array"
}

// IsDefinition returns true if Schema is a complex type that should have a Definition.
func IsDefinition(s spec.Schema) bool {
	return len(s.SchemaProps.Ref.GetPointer().String()) > 0
}

// handle '*', 'a/*', '*/b', '*/*' cases
func EscapeAsterisks(des string) string {
	s := strings.Replace(des, "'*'", `'\*'`, -1)
	s = strings.Replace(s, "/*'", `/\*'`, -1)
	s = strings.Replace(s, "'*/", `'\*/`, -1)
	s = strings.Replace(s, "'*/*'", `'\*/\*'`, -1)
	return s
}

// IsComplex returns true if the schema is for a complex (non-primitive) definitions
func IsComplex(schema spec.Schema) bool {
	_, _, k := GetDefinitionVersionKind(schema)
	return len(k) > 0
}

// GuessGVK makes a guess about the (Group, Version, Kind) tuple based on
// resource name
// TODO: Rework this function because it is ugly
func GuessGVK(name string) (group, version, kind string) {
	parts := strings.Split(name, ".")
	if len(parts) < 4 {
		fmt.Printf("Error: Could not find version and type for definition %s.\n", name)
		return "", "", ""
	}

	if parts[len(parts)-3] == "api" {
		// e.g. "io.k8s.apimachinery.pkg.api.resource.Quantity"
		group = "core"
		version = parts[len(parts)-2]
		kind = parts[len(parts)-1]
	} else if parts[len(parts)-4] == "api" {
		// e.g. "io.k8s.api.core.v1.Pod"
		group = parts[len(parts)-3]
		version = parts[len(parts)-2]
		kind = parts[len(parts)-1]
	} else if parts[len(parts)-4] == "apis" {
		// e.g. "io.k8s.apimachinery.pkg.apis.meta.v1.Status"
		group = parts[len(parts)-3]
		version = parts[len(parts)-2]
		kind = parts[len(parts)-1]
	} else if parts[len(parts)-3] == "util" || parts[len(parts)-3] == "pkg" {
		// e.g. io.k8s.apimachinery.pkg.util.intstr.IntOrString
		// e.g. io.k8s.apimachinery.pkg.runtime.RawExtension
		return "", "", ""
	} else {
		// To report error
		return "error", "", ""
	}
	return group, version, kind
}

func GetInlinedDefinitionNames(parent string) []string {
	names := []string{}
	for _, id := range INLINE_DEFINITIONS {
		name := strings.Replace(id.Match, "${resource}", parent, -1)
		names = append(names, name)
	}
	return names
}
