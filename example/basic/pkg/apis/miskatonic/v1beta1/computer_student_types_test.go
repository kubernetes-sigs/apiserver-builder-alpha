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
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	. "sigs.k8s.io/apiserver-builder-alpha/example/basic/pkg/apis/miskatonic/v1beta1"
	"sigs.k8s.io/apiserver-builder-alpha/pkg/builders"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

var _ = Describe("Student", func() {
	var instance Student

	BeforeEach(func() {
		instance = Student{}
		instance.Name = "instance-1"
		instance.Namespace = "default"

		// expected = instance
	})

	AfterEach(func() {
		cs.Delete(context.TODO(), &instance)
	})

	Describe("when sending a computer request", func() {
		It("should return success", func() {
			err := cs.Create(context.TODO(), &instance)
			Expect(err).ShouldNot(HaveOccurred())

			computer := &StudentComputer{}
			computer.Name = instance.Name
			restClient, err := apiutil.RESTClientForGVK(schema.GroupVersionKind{
				Group:   "miskatonic.k8s.io",
				Version: "v1beta1",
				Kind:    "StudentComputer",
			}, config, serializer.NewCodecFactory(builders.Scheme))
			Expect(err).ShouldNot(HaveOccurred())
			err = restClient.Post().Namespace(instance.Namespace).
				Name(instance.Name).
				Resource("students").
				SubResource("computer").
				Body(computer).Do(context.TODO()).Error()
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
