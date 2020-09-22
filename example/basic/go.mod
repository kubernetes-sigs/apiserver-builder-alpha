module sigs.k8s.io/apiserver-builder-alpha/example/basic

go 1.15

require (
	github.com/go-logr/zapr v0.1.1 // indirect
	github.com/go-openapi/loads v0.19.4
	github.com/go-openapi/runtime v0.19.4
	github.com/go-openapi/spec v0.19.3
	github.com/go-openapi/validate v0.19.5
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.1-0.20190723091251-e0797f438f94 // indirect
	github.com/kubernetes-incubator/reference-docs v0.0.0-20170929004150-fcf65347b256 // indirect
	github.com/markbates/inflect v1.0.4 // indirect
	github.com/onsi/ginkgo v1.13.0 // indirect
	github.com/onsi/gomega v1.10.1
	github.com/pkg/errors v0.9.1 // indirect
	github.com/spf13/cobra v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.7.1
	github.com/tmc/grpc-websocket-proxy v0.0.0-20190109142713-0ad062ec5ee5 // indirect
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	k8s.io/api v0.19.2
	k8s.io/apiextensions-apiserver v0.18.2
	k8s.io/apimachinery v0.19.2
	k8s.io/apiserver v0.19.2
	k8s.io/client-go v0.19.2
	k8s.io/gengo v0.0.0-20200428234225-8167cfdcfc14 // indirect
	k8s.io/klog v1.0.0 // indirect
	k8s.io/klog/v2 v2.2.0
	k8s.io/kube-aggregator v0.18.4 // indirect
	k8s.io/kube-openapi v0.0.0-20200805222855-6aeccd4b50c6
	k8s.io/utils v0.0.0-20200729134348-d5654de09c73 // indirect
	sigs.k8s.io/apiserver-builder-alpha v0.0.0-00010101000000-000000000000
	sigs.k8s.io/apiserver-runtime v0.0.0-20200923110630-6937e37990f9
	sigs.k8s.io/controller-runtime v0.6.0
	sigs.k8s.io/controller-tools v0.1.12 // indirect
	sigs.k8s.io/kubebuilder v1.0.8 // indirect
	sigs.k8s.io/structured-merge-diff v1.0.1-0.20191108220359-b1b620dd3f06 // indirect
	sigs.k8s.io/testing_frameworks v0.1.1 // indirect
)

replace github.com/markbates/inflect => github.com/markbates/inflect v1.0.4

replace github.com/kubernetes-incubator/reference-docs => github.com/kubernetes-sigs/reference-docs v0.0.0-20170929004150-fcf65347b256

replace sigs.k8s.io/apiserver-builder-alpha => ../../../apiserver-builder-alpha
