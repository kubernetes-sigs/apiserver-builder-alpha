{{.BoilerPlate}}

package main

import (
	"k8s.io/klog"
	"sigs.k8s.io/apiserver-runtime/pkg/builder"

	// +kubebuilder:scaffold:resource-imports
)

func main() {
	err := builder.APIServer.
		// +kubebuilder:scaffold:resource-register
		Execute()
	if err != nil {
		klog.Fatal(err)
	}
}