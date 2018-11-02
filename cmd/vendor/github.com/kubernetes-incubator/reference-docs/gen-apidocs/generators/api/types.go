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
	"strings"

	"github.com/go-openapi/spec"
)

type ApiGroup string
type ApiGroups []ApiGroup

type ApiKind string

type ApiVersion string
type ApiVersions []ApiVersion
func (a ApiVersion) String() string {
	return string(a)
}

type SortDefinitionsByName []*Definition

func (a SortDefinitionsByName) Len() int      { return len(a) }
func (a SortDefinitionsByName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortDefinitionsByName) Less(i, j int) bool {
	if a[i].Name == a[j].Name {
		if a[i].Version.String() == a[j].Version.String() {
			return a[i].Group.String() < a[j].Group.String()
		}
		return a[i].Version.LessThan(a[j].Version)
	}
	return a[i].Name < a[j].Name
}

type SortDefinitionsByVersion []*Definition

func (a SortDefinitionsByVersion) Len() int      { return len(a) }
func (a SortDefinitionsByVersion) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortDefinitionsByVersion) Less(i, j int) bool {
	switch {
	case a[i].Version == a[j].Version:
		return strings.Compare(a[i].Group.String(), a[j].Group.String()) < 0
	default:
		return a[i].Version.LessThan(a[j].Version)
	}
}

type Definition struct {
	// open-api schema for the definition
	schema spec.Schema
	// Display name of the definition (e.g. Deployment)
	Name      string
	Group     ApiGroup
	ShowGroup bool

	// Api version of the definition (e.g. v1beta1)
	Version                 ApiVersion
	Kind                    ApiKind
	DescriptionWithEntities string
	GroupFullName           string

	// InToc is true if this definition should appear in the table of contents
	InToc        bool
	IsInlined    bool
	IsOldVersion bool

	FoundInField     bool
	FoundInOperation bool

	// Inline is a list of definitions that should appear inlined with this one in the documentations
	Inline SortDefinitionsByName

	// AppearsIn is a list of definition that this one appears in - e.g. PodSpec in Pod
	AppearsIn SortDefinitionsByName

	OperationCategories []*OperationCategory

	// Fields is a list of fields in this definition
	Fields Fields

	OtherVersions SortDefinitionsByName
	NewerVersions SortDefinitionsByName

	Sample SampleConfig

	FullName string
	Resource string
}

// Definitions indexes open-api definitions
type Definitions struct {
	All    map[string]*Definition
	ByKind map[string]SortDefinitionsByVersion
}

type DefinitionList []*Definition

type Config struct {
	ApiGroups           []ApiGroup          `yaml:"api_groups,omitempty"`
	ExampleLocation     string              `yaml:"example_location,omitempty"`
	OperationCategories []OperationCategory `yaml:"operation_categories,omitempty"`
	ResourceCategories  []ResourceCategory  `yaml:"resource_categories,omitempty"`

	// Used to map the group as the resource sees it to the group as the operation sees it
	GroupMap map[string]string

	Definitions Definitions
	Operations  Operations
}

// InlineDefinition defines a definition that should be inlined when displaying a Concept instead of appearing the in "Definitions"
type InlineDefinition struct {
	Name string `yaml:",omitempty"`
	Match string `yaml:",omitempty"`
}

type Field struct {
	Name                    string
	Type                    string
	Description             string
	DescriptionWithEntities string

	Definition *Definition // Optional Definition for complex types

	PatchStrategy string
	PatchMergeKey string
}

type Fields []*Field

func (a Fields) Len() int           { return len(a) }
func (a Fields) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Fields) Less(i, j int) bool { return a[i].Name < a[j].Name }

func (f Field) Link() string {
	if f.Definition != nil {
		return strings.Replace(f.Type, f.Definition.Name, f.Definition.MdLink(), -1)
	} else {
		return f.Type
	}
}

func (f Field) FullLink() string {
	if f.Definition != nil {
		return strings.Replace(f.Type, f.Definition.Name, f.Definition.HrefLink(), -1)
	} else {
		return f.Type
	}
}

