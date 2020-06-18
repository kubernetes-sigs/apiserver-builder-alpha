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

package v1_test

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "sigs.k8s.io/apiserver-builder-alpha/example/basic/pkg/apis/kingsport/v1"
	. "sigs.k8s.io/apiserver-builder-alpha/example/basic/pkg/client/clientset_generated/clientset/typed/kingsport/v1"
)

var _ = Describe("Festival", func() {
	var instance Festival
	var expected Festival
	var client FestivalInterface

	BeforeEach(func() {
		instance = Festival{}
		instance.Name = "instance-1"
		instance.Spec.Year = 1
		expected = instance
	})

	AfterEach(func() {
		client.Delete(context.TODO(), instance.Name, metav1.DeleteOptions{})
	})

	Describe("when sending a storage request", func() {
		Context("for a valid config", func() {
			It("should provide CRUD access to the object", func() {
				client = cs.KingsportV1().Festivals()

				By("returning success from the create request")
				actual, err := client.Create(context.TODO(), &instance, metav1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				By("defaulting the expected fields")
				Expect(actual.Spec).To(Equal(expected.Spec))

				By("returning the item for list requests")
				result, err := client.List(context.TODO(), metav1.ListOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Items).To(HaveLen(1))
				Expect(result.Items[0].Spec).To(Equal(expected.Spec))

				By("returning the item for get requests")
				actual, err = client.Get(context.TODO(), instance.Name, metav1.GetOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(actual.Spec).To(Equal(expected.Spec))

				By("deleting the item for delete requests")
				err = client.Delete(context.TODO(), instance.Name, metav1.DeleteOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				result, err = client.List(context.TODO(), metav1.ListOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Items).To(HaveLen(0))
			})
		})
		Context("for an invalid config", func() {
			It("should fail", func() {
				instance.Spec.Year = -1
				client = cs.KingsportV1().Festivals()

				By("returning success from the create request")
				_, err := client.Create(context.TODO(), &instance, metav1.CreateOptions{})
				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
