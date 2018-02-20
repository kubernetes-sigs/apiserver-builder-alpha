#!/bin/bash
set -e
set -x

mkdir -p /workspace/vendor/github.com/kubernetes-sigs/kubebuilder/pkg/ || echo ""
cp -r /workspace/pkg/* /workspace/vendor/github.com/kubernetes-sigs/kubebuilder/pkg/

export DEST=/workspace/_output/kubebuilder/bin/
export DEST=/workspace/_output/kubebuilder/bin/
mkdir -p $DEST || echo ""
tar -zcvf $DEST/vendor.tar.gz vendor/
