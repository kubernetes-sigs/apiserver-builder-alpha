package v1beta1

import "log"

func SetDefaults_University(obj *University) {
	log.Printf("Defaulting University %s\n", obj.Name)
	if obj.Spec.MaxStudents == nil {
		n := 15
		obj.Spec.MaxStudents = &n
	}
}