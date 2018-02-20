#!/bin/bash
set -e
set -x

git clone https://github.com/coreos/etcd $GOPATH/src/github.com/coreos/etcd --depth=1

export CGO=0
export DEST=/workspace/_output/kubebuilder/bin/
mkdir -p $DEST || echo ""

go build -o $DEST/etcd github.com/coreos/etcd
