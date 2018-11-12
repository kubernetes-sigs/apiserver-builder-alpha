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
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
)

func NewDefinitions(specs []*loads.Document) Definitions {
	s := Definitions{
		All:	map[string]*Definition{},
		ByKind: map[string]SortDefinitionsByVersion{},
	}

	LoadDefinitions(specs, &s)
	s.initialize()
	return s
}

func (s *Definitions) initialize() {
	// initialize fields for all definitions
	for _, d := range s.All {
		s.InitializeFields(d)
	}

	for _, d := range s.All {
		s.ByKind[d.Name] = append(s.ByKind[d.Name], d)
	}

	// If there are multiple versions for an object.  Mark all by the newest as old
	// Sort the ByKind index in by version with newer versions coming before older versions.
	for k, l := range s.ByKind {
		if len(l) <= 1 {
			continue
		}
		sort.Sort(l)
		// Mark all version as old
		for i, d := range l {
			if len(l) > 1 {
				if i > 0 {
					fmt.Printf("%s.%s.%s", d.Group, d.Version, k)
					if len(l) > i-1 {
						fmt.Printf(",")
					}
				} else {
					fmt.Printf("Current Version: %s.%s.%s", d.Group, d.Version, k)
					if len(l) > i-1 {
						fmt.Printf(" Old Versions: [")
					}
				}
			}
			if i > 0 {
				d.IsOldVersion = true
			}
		}
		if len(l) > 1 {
			fmt.Printf("]\n")
		}
	}

	// Initialize OtherVersions
	for _, d := range s.All {
		defs := s.ByKind[d.Name]
		others := []*Definition{}
		for _, def := range defs {
			if def.Version != d.Version {
				others = append(others, def)
			}
		}
		d.OtherVersions = others
	}

	// Initialize AppearsIn and FoundInField
	for _, d := range s.All {
		for _, r := range s.getReferences(d) {
			r.AppearsIn = append(r.AppearsIn, d)
			r.FoundInField = true
		}
	}

	// Initialize Inline, IsInlined 
	// Note: examples of inline definitions are "Spec", "Status", "List", etc
	for _, d := range s.All {
		for _, name := range GetInlinedDefinitionNames(d.Name) {
			if cr, ok := s.GetByVersionKind(string(d.Group), string(d.Version), name); ok {
				d.Inline = append(d.Inline, cr)
				cr.IsInlined = true
				cr.FoundInField = true
			}
		}
	}
}

func (s *Definitions) getReferences(d *Definition) []*Definition {
	refs := []*Definition{}
	// Find all of the resources referenced by this definition
	for _, p := range d.schema.Properties {
		if !IsComplex(p) {
			// Skip primitive types and collections of primitive types
			continue
		}
		// Look up the definition for the referenced resource
		if schema, ok := s.GetForSchema(p); ok {
			refs = append(refs, schema)
		} else {
			g, v, k := GetDefinitionVersionKind(p)
			fmt.Printf("Could not locate referenced property of %s: %s (%s/%s).\n", d.Name, g, k, v)
		}
	}
	return refs
}

func (s *Definitions) parameterToField(param spec.Parameter) *Field {
	f := &Field{
		Name:        param.Name,
		Description: strings.Replace(param.Description, "\n", " ", -1),
	}
	if param.Schema != nil {
		f.Type = GetTypeName(*param.Schema)
		if fieldType, ok := s.GetForSchema(*param.Schema); ok {
			f.Definition = fieldType
		}
	}
	return f
}

// GetByVersionKind looks up a definition using its primary key (version,kind)
func (s *Definitions) GetByVersionKind(group, version, kind string) (*Definition, bool) {
	key := &Definition{Group: ApiGroup(group), Version: ApiVersion(version), Kind: ApiKind(kind)}
	r, f := s.All[key.Key()]
	return r, f
}

func (s *Definitions) GetForSchema(schema spec.Schema) (*Definition, bool) {
	g, v, k := GetDefinitionVersionKind(schema)
	if len(k) <= 0 {
		return nil, false
	}
	return s.GetByVersionKind(g, v, k)
}

