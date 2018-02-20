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

package initproject

import (
	"github.com/kubernetes-sigs/kubebuilder/cmd/kubebuilder/util"
)

// createBazelWorkspace creates new WORKSPACE and BUILD.bazel files at the root
func createBazelWorkspace() {
	execute("WORKSPACE", "bazel-workspace-template", workspaceTemplate, nil)
	execute(
		"BUILD.bazel",
		"bazel-build-template",
		buildTemplate,
		buildTemplateArguments{util.Repo},
	)
}

type buildTemplateArguments struct {
	Repo string
}

var workspaceTemplate = `
http_archive(
    name = "io_bazel_rules_go",
    url = "https://github.com/bazelbuild/rules_go/releases/download/0.9.0/rules_go-0.9.0.tar.gz",
    sha256 = "4d8d6244320dd751590f9100cf39fd7a4b75cd901e1f3ffdfd6f048328883695",
)
load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")
go_rules_dependencies()
go_register_toolchains()
`

var buildTemplate = `
# gazelle:proto disable
# gazelle:exclude vendor/github.com/json-iterator/go/skip_tests
load("@io_bazel_rules_go//go:def.bzl", "gazelle")

gazelle(
    name = "gazelle",
    command = "fix",
    prefix = "{{.Repo}}",
    external = "vendored",
    args = [
        "-build_file_name",
        "BUILD.bazel",
    ],
)
`
