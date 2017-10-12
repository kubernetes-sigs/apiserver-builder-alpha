
/*
Copyright 2017 The Kubernetes Authors.

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



package mycontroller

import (
	"log"

	"github.com/kubernetes-incubator/apiserver-builder/pkg/builders"
	"k8s.io/client-go/rest"

	"github.com/kubernetes-incubator/apiserver-builder-example/pkg/apis/simple/v1alpha1"
	"github.com/kubernetes-incubator/apiserver-builder-example/pkg/controller/sharedinformers"
	listers "github.com/kubernetes-incubator/apiserver-builder-example/pkg/client/listers_generated/simple/v1alpha1"
)

// +controller:group=simple,version=v1alpha1,kind=MyController,resource=mycontrollers
type MyControllerControllerImpl struct {
	builders.DefaultControllerFns

	// lister indexes properties about MyController
	lister listers.MyControllerLister
}

// Init initializes the controller and is called by the generated code
// Registers eventhandlers to enqueue events
// config - client configuration for talking to the apiserver
// si - informer factory shared across all controllers for listening to events and indexing resource properties
// queue - message queue for handling new events.  unique to this controller.
func (c *MyControllerControllerImpl) Init(
	config *rest.Config,
	si *sharedinformers.SharedInformers,
    reconcileKey func(key string) error) {

	// Set the informer and lister for subscribing to events and indexing mycontrollers labels
	c.lister = si.Factory.Simple().V1alpha1().MyControllers().Lister()
}

// Reconcile handles enqueued messages
func (c *MyControllerControllerImpl) Reconcile(u *v1alpha1.MyController) error {
	// Implement controller logic here
	log.Printf("Running reconcile MyController for %s\n", u.Name)
	return nil
}

func (c *MyControllerControllerImpl) Get(namespace, name string) (*v1alpha1.MyController, error) {
	return c.lister.MyControllers(namespace).Get(name)
}
