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
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-openapi/loads"
)

const (
	patchStrategyKey = "x-kubernetes-patch-strategy"
	patchMergeKeyKey = "x-kubernetes-patch-merge-key"
	resourceNameKey  = "x-kubernetes-resource"
	typeKey          = "x-kubernetes-group-version-kind"
)

// Loads all of the open-api documents
func LoadOpenApiSpec() []*loads.Document {
	dir := filepath.Join(*ConfigDir, "openapi-spec/")
	docs := []*loads.Document{}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		ext := filepath.Ext(path)
		if ext != ".json" {
			return nil
		}
		var d *loads.Document
		d, err = loads.JSONSpec(path)
		if err != nil {
			return fmt.Errorf("Could not load json file %s as api-spec: %v\n", path, err)
		}
		docs = append(docs, d)
		return nil
	})
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("%v", err))
		os.Exit(1)
	}
	return docs
}

// return the map from short group name to full group name
func buildGroupMap(specs []*loads.Document) map[string]string {
	mapping := map[string]string{}
	mapping["apiregistration"] = "apiregistration.k8s.io"
	mapping["apiextensions"] = "apiextensions.k8s.io"
	mapping["certificates"] = "certificates.k8s.io"
	mapping["meta"] = "meta"
	mapping["core"] = "core"
	mapping["extensions"] = "extensions"

	for _, spec := range specs {
		for name, spec := range spec.Spec().Definitions {
			group, _, _ := GuessGVK(name)
			if _, found := mapping[group]; found {
				continue
			}
			// special groups where group name from extension is empty!
			if group == "meta" || group == "core" {
				continue
			}

			// full group not exposed as x-kubernetes- openapi extensions
			// from kube-aggregator project or apiextensions-apiserver project
			if group == "apiregistration" || group == "apiextensions" {
				continue
			}

			if extension, found := spec.Extensions[typeKey]; found {
				gvks, ok := extension.([]interface{})
				if ok {
					for _, item := range gvks {
						gvk, ok := item.(map[string]interface{})
						if ok {
							mapping[group] = gvk["group"].(string)
							break
						}
					}
				}
			}
		}
	}
	return mapping
}

func LoadDefinitions(specs []*loads.Document, s *Definitions) {
	groups := map[string]string{}
	groupMapping := buildGroupMap(specs)
	for _, spec := range specs {
		for name, spec := range spec.Spec().Definitions {
			resource := ""
			if r, ok := spec.Extensions.GetString(resourceNameKey); ok {
				resource = r
			}

			// This actually skips the following groupsi, i.e. old definitions
			//  'io.k8s.kubernetes.pkg.api.*'
			//  'io.k8s.kubernetes.pkg.apis.*'
			if strings.HasPrefix(spec.Description, "Deprecated. Please use") {
				continue
			}

			// NOTE:
			if strings.Contains(name, "JSONSchemaPropsOrStringArray") {
				continue
			}

			group, version, kind := GuessGVK(name)
			if group == "" {
				continue
			} else if group == "error" {
				panic(errors.New(fmt.Sprintf("Could not locate group for %s", name)))
			}
			groups[group] = ""

			full_group, found := groupMapping[group]
			if !found {
				// fall back to group name if no mapping found
				full_group = group
			}

			d := &Definition{
				schema:        spec,
				Name:          kind,
				Version:       ApiVersion(version),
				Kind:          ApiKind(kind),
				Group:         ApiGroup(group),
				GroupFullName: full_group,
				ShowGroup:     true,
				Resource:      resource,
			}

			s.All[d.Key()] = d
		}
	}
}
