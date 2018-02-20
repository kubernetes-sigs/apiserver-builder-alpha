#!/bin/bash
set -e
set -x

cd /workspace/_output/
tar -zcvf /workspace/kubebuilder-$VERSION-$GOOS-$GOARCH.tar.gz kubebuilder