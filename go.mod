module sigs.k8s.io/apiserver-builder-alpha

go 1.15

require (
	github.com/evanphx/json-patch v4.9.0+incompatible // indirect
	github.com/go-logr/logr v0.2.1 // indirect
	github.com/go-logr/zapr v0.2.0 // indirect
	github.com/go-openapi/loads v0.19.4
	github.com/kubernetes-incubator/reference-docs v0.0.0-20170929004150-fcf65347b256
	github.com/markbates/inflect v1.0.4
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.2
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.6.1
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/apiserver v0.19.2
	k8s.io/client-go v0.19.2
	k8s.io/gengo v0.0.0-20200428234225-8167cfdcfc14
	k8s.io/klog v1.0.0
	k8s.io/kube-aggregator v0.19.2
	k8s.io/kube-openapi v0.0.0-20200805222855-6aeccd4b50c6
	k8s.io/utils v0.0.0-20200912215256-4140de9c8800
	sigs.k8s.io/apiserver-runtime v0.0.0-20200925141712-5fcfc91568ad // indirect
	sigs.k8s.io/controller-runtime v0.6.0
	sigs.k8s.io/controller-tools v0.1.12
	sigs.k8s.io/kubebuilder v1.0.9-0.20200925141511-a2f239880b04
	sigs.k8s.io/testing_frameworks v0.1.1
)

replace sigs.k8s.io/apiserver-builder-alpha/test => ./test

replace sigs.k8s.io/apiserver-builder-alpha/example/basic => ./example/basic
