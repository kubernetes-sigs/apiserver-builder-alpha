package v1beta1

import (
	v1 "k8s.io/api/autoscaling/v1"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource"
)

var _ resource.ObjectWithScaleSubResource = &University{}

func (in *University) SetScale(scaleSubResource *v1.Scale) {
	in.Spec.FacultySize = int(scaleSubResource.Spec.Replicas)
}

func (in *University) GetScale() (scaleSubResource *v1.Scale) {
	return &v1.Scale{
		Spec: v1.ScaleSpec{
			Replicas: int32(in.Spec.FacultySize),
		},
	}
}
