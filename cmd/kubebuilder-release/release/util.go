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

package release

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func RunCmd(cmd *exec.Cmd, gopath string) {
	setgopath := len(gopath) > 0
	gopath, err := filepath.Abs(gopath)
	if err != nil {
		log.Fatal(err)
	}
	gopath, err = filepath.EvalSymlinks(gopath)
	if err != nil {
		log.Fatal(err)
	}
	if setgopath {
		cmd.Env = append(cmd.Env, fmt.Sprintf("GOPATH=%s", gopath))
	}
	for _, v := range os.Environ() {
		if strings.HasPrefix(v, "GOPATH=") && setgopath {
			continue
		}
		cmd.Env = append(cmd.Env, v)
	}

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if len(cmd.Dir) == 0 && len(gopath) > 0 {
		cmd.Dir = gopath
	}
	fmt.Printf("%s\n", strings.Join(cmd.Args, " "))
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
