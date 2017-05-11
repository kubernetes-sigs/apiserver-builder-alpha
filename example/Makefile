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

# NOTE: Project must live under $GOPATH/src/github.com/kubernetes-incubator/apiserver-builder/example for the code to compile
REPO=github.com/kubernetes-incubator/apiserver-builder/example
BOILER_PLATE="boilerplate.go.txt"
UNVERSIONED="miskatonic/,innsmouth/"
VERSIONED="miskatonic/v1beta1,innsmouth/v1"
UNVERSIONED_DIRS="--input-dirs=github.com/kubernetes-incubator/apiserver-builder/example/pkg/apis/miskatonic --input-dirs=github.com/kubernetes-incubator/apiserver-builder/example/pkg/apis/innsmouth"
VERSIONED_DIRS="--input-dirs=github.com/kubernetes-incubator/apiserver-builder/example/pkg/apis/miskatonic/v1beta1 --input-dirs=github.com/kubernetes-incubator/apiserver-builder/example/pkg/apis/innsmouth/v1"

REGISTER_GEN="apiregister-gen"
REGISTER_GEN="../cmd/apiregister-gen/main.go"

SRC="$(UNVERSIONED_DIRS) $(VERSIONED_DIRS)"
DIR=pkg

GENERIC_API="k8s.io/apimachinery/pkg/apis/meta/v1,k8s.io/apimachinery/pkg/api/resource/,k8s.io/apimachinery/pkg/version/,k8s.io/apimachinery/pkg/runtime/,k8s.io/apimachinery/pkg//util/intstr/"
OUT=$(GOPATH)/src
CLIENT_PKG=$(REPO)/pkg/client
CLIENT_PATH=$(CLIENT_PKG)/clientset_generated
CLIENT=clientset
INTERNAL_CLIENT=internalclientset
LISTERS_PKG="$(CLIENT_PKG)/listers_generated"
INFORMERS_PKG="$(CLIENT_PKG)/informers_generated"

all: test

# Run integration tests
test: build
	go test github.com/kubernetes-incubator/apiserver-builder/example/pkg/apis/miskatonic/v1beta1
	go test github.com/kubernetes-incubator/apiserver-builder/example/pkg/apis/innsmouth/v1
	bash -c "find $(DIR)/apis/ -name apiserver.local.config | xargs rm -rf"

# Build the main binary
build: generate
	go build -o bin/apiserver main.go

# Generate code
generate: cmds cleangenerated
	bin/apiserver-boot generate

# Build docs
docs: cleandocs build
	apiserver-boot docs --server bin/apiserver

# Build the tools needed by apiserver-boot
# Note: These are available in the release tar
cmds: deps
	go build -o bin/apiserver-boot ../cmd/apiserver-boot/main.go
	go build -o bin/apiregister-gen ../cmd/apiregister-gen/main.go
	go build -o bin/client-gen k8s.io/kubernetes/cmd/libs/go2idl/client-gen
	go build -o bin/conversion-gen k8s.io/kubernetes/cmd/libs/go2idl/conversion-gen
	go build -o bin/deepcopy-gen k8s.io/kubernetes/cmd/libs/go2idl/deepcopy-gen
	go build -o bin/openapi-gen k8s.io/kubernetes/cmd/libs/go2idl/openapi-gen
	go build -o bin/defaulter-gen k8s.io/kubernetes/cmd/libs/go2idl/defaulter-gen
	go build -o bin/lister-gen k8s.io/kubernetes/cmd/libs/go2idl/lister-gen
	go build -o bin/informer-gen k8s.io/kubernetes/cmd/libs/go2idl/informer-gen
	go build -o bin/apiregister-gen github.com/kubernetes-incubator/apiserver-builder/cmd/apiregister-gen
	go build -o bin/gen-apidocs github.com/kubernetes-incubator/reference-docs/gen-apidocs

# Download all of the code generators
# NOTE: These are available in the release tar
deps:
	go get k8s.io/kubernetes/cmd/libs/go2idl/client-gen
	go get k8s.io/kubernetes/cmd/libs/go2idl/conversion-gen
	go get k8s.io/kubernetes/cmd/libs/go2idl/deepcopy-gen
	go get k8s.io/kubernetes/cmd/libs/go2idl/openapi-gen
	go get k8s.io/kubernetes/cmd/libs/go2idl/defaulter-gen
	go get k8s.io/kubernetes/cmd/libs/go2idl/lister-gen
	go get k8s.io/kubernetes/cmd/libs/go2idl/informer-gen
	go get github.com/kubernetes-incubator/apiserver-builder/cmd/apiregister-gen
	go get github.com/kubernetes-incubator/reference-docs/gen-apidocs

# Clean up functions

clean: cleanbin cleangenerated cleandocs

cleanbin:
	rm -f main

cleangenerated: cmds
	bin/apiserver-boot generate clean

cleandocs: cmds
	bin/apiserver-boot docs clean