// Operation defines a highlevel operation type such as Read, Replace, Patch
type OperationType struct {
	// Name is the display name of this operation
	Name string `yaml:",omitempty"`
	// Match is the regular expression of operation IDs that match this group where '${resource}' matches the resource name.
	Match string `yaml:",omitempty"`
}

type ExampleText struct {
	Tab  string
	Type string
	Text string
	Msg  string
}

type HttpResponse struct {
	Field
	Code string
}

type HttpResponses []*HttpResponse

type Operation struct {
	item          spec.PathItem
	op            *spec.Operation
	ID            string
	Type          OperationType
	Path          string
	HttpMethod    string
	Definition    *Definition
	BodyParams    Fields
	QueryParams   Fields
	PathParams    Fields
	HttpResponses HttpResponses

	ExampleConfig ExampleConfig
}

type Operations map[string]*Operation

// OperationCategory defines a group of related operations
type OperationCategory struct {
	// Name is the display name of this group
	Name string `yaml:",omitempty"`
	// Operations are the collection of Operations in this group
	OperationTypes []OperationType `yaml:"operation_types,omitempty"`
	// Default is true if this is the default operation group for operations that do not match any other groups
	Default bool `yaml:",omitempty"`

	Operations []*Operation
}

type ExampleProvider interface {
	GetTab() string
	GetRequestMessage() string
	GetResponseMessage() string
	GetRequestType() string
	GetResponseType() string
	GetSampleType() string
	GetSample(d *Definition) string
	GetRequest(o *Operation) string
	GetResponse(o *Operation) string
}

type EmptyExample struct{}
type CurlExample struct{}
type KubectlExample struct{}

type Resource struct {
	// Name is the display name of this Resource
	Name    string `yaml:",omitempty"`
	Version string `yaml:",omitempty"`
	Group   string `yaml:",omitempty"`

	// InlineDefinition is a list of definitions to show along side this resource when displaying it
	InlineDefinition []string `yaml:inline_definition",omitempty"`
	// DescriptionWarning is a warning message to show along side this resource when displaying it
	DescriptionWarning string `yaml:"description_warning,omitempty"`
	// DescriptionNote is a note message to show along side this resource when displaying it
	DescriptionNote string `yaml:"description_note,omitempty"`
	// ConceptGuide is a link to the concept guide for this resource if it exists
	ConceptGuide string `yaml:"concept_guide,omitempty"`
	// RelatedTasks is as list of tasks related to this concept
	RelatedTasks []string `yaml:"related_tasks,omitempty"`
	// IncludeDescription is the path to an md file to incline into the description
	IncludeDescription string `yaml:"include_description,omitempty"`
	// LinkToMd is the relative path to the md file containing the contents that clicking on this should link to
	LinkToMd string `yaml:"link_to_md,omitempty"`

	// Definition of the object
	Definition *Definition
}

type Resources []*Resource

// ResourceCategory defines a category of Concepts
type ResourceCategory struct {
	// Name is the display name of this group
	Name string `yaml:",omitempty"`
	// Include is the name of the _resource.md file to include in the index.html.md
	Include string `yaml:",omitempty"`
	// Resources are the collection of Resources in this group
	Resources Resources `yaml:",omitempty"`
	// LinkToMd is the relative path to the md file containing the contents that clicking on this should link to
	LinkToMd string `yaml:"link_to_md,omitempty"`
}

type ExampleConfig struct {
	Name         string `yaml:",omitempty"`
	Namespace    string `yaml:",omitempty"`
	Request      string `yaml:",omitempty"`
	Response     string `yaml:",omitempty"`
	RequestNote  string `yaml:",omitempty"`
	ResponseNote string `yaml:",omitempty"`
}

type SampleConfig struct {
	Note   string `yaml:",omitempty"`
	Sample string `yaml:",omitempty"`
}

type ResourceVisitor func(resource *Resource, d *Definition)
