#!/usr/bin/env bash

version=1.0.8 # latest stable version
arch=amd64
kubernetes_version=v1.15.0

# download the release
curl -L -O "https://github.com/kubernetes-sigs/kubebuilder/releases/download/v${version}/kubebuilder_${version}_linux_${arch}.tar.gz"

# extract the archive
tar -zxvf kubebuilder_${version}_linux_${arch}.tar.gz
mv kubebuilder_${version}_linux_${arch} kubebuilder && sudo mv kubebuilder /usr/local/


curl -L -O "https://dl.k8s.io/${kubernetes_version}/kubernetes-server-linux-amd64.tar.gz"
tar -zxvf kubernetes-server-linux-amd64.tar.gz
mv kubernetes/server/bin/kube-apiserver /usr/local/kubebuilder/bin/

# update your PATH to include /usr/local/kubebuilder/bin
export PATH=$PATH:/usr/local/kubebuilder/bin
