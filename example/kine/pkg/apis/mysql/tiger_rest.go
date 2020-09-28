package mysql

import (
	"context"
	"fmt"

	"github.com/rancher/kine/pkg/endpoint"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"sigs.k8s.io/apiserver-builder-alpha/pkg/builders"
)

func NewTigerREST(getter generic.RESTOptionsGetter) rest.Storage {
	groupResource := schema.GroupResource{
		Group:    "mysql.example.com",
		Resource: "tigers",
	}
	AddToScheme(builders.Scheme) // Fix it
	strategy := &TigerStrategy{builders.StorageStrategySingleton}
	store := &genericregistry.Store{
		NewFunc:                  func() runtime.Object { return &Tiger{} },
		NewListFunc:              func() runtime.Object { return &TigerList{} },
		DefaultQualifiedResource: groupResource,
		TableConvertor:           rest.NewDefaultTableConvertor(groupResource),

		CreateStrategy: strategy, // TODO: specify create strategy
		UpdateStrategy: strategy, // TODO: specify update strategy
		DeleteStrategy: strategy, // TODO: specify delete strategy
	}
	options := &generic.StoreOptions{RESTOptions: NewKineRESTOptionsGetter(getter)}
	if err := store.CompleteWithOptions(options); err != nil {
		panic(err) // TODO: Propagate error up
	}
	return &TigerREST{store}
}

// +k8s:deepcopy-gen=false
type TigerREST struct {
	*genericregistry.Store
}

func NewKineRESTOptionsGetter(getter generic.RESTOptionsGetter) generic.RESTOptionsGetter {
	return &kineProxiedRESTOptionsGetter{
		delegate: getter,
	}
}

type kineProxiedRESTOptionsGetter struct {
	delegate generic.RESTOptionsGetter
}

func (g *kineProxiedRESTOptionsGetter) GetRESTOptions(resource schema.GroupResource) (generic.RESTOptions, error) {
	restOptions, err := g.delegate.GetRESTOptions(resource)
	if err != nil {
		return generic.RESTOptions{}, err
	}

	if len(restOptions.StorageConfig.Transport.ServerList) != 1 {
		return generic.RESTOptions{}, fmt.Errorf("no valid mysql dsn found")
	}

	etcdConfig, err := endpoint.Listen(context.TODO(), endpoint.Config{
		Endpoint: restOptions.StorageConfig.Transport.ServerList[0],
	})
	if err != nil {
		//return restOptions, nil
		return generic.RESTOptions{}, err
	}

	restOptions.StorageConfig.Transport.ServerList = etcdConfig.Endpoints
	restOptions.StorageConfig.Transport.TrustedCAFile = etcdConfig.TLSConfig.CAFile
	restOptions.StorageConfig.Transport.CertFile = etcdConfig.TLSConfig.CertFile
	restOptions.StorageConfig.Transport.KeyFile = etcdConfig.TLSConfig.KeyFile
	return restOptions, nil
}
