
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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "sigs.k8s.io/apiserver-builder-alpha/example/podlogs/pkg/apis/podlogs/v1"
	. "sigs.k8s.io/apiserver-builder-alpha/example/podlogs/pkg/client/clientset_generated/clientset/typed/podlogs/v1"
)

var _ = Describe("Pod", func() {
	var instance Pod
	var expected Pod
	var client PodInterface

	BeforeEach(func() {
		instance = Pod{}
		instance.Name = "instance-1"

		expected = instance
	})

	AfterEach(func() {
		client.Delete(instance.Name, &metav1.DeleteOptions{})
	})

	Describe("when sending a logs request", func() {
		It("should return success", func() {
			client = cs.PodlogsV1().Pods("pod-test-logs")
			_, err := client.Create(&instance)
			Expect(err).ShouldNot(HaveOccurred())

			logs := &PodLogs{}
			logs.Name = instance.Name
			restClient := cs.PodlogsV1().RESTClient()
			err = restClient.Post().Namespace("pod-test-logs").
				Name(instance.Name).
				Resource("pods").
				SubResource("logs").
				Body(logs).Do().Error()
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
