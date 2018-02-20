#!/bin/bash
set -e
set -x

apt update
apt install rsync -y
go get github.com/jteeuwen/go-bindata/go-bindata

git clone https://github.com/kubernetes/kubernetes $GOPATH/src/k8s.io/kubernetes --depth=1 -b release-1.9

export CGO=0
export KUBE_BUILD_PLATFORMS=$GOOS/$GOARCH
export DEST=/workspace/_output/kubebuilder/bin/
mkdir -p $DEST || echo ""

cd $GOPATH/src/k8s.io/kubernetes
make clean
WHAT=cmd/kube-apiserver make
WHAT=cmd/kube-controller-manager make
WHAT=cmd/kubectl make

cp _output/local/bin/$GOOS/$GOARCH/kube-apiserver $DEST
cp _output/local/bin/$GOOS/$GOARCH/kube-controller-manager $DEST
cp _output/local/bin/$GOOS/$GOARCH/kubectl $DEST

