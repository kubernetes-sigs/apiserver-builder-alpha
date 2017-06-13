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
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var generateForBuild bool

var createBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Creates an API resource",
	Long:  `Creates an API resource`,
	Run:   RunBuild,
}

func AddBuild(cmd *cobra.Command) {
	cmd.AddCommand(createBuildCmd)

	createBuildCmd.Flags().BoolVar(&generateForBuild, "generate", true, "if true, generate code before building")
}

func RunBuild(cmd *cobra.Command, args []string) {
	if generateForBuild {
		RunGenerate(cmd, args)
	}

	path := filepath.Join("cmd", "apiserver", "main.go")
	c := exec.Command("go", "build", "-o", "bin/apiserver", path)
	fmt.Printf("%s\n", strings.Join(c.Args, " "))
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	err := c.Run()
	if err != nil {
		log.Fatal(err)
	}

	path = filepath.Join("cmd", "controller", "main.go")
	c = exec.Command("go", "build", "-o", "bin/controller-manager", path)
	fmt.Println(strings.Join(c.Args, " "))
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	err = c.Run()
	if err != nil {
		log.Fatal(err)
	}
}
