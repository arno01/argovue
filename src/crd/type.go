package crd

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Service struct {
	Name              string
	Namespace         string
	Image             string
	SharedVolume      string
	PrivateVolumeSize string
}

func Parse(obj interface{}) *Service {
	m := new(Service)
	object := obj.(*unstructured.Unstructured)
	spec := object.Object["spec"].(map[string]interface{})
	m.Name = object.GetName()
	m.Namespace = object.GetNamespace()
	m.Image = spec["image"].(string)
	if sharedVolume, ok := spec["sharedVolume"].(string); ok {
		m.SharedVolume = sharedVolume
	}
	if privateVolumeSize, ok := spec["privateVolumeSize"].(string); ok {
		m.PrivateVolumeSize = privateVolumeSize
	}
	return m
}
