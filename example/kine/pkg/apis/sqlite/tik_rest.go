package sqlite

import (
	"context"
	"github.com/rancher/kine/pkg/endpoint"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"sigs.k8s.io/apiserver-builder-alpha/pkg/builders"
)

func NewTikREST(getter generic.RESTOptionsGetter) rest.Storage {
	groupResource := schema.GroupResource{
		Group:    "sqlite",
		Resource: "tiks",
	}
	strategy := &TikStrategy{builders.StorageStrategySingleton}
	store := &genericregistry.Store{
		NewFunc:                  func() runtime.Object { return &Tik{} },
		NewListFunc:              func() runtime.Object { return &TikList{} },
		DefaultQualifiedResource: groupResource,

		CreateStrategy: strategy, // TODO: specify create strategy
		UpdateStrategy: strategy, // TODO: specify update strategy
		DeleteStrategy: strategy, // TODO: specify delete strategy
	}
	options := &generic.StoreOptions{RESTOptions: NewKineRESTOptionsGetter(getter, endpoint.Config{
		Listener: endpoint.KineSocket,
	})}
	if err := store.CompleteWithOptions(options); err != nil {
		panic(err) // TODO: Propagate error up
	}
	return &TikREST{store}
}

// +k8s:deepcopy-gen=false
type TikREST struct {
	*genericregistry.Store
}

func NewKineRESTOptionsGetter(getter generic.RESTOptionsGetter, config endpoint.Config) generic.RESTOptionsGetter {
	return &kineProxiedRESTOptionsGetter{
		delegate:   getter,
		kineConfig: config,
	}
}

type kineProxiedRESTOptionsGetter struct {
	delegate   generic.RESTOptionsGetter
	kineConfig endpoint.Config
}

func (g *kineProxiedRESTOptionsGetter) GetRESTOptions(resource schema.GroupResource) (generic.RESTOptions, error) {
	restOptions, err := g.delegate.GetRESTOptions(resource)
	if err != nil {
		return generic.RESTOptions{}, err
	}
	etcdConfig, err := endpoint.Listen(context.TODO(), g.kineConfig)
	if err != nil {
		//return restOptions, nil
		return generic.RESTOptions{}, err
	}

	restOptions.StorageConfig.Transport.ServerList = etcdConfig.Endpoints
	restOptions.StorageConfig.Transport.CAFile = etcdConfig.TLSConfig.CAFile
	restOptions.StorageConfig.Transport.CertFile = etcdConfig.TLSConfig.CertFile
	restOptions.StorageConfig.Transport.KeyFile = etcdConfig.TLSConfig.KeyFile
	return restOptions, nil
}
