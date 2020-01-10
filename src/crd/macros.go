package crd

import (
	"encoding/json"
	"fmt"

	v1 "argovue/apis/argovue.io/v1"
)

func GetByKind(kind, namespace, name string) (crd *Crd, err error) {
	err = nil
	switch kind {
	case "pod":
		crd = New("", "v1", "pods")
	case "pvc":
		crd = New("", "v1", "persistentvolumeclaims")
	case "service":
		crd = New("", "v1", "services")
	case "workflow":
		crd = New("argoproj.io", "v1alpha1", "workflows")
	case "catalogue":
		crd = New("argovue.io", "v1", "services")
	case "instance":
		crd = New("helm.fluxcd.io", "v1", "helmreleases")
	default:
		return nil, fmt.Errorf("Can't create crd by kind:%s", kind)
	}
	crd.SetFieldSelector(fmt.Sprintf("metadata.name=%s,metadata.namespace=%s", name, namespace))
	return
}

func WorkflowPods(wfName, pod string) *Crd {
	return New("", "v1", "pods").
		SetLabelSelector("workflows.argoproj.io/workflow=" + wfName).
		SetFieldSelector("metadata.name=" + pod)
}

func WorkflowServices(wfName string) *Crd {
	return New("helm.fluxcd.io", "v1", "helmreleases").
		SetLabelSelector("workflows.argoproj.io/workflow=" + wfName)
}

func WorkflowMounts(wfName string) *Crd {
	return New("", "v1", "services").
		SetLabelSelector("workflows.argoproj.io/workflow=" + wfName)
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
	return New("helm.fluxcd.io", "v1", "helmreleases").
		SetLabelSelector("service.argovue.io/name=" + name)
}

func CatalogueResources(name string) *Crd {
	return New("", "v1", "pods").
		SetLabelSelector("app.kubernetes.io/name=" + name)
}

func CatalogueInstancePods(name string) *Crd {
	return New("", "v1", "pods").
		SetLabelSelector("app.kubernetes.io/instance=" + name)
}

func CatalogueInstanceServices(name string) *Crd {
	return New("", "v1", "services").
		SetLabelSelector("app.kubernetes.io/instance=" + name)
}

func CatalogueInstancePvcs(name string) *Crd {
	return New("", "v1", "persistentvolumeclaims").
		SetLabelSelector("app.kubernetes.io/instance=" + name)
}

func CatalogueInstance(name, instance string) *Crd {
	return New("helm.fluxcd.io", "v1", "helmreleases").
		SetLabelSelector("service.argovue.io/name=" + name).
		SetFieldSelector("metadata.name=" + instance)
}

func Typecast(thing interface{}) (*v1.Service, error) {
	if thing == nil {
		return nil, fmt.Errorf("Service typecast nil input")
	}
	buf, err := json.Marshal(thing)
	if err != nil {
		return nil, err
	}
	svc := new(v1.Service)
	err = json.Unmarshal(buf, svc)
	if err != nil {
		return nil, err
	}
	return svc, nil
}

func TypecastConfig(thing interface{}) (*v1.AppConfig, error) {
	if thing == nil {
		return nil, fmt.Errorf("Service typecast nil input")
	}
	buf, err := json.Marshal(thing)
	if err != nil {
		return nil, err
	}
	cfg := new(v1.AppConfig)
	err = json.Unmarshal(buf, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
