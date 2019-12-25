package crd

func WorkflowPods(wfName, pod string) *Crd {
	return New("", "v1", "pods").
		SetLabelSelector("workflows.argoproj.io/workflow=" + wfName).
		SetFieldSelector("metadata.name=" + pod)
}

func Workflow(wfName string) *Crd {
	return New("argoproj.io", "v1alpha1", "workflows").
		SetFieldSelector("metadata.name=" + wfName)
}
