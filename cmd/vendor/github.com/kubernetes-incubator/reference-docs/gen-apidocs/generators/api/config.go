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
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"html"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"github.com/go-openapi/loads"
)

var AllowErrors = flag.Bool("allow-errors", false, "If true, don't fail on errors.")
var ConfigDir = flag.String("config-dir", "", "Directory contain api files.")
var UseTags = flag.Bool("use-tags", false, "If true, use the openapi tags instead of the config yaml.")
var MungeGroups = flag.Bool("munge-groups", true, "If true, munge the group names for the operations to match.")

func NewConfig() *Config {
	config := LoadConfigFromYAML()
	specs := LoadOpenApiSpec()

	// Initialize all of the operations
	config.Definitions = NewDefinitions(specs)

	if *UseTags {
		// Initialize the config and ToC from the tags on definitions
		config.genConfigFromTags(specs)
	} else {
		// Initialization for ToC resources only
		config.visitResourcesInToc()
	}

	config.initOperations(specs)

	// replace unicode escape sequences with HTML entities.
	config.escapeDescriptions()

	config.CleanUp()

	// Prune anything that shouldn't be in the ToC
	if *UseTags {
		categories := []ResourceCategory{}
		for _, c := range config.ResourceCategories {
			resources := Resources{}
			for _, r := range c.Resources {
				if d, f := config.Definitions.GetByVersionKind(r.Group, r.Version, r.Name); f {
					if d.InToc {
						resources = append(resources, r)
					}
				}
			}
			c.Resources = resources
			if len(resources) > 0 {
				categories = append(categories, c)
			}
		}
		config.ResourceCategories = categories
	}

	return config
}

func (c *Config) genConfigFromTags(specs []*loads.Document) {
	log.Printf("Using OpenAPI extension tags to configure.")

	c.ExampleLocation = "examples"
	// build the apis from the observed groups
	groupsMap := map[ApiGroup]DefinitionList{}
	for _, d := range c.Definitions.All {
		if strings.HasSuffix(d.Name, "List") {
			continue
		}
		if strings.HasSuffix(d.Name, "Status") {
			continue
		}
		if strings.HasPrefix(d.Description(), "Deprecated. Please use") {
			// Don't look at deprecated types
			continue
		}
		d.initExample(c)
		g := d.Group
		groupsMap[g] = append(groupsMap[g], d)
	}

	groupsList := ApiGroups{}
	for g := range groupsMap {
		groupsList = append(groupsList, g)
	}

	sort.Sort(groupsList)

	for _, g := range groupsList {
		groupName := strings.Title(string(g))
		c.ApiGroups = append(c.ApiGroups, ApiGroup(groupName))
		rc := ResourceCategory{
			Include: string(g),
			Name: groupName,
		}
		defList := groupsMap[g]
		sort.Sort(defList)
		for _, d := range defList {
			r := &Resource{
				Name: d.Name,
				Group: string(d.Group),
				Version: string(d.Version),
				Definition: d,
			}
			rc.Resources = append(rc.Resources, r)
		}
		c.ResourceCategories = append(c.ResourceCategories, rc)
	}
}

