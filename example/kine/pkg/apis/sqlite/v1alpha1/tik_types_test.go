
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



package v1alpha1_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "sigs.k8s.io/apiserver-builder-alpha/example/kine/pkg/apis/sqlite/v1alpha1"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Tik", func() {
	var instance Tik
	var expected Tik

	BeforeEach(func() {
		instance = Tik{}
		instance.Name = "instance-1"
		instance.Namespace = "default"

		expected = instance
	})

	AfterEach(func() {
		cs.Delete(context.TODO(), &instance)
	})

	Describe("when sending a storage request", func() {
		Context("for a valid config", func() {
			It("should provide CRUD access to the object", func() {
				actual := instance.DeepCopy()
				By("returning success from the create request")
				err := cs.Create(context.TODO(), actual)
				Expect(err).ShouldNot(HaveOccurred())

				By("defaulting the expected fields")
				Expect(actual.Spec).To(Equal(expected.Spec))

				By("returning the item for list requests")
				var result TikList
				err = cs.List(context.TODO(), &result)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Items).To(HaveLen(1))
				Expect(result.Items[0].Spec).To(Equal(expected.Spec))

				By("returning the item for get requests")
				err = cs.Get(context.TODO(), client.ObjectKey{
					Name:instance.Name,
					Namespace:instance.Namespace,
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
})
