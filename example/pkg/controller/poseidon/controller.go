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

package poseidon

import (
	"log"

	extensionsv1beta1listers "k8s.io/client-go/listers/extensions/v1beta1"

	"github.com/kubernetes-incubator/apiserver-builder/pkg/builders"
	"k8s.io/client-go/rest"

	"fmt"
	olympusv1beta1 "github.com/kubernetes-incubator/apiserver-builder/example/pkg/apis/olympus/v1beta1"
	listers "github.com/kubernetes-incubator/apiserver-builder/example/pkg/client/listers_generated/olympus/v1beta1"
	"github.com/kubernetes-incubator/apiserver-builder/example/pkg/controller/sharedinformers"
	"github.com/kubernetes-incubator/apiserver-builder/pkg/controller"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// +controller:group=olympus,version=v1beta1,kind=Poseidon,resource=poseidons
type PoseidonControllerImpl struct {
	builders.DefaultControllerFns

	// lister indexes properties about Poseidon
	lister listers.PoseidonLister

	deploymentLister extensionsv1beta1listers.DeploymentLister

	cs *kubernetes.Clientset
}

// Init initializes the controller and is called by the generated code
// Registers eventhandlers to enqueue events
// config - client configuration for talking to the apiserver
// si - informer factory shared across all controllers for listening to events and indexing resource properties
// queue - message queue for handling new events.  unique to this controller.
func (c *PoseidonControllerImpl) Init(
	config *rest.Config,
	si *sharedinformers.SharedInformers,
	r func(key string) error) {

	// Set the informer and lister for subscribing to events and indexing poseidons labels
	i := si.Factory.Olympus().V1beta1().Poseidons()
	c.lister = i.Lister()

	// For watching Deployments
	log.Printf("Register Poseidon controller for Deployment events")
	di := si.KubernetesFactory.Extensions().V1beta1().Deployments()
	c.deploymentLister = di.Lister()
	si.Watch("PoseidonPod", di.Informer(), c.DeploymentToPoseidon, r)

	c.cs = si.KubernetesClientSet
}

func (c *PoseidonControllerImpl) DeploymentToPoseidon(i interface{}) (string, error) {
	d, _ := i.(*v1beta1.Deployment)
	log.Printf("Deployment update: %v", d.Name)
	if len(d.OwnerReferences) == 1 && d.OwnerReferences[0].Kind == "Poseidon" {
		return d.Namespace + "/" + d.OwnerReferences[0].Name, nil
	} else {
		// Not owned
		return "", nil
	}
}

// Reconcile handles enqueued messages
func (c *PoseidonControllerImpl) Reconcile(u *olympusv1beta1.Poseidon) error {
	// TODO: Instead of using the same name, include a hash function against the PodTemplate to uniquely identify multiple
	// Deployments

	d, err := c.deploymentLister.Deployments(u.Namespace).Get(u.Name)
	hash := fmt.Sprintf("%d", controller.GetHash(u.Spec.Template))

	if d == nil || err != nil {
		log.Printf("Creating Deployment for Poseidon %s", u.Name)
		d = &v1beta1.Deployment{}
		// Note, these may be defaulted on when the Deployment is created, so don't count on them
		// always being the same.  Will need to save a hash copy of the template to check if it has
		// changed on the resource since the Deployment was updated
		d.Name = u.Name
		d.Spec.Template.Spec = u.Spec.Template
		d.Spec.Selector = &v1.LabelSelector{}

		// Use the labels from the Poseidon object - TODO: Validate that the Poseidon labels are
		// specified, and make them immutable
		d.Spec.Selector.MatchLabels = u.Labels
		d.Spec.Template.ObjectMeta.Labels = u.Labels
		d.Annotations = u.Annotations
		if d.Annotations == nil {
			d.Annotations = map[string]string{}
		}
		// TODO: incorporate the name of the API groupversionkind into the annotation label
		d.Annotations["pod-hash"] = hash
		owner := v1.OwnerReference{
			Name:       u.Name,
			UID:        u.UID,
			Kind:       "Poseidon",
			APIVersion: olympusv1beta1.SchemeGroupVersion.String(),
		}
		d.OwnerReferences = append(d.OwnerReferences, owner)
		_, err = c.cs.ExtensionsV1beta1().Deployments(u.Namespace).Create(d)
		if err != nil {
			log.Printf("Error: %v", err)
		}
		return err
	} else if hash != d.Annotations["pod-hash"] {
		log.Printf("Updating Deployment for Poseidon %s", u.Name)
		d.Spec.Template.Spec = u.Spec.Template
		_, err = c.cs.ExtensionsV1beta1().Deployments(u.Namespace).Update(d)
		if err != nil {
			log.Printf("Error: %v", err)

		}
		return err
	} else {
		log.Printf("No changes to Deployment for Poseidon %s", u.Name)
	}

	return nil
}

func (c *PoseidonControllerImpl) Get(namespace, name string) (*olympusv1beta1.Poseidon, error) {
	return c.lister.Poseidons(namespace).Get(name)
}
