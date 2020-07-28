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
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"sigs.k8s.io/apiserver-builder-alpha/pkg/builders"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"

	. "sigs.k8s.io/apiserver-builder-alpha/example/basic/pkg/apis/miskatonic/v1beta1"
)

var _ = Describe("University", func() {
	var instance University
	var expected University

	BeforeEach(func() {
		instance = University{}
		instance.Name = "instance-1"
		instance.Namespace = "default"
		instance.Spec.FacultySize = 7

		expected = instance
		val := 15
		expected.Spec.MaxStudents = &val
	})

	AfterEach(func() {
		cs.Delete(context.TODO(), &instance)
	})

	Describe("when sending a campus request", func() {
		It("should return success", func() {
			err := cs.Create(context.TODO(), &instance)
			Expect(err).ShouldNot(HaveOccurred())

			campus := &UniversityCampus{
				Faculty: 30,
			}
			campus.Name = instance.Name
			restClient, err := apiutil.RESTClientForGVK(schema.GroupVersionKind{
				Group:   "miskatonic.k8s.io",
				Version: "v1beta1",
				Kind:    "UniversityCampus",
			}, config, serializer.NewCodecFactory(builders.Scheme))
			Expect(err).ShouldNot(HaveOccurred())
			err = restClient.Post().Namespace(instance.Namespace).
				Name(instance.Name).
				Resource("universities").
				SubResource("campus").
				Body(campus).Do(context.TODO()).Error()
			Expect(err).ShouldNot(HaveOccurred())

			expected.Spec.FacultySize = 30
			actual := instance.DeepCopy()
			err = cs.Get(context.TODO(), client.ObjectKey{
				Namespace: instance.Namespace,
				Name:      instance.Name,
			}, actual)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(actual.Spec).Should(Equal(expected.Spec))
		})
	})
})
