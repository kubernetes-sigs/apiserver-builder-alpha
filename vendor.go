package apiserver_builder_alpha

import (
	// these imports are for ensuring the packages present in the vendor by `go mod vendor`.
	_ "github.com/go-openapi/loads"
	_ "github.com/kubernetes-incubator/reference-docs/gen-apidocs/generators"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	_ "sigs.k8s.io/controller-runtime/pkg/client"
	_ "sigs.k8s.io/controller-runtime/pkg/controller"
	_ "sigs.k8s.io/controller-runtime/pkg/handler"
	_ "sigs.k8s.io/controller-runtime/pkg/manager"
	_ "sigs.k8s.io/controller-runtime/pkg/reconcile"
	_ "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	_ "sigs.k8s.io/controller-runtime/pkg/runtime/signals"
	_ "sigs.k8s.io/controller-runtime/pkg/source"
	_ "sigs.k8s.io/controller-tools/pkg/util"
)
