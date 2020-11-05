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

# from cmd/
#  bash -c "find vendor/sigs.k8s.io/apiserver-builder-alpha -name BUILD.bazel| xargs sed -i='' s'|//pkg|//vendor/sigs.k8s.io/apiserver-builder-alpha/pkg|g'"

# from /

gazelle:
	bazel run //:gazelle

NAME=apiserver-builder-alpha
VENDOR=kubernetes-sigs
VERSION=$(shell git describe --always --tags HEAD)
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
	go test ./pkg/... ./cmd/...

.PHONY: install
install: build
	@echo "Installing apiserver-builder suite tools..."
	@echo "GOOS: $(GOOS)"
	@echo "GOARCH: $(GOARCH)"
	@echo "ARCH: $(ARCH)"
	cp release/$(VERSION)/bin/* $(GOPATH)/bin/

.PHONY: clean
clean:
	rm -rf *.deb *.rpm *.tar.gz ./release

.PHONY: build
build: clean ## Create release artefacts for darwin:amd64, linux:amd64 and windows:amd64. Requires etcd, glide, hg.
	mkdir -p release/$(VERSION)/src
	bazel build --platforms=@io_bazel_rules_go//go/toolchain:$(GOOS)_$(GOARCH) cmd:apiserver-builder
	ls -lh bazel-bin/cmd
	cp bazel-bin/cmd/apiserver-builder.tar.gz apiserver-builder-alpha-$(VERSION)-$(GOOS)-$(GOARCH).tar.gz
	tar xzf apiserver-builder-alpha-$(VERSION)-$(GOOS)-$(GOARCH).tar.gz -C release/$(VERSION)

.PHONY: package
package: package-linux-amd64

.PHONY: package-linux-amd64
package-linux-amd64: package-linux-amd64-deb package-linux-amd64-rpm

.PHONY: package-linux-amd64-deb
package-linux-amd64-deb: ## Create a Debian package. Requires jordansissel/fpm.
	fpm --force --name '$(NAME)' --version '$(VERSION)' \
	  --input-type tar \
	  --output-type deb \
	  --vendor '$(VENDOR)' \
	  --description '$(DESCRIPTION)' \
	  --url '$(URL)' \
	  --maintainer '$(MAINTAINER)' \
	  --license '$(LICENSE)' \
	  --package $(NAME)-$(VERSION)-amd64.deb \
	  --prefix /usr/local/apiserver-builder \
	  $(NAME)-$(VERSION)-linux-amd64.tar.gz

.PHONY: package-linux-amd64-rpm
package-linux-amd64-rpm: ## Create an RPM package. Requires jordansissel/fpm, rpm.
	fpm --force --name '$(NAME)' --version '$(VERSION)' \
	  --input-type tar \
	  --output-type rpm \
	  --vendor '$(VENDOR)' \
	  --description '$(DESCRIPTION)' \
	  --url '$(URL)' \
	  --maintainer '$(MAINTAINER)' \
	  --license '$(LICENSE)' \
	  --rpm-os linux \
	  --package $(NAME)-$(VERSION)-amd64.rpm \
	  --prefix /usr/local/apiserver-builder \
	  $(NAME)-$(VERSION)-linux-amd64.tar.gz

gazelle-reset:
	bazel run \
		//:gazelle -- \
		update-repos \
		--from_file=go.mod \
		--to_macro=repos.bzl%go_repositories \
		--build_file_generation=on \
		--build_file_proto_mode=disable \
		--prune
