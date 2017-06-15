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
var goos string
var goarch string
var outputdir string

var createBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds the source into executables",
	Long:  `Builds the source into executables`,
	Run:   RunBuild,
}

func AddBuild(cmd *cobra.Command) {
	cmd.AddCommand(createBuildCmd)

	createBuildCmd.Flags().BoolVar(&generateForBuild, "generate", true, "if true, generate code before building")
	createBuildCmd.Flags().StringVar(&goos, "goos", "", "if set, use this GOOS")
	createBuildCmd.Flags().StringVar(&goarch, "goarch", "", "if set, use this GOARCH")
	createBuildCmd.Flags().StringVar(&outputdir, "outputdir", "bin", "if set, write the binary to this directory")
}

func RunBuild(cmd *cobra.Command, args []string) {
	if generateForBuild {
		log.Printf("regenerating generated code.  To disable regeneration, run with --generate=false.")
		RunGenerate(cmd, args)
	}

	// Build the apiserver
	path := filepath.Join("cmd", "apiserver", "main.go")
	c := exec.Command("go", "build", "-o", filepath.Join(outputdir, "apiserver"), path)
	c.Env = append(os.Environ(), "CGO_ENABLED=0")
	if len(goos) > 0 {
		c.Env = append(c.Env, fmt.Sprintf("GOOS=%s", goos))
	}
	if len(goarch) > 0 {
		c.Env = append(c.Env, fmt.Sprintf("GOARCH=%s", goarch))
	}

	fmt.Printf("%s\n", strings.Join(c.Args, " "))
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	err := c.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Build the controller manager
	path = filepath.Join("cmd", "controller", "main.go")
	c = exec.Command("go", "build", "-o", filepath.Join(outputdir, "controller-manager"), path)
	c.Env = append(os.Environ(), "CGO_ENABLED=0")
	if len(goos) > 0 {
		c.Env = append(c.Env, fmt.Sprintf("GOOS=%s", goos))
	}
	if len(goarch) > 0 {
		c.Env = append(c.Env, fmt.Sprintf("GOARCH=%s", goarch))
	}

	fmt.Println(strings.Join(c.Args, " "))
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	err = c.Run()
	if err != nil {
		log.Fatal(err)
	}
}
