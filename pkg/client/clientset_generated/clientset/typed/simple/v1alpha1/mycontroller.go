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

package v1alpha1

import (
	v1alpha1 "github.com/kubernetes-incubator/apiserver-builder-example/pkg/apis/simple/v1alpha1"
	scheme "github.com/kubernetes-incubator/apiserver-builder-example/pkg/client/clientset_generated/clientset/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// MyControllersGetter has a method to return a MyControllerInterface.
// A group's client should implement this interface.
type MyControllersGetter interface {
	MyControllers(namespace string) MyControllerInterface
}

// MyControllerInterface has methods to work with MyController resources.
type MyControllerInterface interface {
	Create(*v1alpha1.MyController) (*v1alpha1.MyController, error)
	Update(*v1alpha1.MyController) (*v1alpha1.MyController, error)
	UpdateStatus(*v1alpha1.MyController) (*v1alpha1.MyController, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.MyController, error)
	List(opts v1.ListOptions) (*v1alpha1.MyControllerList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.MyController, err error)
	MyControllerExpansion
}

// myControllers implements MyControllerInterface
type myControllers struct {
	client rest.Interface
	ns     string
}

// newMyControllers returns a MyControllers
func newMyControllers(c *SimpleV1alpha1Client, namespace string) *myControllers {
	return &myControllers{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Create takes the representation of a myController and creates it.  Returns the server's representation of the myController, and an error, if there is any.
func (c *myControllers) Create(myController *v1alpha1.MyController) (result *v1alpha1.MyController, err error) {
	result = &v1alpha1.MyController{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("mycontrollers").
		Body(myController).
		Do().
		Into(result)
	return
}

// Update takes the representation of a myController and updates it. Returns the server's representation of the myController, and an error, if there is any.
func (c *myControllers) Update(myController *v1alpha1.MyController) (result *v1alpha1.MyController, err error) {
	result = &v1alpha1.MyController{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("mycontrollers").
		Name(myController.Name).
		Body(myController).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclientstatus=false comment above the type to avoid generating UpdateStatus().

func (c *myControllers) UpdateStatus(myController *v1alpha1.MyController) (result *v1alpha1.MyController, err error) {
	result = &v1alpha1.MyController{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("mycontrollers").
		Name(myController.Name).
		SubResource("status").
		Body(myController).
		Do().
		Into(result)
	return
}

// Delete takes name of the myController and deletes it. Returns an error if one occurs.
func (c *myControllers) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("mycontrollers").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *myControllers) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("mycontrollers").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Get takes name of the myController, and returns the corresponding myController object, and an error if there is any.
func (c *myControllers) Get(name string, options v1.GetOptions) (result *v1alpha1.MyController, err error) {
	result = &v1alpha1.MyController{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("mycontrollers").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of MyControllers that match those selectors.
func (c *myControllers) List(opts v1.ListOptions) (result *v1alpha1.MyControllerList, err error) {
	result = &v1alpha1.MyControllerList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("mycontrollers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested myControllers.
func (c *myControllers) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("mycontrollers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Patch applies the patch and returns the patched myController.
func (c *myControllers) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.MyController, err error) {
	result = &v1alpha1.MyController{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("mycontrollers").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
