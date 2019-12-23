package crd

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
)

type Service struct {
	Name              string
	Namespace         string
	Image             string
	SharedVolume      string
	PrivateVolumeSize string
	UID               types.UID
}

func Parse(obj interface{}) *Service {
	m := new(Service)
	object := obj.(*unstructured.Unstructured)
	spec := object.Object["spec"].(map[string]interface{})
	m.Name = object.GetName()
	m.Namespace = object.GetNamespace()
	m.UID = object.GetUID()
	m.Image = spec["image"].(string)
	if sharedVolume, ok := spec["sharedVolume"].(string); ok {
		m.SharedVolume = sharedVolume
	}
	if privateVolumeSize, ok := spec["privateVolumeSize"].(string); ok {
		m.PrivateVolumeSize = privateVolumeSize
	}
	return m
}
