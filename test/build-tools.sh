#!/usr/bin/env bash

set -x -e

# NOTE: Do not copy this file unless you need to use apiserver-builder at HEAD.
# Otherwise, download the pre-built apiserver-builder tar release from
# https://sigs.k8s.io/apiserver-builder-alpha/releases instead.

(
    mkdir -p /home/travis/gopath/src/github.com/Masterminds
	cd /home/travis/gopath/src/github.com/Masterminds
	git clone https://github.com/Masterminds/glide.git
	cd glide
	make build
)

export PATH=/home/travis/gopath/src/github.com/Masterminds/glide:$PATH

# Install generators from this repo
cd ..
go build -o bin/apiserver-builder-release cmd/apiserver-builder-release/main.go
./bin/apiserver-builder-release vendor --version 1.0
./bin/apiserver-builder-release build --version 1.0 --targets linux:amd64

tar -xzf apiserver-builder-alpha-1.0-linux-amd64.tar.gz -C test
