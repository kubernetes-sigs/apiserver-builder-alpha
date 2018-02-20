#!/bin/bash
set -e
set -x

git clone https://github.com/kubernetes-incubator/reference-docs $GOPATH/src/github.com/kubernetes-incubator/reference-docs --depth=1

export CGO=0
export DEST=/workspace/_output/kubebuilder/bin/
mkdir -p $DEST || echo ""

go build -o $DEST/gen-apidocs github.com/kubernetes-incubator/reference-docs/gen-apidocs
