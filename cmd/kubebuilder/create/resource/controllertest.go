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

package resource

import (
	"fmt"
	"path/filepath"
	"strings"

	createutil "github.com/kubernetes-sigs/kubebuilder/cmd/kubebuilder/create/util"
	"github.com/kubernetes-sigs/kubebuilder/cmd/kubebuilder/util"
)

func doControllerTest(dir string, args resourceTemplateArgs) bool {
	path := filepath.Join(dir, "pkg", "controller", strings.ToLower(createutil.KindName),
		fmt.Sprintf("%s_suite_test.go",
			strings.ToLower(createutil.KindName)))
	util.WriteIfNotFound(path, "resource-controller-suite-test-template", controllerSuiteTestTemplate, args)

	path = filepath.Join(dir, "pkg", "controller", strings.ToLower(createutil.KindName), "controller_test.go")
	fmt.Printf("\t%s\n", filepath.Join(
		"pkg", "controller", strings.ToLower(createutil.KindName), "controller_test.go"))
	return util.WriteIfNotFound(path, "controller-test-template", controllerTestTemplate, args)
}

var controllerSuiteTestTemplate = `
{{.BoilerPlate}}

package {{lower .Kind}}_test

import (
    "testing"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    "k8s.io/client-go/rest"
    "github.com/kubernetes-sigs/kubebuilder/pkg/test"

    "{{ .Repo }}/pkg/apis"
    "{{ .Repo }}/pkg/client/clientset_generated/clientset"
    "{{ .Repo }}/pkg/controller/sharedinformers"
    "{{ .Repo }}/pkg/controller/{{lower .Kind}}"
)

var testenv *test.TestEnvironment
var config *rest.Config
var cs *clientset.Clientset
var shutdown chan struct{}
var controller *{{ lower .Kind }}.{{ .Kind }}Controller
var si *sharedinformers.SharedInformers

func Test{{.Kind}}(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecsWithDefaultAndCustomReporters(t, "{{ .Kind }} Suite", []Reporter{test.NewlineReporter{}})
}

var _ = BeforeSuite(func() {
    testenv = &test.TestEnvironment{CRDs: apis.APIMeta.GetCRDs()}
    var err error
    config, err = testenv.Start()
    Expect(err).NotTo(HaveOccurred())
    cs = clientset.NewForConfigOrDie(config)

    shutdown = make(chan struct{})
    si = sharedinformers.NewSharedInformers(config, shutdown)
    controller = {{ lower .Kind }}.New{{ .Kind}}Controller(config, si)
    controller.Run(shutdown)
})

var _ = AfterSuite(func() {
    testenv.Stop()
})
`

var controllerTestTemplate = `
{{.BoilerPlate}}

package {{ lower .Kind }}_test

import (
    . "{{ .Repo }}/pkg/apis/{{ .Group }}/{{ .Version }}"
    . "{{ .Repo }}/pkg/client/clientset_generated/clientset/typed/{{ .Group }}/{{ .Version }}"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!
// Created by "kubebuilder create resource" for you to implement controller logic tests

var _ = Describe("{{ .Kind }} controller", func() {
    var instance {{ .Kind }}
    var expectedKey string
    var client {{ .Kind }}Interface

    BeforeEach(func() {
        instance = {{ .Kind }}{}
        instance.Name = "instance-1"
        expectedKey = "{{ if not .NonNamespacedKind }}default/{{ end }}instance-1"
    })

    AfterEach(func() {
        client.Delete(instance.Name, &metav1.DeleteOptions{})
    })

    Describe("when creating a new object", func() {
        It("invoke the reconcile method", func() {
            after := make(chan struct{})
            controller.AfterReconcile = func(key string, err error) {
                defer func() {
                    // Recover in case the key is reconciled multiple times
                    defer func() { recover() }()
                    close(after)
                }()
                Expect(key).To(Equal(expectedKey))
                Expect(err).ToNot(HaveOccurred())
            }

            // Create the instance
            client = cs.{{title .Group}}{{title .Version}}().{{ plural .Kind }}({{ if not .NonNamespacedKind }}"default"{{ end }})
            _, err := client.Create(&instance)
            Expect(err).ShouldNot(HaveOccurred())

            // Wait for reconcile to happen
            Eventually(after).Should(BeClosed())

            // INSERT YOUR CODE HERE - test conditions post reconcile
        })
    })
})
`
