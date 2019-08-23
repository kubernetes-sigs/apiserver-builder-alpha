
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



package deeponeadmission

import (
	"fmt"
	aggregatedadmission "sigs.k8s.io/apiserver-builder-alpha/example/basic/plugin/admission"
	aggregatedinformerfactory "sigs.k8s.io/apiserver-builder-alpha/example/basic/pkg/client/informers_generated/externalversions"
	aggregatedclientset "sigs.k8s.io/apiserver-builder-alpha/example/basic/pkg/client/clientset_generated/clientset"
	genericadmissioninitializer "k8s.io/apiserver/pkg/admission/initializer"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/apiserver/pkg/admission"
)

var _ admission.Interface 											= &deeponePlugin{}
var _ admission.MutationInterface 									= &deeponePlugin{}
var _ admission.ValidationInterface 								= &deeponePlugin{}
var _ genericadmissioninitializer.WantsExternalKubeInformerFactory 	= &deeponePlugin{}
var _ genericadmissioninitializer.WantsExternalKubeClientSet 		= &deeponePlugin{}
var _ aggregatedadmission.WantsAggregatedResourceInformerFactory 	= &deeponePlugin{}
var _ aggregatedadmission.WantsAggregatedResourceClientSet 			= &deeponePlugin{}

func NewDeepOnePlugin() *deeponePlugin {
	return &deeponePlugin{
		Handler: admission.NewHandler(admission.Create, admission.Update),
	}
}

type deeponePlugin struct {
	*admission.Handler
}

func (p *deeponePlugin) ValidateInitialization() error {
	return nil
}

func (p *deeponePlugin) Admit(a admission.Attributes, o admission.ObjectInterfaces) error {
	fmt.Println("admitting deepones")
	return nil
}

func (p *deeponePlugin) Validate(a admission.Attributes, o admission.ObjectInterfaces) error {
	return nil
}

func (p *deeponePlugin) SetAggregatedResourceInformerFactory(aggregatedinformerfactory.SharedInformerFactory) {}

func (p *deeponePlugin) SetAggregatedResourceClientSet(aggregatedclientset.Interface) {}

func (p *deeponePlugin) SetExternalKubeInformerFactory(informers.SharedInformerFactory) {}

func (p *deeponePlugin) SetExternalKubeClientSet(kubernetes.Interface) {}
