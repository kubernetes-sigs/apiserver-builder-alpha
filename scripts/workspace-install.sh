#!/usr/bin/env bash

[[ -f bazel_0.29.1-linux-x86_64.deb ]] || wget https://github.com/bazelbuild/bazel/releases/download/0.29.1/bazel_0.29.1-linux-x86_64.deb
[[ -f kubernetes-server-linux-amd64.tar.gz ]] || wget https://dl.k8s.io/v1.16.0/kubernetes-server-linux-amd64.tar.gz
[[ -f etcd-v3.2.0-linux-amd64.tar.gz ]] || wget https://github.com/coreos/etcd/releases/download/v3.2.0/etcd-v3.2.0-linux-amd64.tar.gz

# bazel installation
sudo dpkg -i bazel_0.29.1-linux-x86_64.deb

# kubebuilder installation
sudo mkdir -p /usr/local/kubebuilder/bin/
tar -zxvf kubernetes-server-linux-amd64.tar.gz
sudo mv kubernetes/server/bin/kube-apiserver /usr/local/kubebuilder/bin/

# etcd installation
mkdir -p /tmp/test-etcd/
tar xzvf etcd-v3.2.0-linux-amd64.tar.gz -C /tmp/test-etcd/ --strip-components=1
sudo cp /tmp/test-etcd/{etcd,etcdctl} /usr/local/kubebuilder/bin/
sudo cp /tmp/test-etcd/{etcd,etcdctl} /usr/local/bin/

# install apiserver-boots
mkdir -p $(go env GOPATH)/bin/
