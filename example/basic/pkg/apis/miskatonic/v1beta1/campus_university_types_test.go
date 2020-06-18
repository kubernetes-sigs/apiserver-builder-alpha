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

package v1beta1_test

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "sigs.k8s.io/apiserver-builder-alpha/example/basic/pkg/apis/miskatonic/v1beta1"
	. "sigs.k8s.io/apiserver-builder-alpha/example/basic/pkg/client/clientset_generated/clientset/typed/miskatonic/v1beta1"
)

var _ = Describe("University", func() {
	var instance University
	var expected University
	var client UniversityInterface

	BeforeEach(func() {
		instance = University{}
		instance.Name = "instance-1"
		instance.Spec.FacultySize = 7

		expected = instance
		val := 15
		expected.Spec.MaxStudents = &val
	})

	AfterEach(func() {
		client.Delete(context.TODO(), instance.Name, metav1.DeleteOptions{})
	})

	Describe("when sending a campus request", func() {
		It("should return success", func() {
			client = cs.MiskatonicV1beta1().Universities("university-test-campus")
			_, err := client.Create(context.TODO(), &instance, metav1.CreateOptions{})
			Expect(err).ShouldNot(HaveOccurred())

			campus := &UniversityCampus{
				Faculty: 30,
			}
			campus.Name = instance.Name
			restClient := cs.MiskatonicV1beta1().RESTClient()
			err = restClient.Post().Namespace("university-test-campus").
				Name(instance.Name).
				Resource("universities").
				SubResource("campus").
				Body(campus).Do(context.TODO()).Error()
			Expect(err).ShouldNot(HaveOccurred())

			expected.Spec.FacultySize = 30
			actual, err := client.Get(context.TODO(), instance.Name, metav1.GetOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(actual.Spec).Should(Equal(expected.Spec))
		})
	})
})
