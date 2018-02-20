#!/bin/bash
set -e
set -x

git clone https://github.com/kubernetes/code-generator $GOPATH/src/k8s.io/code-generator --depth=1 -b release-1.9

export CGO=0
export DEST=/workspace/_output/kubebuilder/bin/
mkdir -p $DEST || echo ""

go build -o $DEST/client-gen k8s.io/code-generator/cmd/client-gen
go build -o $DEST/conversion-gen k8s.io/code-generator/cmd/conversion-gen
go build -o $DEST/deepcopy-gen k8s.io/code-generator/cmd/deepcopy-gen
go build -o $DEST/defaulter-gen k8s.io/code-generator/cmd/defaulter-gen
go build -o $DEST/informer-gen k8s.io/code-generator/cmd/informer-gen
go build -o $DEST/lister-gen k8s.io/code-generator/cmd/lister-gen
go build -o $DEST/openapi-gen k8s.io/code-generator/cmd/openapi-gen
