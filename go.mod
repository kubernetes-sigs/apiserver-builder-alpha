module sigs.k8s.io/apiserver-builder-alpha

go 1.13

require (
	github.com/go-openapi/loads v0.19.4
	github.com/kubernetes-incubator/reference-docs v0.0.0-20170929004150-fcf65347b256
	github.com/markbates/inflect v1.0.4
	github.com/onsi/ginkgo v1.11.0
	github.com/onsi/gomega v1.8.1
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.4.0
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/apiserver v0.19.2
	k8s.io/client-go v0.19.2
	k8s.io/gengo v0.0.0-20200428234225-8167cfdcfc14
	k8s.io/klog v1.0.0
	k8s.io/kube-aggregator v0.19.2
	k8s.io/kube-openapi v0.0.0-20200805222855-6aeccd4b50c6
	k8s.io/utils v0.0.0-20200729134348-d5654de09c73
	sigs.k8s.io/controller-runtime v0.6.0
	sigs.k8s.io/controller-tools v0.1.12
	sigs.k8s.io/kubebuilder v1.0.8
	sigs.k8s.io/testing_frameworks v0.1.1
)

replace sigs.k8s.io/apiserver-builder-alpha/test => ./test

replace sigs.k8s.io/apiserver-builder-alpha/example/basic => ./example/basic
