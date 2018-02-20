#!/bin/bash
set -e
set -x

mkdir -p /go/src/github.com/kubernetes-sigs
ln -s $(pwd) /go/src/github.com/kubernetes-sigs/kubebuilder

export CGO=0
export DEST=/workspace/_output/kubebuilder/bin/
mkdir -p $DEST || echo ""

export X=github.com/kubernetes-sigs/kubebuilder/cmd/kubebuilder/version

go build -o $DEST/kubebuilder \
 -ldflags "-X $X.kubeBuilderVersion=$VERSION -X $X.goos=$GOOS -X $X.goarch=$GOARCH -X $X.kubernetesVendorVersion=$KUBERNETES_VERSION" \
 github.com/kubernetes-sigs/kubebuilder/cmd/kubebuilder

go build -o $DEST/kubebuilder-gen github.com/kubernetes-sigs/kubebuilder/cmd/kubebuilder-gen
