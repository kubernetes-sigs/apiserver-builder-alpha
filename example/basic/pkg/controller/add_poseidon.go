package controller

import (
	"sigs.k8s.io/apiserver-builder-alpha/example/basic/pkg/controller/poseidon"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, poseidon.Add)
}
