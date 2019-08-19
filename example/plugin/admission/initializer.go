
/*
Copyright YEAR The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/



package admission

import (
	aggregatedclientset "sigs.k8s.io/apiserver-builder-alpha/example/pkg/client/clientset_generated/clientset"
	aggregatedinformerfactory "sigs.k8s.io/apiserver-builder-alpha/example/pkg/client/informers_generated/externalversions"
	"k8s.io/apiserver/pkg/admission"
)

// WantsAggregatedResourceClientSet defines a function which sets external ClientSet for admission plugins that need it
type WantsAggregatedResourceClientSet interface {
	SetAggregatedResourceClientSet(aggregatedclientset.Interface)
	admission.InitializationValidator
}

// WantsAggregatedResourceInformerFactory defines a function which sets InformerFactory for admission plugins that need it
type WantsAggregatedResourceInformerFactory interface {
	SetAggregatedResourceInformerFactory(aggregatedinformerfactory.SharedInformerFactory)
	admission.InitializationValidator
}

// New creates an instance of admission plugins initializer.
func New(
	clientset aggregatedclientset.Interface,
	informers aggregatedinformerfactory.SharedInformerFactory,
) pluginInitializer {
	return pluginInitializer{
		aggregatedResourceClient:    clientset,
		aggregatedResourceInformers: informers,
	}
}

type pluginInitializer struct {
	aggregatedResourceClient    aggregatedclientset.Interface
	aggregatedResourceInformers aggregatedinformerfactory.SharedInformerFactory
}

func (i pluginInitializer) Initialize(plugin admission.Interface) {
	if wants, ok := plugin.(WantsAggregatedResourceClientSet); ok {
		wants.SetAggregatedResourceClientSet(i.aggregatedResourceClient)
	}
	if wants, ok := plugin.(WantsAggregatedResourceInformerFactory); ok {
		wants.SetAggregatedResourceInformerFactory(i.aggregatedResourceInformers)
	}
}

