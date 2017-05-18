# Copyright 2017 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

all: test

test: build
	go test github.com/kubernetes-incubator/apiserver-builder/example/pkg/apis/miskatonic/v1beta1
	go test github.com/kubernetes-incubator/apiserver-builder/example/pkg/apis/innsmouth/v1
	go test github.com/kubernetes-incubator/apiserver-builder/example/pkg/controller/university
	bash -c "find pkg/apis/ -name apiserver.local.config | xargs rm -rf"

build: generate
	go build -o bin/apiserver cmd/apiserver/main.go
	go build -o bin/controller-manager cmd/controller-manager/main.go

# Generate code
generate: cmds cleangenerated
	bin/apiserver-boot generate

# Build docs
docs: cleandocs build
	apiserver-boot docs --server bin/apiserver

# Build the tools
# Note: These are available in the apiserver-builder release tar.gz
cmds:
	./build-tools.sh

# Clean up functions
clean: cleangenerated cleandocs

cleangenerated: cmds
	bin/apiserver-boot generate clean

cleandocs: cmds
	bin/apiserver-boot docs clean
