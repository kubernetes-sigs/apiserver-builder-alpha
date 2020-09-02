package filepath

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/klog"
	"sigs.k8s.io/apiserver-builder-alpha/pkg/storage/filepath"
)

func NewBurgerREST(getter generic.RESTOptionsGetter) rest.Storage {
	gr := schema.GroupResource{
		Group:    "filepath.example.com",
		Resource: "burgers",
	}
	opt, err := getter.GetRESTOptions(gr)
	if err != nil {
		klog.Fatal(err)
	}
	return filepath.NewFilepathREST(
		gr,
		opt.StorageConfig.Codec,
		"/data/",
		true,
		func() runtime.Object { return &Burger{} },
		func() runtime.Object { return &BurgerList{} },
	)
}
