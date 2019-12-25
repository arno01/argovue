package crd

import "fmt"

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
		SetLabelSelector("argovue.io/service=" + name)
}

func CatalogueInstance(name, instance string) *Crd {
	return New("", "v1", "services").
		SetLabelSelector(fmt.Sprintf("argovue.io/service=%s,service=%s", name, instance))
}
