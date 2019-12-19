#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

apiserver_builder_version="v1.16.alpha.0"
k8s_vendor="v1.16.0"
git_commit="$(git rev-parse HEAD)"
build_date="$(date +%Y-%m-%d-%H:%M:%S)"

cat <<EOF
APISERVER_BUILDER_VERSION ${apiserver_builder_version}
GIT_COMMIT ${git_commit}
BUILD_DATE ${build_date}
K8S_VENDOR ${k8s_vendor}
EOF
