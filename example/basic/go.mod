module sigs.k8s.io/apiserver-builder-alpha/example/basic

go 1.13

require (
	github.com/markbates/inflect v0.0.0-00010101000000-000000000000
	github.com/onsi/ginkgo v1.11.0
	github.com/onsi/gomega v1.8.1
	github.com/pkg/errors v0.8.1
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	k8s.io/api v0.17.3
	k8s.io/apimachinery v0.17.3
	k8s.io/apiserver v0.17.3
	k8s.io/client-go v0.17.3
	k8s.io/gengo v0.0.0-20190822140433-26a664648505
	k8s.io/klog v1.0.0
	k8s.io/kube-aggregator v0.17.3
	k8s.io/kube-openapi v0.0.0-20191107075043-30be4d16710a
	k8s.io/utils v0.0.0-20191114184206-e782cd3c129f
	sigs.k8s.io/apiserver-builder-alpha v0.0.0-00010101000000-000000000000
	sigs.k8s.io/controller-runtime v0.5.1
	sigs.k8s.io/kubebuilder v1.0.8
	sigs.k8s.io/testing_frameworks v0.1.1
)

replace sigs.k8s.io/controller-tools => sigs.k8s.io/controller-tools v0.1.12

replace sigs.k8s.io/kubebuilder => sigs.k8s.io/kubebuilder v1.0.8

replace github.com/markbates/inflect => github.com/markbates/inflect v1.0.4

replace github.com/kubernetes-incubator/reference-docs => github.com/kubernetes-sigs/reference-docs v0.0.0-20170929004150-fcf65347b256

replace sigs.k8s.io/apiserver-builder-alpha => ../../../apiserver-builder-alpha
