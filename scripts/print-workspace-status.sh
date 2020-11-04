#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

apiserver_builder_version=$(git rev-parse --verify --short HEAD)
k8s_vendor=kubernetes-1.19.2
git_commit="$(git rev-parse HEAD)"
build_date="$(date +%Y-%m-%d-%H:%M:%S)"

cat <<EOF
APISERVER_BUILDER_VERSION ${apiserver_builder_version}
GIT_COMMIT ${git_commit}
BUILD_DATE ${build_date}
K8S_VENDOR ${k8s_vendor}
EOF