func (config *Config) initOperationsFromTags(specs []*loads.Document) {
	if *UseTags {
		ops := map[string]map[string][]*Operation{}
		defs := map[string]*Definition{}
		for _, d := range config.Definitions.All {
			name := fmt.Sprintf("%s.%s.%s", d.Group, d.Version, d.GetResourceName())
			defs[name] = d
		}

		VisitOperations(specs, func(operation Operation) {
			if o, found := config.Operations[operation.ID]; found && o.Definition != nil {
				return
			}
			op := operation
			o := &op
			config.Operations[operation.ID] = o
			group, version, kind, sub := o.GetGroupVersionKindSub()
			if sub == "status" {
				return
			}
			if len(group) == 0 {
				return
			}
			key := fmt.Sprintf("%s.%s.%s", group, version, kind)
			o.Definition = defs[key]

			// Index by group and subresource
			if _, f := ops[key]; !f {
				ops[key] = map[string][]*Operation{}
			}
			ops[key][sub] = append(ops[key][sub], o)
		})

		for key, subMap := range ops {
			def := defs[key]
			if def == nil {
				panic(fmt.Errorf("Unable to locate resource %s in resource map\n%v\n", key, defs))
			}
			subs := []string{}
			for s := range subMap {
				subs = append(subs, s)
			}
			sort.Strings(subs)
			for _, s := range subs {
				cat := &OperationCategory{}
				cat.Name = strings.Title(s) + " Operations"
				for _, op := range subMap[s] {
					ot := OperationType{}
					ot.Name = op.GetMethod() + " " + strings.Title(s)
					op.Type = ot
					cat.Operations = append(cat.Operations, op)
				}
				def.OperationCategories = append(def.OperationCategories, cat)
			}
		}
	}
}

// initOperations returns all Operations found in the Documents
func (c *Config) initOperations(specs []*loads.Document) {
	ops := Operations{}

	c.GroupMap = map[string]string{}
	VisitOperations(specs, func(op Operation) {
		ops[op.ID] = &op

		// Build a map of the group names to the group name appearing in operation ids
		// This is necessary because the group will appear without the domain
		// in the resource, but with the domain in the operationID, and we
		// will be unable to match the operationID to the resource because they
		// don't agree on the name of the group.
		// TODO: Fix this by getting the group-version-kind in the resource
		if *MungeGroups {
			if v, ok := op.op.Extensions[typeKey]; ok {
				gvk := v.(map[string]interface{})
				group, ok := gvk["group"].(string)
				if !ok {
					log.Fatalf("group not type string %v", v)
				}
				groupId := ""
				for _, s := range strings.Split(group, ".") {
					groupId = groupId + strings.Title(s)
				}
				c.GroupMap[strings.Title(strings.Split(group, ".")[0])] = groupId
			}
		}
	})

	c.Operations = ops
	c.mapOperationsToDefinitions()
	c.initOperationsFromTags(specs)

	VisitOperations(specs, func(target Operation) {
		if op, ok := c.Operations[target.ID]; !ok || op.Definition == nil {
			op.VerifyBlackListed()
		}
	})
	c.initOperationParameters()

	// Clear the operations.  We still have to calculate the operations because that is how we determine
	// the API Group for each definition.
	if !*BuildOps {
		c.Operations = Operations{}
		c.OperationCategories = []OperationCategory{}
		for _, d := range c.Definitions.All {
			d.OperationCategories = []*OperationCategory{}
		}
	}
}

// CleanUp sorts and dedups fields
func (c *Config) CleanUp() {
	for _, d := range c.Definitions.All {
		sort.Sort(d.AppearsIn)
		sort.Sort(d.Fields)
		dedup := SortDefinitionsByName{}
		var last *Definition
		for _, i := range d.AppearsIn {
			if last != nil &&
				i.Name == last.Name &&
				i.Group.String() == last.Group.String() &&
				i.Version.String() == last.Version.String() {
				continue
			}
			last = i
			dedup = append(dedup, i)
		}
		d.AppearsIn = dedup
	}
}

