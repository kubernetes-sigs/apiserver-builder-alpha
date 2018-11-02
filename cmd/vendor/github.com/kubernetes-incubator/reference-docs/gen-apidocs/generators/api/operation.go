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
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
)

var BuildOps = flag.Bool("build-operations", true, "If true build operations in the docs.")

// GetOperationId returns the ID of the operation for the given definition
func (ot OperationType) GetOperationId(d string) string {
	return strings.Replace(ot.Match, "${resource}", d, -1)
}

func (o *Operation) GetExampleRequests() []ExampleText {
	r := []ExampleText{}
	for _, p := range GetExampleProviders() {
		text := p.GetRequest(o)
		if len(text) > 0 {
			r = append(r, ExampleText{
				Tab:  p.GetTab(),
				Type: p.GetRequestType(),
				Text: p.GetRequest(o),
				Msg:  p.GetRequestMessage(),
			})
		}
	}
	return r
}

func (o *Operation) GetExampleResponses() []ExampleText {
	r := []ExampleText{}
	for _, p := range GetExampleProviders() {
		text := p.GetResponse(o)
		if len(text) > 0 {
			r = append(r, ExampleText{
				Tab:  p.GetTab(),
				Type: p.GetResponseType(),
				Text: p.GetResponse(o),
				Msg:  p.GetResponseMessage(),
			})
		}
	}
	return r
}

func (o *Operation) Description() string {
	return o.op.Description
}

func (a HttpResponses) Len() int      { return len(a) }
func (a HttpResponses) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a HttpResponses) Less(i, j int) bool { return a[i].Code < a[j].Code }

// VisitOperations calls fn once for each operation found in the collection of Documents
// VisitOperations calls fn once for each operation found in the collection of Documents
func VisitOperations(specs []*loads.Document, fn func(operation Operation)) {
	for _, d := range specs {
		for path, item := range d.Spec().Paths.Paths {
			for method, operation := range getOperationsForItem(item) {
				if operation != nil && !IsBlacklistedOperation(operation) {
					fn(Operation{
						item:       item,
						op:         operation,
						Path:       path,
						HttpMethod: method,
						ID:         operation.ID,
					})
				}
			}
		}
	}
}

func IsBlacklistedOperation(o *spec.Operation) bool {
	return strings.HasSuffix(o.ID, "APIGroup") || // These are just the API group meta datas.  Ignore for now.
		strings.HasSuffix(o.ID, "APIResources") || // These are just the API group meta datas.  Ignore for now.
		strings.HasSuffix(o.ID, "APIVersions") // || // These are just the API group meta datas.  Ignore for now.
		//strings.HasPrefix(o.ID, "connect") || // Skip pod connect apis for now.  There are too many.
		//strings.HasPrefix(o.ID, "proxy")
}

// Get all operations from the pathitem so we cacn iterate over them
func getOperationsForItem(pathItem spec.PathItem) map[string]*spec.Operation {
	return map[string]*spec.Operation{
		"GET":    pathItem.Get,
		"DELETE": pathItem.Delete,
		"PATCH":  pathItem.Patch,
		"PUT":    pathItem.Put,
		"POST":   pathItem.Post,
		"HEAD":   pathItem.Head,
	}
}

func (o *Operation) GetDisplayHttp() string {
	return fmt.Sprintf("%s %s", o.HttpMethod, o.Path)
}

func (o *Operation) VerifyBlackListed() {
	switch {
	case strings.Contains(o.ID, "connectCoreV1Patch"):
	case strings.Contains(o.ID, "createCoreV1NamespacedPodBinding"):
	case strings.Contains(o.ID, "getCodeVersion"):
	case strings.Contains(o.ID, "logFileHandler"):
	case strings.Contains(o.ID, "logFileListHandler"):
	case strings.Contains(o.ID, "NamespacedPodAttach"):
	case strings.Contains(o.ID, "NamespacedPodExec"):
	case strings.Contains(o.ID, "replaceCoreV1NamespaceFinalize"):
	case strings.Contains(o.ID, "V1beta1CertificateSigningRequestApproval"):
	case strings.Contains(o.ID, "V1beta1NamespacedReplicationControllerDummyScale"):
	default:
		fmt.Printf("No Definition found for %s [%s].  \n", o.ID, o.Path)
	}
}

func (o *Operation) GetMethod() string {
	switch o.HttpMethod {
	case "GET":
		return "List"
	case "POST":
		return "Create"
	case "PATCH":
		return "Patch"
	case "DELETE":
		return "Delete"
	case "PUT":
		return "Update"
	}
	return ""
}

// /apis/<group>/<version>/namespaces/{namespace}/<resources>/{name}/<subresource>
var matchNamespaced = regexp.MustCompile(
	`^/apis/([A-Za-z0-9\.]+)/([A-Za-z0-9]+)/namespaces/\{namespace\}/([A-Za-z0-9\.]+)/\{name\}/([A-Za-z0-9\.]+)$`)
var matchUnnamespaced = regexp.MustCompile(
	`^/apis/([A-Za-z0-9\.]+)/([A-Za-z0-9]+)/([A-Za-z0-9\.]+)/\{name\}/([A-Za-z0-9\.]+)$`)

func (o *Operation) GetGroupVersionKindSub() (string, string, string, string) {
	if matchNamespaced.MatchString(o.Path) {
		m := matchNamespaced.FindStringSubmatch(o.Path)
		return strings.Split(m[1], ".")[0], m[2], m[3], m[4]
	} else if matchUnnamespaced.MatchString(o.Path) {
		m := matchUnnamespaced.FindStringSubmatch(o.Path)
		return m[1], m[2], m[3], m[4]
	}
	return "", "", "", ""
}

// initExample reads the example config for an operation
func (o *Operation) initExample(config *Config) {
	path := o.Type.Name + ".yaml"
	path = filepath.Join(*ConfigDir, config.ExampleLocation, o.Definition.Name, path)
	path = strings.Replace(path, " ", "_", -1)
	path = strings.ToLower(path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(content, &o.ExampleConfig)
	if err != nil {
		panic(fmt.Sprintf("Could not Unmarshal ExampleConfig yaml: %s\n", content))
	}
}
