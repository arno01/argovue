package crd

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
)

type Service struct {
	Name              string
	Namespace         string
	Image             string
	Args              []string
	SharedVolume      string
	PrivateVolumeSize string
	Port              int32
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
	if port, ok := spec["port"].(int64); ok {
		m.Port = int32(port)
	} else {
		m.Port = 80
	}
	if args, ok := spec["args"].([]interface{}); ok {
		for _, arg := range args {
			m.Args = append(m.Args, arg.(string))
		}
	}
	return m
}
