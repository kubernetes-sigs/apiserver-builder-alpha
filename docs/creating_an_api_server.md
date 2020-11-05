# Creating an API Server

**Note:** This document explains how to manually create the files generated
by `apiserver-boot`.  It is recommended to automatically create these files
instead.

## Create the apiserver command

Create a file called `main.go` under `cmd/apiserver` in your project. This
file bootstraps the apiserver by invoking the apiserver start function
with the generated API code.

```go

package main

import (
	"k8s.io/klog"
	"sigs.k8s.io/apiserver-runtime/pkg/builder"

	// +kubebuilder:scaffold:resource-imports
	storagev1 "example.io/kimmin/pkg/apis/storage/v1"
)

const storagePath = "/registry/YOUR.DOMAIN" // Change this

func main() {
	err := builder.APIServer.
		// +kubebuilder:scaffold:resource-register
		WithResource(&storagev1.VolumeClaim{}).
		SetDelegateAuthOptional().
		WithLocalDebugExtension().
		WithOptionsFns(func(options *builder.ServerOptions) *builder.ServerOptions {
			options.RecommendedOptions.CoreAPI = nil
			options.RecommendedOptions.Admission = nil
			options.RecommendedOptions.Authorization = nil
			return options
		}).
		Execute()
	if err != nil {
		klog.Fatal(err)
	}
}
```

## Create the API root package

Create your API root under `pkg/apis`

*Location:*  `GOPATH/src/YOUR/GO/PACKAGE/pkg/apis/doc.go`

- Change `YOUR.DOMAIN` to the domain you want your API groups to appear under.

```go
// +domain=YOUR.DOMAIN

package apis
```

## Create an API group

Create your API group under `pkg/apis/GROUP`

*Location:* `GOPATH/src/YOUR/GO/PACKAGE/pkg/apis/GROUP/doc.go`

- Change GROUP to be the group name.
- Change YOUR.DOMAIN to be your domain.

```go
// +k8s:deepcopy-gen=package,register
// +groupName=GROUP.YOUR.DOMAIN

// Package api is the internal version of the API.
package GROUP
```

## Create an API version

Create your API group under `pkg/apis/GROUP/VERSION`

*Location:* `GOPATH/src/YOUR/GO/PACKAGE/pkg/apis/GROUP/VERSION/doc.go`

- Change GROUP to be the group name.
- Change VERSION to the be the api version name.
- Change YOUR/GO/PACKAGE to be the go package of you project.

```go
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen=package,register
// +k8s:conversion-gen=YOUR/GO/PACKAGE/pkg/apis/GROUP
// +k8s:defaulter-gen=TypeMeta

// +groupName=GROUP.VERSION
package VERSION // import "YOUR/GO/PACKAGE/pkg/apis/GROUP/VERSION"
```

## Create the API type definitions

## Generate the code

## Create an integration test

## Start the server locally