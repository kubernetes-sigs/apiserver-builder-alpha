package v1beta1

import "k8s.io/klog"

func SetDefaults_University(obj *University) {
	klog.Infof("Defaulting University %s", obj.Name)
	if obj.Spec.MaxStudents == nil {
		n := 15
		obj.Spec.MaxStudents = &n
	}
}
