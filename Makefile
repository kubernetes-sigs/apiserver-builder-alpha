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

NAME=apiserver-builder-alpha
VENDOR=kubernetes-sigs
VERSION=$(shell git describe --always --tags HEAD)
COMMIT=$(shell git rev-parse --short HEAD)
KUBE_VERSION?=v0.23.5
DESCRIPTION=apiserver-builder implements libraries and tools to quickly and easily build Kubernetes apiservers to support custom resource types.
MAINTAINER=The Kubernetes Authors
URL=https://github.com/$(VENDOR)/$(NAME)
LICENSE=Apache-2.0
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)
GOPATH?=$(shell go env GOPATH)

.PHONY: default
default: install

.PHONY: test
test:
	go test ./cmd/... ./pkg/...

.PHONY: install
install: build
	@echo "Installing apiserver-builder suite tools..."
	@echo "GOOS: $(GOOS)"
	@echo "GOARCH: $(GOARCH)"
	@echo "ARCH: $(ARCH)"
	go install ./cmd/apiserver-boot

.PHONY: clean
clean:
	rm -rf *.deb *.rpm *.tar.gz ./release bin/*

.PHONY: build
build: clean ## Create release artefacts for darwin:amd64, linux:amd64 and windows:amd64. Requires etcd, glide, hg.
	mkdir -p bin
	go build -o bin/apiserver-boot ./cmd/apiserver-boot

release-binary:
	mkdir -p bin
	go build \
		-ldflags=" \
			-X 'sigs.k8s.io/apiserver-builder-alpha/cmd/apiserver-boot/boot/version.goos=${GOOS}' \
			-X 'sigs.k8s.io/apiserver-builder-alpha/cmd/apiserver-boot/boot/version.goarch=${GOARCH}' \
			-X 'sigs.k8s.io/apiserver-builder-alpha/cmd/apiserver-boot/boot/version.kubernetesVendorVersion=${KUBE_VERSION}' \
			-X 'sigs.k8s.io/apiserver-builder-alpha/cmd/apiserver-boot/boot/version.apiserverBuilderVersion=${VERSION}' \
			-X 'sigs.k8s.io/apiserver-builder-alpha/cmd/apiserver-boot/boot/version.gitCommit=${COMMIT}' \
			" \
 		-o bin/apiserver-boot ./cmd/apiserver-boot
	tar czvf apiserver-boot-${GOOS}-${GOARCH}.tar.gz bin/apiserver-boot