// LoadConfigFromYAML reads the config yaml file into a struct
func LoadConfigFromYAML() *Config {
	f := filepath.Join(*ConfigDir, "config.yaml")

	config := &Config{}
	contents, err := ioutil.ReadFile(f)
	if err != nil {
		if !*UseTags {
			fmt.Printf("Failed to read yaml file %s: %v", f, err)
			os.Exit(2)
		}
	} else {
		err = yaml.Unmarshal(contents, config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	writeCategory := OperationCategory{
		Name: "Write Operations",
		OperationTypes: []OperationType{
			{
				Name:  "Create",
				Match: "create${group}${version}(Namespaced)?${resource}",
			},
			{
				Name:  "Create Eviction",
				Match: "create${group}${version}(Namespaced)?${resource}Eviction",
			},
			{
				Name:  "Patch",
				Match: "patch${group}${version}(Namespaced)?${resource}",
			},
			{
				Name:  "Replace",
				Match: "replace${group}${version}(Namespaced)?${resource}",
			},
			{
				Name:  "Delete",
				Match: "delete${group}${version}(Namespaced)?${resource}",
			},
			{
				Name:  "Delete Collection",
				Match: "delete${group}${version}Collection(Namespaced)?${resource}",
			},
		},
	}

	readCategory := OperationCategory{
		Name: "Read Operations",
		OperationTypes: []OperationType{
			{
				Name:  "Read",
				Match: "read${group}${version}(Namespaced)?${resource}",
			},
			{
				Name:  "List",
				Match: "list${group}${version}(Namespaced)?${resource}",
			},
			{
				Name:  "List All Namespaces",
				Match: "list${group}${version}(Namespaced)?${resource}ForAllNamespaces",
			},
			{
				Name:  "Watch",
				Match: "watch${group}${version}(Namespaced)?${resource}",
			},
			{
				Name:  "Watch List",
				Match: "watch${group}${version}(Namespaced)?${resource}List",
			},
			{
				Name:  "Watch List All Namespaces",
				Match: "watch${group}${version}(Namespaced)?${resource}ListForAllNamespaces",
			},
		},
	}

	statusCategory := OperationCategory{
		Name: "Status Operations",
		OperationTypes: []OperationType{
			{
				Name:  "Patch Status",
				Match: "patch${group}${version}(Namespaced)?${resource}Status",
			},
			{
				Name:  "Read Status",
				Match: "read${group}${version}(Namespaced)?${resource}Status",
			},
			{
				Name:  "Replace Status",
				Match: "replace${group}${version}(Namespaced)?${resource}Status",
			},
		},
	}

	config.OperationCategories = append([]OperationCategory{writeCategory, readCategory, statusCategory}, config.OperationCategories...)

	return config
}


const (
	PATH  = "path"
	QUERY = "query"
	BODY  = "body"
)

func (c *Config) initOperationParameters() {
	s := c.Definitions
	for _, op := range c.Operations {
		pathItem := op.item

		// Path parameters
		for _, p := range pathItem.Parameters {
			switch p.In {
			case PATH:
				op.PathParams = append(op.PathParams, s.parameterToField(p))
			case QUERY:
				op.QueryParams = append(op.QueryParams, s.parameterToField(p))
			case BODY:
				op.BodyParams = append(op.BodyParams, s.parameterToField(p))
			default:
				panic("")
			}
		}

		// Query parameters
		for _, p := range op.op.Parameters {
			switch p.In {
			case PATH:
				op.PathParams = append(op.PathParams, s.parameterToField(p))
			case QUERY:
				op.QueryParams = append(op.QueryParams, s.parameterToField(p))
			case BODY:
				op.BodyParams = append(op.BodyParams, s.parameterToField(p))
			default:
				panic("")
			}
		}

		for code, response := range op.op.Responses.StatusCodeResponses {
			if response.Schema == nil {
				// fmt.Printf("Nil Schema for response: %+v\n", op.Path)
				continue
			}
			r := &HttpResponse{
				Field: Field{
					Description: strings.Replace(response.Description, "\n", " ", -1),
					Type:        GetTypeName(*response.Schema),
					Name:        fmt.Sprintf("%d", code),
				},
				Code: fmt.Sprintf("%d", code),
			}
			if IsComplex(*response.Schema) {
				r.Definition, _ = s.GetForSchema(*response.Schema)
				if r.Definition != nil {
					r.Definition.FoundInOperation = true
				}
			}
			op.HttpResponses = append(op.HttpResponses, r)
		}
	}
}

func (c *Config) getOperationId(match string, group string, version ApiVersion, kind string) string {
	// Lookup the name of the group as the operation expects it (different than the resource)
	if g, ok := c.GroupMap[group]; ok {
		group = g
	}

	ver := []rune(string(version))
	ver[0] = unicode.ToUpper(ver[0])

	match = strings.Replace(match, "${group}", string(group), -1)
	match = strings.Replace(match, "${version}", string(ver), -1)
	match = strings.Replace(match, "${resource}", kind, -1)
	return match
}

func (c *Config) setOperation(match, namespace string, ot *OperationType, oc *OperationCategory, d *Definition) {

	key := strings.Replace(match, "(Namespaced)?", namespace, -1)
	if o, ok := c.Operations[key]; ok {
		// Each operation should have exactly 1 definition
		if o.Definition != nil {
			panic(fmt.Sprintf(
				"Found multiple matching definitions [%s/%s/%s, %s/%s/%s] for operation key: %s",
				d.Group, d.Version, d.Name, o.Definition.Group, o.Definition.Version, o.Definition.Name, key))
		}
		o.Type = *ot
		o.Definition = d
		o.initExample(c)
		oc.Operations = append(oc.Operations, o)

		// When using tags for the configuration, everything with an operation goes in the ToC
		if *UseTags && !o.Definition.IsOldVersion {
			o.Definition.InToc = true
		}
	}
}

// mapOperationsToDefinitions adds operations to the definitions they operate
func (c *Config) mapOperationsToDefinitions() {
	for _, d := range c.Definitions.All {
		if d.IsInlined {
			continue
		}

		for i := range c.OperationCategories {
			oc := c.OperationCategories[i]
			for j := range oc.OperationTypes {
				ot := oc.OperationTypes[j]
				operationId := c.getOperationId(ot.Match, d.GetOperationGroupName(), d.Version, d.Name)
				c.setOperation(operationId, "Namespaced", &ot, &oc, d)
				c.setOperation(operationId, "", &ot, &oc, d)
			}

			if len(oc.Operations) > 0 {
				d.OperationCategories = append(d.OperationCategories, &oc)
			}
		}
	}
}

// The OpenAPI spec has escape sequences like \u003c. When the spec is unmarshaled,
// the escape sequences get converted to ordinary characters. For example,
// \u003c gets converted to a regular < character. But we can't use  regular <
// and > characters in our HTML document. This function replaces these regular
// characters with HTML entities: <, >, &, ', and ".
func (c *Config) escapeDescriptions() {
	for _, d := range c.Definitions.All {
		d.DescriptionWithEntities = html.EscapeString(d.Description())

		for _, f := range d.Fields {
			f.DescriptionWithEntities = html.EscapeString(f.Description)
		}
	}

	for _, op := range c.Operations {
		for _, p := range op.BodyParams {
			p.DescriptionWithEntities = html.EscapeString(p.Description)
		}
		for _, p := range op.QueryParams {
			p.DescriptionWithEntities = html.EscapeString(p.Description)
		}
		for _, p := range op.PathParams {
			p.DescriptionWithEntities = html.EscapeString(p.Description)
		}
		for _, r := range op.HttpResponses {
			r.DescriptionWithEntities = html.EscapeString(r.Description)
		}
	}
}

// For each resource in the ToC, look up its definition and visit it.
func (c *Config) visitResourcesInToc() {
	missing := false
	for _, cat := range c.ResourceCategories {
		for _, r := range cat.Resources {
			if d, ok := c.Definitions.GetByVersionKind(r.Group, r.Version, r.Name); ok {
				d.InToc = true // Mark as in Toc
				d.initExample(c)
				r.Definition = d
			} else {
				fmt.Printf("Could not find definition for resource in TOC: %s %s %s.\n", r.Group, r.Version, r.Name)
				missing = true
			}
		}
	}
	if missing {
		fmt.Printf("All known definitions: %v\n", c.Definitions.All)
	}
}
