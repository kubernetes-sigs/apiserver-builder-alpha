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

package university_test

import (
	v1beta1miskatonic "github.com/kubernetes-incubator/apiserver-builder/example/pkg/apis/miskatonic/v1beta1"
	"testing"
	"time"
)

func TestReconcileUniversity(t *testing.T) {
	beforeChan := make(chan struct{})
	afterChan := make(chan struct{})
	expectedKey := "test-controller-universities/miskatonic-university"
	controller.BeforeReconcile = func(key string) {
		defer close(beforeChan)
		if key != expectedKey {
			t.Fatalf("Expected reconcile before university %s got %s", expectedKey, key)
		}
	}
	controller.AfterReconcile = func(key string, err error) {
		defer close(afterChan)
		if key != expectedKey {
			t.Fatalf("Expected reconcile after university %s got %s", expectedKey, key)
		}
		if err != nil {
			t.Fatalf("Expected no error on reconcile university %s", key)
		}
	}

	client := cs.MiskatonicV1beta1Client
	intf := client.Universities("test-controller-universities")

	univ := &v1beta1miskatonic.University{}
	univ.Name = "miskatonic-university"
	univ.Spec.FacultySize = 7

	// Make sure we can create the resource
	if _, err := intf.Create(univ); err != nil {
		t.Fatalf("Failed to create %T %v", univ, err)
	}

	select {
	case <-beforeChan:
	case <-time.After(time.Second * 2):
		t.Fatalf("Create University event never reconciled")
	}

	select {
	case <-afterChan:
	case <-time.After(time.Second * 2):
		t.Fatalf("Create University event never finished")
	}
}
