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

package boot

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

var groupName string
var kindName string
var resourceName string
var versionName string
var copyright string
var domain string
var Repo string
var GoSrc string

// writeIfNotFound returns true if the file was created and false if it already exists
func writeIfNotFound(path, templateName, templateValue string, data interface{}) bool {
	// Make sure the directory exists
	exec.Command("mkdir", "-p", filepath.Dir(path)).CombinedOutput()

	// Don't create the doc.go if it exists
	if _, err := os.Stat(path); err == nil {
		return false
	} else if !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Could not stat %s: %v\n", path, err)
		os.Exit(-1)

	}
	create(path)

	t := template.Must(template.New(templateName).Parse(templateValue))

	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create %s: %v\n", path, err)
		os.Exit(-1)
	}
	defer f.Close()

	err = t.Execute(f, data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create %s: %v\n", path, err)
		os.Exit(-1)
	}

	return true
}

func create(path string) {
	f, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
}
