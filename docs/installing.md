# Installing the kubebuilder development tools

Installation instructions:

- Download the latest [release](https://github.com/najena/kubebuilder/releases)
- Extract to `/usr/local/kubebuilder/`
  - Create this directory if it does not already exist
- Add `/usr/local/kubebuilder/bin` to your path
  `export PATH=$PATH:/usr/local/kubebuilder/bin`
- Test things are working by running `kubebuilder version`
- Set test framework environment variables
  - export `TEST_ASSET_KUBECTL=/usr/local/kubebuilder/bin/kubectl`
  - export `TEST_ASSET_KUBE_APISERVER=/usr/local/kubebuilder/bin/kube-apiserver`
  - export `TEST_ASSET_ETCD=/usr/local/kubebuilder/bin/etcd`
