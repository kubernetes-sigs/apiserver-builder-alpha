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

package util

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/markbates/inflect"
	"k8s.io/apiserver/pkg/server"
	"k8s.io/klog/v2"
)

var Domain string

// writeIfNotFound returns true if the file was created and false if it already exists
func WriteIfNotFound(path, templateName, templateValue string, data interface{}) bool {
	// Make sure the directory exists
	os.MkdirAll(filepath.Dir(path), 0700)

	// Don't create the doc.go if it exists
	if _, err := os.Stat(path); err == nil {
		return false
	} else if !os.IsNotExist(err) {
		klog.Fatalf("Could not stat %s: %v", path, err)

	}
	create(path)

	t := template.Must(template.New(templateName).Funcs(
		template.FuncMap{
			"title":  strings.Title,
			"lower":  strings.ToLower,
			"plural": inflect.NewDefaultRuleset().Pluralize,
		},
	).Parse(templateValue))

	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		klog.Fatalf("Failed to create %s: %v", path, err)
	}
	defer f.Close()

	err = t.Execute(f, data)
	if err != nil {
		klog.Fatalf("Failed to create %s: %v", path, err)
	}

	return true
}

// Overwrite always updates the target file with the new content.
func Overwrite(path, templateName, templateValue string, data interface{}) bool {
	// Make sure the directory exists
	os.MkdirAll(filepath.Dir(path), 0700)

	create(path)
	t := template.Must(template.New(templateName).Funcs(
		template.FuncMap{
			"title":  strings.Title,
			"lower":  strings.ToLower,
			"plural": inflect.NewDefaultRuleset().Pluralize,
		},
	).Parse(templateValue))

	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		klog.Fatalf("Failed to create %s: %v", path, err)
	}
	defer f.Close()

	err = t.Execute(f, data)
	if err != nil {
		klog.Fatalf("Failed to create %s: %v", path, err)
	}

	return true
}

func GetCopyright(file string) string {
	if len(file) == 0 {
		wd, err := os.Getwd()
		if err != nil {
			klog.Fatal(err)
		}
		file = filepath.Join(wd, "hack", "boilerplate.go.txt")
	}
	cr, err := ioutil.ReadFile(file)
	if err != nil {
		klog.Fatalf("Must create boilerplate.go.txt file with copyright and file headers: %v", err)
	}
	return string(cr)
}

func GetDomain() string {
	b, err := ioutil.ReadFile(filepath.Join("pkg", "apis", "doc.go"))
	if err != nil {
		klog.Fatalf("Could not find pkg/apis/doc.go.  First run `apiserver-boot init --domain <domain>`.")
	}
	r := regexp.MustCompile("\\+domain=(.*)")
	l := r.FindSubmatch(b)
	if len(l) < 2 {
		klog.Fatalf("pkg/apis/doc.go does not contain the domain (// +domain=.*)")
	}
	Domain = string(l[1])
	return Domain
}

func DoCmdWithInfo(s string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", s)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		klog.Fatalf("command failed %v", err)
	}
	return out.String(), err
}

func CompareVersion(version1 string, version2 string) int {
	var res int
	ver1Strs := strings.Split(version1, ".")
	ver2Strs := strings.Split(version2, ".")
	ver1Len := len(ver1Strs)
	ver2Len := len(ver2Strs)
	verLen := ver1Len
	if len(ver1Strs) < len(ver2Strs) {
		verLen = ver2Len
	}
	for i := 0; i < verLen; i++ {
		var ver1Int, ver2Int int
		if i < ver1Len {
			ver1Int, _ = strconv.Atoi(ver1Strs[i])
		}
		if i < ver2Len {
			ver2Int, _ = strconv.Atoi(ver2Strs[i])
		}
		if ver1Int < ver2Int {
			res = -1
			break
		}
		if ver1Int > ver2Int {
			res = 1
			break
		}
	}
	return res
}

func create(path string) {
	f, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
}

func DoCmd(cmd string, args ...string) {
	c := exec.Command(cmd, args...)
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	klog.Infof("%s", strings.Join(c.Args, " "))
	err := c.Run()
	if err != nil {
		klog.Fatalf("command failed %v", err)
	}
}

func CheckInstall() {
	//bins := []string{"apiregister-gen", "client-gen", "deepcopy-gen", "gen-apidocs", "informer-gen",
	//	"openapi-gen", "apiserver-boot", "conversion-gen", "defaulter-gen", "lister-gen"}
	bins := []string{}
	missing := []string{}

	e, err := os.Executable()
	if err != nil {
		klog.Fatal("unable to get directory of apiserver-builder tools")
	}

	dir := filepath.Dir(e)
	for _, b := range bins {
		_, err = os.Stat(filepath.Join(dir, b))
		if err != nil {
			missing = append(missing, b)
		}
	}
	if len(missing) > 0 {
		klog.Fatalf("Error running apiserver-boot."+
			"\nThe following files are missing [%s]"+
			"\napiserver-boot must be installed using a release tar.gz downloaded from the git repo.",
			strings.Join(missing, ","))
	}
}

func CancelWhenSignaled(parent context.Context) context.Context {
	ctx, cancel := context.WithCancel(parent)

	go func() {
		signalChannel := server.SetupSignalHandler()
		<-signalChannel
		cancel()
	}()

	return ctx
}
