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

NAME=apiserver-builder
VENDOR=kubernetes-incubator
VERSION=$(shell cat VERSION)
DESCRIPTION=apiserver-builder implements libraries and tools to quickly and easily build Kubernetes apiservers to support custom resource types.
MAINTAINER=The Kubernetes Authors
URL=https://github.com/$(VENDOR)/$(NAME)
LICENSE=Apache-2.0

BUILD_DIR=$(shell pwd)/build
DARWIN_AMD64_BUILD_BIN_DIR=$(BUILD_DIR)/darwin-amd64/bin
LINUX_AMD64_BUILD_BIN_DIR=$(BUILD_DIR)/linux-amd64/bin
LINUX_AMD64_BUILD_PKG_DIR=$(BUILD_DIR)/linux-amd64/pkg
WINDOWS_AMD64_BUILD_BIN_DIR=$(BUILD_DIR)/windows-amd64/bin

.PHONY: default
default: install

.PHONY: test
test:
	go test ./pkg/... ./cmd/...

.PHONY: install
install:
	go install -v ./pkg/... ./cmd/...

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

.PHONY: build
build: clean build-darwin-amd64 build-linux-amd64 build-windows-amd64

.PHONY: build-darwin-amd64
build-darwin-amd64: build-apiregister-gen-darwin-amd64 build-apiserver-boot-darwin-amd64 build-apiserver-builder-release-darwin-amd64

.PHONY: build-apiregister-gen-darwin-amd64 build-apiserver-boot-darwin-amd64 build-apiserver-builder-release-darwin-amd64
build-apiregister-gen-darwin-amd64 build-apiserver-boot-darwin-amd64 build-apiserver-builder-release-darwin-amd64: build-%-darwin-amd64:
	mkdir -p $(DARWIN_AMD64_BUILD_BIN_DIR)
	GOOS=darwin GOARCH=amd64 go build -o $(DARWIN_AMD64_BUILD_BIN_DIR)/$* ./cmd/$*/main.go

.PHONY: build-linux-amd64
build-linux-amd64: build-apiregister-gen-linux-amd64 build-apiserver-boot-linux-amd64 build-apiserver-builder-release-linux-amd64

.PHONY: build-apiregister-gen-linux-amd64 build-apiserver-boot-linux-amd64 build-apiserver-builder-release-linux-amd64
build-apiregister-gen-linux-amd64 build-apiserver-boot-linux-amd64 build-apiserver-builder-release-linux-amd64: build-%-linux-amd64:
	mkdir -p $(LINUX_AMD64_BUILD_BIN_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(LINUX_AMD64_BUILD_BIN_DIR)/$* ./cmd/$*/main.go

.PHONY: build-windows-amd64
build-darwin-amd64: build-apiregister-gen-windows-amd64 build-apiserver-boot-windows-amd64 build-apiserver-builder-release-windows-amd64

.PHONY: build-apiregister-gen-windows-amd64 build-apiserver-boot-windows-amd64 build-apiserver-builder-release-windows-amd64
build-apiregister-gen-windows-amd64 build-apiserver-boot-windows-amd64 build-apiserver-builder-release-windows-amd64: build-%-windows-amd64:
	mkdir -p $(WINDOWS_AMD64_BUILD_BIN_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(WINDOWS_AMD64_BUILD_BIN_DIR)/$* ./cmd/$*/main.go

.PHONY: package
package: build package-linux-amd64

.PHONY: package-linux-amd64
package-linux-amd64: build-linux-amd64 package-linux-amd64-deb package-linux-amd64-rpm

.PHONY: package-linux-amd64-deb
package-linux-amd64-deb: ## Build a Debian package. Requires jordansissel/fpm.
	mkdir -p $(LINUX_AMD64_BUILD_PKG_DIR)

	fpm --name '$(NAME)' --version '$(VERSION)' \
	  --input-type dir \
	  --output-type deb \
	  --vendor '$(VENDOR)' \
	  --description '$(DESCRIPTION)' \
	  --url '$(URL)' \
	  --maintainer '$(MAINTAINER)' \
	  --license '$(LICENSE)' \
	  --package $(LINUX_AMD64_BUILD_PKG_DIR)/$(NAME)_$(VERSION)_amd64.deb \
	  $(LINUX_AMD64_BUILD_BIN_DIR)/=/usr/local/bin

.PHONY: package-linux-amd64-rpm
package-linux-amd64-rpm: ## Build a Debian package. Requires jordansissel/fpm and rpmbuild.
	mkdir -p $(LINUX_AMD64_BUILD_PKG_DIR)

	fpm --name '$(NAME)' --version '$(VERSION)' \
	  --input-type dir \
	  --output-type rpm \
	  --vendor '$(VENDOR)' \
	  --description '$(DESCRIPTION)' \
	  --url '$(URL)' \
	  --maintainer '$(MAINTAINER)' \
	  --license '$(LICENSE)' \
	  --rpm-os linux \
	  --package $(LINUX_AMD64_BUILD_PKG_DIR)/$(NAME)_$(VERSION)_amd64.rpm \
	  $(LINUX_AMD64_BUILD_BIN_DIR)/=/usr/local/bin
