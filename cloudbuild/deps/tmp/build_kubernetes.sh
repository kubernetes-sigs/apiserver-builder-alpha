#!/bin/bash
set -e
set -x

export CGO=0
export KUBE_BUILD_PLATFORMS=$GOOS/$GOARCH
export DEST=/user/local/kubebuilder/bin/
mkdir -p $DEST || echo ""

cd $GOPATH/src/k8s.io/kubernetes
make clean
mkdir -p  _output/bin
chmod +777 -R _output
WHAT=cmd/kube-apiserver make
WHAT=cmd/kube-controller-manager make
WHAT=cmd/kubectl make

#cp _output/local/bin/$GOOS/$GOARCH/kube-apiserver $DEST
#cp _output/local/bin/$GOOS/$GOARCH/kube-controller-manager $DEST
#cp _output/local/bin/$GOOS/$GOARCH/kubectl $DEST

