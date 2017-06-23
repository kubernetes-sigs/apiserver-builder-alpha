
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



package festival

import (
	"log"

	"github.com/kubernetes-incubator/apiserver-builder/pkg/controller"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	"github.com/kubernetes-incubator/apiserver-builder/example/pkg/apis/kingsport/v1"
	"github.com/kubernetes-incubator/apiserver-builder/example/pkg/controller/sharedinformers"
	listers "github.com/kubernetes-incubator/apiserver-builder/example/pkg/client/listers_generated/kingsport/v1"
)

// +controller:group=kingsport,version=v1,kind=Festival,resource=festivals
type FestivalControllerImpl struct {
	// informer listens for events about Festival
	informer cache.SharedIndexInformer

	// lister indexes properties about Festival
	lister listers.FestivalLister
}

// Init initializes the controller and is called by the generated code
// Registers eventhandlers to enqueue events
// config - client configuration for talking to the apiserver
// si - informer factory shared across all controllers for listening to events and indexing resource properties
// queue - message queue for handling new events.  unique to this controller.
func (c *FestivalControllerImpl) Init(
	config *rest.Config,
	si *sharedinformers.SharedInformers,
	queue workqueue.RateLimitingInterface) {

	// Set the informer and lister for subscribing to events and indexing festivals labels
	i := si.Factory.Kingsport().V1().Festivals()
	c.informer = i.Informer()
	c.lister = i.Lister()

	// Add an event handler to enqueue a message for festivals adds / updates
	c.informer.AddEventHandler(&controller.QueueingEventHandler{queue})
}

// Reconcile handles enqueued messages
func (c *FestivalControllerImpl) Reconcile(u *v1.Festival) error {
	// Implement controller logic here
	log.Printf("Running reconcile Festival for %s\n", u.Name)
	return nil
}

func (c *FestivalControllerImpl) Get(namespace, name string) (*v1.Festival, error) {
	return c.lister.Get(name)
}
