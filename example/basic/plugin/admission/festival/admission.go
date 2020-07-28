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

package festivaladmission

import (
	"context"
	"fmt"
	"k8s.io/apiserver/pkg/admission"
	genericadmissioninitializer "k8s.io/apiserver/pkg/admission/initializer"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	aggregatedadmission "sigs.k8s.io/apiserver-builder-alpha/example/basic/plugin/admission"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ admission.Interface = &festivalPlugin{}
var _ admission.MutationInterface = &festivalPlugin{}
var _ admission.ValidationInterface = &festivalPlugin{}
var _ genericadmissioninitializer.WantsExternalKubeInformerFactory = &festivalPlugin{}
var _ genericadmissioninitializer.WantsExternalKubeClientSet = &festivalPlugin{}
var _ aggregatedadmission.WantsAggregatedResourceInformerFactory = &festivalPlugin{}
var _ aggregatedadmission.WantsAggregatedResourceClientSet = &festivalPlugin{}

func NewFestivalPlugin() *festivalPlugin {
	return &festivalPlugin{
		Handler: admission.NewHandler(admission.Create, admission.Update),
	}
}

type festivalPlugin struct {
	*admission.Handler
}

func (p *festivalPlugin) ValidateInitialization() error {
	return nil
}

func (p *festivalPlugin) Admit(ctx context.Context, a admission.Attributes, o admission.ObjectInterfaces) error {
	fmt.Println("admitting festivals")
	return nil
}

func (p *festivalPlugin) Validate(ctx context.Context, a admission.Attributes, o admission.ObjectInterfaces) error {
	return nil
}

func (p *festivalPlugin) SetAggregatedResourceInformerFactory(cache.Cache) {}

func (p *festivalPlugin) SetAggregatedResourceClientSet(client.Client) {}

func (p *festivalPlugin) SetExternalKubeInformerFactory(informers.SharedInformerFactory) {}

func (p *festivalPlugin) SetExternalKubeClientSet(kubernetes.Interface) {}
