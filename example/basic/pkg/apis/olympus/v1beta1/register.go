package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

//import (
//	"fmt"
//	"sigs.k8s.io/apiserver-builder-alpha/pkg/builders"
//)
//
//func init() {
//	builders.Scheme.AddFieldLabelConversionFunc(SchemeGroupVersion.WithKind("Poseidon"), PoseidonFieldSelectorConversion)
//}
//
//// All field selector fields must appear in this function
//func PoseidonFieldSelectorConversion(label, value string) (string, string, error) {
//	switch label {
//	case "metadata.name":
//		return label, value, nil
//	case "metadata.namespace":
//		return label, value, nil
//	case "spec.deployment.name":
//		return label, value, nil
//	default:
//		return "", "", fmt.Errorf("%q is not a known field selector: only %q, %q, %q", label, "metadata.name", "metadata.namespace", "spec.deployment.name")
//	}
//}

var AddToScheme = func(scheme *runtime.Scheme) error {
	metav1.AddToGroupVersion(scheme, schema.GroupVersion{
		Group:   "olympus.k8s.io",
		Version: "v1beta1",
	})
	// +kubebuilder:scaffold:install
	scheme.AddKnownTypes(schema.GroupVersion{
		Group:   "olympus.k8s.io",
		Version: "v1beta1",
	}, &Poseidon{}, &PoseidonList{})
	return nil
}
