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

package v1beta1_test

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "sigs.k8s.io/apiserver-builder-alpha/example/basic/pkg/apis/miskatonic/v1beta1"
	. "sigs.k8s.io/apiserver-builder-alpha/example/basic/pkg/client/clientset_generated/clientset/typed/miskatonic/v1beta1"
)

var _ = Describe("Student", func() {
	var instance Student
	// var expected Student
	var client StudentInterface

	BeforeEach(func() {
		instance = Student{}
		instance.Name = "instance-1"

		// expected = instance
	})

	AfterEach(func() {
		client.Delete(context.TODO(), instance.Name, metav1.DeleteOptions{})
	})

	Describe("when sending a computer request", func() {
		It("should return success", func() {
			client = cs.MiskatonicV1beta1().Students("student-test-computer")
			_, err := client.Create(context.TODO(), &instance, metav1.CreateOptions{})
			Expect(err).ShouldNot(HaveOccurred())

			computer := &StudentComputer{}
			computer.Name = instance.Name
			restClient := cs.MiskatonicV1beta1().RESTClient()
			err = restClient.Post().Namespace("student-test-computer").
				Name(instance.Name).
				Resource("students").
				SubResource("computer").
				Body(computer).Do(context.TODO()).Error()
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
