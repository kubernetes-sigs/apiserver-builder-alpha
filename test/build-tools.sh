#!/usr/bin/env bash

set -x -e

# NOTE: Do not copy this file unless you need to use apiserver-builder at HEAD.
# Otherwise, download the pre-built apiserver-builder tar release from
# https://sigs.k8s.io/apiserver-builder-alpha/releases instead.

if [[ -n "${BUILD_WORKSPACE_DIRECTORY:-}" ]]; then # Running inside bazel
  echo "Updating codegen files..." >&2
elif ! command -v bazel &>/dev/null; then
  echo "Install bazel at https://bazel.build" >&2
  exit 1
else
  (
    set -o xtrace
    bazel run //test:build-tools
  )
  exit 0
fi

out_dir=$BUILD_WORKSPACE_DIRECTORY/$(dirname "$1")
out_tar=$BUILD_WORKSPACE_DIRECTORY/$2

tar -zxvf "$out_tar" -C "$out_dir"
