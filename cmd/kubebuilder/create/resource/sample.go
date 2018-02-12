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

package resource

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	createutil "github.com/najena/kubebuilder/cmd/kubebuilder/create/util"
	"github.com/najena/kubebuilder/cmd/kubebuilder/util"
)

func doSample(dir string, args resourceTemplateArgs) bool {
	os.MkdirAll("sample", 0700)
	samplepath := filepath.Join("sample", fmt.Sprintf("%s.yaml", strings.ToLower(createutil.KindName)))
	return util.WriteIfNotFound(samplepath, "sample-template", sampleTemplate, args)
}

var sampleTemplate = `apiVersion: {{ .Group }}.{{ .Domain }}/{{ .Version }}
kind: {{ .Kind }}
metadata:
  name: {{ lower .Kind }}-example
spec:
`
