package crd

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Object struct {
	Group     string
	Version   string
	Name      string
	Namespace string
}

func Parse(obj interface{}) *Object {
	m := new(Object)
	spec := obj.(*unstructured.Unstructured).Object["spec"].(map[string]interface{})
	m.Group = spec["group"].(string)
	m.Version = spec["version"].(string)
	m.Name = spec["name"].(string)
	m.Namespace = spec["namespace"].(string)
	return m
}
