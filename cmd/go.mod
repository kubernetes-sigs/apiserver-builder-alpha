module sigs.k8s.io/apiserver-builder-alpha/cmd

go 1.12

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/emicklei/go-restful v0.0.0-20170410110728-ff4f55a20633
	github.com/evanphx/json-patch v4.2.0+incompatible // indirect
	github.com/gobuffalo/envy v1.6.10 // indirect
	github.com/markbates/inflect v1.0.4
	github.com/pkg/errors v0.8.0
	github.com/rogpeppe/go-internal v1.0.1-alpha.3 // indirect
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cobra v0.0.3
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2 // indirect
	golang.org/x/tools v0.0.0-20190328211700-ab21143f2384 // indirect
	k8s.io/apimachinery v0.0.0-20190817020851-f2f3a405f61d
	k8s.io/apiserver v0.0.0-20190819142446-92cc630367d0
	k8s.io/client-go v0.0.0-20190819141724-e14f31a72a77
	k8s.io/gengo v0.0.0-20190907103519-ebc107f98eab
	k8s.io/klog v0.3.1
	sigs.k8s.io/controller-tools v0.1.10 // indirect
	sigs.k8s.io/kubebuilder v0.0.0-20190320190143-2621a6fdb324
)

replace github.com/golang/glog => k8s.io/klog v0.3.1
