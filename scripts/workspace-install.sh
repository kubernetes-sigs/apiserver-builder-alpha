#!/usr/bin/env bash

[[ -f bazel_0.29.1-linux-x86_64.deb ]] || wget https://github.com/bazelbuild/bazel/releases/download/0.29.1/bazel_0.29.1-linux-x86_64.deb

# bazel installation
sudo dpkg -i bazel_0.29.1-linux-x86_64.deb

# install apiserver-boots
mkdir -p $(go env GOPATH)/bin/
