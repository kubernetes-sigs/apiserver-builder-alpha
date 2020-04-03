#!/usr/bin/env bash

set -x -e

if [ "$TEST" == "example" ]; then
	cd example/basic
	PATH=$GOPATH/bin:$PATH:/tmp/test-etcd make test
elif [ "$TEST" == "test" ]; then
	cd test
	PATH=$GOPATH/bin:$PATH:/tmp/test-etcd make test
fi
