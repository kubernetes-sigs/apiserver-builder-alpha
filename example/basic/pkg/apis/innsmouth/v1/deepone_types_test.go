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

package v1_test

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	. "sigs.k8s.io/apiserver-builder-alpha/example/basic/pkg/apis/innsmouth/v1"
)

var _ = Describe("Deepone", func() {
	var instance DeepOne
	var expected DeepOne

	BeforeEach(func() {
		instance = DeepOne{}
		instance.Name = "deepone-1"
		instance.Namespace = "default"
		instance.Spec.FishRequired = 150

		expected = instance
	})

	AfterEach(func() {
		cs.Delete(context.TODO(), &instance)
	})

	Describe("when sending a storage request", func() {
		Context("for a valid config", func() {
			It("should provide CRUD access to the object", func() {
				By("returning success from the create request")
				actual := instance.DeepCopy()
				err := cs.Create(context.TODO(), actual)
				Expect(err).ShouldNot(HaveOccurred())

				By("defaulting the expected fields")
				Expect(actual.Spec).To(Equal(expected.Spec))

				By("returning the item for list requests")
				var result DeepOneList
				err = cs.List(context.TODO(), &result)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Items).To(HaveLen(1))
				Expect(result.Items[0].Spec).To(Equal(expected.Spec))

				By("returning the item for get requests")
				err = cs.Get(context.TODO(), client.ObjectKey{
					Namespace: instance.Namespace,
					Name:      instance.Name,
				}, actual)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(actual.Spec).To(Equal(expected.Spec))

				By("deleting the item for delete requests")
				err = cs.Delete(context.TODO(), &instance)
				Expect(err).ShouldNot(HaveOccurred())
				err = cs.List(context.TODO(), &result)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Items).To(HaveLen(0))
			})
		})
	})

	Describe("when listing a resource", func() {
		Context("using labels", func() {
			It("should find the matchings objects", func() {
				instance1 := DeepOne{}
				instance1.Name = "deepone-1"
				instance1.Namespace = "default"
				instance1.Spec.FishRequired = 150
				instance1.Labels = map[string]string{"foo": "1"}

				expected1 := instance1

				instance2 := DeepOne{}
				instance2.Name = "deepone-2"
				instance2.Namespace = "default"
				instance2.Spec.FishRequired = 140
				instance2.Labels = map[string]string{"foo": "2"}

				expected2 := instance2

				By("returning success from the create request")
				err := cs.Create(context.TODO(), &instance1)
				Expect(err).ShouldNot(HaveOccurred())

				By("returning success from the create request")
				err = cs.Create(context.TODO(), &instance2)
				Expect(err).ShouldNot(HaveOccurred())

				By("returning the item for list requests")
				var result DeepOneList
				err = cs.List(context.TODO(), &result, &client.ListOptions{Raw: &metav1.ListOptions{LabelSelector: "foo=1"}})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Items).To(HaveLen(1))
				Expect(result.Items[0].Spec).To(Equal(expected1.Spec))

				By("returning the item for list requests")
				err = cs.List(context.TODO(), &result, &client.ListOptions{Raw: &metav1.ListOptions{LabelSelector: "foo=2"}})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Items).To(HaveLen(1))
				Expect(result.Items[0].Spec).To(Equal(expected2.Spec))

				By("returning the item for list requests")
				err = cs.List(context.TODO(), &result)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Items).To(HaveLen(2))
			})
		})
	})
})
