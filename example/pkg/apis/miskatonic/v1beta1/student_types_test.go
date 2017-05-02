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
	"testing"

	v1beta1miskatonic "github.com/kubernetes-incubator/apiserver-builder/example/pkg/apis/miskatonic/v1beta1"
)

func TestCreateStudents(t *testing.T) {
	client := cs.MiskatonicV1beta1Client
	sclient := client.Students("test-create-delete-students")

	student := &v1beta1miskatonic.Student{}
	student.Name = "joe"
	student.Spec.ID = 3
	if actual, err := sclient.Create(student); err != nil {
		t.Fatalf("Failed to create %T %v", student, err)
	} else {
		if actual.Spec.ID != student.Spec.ID+1 {
			t.Fatalf("Expected to find ID %d, found %d", student.Spec.ID, actual.Spec.ID+1)
		}
	}
}
