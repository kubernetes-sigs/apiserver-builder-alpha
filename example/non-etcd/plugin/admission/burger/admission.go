
/*
Copyright 2020 The Kubernetes Authors.

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



package burgeradmission

import (
	"context"
	aggregatedadmission "sigs.k8s.io/apiserver-builder-alpha/example/non-etcd/plugin/admission"
	aggregatedinformerfactory "sigs.k8s.io/apiserver-builder-alpha/example/non-etcd/pkg/client/informers_generated/externalversions"
	aggregatedclientset "sigs.k8s.io/apiserver-builder-alpha/example/non-etcd/pkg/client/clientset_generated/clientset"
	genericadmissioninitializer "k8s.io/apiserver/pkg/admission/initializer"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/apiserver/pkg/admission"
)

var _ admission.Interface 											= &burgerPlugin{}
var _ admission.MutationInterface 									= &burgerPlugin{}
var _ admission.ValidationInterface 								= &burgerPlugin{}
var _ genericadmissioninitializer.WantsExternalKubeInformerFactory 	= &burgerPlugin{}
var _ genericadmissioninitializer.WantsExternalKubeClientSet 		= &burgerPlugin{}
var _ aggregatedadmission.WantsAggregatedResourceInformerFactory 	= &burgerPlugin{}
var _ aggregatedadmission.WantsAggregatedResourceClientSet 			= &burgerPlugin{}

func NewBurgerPlugin() *burgerPlugin {
	return &burgerPlugin{
		Handler: admission.NewHandler(admission.Create, admission.Update),
	}
}

type burgerPlugin struct {
	*admission.Handler
}

func (p *burgerPlugin) ValidateInitialization() error {
	return nil
}

func (p *burgerPlugin) Admit(ctx context.Context, a admission.Attributes, o admission.ObjectInterfaces) error {
	return nil
}

func (p *burgerPlugin) Validate(ctx context.Context, a admission.Attributes, o admission.ObjectInterfaces) error {
	return nil
}

func (p *burgerPlugin) SetAggregatedResourceInformerFactory(aggregatedinformerfactory.SharedInformerFactory) {}

func (p *burgerPlugin) SetAggregatedResourceClientSet(aggregatedclientset.Interface) {}

func (p *burgerPlugin) SetExternalKubeInformerFactory(informers.SharedInformerFactory) {}

func (p *burgerPlugin) SetExternalKubeClientSet(kubernetes.Interface) {}
