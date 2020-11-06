# Manually Installing the apiserver build tools


Requires the following to already be installed:
- kubectl
- go
- docker
- openssl
- base64
- tar
- cp

Installing by Go Get

```go
GO111MODULE=on go get sigs.k8s.io/apiserver-builder-alpha/cmd/apiserver-boot
```

Alternative Manual Installation instructions:

- Download the latest [release](https://github.com/kubernetes-sigs/apiserver-builder-alpha/releases)
- Extract to `/usr/local/apiserver-builder/`
  - Create this directory if it does not already exist
- Add `/usr/local/apiserver-builder/bin` to your path
  `export PATH=$PATH:/usr/local/apiserver-builder/bin`
- Test things are working by running `apiserver-boot -h`.