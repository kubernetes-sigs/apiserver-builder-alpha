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

package parse

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"bytes"
	"encoding/json"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/gengo/types"
	"text/template"
)

func (b *APIs) parseJSONSchemaProps() {
	for _, group := range b.APIs.Groups {
		for _, version := range group.Versions {
			for _, resource := range version.Resources {
				if IsAPIResource(resource.Type) {
					fmt.Printf("Resource: %+v\n\n", resource.Kind)
					resource.JSONSchemaProps, resource.Validation = b.typeToJSONSchemaProps(resource.Type)
					j, err := json.MarshalIndent(resource.JSONSchemaProps, "", "    ")
					if err != nil {
						log.Fatalf("Could not Marshall validation %v\n", err)
					}
					fmt.Printf("JSON RESOURCE: %s\n\n", string(j))
					resource.ValidationComments = string(j)
				}
			}
		}
	}
}

func (b *APIs) getTime() string {
	return `v1beta1.JSONSchemaProps{
    Type:   "string",
    Format: "date-time",
}`
}

func (b *APIs) getMeta() string {
	return `v1beta1.JSONSchemaProps{
    Type:   "object",
}`
}

func (b *APIs) typeToJSONSchemaProps(t *types.Type) (v1beta1.JSONSchemaProps, string) {
	// Special cases
	fmt.Printf("(%s) (%s) (%s)\n", t.Name.Name, t.Name.Path, t.Name.Package)
	time := types.Name{Name: "Time", Package: "k8s.io/apimachinery/pkg/apis/meta/v1"}
	meta := types.Name{Name: "ObjectMeta", Package: "k8s.io/apimachinery/pkg/apis/meta/v1"}
	switch t.Name {
	case time:
		return v1beta1.JSONSchemaProps{
			Type:   "string",
			Format: "date-time",
		}, b.getTime()
	case meta:
		return v1beta1.JSONSchemaProps{
			Type: "object",
		}, b.getMeta()
	}

	switch t.Kind {
	case types.Builtin:
		return b.parsePrimitiveValidation(t)
	case types.Struct:
		return b.parseObjectValidation(t)
	case types.Map:
		return b.parseMapValidation(t)
		return b.parseMapValidation(t)
	case types.Slice:
		return b.parseArrayValidation(t)
	case types.Array:
		return b.parseArrayValidation(t)
	case types.Pointer:
		return b.typeToJSONSchemaProps(t.Elem)
	case types.Alias:
		return b.typeToJSONSchemaProps(t.Underlying)
	default:
		log.Fatalf("Unknown supported Kind %v\n", t.Kind)
	}
	// Unreachable
	return v1beta1.JSONSchemaProps{}, ""
}

var jsonRegex = regexp.MustCompile("json:\"([a-zA-Z,]+)\"")

var primitiveTemplate = template.Must(template.New("map-template").Parse(
	`v1beta1.JSONSchemaProps{Type: "{{.}}"}`))

func (b *APIs) parsePrimitiveValidation(t *types.Type) (v1beta1.JSONSchemaProps, string) {
	props := v1beta1.JSONSchemaProps{Type: string(t.Name.Name)}

	buff := &bytes.Buffer{}
	if err := primitiveTemplate.Execute(buff, t.Name.Name); err != nil {
		log.Fatalf("%v", err)
	}

	return props, buff.String()
}

var mapTemplate = template.Must(template.New("map-template").Parse(
	`v1beta1.JSONSchemaProps{
    Type:                 "object",
    AdditionalProperties: &v1beta1.JSONSchemaPropsOrBool{
        Allows: true,
        //Schema: &{{.}},
    },
}`))

func (b *APIs) parseMapValidation(t *types.Type) (v1beta1.JSONSchemaProps, string) {
	additionalProps, _ := b.typeToJSONSchemaProps(t.Elem)
	props := v1beta1.JSONSchemaProps{
		Type: "object",
		AdditionalProperties: &v1beta1.JSONSchemaPropsOrBool{
			Allows: true,
			Schema: &additionalProps},
	}

	buff := &bytes.Buffer{}
	if err := mapTemplate.Execute(buff, ""); err != nil {
		log.Fatalf("%v", err)
	}
	return props, buff.String()
}

var arrayTemplate = template.Must(template.New("array-template").Parse(
	`v1beta1.JSONSchemaProps{
    Type:                 "array",
    Items: &v1beta1.JSONSchemaPropsOrArray{
        Schema: &{{.}},
    },
}`))

func (b *APIs) parseArrayValidation(t *types.Type) (v1beta1.JSONSchemaProps, string) {
	items, result := b.typeToJSONSchemaProps(t.Elem)
	props := v1beta1.JSONSchemaProps{
		Type:  "array",
		Items: &v1beta1.JSONSchemaPropsOrArray{Schema: &items},
	}

	buff := &bytes.Buffer{}
	if err := arrayTemplate.Execute(buff, result); err != nil {
		log.Fatalf("%v", err)
	}
	return props, buff.String()
}

var objectTemplate = template.Must(template.New("object-template").Parse(
	`v1beta1.JSONSchemaProps{
    Type:                 "object",
    Properties: map[string]v1beta1.JSONSchemaProps{
        {{ range $k, $v := . -}}
        "{{ $k }}": {{ $v }},
        {{ end -}}
    },
}`))

func (b *APIs) parseObjectValidation(t *types.Type) (v1beta1.JSONSchemaProps, string) {
	m, result := b.getMembers(t)
	props := v1beta1.JSONSchemaProps{
		Type:       "object",
		Properties: m,
	}

	buff := &bytes.Buffer{}
	if err := objectTemplate.Execute(buff, result); err != nil {
		log.Fatalf("%v", err)
	}
	return props, buff.String()
}

func (b *APIs) getMembers(t *types.Type) (map[string]v1beta1.JSONSchemaProps, map[string]string) {
	members := map[string]v1beta1.JSONSchemaProps{}
	result := map[string]string{}
	for _, member := range t.Members {
		fmt.Printf("Field %s - ", member.Name)

		tags := jsonRegex.FindStringSubmatch(member.Tags)
		if len(tags) == 0 {
			// Skip fields without json tags
			//fmt.Printf("Skipping member %s %s\n", member.Name, member.Type.Name.String())
			continue
		}
		ts := strings.Split(tags[1], ",")
		name := member.Name
		strat := ""
		if len(ts) > 0 && len(ts[0]) > 0 {
			name = ts[0]
		}
		if len(ts) > 1 {
			strat = ts[1]
		}

		// Inline "inline" structs
		if strat == "inline" {
			m, r := b.getMembers(member.Type)
			for n, v := range m {
				members[n] = v
			}
			for n, v := range r {
				result[n] = v
			}
		} else {
			m, r := b.typeToJSONSchemaProps(member.Type)
			members[name] = m
			result[name] = r
		}
	}
	return members, result
}
