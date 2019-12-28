package crd

import (
	"fmt"

	"argovue/kube"

	argovuev1 "argovue/apis/argovue.io/v1"

	wfv1alpha1 "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

func int32Ptr(i int32) *int32 { return &i }

func maybeReadOnly(volumeName string) bool {
	if volumeName == "private" {
		return false
	}
	return true
}

func makeVolumeMounts(volumes []apiv1.Volume) []apiv1.VolumeMount {
	mounts := []apiv1.VolumeMount{}
	for _, v := range volumes {
		mounts = append(mounts, apiv1.VolumeMount{Name: v.Name, MountPath: "/mnt/" + v.Name, ReadOnly: maybeReadOnly(v.Name)})
	}
	return mounts
}

func createPVC(s *argovuev1.Service, clientset *kubernetes.Clientset, instance, owner string) (*apiv1.PersistentVolumeClaim, error) {
	log.Debugf("Kube: create pvc %s/%s, owner:%s", s.Namespace, instance, owner)
	pvc := &apiv1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance,
			Namespace: s.Namespace,
			Labels: map[string]string{
				"service.argovue.io/name":     s.Name,
				"service.argovue.io/instance": instance,
				"oidc.argovue.io/id":          owner,
			},
		},
		Spec: apiv1.PersistentVolumeClaimSpec{
			AccessModes: []apiv1.PersistentVolumeAccessMode{apiv1.ReadWriteOnce},
			Resources: apiv1.ResourceRequirements{
				Requests: apiv1.ResourceList{
					apiv1.ResourceName(apiv1.ResourceStorage): resource.MustParse(s.Spec.PrivateVolumeSize),
				},
			},
		},
	}
	return clientset.CoreV1().PersistentVolumeClaims(s.Namespace).Create(pvc)
}

func createService(s *argovuev1.Service, clientset *kubernetes.Clientset, instance, owner string) (*apiv1.Service, error) {
	log.Debugf("Kube: create service %s/%s, owner:%s", s.Namespace, instance, owner)
	svc := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance,
			Namespace: s.Namespace,
			Labels: map[string]string{
				"service.argovue.io/name":     s.Name,
				"service.argovue.io/instance": instance,
				"oidc.argovue.io/id":          owner,
			},
		},
		Spec: apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{{Port: 80, Protocol: apiv1.ProtocolTCP, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: s.Spec.Port}}},
			Type:  apiv1.ServiceTypeClusterIP,
			Selector: map[string]string{
				"service.argovue.io/instance": instance,
				"oidc.argovue.io/id":          owner,
			},
		},
	}
	return clientset.CoreV1().Services(s.Namespace).Create(svc)
}

func createDeployment(s *argovuev1.Service, clientset *kubernetes.Clientset, instance, owner, baseUrl string, volumes []apiv1.Volume) (*appsv1.Deployment, error) {
	log.Debugf("Kube: create deployment %s/%s, owner:%s", s.Namespace, instance, owner)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance,
			Namespace: s.Namespace,
			Labels: map[string]string{
				"service.argovue.io/name":     s.Name,
				"service.argovue.io/instance": instance,
				"oidc.argovue.io/id":          owner,
			},
			OwnerReferences: []metav1.OwnerReference{{APIVersion: "argovue.io/v1", Kind: "Service", Name: s.Name, UID: s.UID}},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{
				"service.argovue.io/name":     s.Name,
				"service.argovue.io/instance": instance,
				"oidc.argovue.io/id":          owner,
			}},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{
					"service.argovue.io/name":     s.Name,
					"service.argovue.io/instance": instance,
					"oidc.argovue.io/id":          owner,
				}},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:         s.Name,
							Image:        s.Spec.Image,
							Args:         s.Spec.Args,
							VolumeMounts: makeVolumeMounts(volumes),
							Env:          []apiv1.EnvVar{{Name: "BASE_URL", Value: baseUrl}},
							Ports:        []apiv1.ContainerPort{{Name: "http", Protocol: apiv1.ProtocolTCP, ContainerPort: s.Spec.Port}},
						},
					},
					Volumes: volumes,
				},
			},
		},
	}
	return clientset.AppsV1().Deployments(s.Namespace).Create(deployment)
}

