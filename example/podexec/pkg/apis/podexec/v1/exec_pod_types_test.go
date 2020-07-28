
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



package v1_test

import (
	"context"
	"k8s.io/apimachinery/pkg/runtime/serializer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "sigs.k8s.io/apiserver-builder-alpha/example/podexec/pkg/apis/podexec/v1"

	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/apiserver-builder-alpha/pkg/builders"
)

var _ = Describe("Pod", func() {
	var instance Pod
	var expected Pod

	BeforeEach(func() {
		instance = Pod{}
		instance.Name = "instance-1"
		instance.Namespace = "default"

		expected = instance
	})

	AfterEach(func() {
		cs.Delete(context.TODO(), &instance)
	})

	Describe("when sending a exec request", func() {
		It("should return success", func() {
			err := cs.Create(context.TODO(), &instance)
			Expect(err).ShouldNot(HaveOccurred())

			exec := &PodExec{}
			restClient, err := apiutil.RESTClientForGVK(schema.GroupVersionKind{
				Group:   "podexec.example.com",
				Version: "v1",
				Kind:    "PodPodExec",
			}, config, serializer.NewCodecFactory(builders.Scheme))
			Expect(err).ShouldNot(HaveOccurred())
			err = restClient.Post().Namespace("pod-test-exec").
				Name(instance.Name).
				Resource("pods").
				SubResource("exec").
				Body(exec).
				Do(context.TODO()).
				Error()
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
