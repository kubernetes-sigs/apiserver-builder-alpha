# Getting started

## Create your Go project

Create a Go project under GOPATH/src/

For example

> GOPATH/src/github.com/my-org/my-project

## Download and install the code generators

Make sure the GOPATH/bin directory is on your path, and then use
 `go get` to download and compile the code generators:

```sh
go get k8s.io/kubernetes/cmd/libs/go2idl/client-gen
go get k8s.io/kubernetes/cmd/libs/go2idl/conversion-gen
go get k8s.io/kubernetes/cmd/libs/go2idl/deepcopy-gen
go get k8s.io/kubernetes/cmd/libs/go2idl/openapi-gen
go get k8s.io/kubernetes/cmd/libs/go2idl/defaulter-gen
go get k8s.io/kubernetes/cmd/libs/go2idl/lister-gen
go get k8s.io/kubernetes/cmd/libs/go2idl/informer-gen
go get github.com/kubernetes-incubator/apiserver-builder/cmd/apiregister-gen
go get github.com/kubernetes-incubator/reference-docs/gen-apidocs
```

Verify the downloaded code generators can be found on the path by running
`client-gen`

## Bootstrap `main.go`

Create a `main.go` file with the following contents that will be used
by Glide to bootstrap fetching dependencies.

```go
package main

import (
        _ "k8s.io/client-go/plugin/pkg/client/auth" // Enable cloud provider auth

	// TODO: Delete these after running glide to fetch the libraries
	_ "github.com/kubernetes-incubator/apiserver-builder/pkg/cmd/server"
	_ "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/discovery"
)

func main() {}
```

## Use Glide to vendor the apiserver-builder go libraries

*Note:* Other dependency managers may also be used.  The insructions
here are for using glide.

1. Install `glide` if you have not already.  [Instructions](https://github.com/Masterminds/glide)

2. Create a `glide.yaml` with the follow contents:

```yaml
package: YOUR/GO/PACKAGE
ignore:
# prevent glide from trying to vendor ourselves
- YOUR/GO/PACKAGE
import:
- package: github.com/kubernetes-incubator/apiserver-builder
- package: k8s.io/apimachinery
- package: k8s.io/apiserver
- package: k8s.io/client-go
```

3. Use glide to install the apiserver-builder libraries.

*Note:* The `--strip-vendor` flag must be supplied to flatten the dependencies
so that vendored libraries share the same instances of shared libaries.

```sh
glide install --strip-vendor
```

4.