// Initializes the fields for a definition
func (s *Definitions) InitializeFields(d *Definition) {
	for fieldName, property := range d.schema.Properties {
		des := strings.Replace(property.Description, "\n", " ", -1)
		f := &Field{
			Name:        fieldName,
			Type:        GetTypeName(property),
			Description: EscapeAsterisks(des),
		}
		if len(property.Extensions) > 0 {
			if ps, ok := property.Extensions.GetString(patchStrategyKey); ok {
				f.PatchStrategy = ps
			}
			if pmk, ok := property.Extensions.GetString(patchMergeKeyKey); ok {
				f.PatchMergeKey = pmk
			}
		}

		if fd, ok := s.GetForSchema(property); ok {
			f.Definition = fd
		}
		d.Fields = append(d.Fields, f)
	}
}

func (d *Definition) GroupDisplayName() string {
	if len(d.GroupFullName) > 0 {
		return d.GroupFullName
	}
	if len(d.Group) <= 0 || d.Group == "core" {
		return "Core"
	}
	return string(d.Group)
}
func (d *Definition) GetOperationGroupName() string {
	if strings.ToLower(d.Group.String()) == "rbac" {
		return "RbacAuthorization"
	}
	return strings.Title(d.Group.String())
}

func (d *Definition) Key() string {
	return fmt.Sprintf("%s.%s.%s", d.Group, d.Version, d.Kind)
}

func (d *Definition) LinkID() string {
	groupName := strings.Replace(strings.ToLower(d.GroupFullName), ".", "-", -1)
	link := fmt.Sprintf("%s-%s-%s", d.Name, d.Version, groupName)
	return strings.ToLower(link)
}

func (d *Definition) MdLink() string {
	groupName := strings.Replace(strings.ToLower(d.GroupFullName), ".", "-", -1)
	return fmt.Sprintf("[%s](#%s-%s-%s)", d.Name, strings.ToLower(d.Name), d.Version, groupName)
}

func (d *Definition) HrefLink() string {
	groupName := strings.Replace(strings.ToLower(d.GroupFullName), ".", "-", -1)
	return fmt.Sprintf("<a href=\"#%s-%s-%s\">%s</a>", strings.ToLower(d.Name), d.Version, groupName, d.Name)
}

func (d *Definition) FullHrefLink() string {
	groupName := strings.Replace(strings.ToLower(d.GroupFullName), ".", "-", -1)
	return fmt.Sprintf("<a href=\"#%s-%s-%s\">%s [%s/%s]</a>", strings.ToLower(d.Name),
		d.Version, groupName, d.Name, d.Group, d.Version)
}

func (d *Definition) VersionLink() string {
	groupName := strings.Replace(strings.ToLower(d.GroupFullName), ".", "-", -1)
	return fmt.Sprintf("<a href=\"#%s-%s-%s\">%s</a>", strings.ToLower(d.Name), d.Version, groupName, d.Version)
}

func (d *Definition) Description() string {
	return EscapeAsterisks(d.schema.Description)
}

func (d *Definition) GetResourceName() string {
	if len(d.Resource) > 0 {
		return d.Resource
	}
	resource := strings.ToLower(d.Name)
	if strings.HasSuffix(resource, "y") {
		return strings.TrimSuffix(resource, "y") + "ies"
	}
	return resource + "s"
}

func (d *Definition) initExample(config *Config) {
	path := filepath.Join(*ConfigDir, config.ExampleLocation, d.Name, d.Name + ".yaml")
	file := strings.Replace(strings.ToLower(path), " ", "_", -1)
	content, err := ioutil.ReadFile(file)
	if err != nil || len(content) <= 0 {
		return
	}
	err = yaml.Unmarshal(content, &d.Sample)
	if err != nil {
		panic(fmt.Sprintf("Could not Unmarshal SampleConfig yaml: %s\n", content))
	}
}

func (d *Definition) GetSamples() []ExampleText {
	r := []ExampleText{}
	for _, p := range GetExampleProviders() {
		r = append(r, ExampleText{
			Tab:  p.GetTab(),
			Type: p.GetSampleType(),
			Text: p.GetSample(d),
		})
	}
	return r
}

func (a DefinitionList) Len() int      { return len(a) }
func (a DefinitionList) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a DefinitionList) Less(i, j int) bool {
	return strings.Compare(a[i].Name, a[j].Name) < 0
}