func Deploy(s *argovuev1.Service, owner string) error {
	clientset, err := kube.GetClient()
	if err != nil {
		return err
	}
	instance := fmt.Sprintf("%s-%s", s.Name, xid.New().String())
	baseUrl := fmt.Sprintf("/proxy/%s/%s/%d", s.Namespace, instance, 80)
	log.Debugf("Kube: create service %s/%s owner:%s instance:%s args:%s", s.Namespace, s.Name, owner, instance, s.Spec.Args)
	for i, arg := range s.Spec.Args {
		if arg == "BASE_URL" {
			s.Spec.Args[i] = baseUrl
		}
	}

	volumes := []apiv1.Volume{}

	if len(s.Spec.PrivateVolumeSize) > 0 {
		pvc, err := createPVC(s, clientset, instance, owner)
		if err != nil {
			log.Errorf("Kube: can't create pvc, error:%s", err)
		}
		volumes = append(volumes, apiv1.Volume{
			Name:         "private",
			VolumeSource: apiv1.VolumeSource{PersistentVolumeClaim: &apiv1.PersistentVolumeClaimVolumeSource{ClaimName: pvc.GetName()}},
		})
	}
	if len(s.Spec.SharedVolume) > 0 {
		volumes = append(volumes, apiv1.Volume{
			Name:         "shared",
			VolumeSource: apiv1.VolumeSource{PersistentVolumeClaim: &apiv1.PersistentVolumeClaimVolumeSource{ClaimName: s.Spec.SharedVolume}},
		})
	}
	_, err = createService(s, clientset, instance, owner)
	if err != nil {
		log.Errorf("Kube: can't create service, error:%s", err)
		return err
	}
	_, err = createDeployment(s, clientset, instance, owner, baseUrl, volumes)
	if err != nil {
		log.Errorf("Kube: can't create deployment, error:%s", err)
		return err
	}
	return nil
}

func Delete(s *argovuev1.Service, instance string) error {
	clientset, err := kube.GetClient()
	if err != nil {
		return err
	}
	deletePolicy := metav1.DeletePropagationForeground
	opts := &metav1.DeleteOptions{PropagationPolicy: &deletePolicy}
	clientset.CoreV1().PersistentVolumeClaims(s.Namespace).Delete(instance, opts)
	clientset.CoreV1().Services(s.Namespace).Delete(instance, opts)
	clientset.AppsV1().Deployments(s.Namespace).Delete(instance, opts)
	return nil
}

func DeployFilebrowser(wf *wfv1alpha1.Workflow, owner string) error {
	clientset, err := kube.GetV1Clientset()
	if err != nil {
		return err
	}
	for _, pvc := range wf.Status.PersistentVolumeClaims {
		log.Infof("CRD: Deploy filebrowser service for workflow:%s, pvc:%s", wf.GetName(), pvc.Name)
		svc := argovuev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: wf.GetNamespace(),
				Name:      wf.GetName(),
				Labels: map[string]string{
					"oidc.argovue.io/id":        owner,
					"workflow.argoproj.io/name": wf.GetName(),
				},
			},
			Spec: argovuev1.ServiceSpec{
				Image:        "filebrowser/filebrowser:latest",
				Args:         []string{"--noauth", "--root", "/mnt", "--baseurl", "BASE_URL"},
				Port:         80,
				SharedVolume: pvc.PersistentVolumeClaim.ClaimName,
			},
		}
		if obj, err := clientset.ArgovueV1().Services(wf.GetNamespace()).Create(&svc); err != nil {
			log.Errorf("CRD: DeployFilebrowser error:%s", err)
		} else {
			Deploy(obj, owner)
		}
	}
	return nil
}
