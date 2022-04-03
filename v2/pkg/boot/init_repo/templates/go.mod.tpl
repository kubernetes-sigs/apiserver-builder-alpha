module {{.Repo}}

go 1.17

require (
	github.com/go-logr/logr v0.2.1 // indirect
	github.com/go-logr/zapr v0.2.0 // indirect
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v0.19.2
	k8s.io/klog v1.0.0
	sigs.k8s.io/apiserver-runtime v1.0.3
	sigs.k8s.io/controller-runtime v0.11.1
)