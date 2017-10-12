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

package fake

import (
	v1alpha1 "github.com/kubernetes-incubator/apiserver-builder-example/pkg/apis/simple/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeMyControllers implements MyControllerInterface
type FakeMyControllers struct {
	Fake *FakeSimpleV1alpha1
	ns   string
}

var mycontrollersResource = schema.GroupVersionResource{Group: "simple.example.com", Version: "v1alpha1", Resource: "mycontrollers"}

var mycontrollersKind = schema.GroupVersionKind{Group: "simple.example.com", Version: "v1alpha1", Kind: "MyController"}

func (c *FakeMyControllers) Create(myController *v1alpha1.MyController) (result *v1alpha1.MyController, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(mycontrollersResource, c.ns, myController), &v1alpha1.MyController{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MyController), err
}

func (c *FakeMyControllers) Update(myController *v1alpha1.MyController) (result *v1alpha1.MyController, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(mycontrollersResource, c.ns, myController), &v1alpha1.MyController{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MyController), err
}

func (c *FakeMyControllers) UpdateStatus(myController *v1alpha1.MyController) (*v1alpha1.MyController, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(mycontrollersResource, "status", c.ns, myController), &v1alpha1.MyController{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MyController), err
}

func (c *FakeMyControllers) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(mycontrollersResource, c.ns, name), &v1alpha1.MyController{})

	return err
}

func (c *FakeMyControllers) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(mycontrollersResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.MyControllerList{})
	return err
}

func (c *FakeMyControllers) Get(name string, options v1.GetOptions) (result *v1alpha1.MyController, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(mycontrollersResource, c.ns, name), &v1alpha1.MyController{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MyController), err
}

func (c *FakeMyControllers) List(opts v1.ListOptions) (result *v1alpha1.MyControllerList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(mycontrollersResource, mycontrollersKind, c.ns, opts), &v1alpha1.MyControllerList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.MyControllerList{}
	for _, item := range obj.(*v1alpha1.MyControllerList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested myControllers.
func (c *FakeMyControllers) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(mycontrollersResource, c.ns, opts))

}

// Patch applies the patch and returns the patched myController.
func (c *FakeMyControllers) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.MyController, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(mycontrollersResource, c.ns, name, data, subresources...), &v1alpha1.MyController{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MyController), err
}
