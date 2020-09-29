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

.PHONY: all test build generate docs cmds clean cleangenerated cleandocs skeleton

all: test

NON_INTERACTIVE_FLAG=--skip-resource=false --skip-controller=false --skip-admission-controller=false

test: build check
	go test ./pkg/...
	bash -c "find pkg/apis/ -name apiserver.local.config | xargs rm -rf"

check: build
	go vet $$(go list ./... | grep -vE '(clientset|listers|informers)_generated')

skeleton:
	apiserver-boot init repo --domain sample.kubernetes.io
	apiserver-boot create group version resource --group storage --version v1 --kind VolumeClaim $(NON_INTERACTIVE_FLAG)
	apiserver-boot create group version resource --group storage --version v1 --kind SnapshotClaim $(NON_INTERACTIVE_FLAG)
	apiserver-boot create group version resource --group storage --version v1 --kind Volume --non-namespaced $(NON_INTERACTIVE_FLAG)
	apiserver-boot create group version resource --group storage --version v1 --kind Snapshot --non-namespaced $(NON_INTERACTIVE_FLAG)

build: cmds skeleton
	GO111MODULE=off apiserver-boot build executables

# Build docs
docs: cleandocs build
	apiserver-boot build docs --server bin/apiserver

# Clean up functions
clean: cleangenerated cleandocs

cleangenerated: cmds
	apiserver-boot build generated clean

cleandocs: cmds
	apiserver-boot build docs clean
