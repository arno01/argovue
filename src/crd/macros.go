package crd

import (
	"encoding/json"
	"fmt"

	v1 "argovue/apis/argovue.io/v1"
)

func WorkflowPods(wfName, pod string) *Crd {
	return New("", "v1", "pods").
		SetLabelSelector("workflows.argoproj.io/workflow=" + wfName).
		SetFieldSelector("metadata.name=" + pod)
}

func Workflow(wfName string) *Crd {
	return New("argoproj.io", "v1alpha1", "workflows").
		SetFieldSelector("metadata.name=" + wfName)
}

func Catalogue(name string) *Crd {
	return New("argovue.io", "v1", "services").
		SetFieldSelector("metadata.name=" + name)
}

func CatalogueInstances(name string) *Crd {
	return New("", "v1", "services").
		SetLabelSelector("service.argovue.io/name=" + name)
}

func CatalogueInstance(name, instance string) *Crd {
	return New("", "v1", "services").
		SetLabelSelector(fmt.Sprintf("service.argovue.io/name=%s,service.argovue.io/instance=%s", name, instance))
}

func Typecast(thing interface{}) (*v1.ServiceType, error) {
	if thing == nil {
		return nil, fmt.Errorf("Service typecast nil input")
	}
	buf, err := json.Marshal(thing)
	if err != nil {
		return nil, err
	}
	svc := new(v1.ServiceType)
	err = json.Unmarshal(buf, svc)
	if err != nil {
		return nil, err
	}
	return svc, nil
}
