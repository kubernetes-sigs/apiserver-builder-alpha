
/*
Copyright 2019 The Kubernetes Authors.

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



package tigeradmission

import (
	"context"
	aggregatedadmission "sigs.k8s.io/apiserver-builder-alpha/example/kine/plugin/admission"
	aggregatedinformerfactory "sigs.k8s.io/apiserver-builder-alpha/example/kine/pkg/client/informers_generated/externalversions"
	aggregatedclientset "sigs.k8s.io/apiserver-builder-alpha/example/kine/pkg/client/clientset_generated/clientset"
	genericadmissioninitializer "k8s.io/apiserver/pkg/admission/initializer"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/apiserver/pkg/admission"
)

var _ admission.Interface 											= &tigerPlugin{}
var _ admission.MutationInterface 									= &tigerPlugin{}
var _ admission.ValidationInterface 								= &tigerPlugin{}
var _ genericadmissioninitializer.WantsExternalKubeInformerFactory 	= &tigerPlugin{}
var _ genericadmissioninitializer.WantsExternalKubeClientSet 		= &tigerPlugin{}
var _ aggregatedadmission.WantsAggregatedResourceInformerFactory 	= &tigerPlugin{}
var _ aggregatedadmission.WantsAggregatedResourceClientSet 			= &tigerPlugin{}

func NewTigerPlugin() *tigerPlugin {
	return &tigerPlugin{
		Handler: admission.NewHandler(admission.Create, admission.Update),
	}
}

type tigerPlugin struct {
	*admission.Handler
}

func (p *tigerPlugin) ValidateInitialization() error {
	return nil
}

func (p *tigerPlugin) Admit(ctx context.Context, a admission.Attributes, o admission.ObjectInterfaces) error {
	return nil
}

func (p *tigerPlugin) Validate(ctx context.Context, a admission.Attributes, o admission.ObjectInterfaces) error {
	return nil
}

func (p *tigerPlugin) SetAggregatedResourceInformerFactory(aggregatedinformerfactory.SharedInformerFactory) {}

func (p *tigerPlugin) SetAggregatedResourceClientSet(aggregatedclientset.Interface) {}

func (p *tigerPlugin) SetExternalKubeInformerFactory(informers.SharedInformerFactory) {}

func (p *tigerPlugin) SetExternalKubeClientSet(kubernetes.Interface) {}
