package crd

import (
	"fmt"

	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func int32Ptr(i int32) *int32 { return &i }

func (s *Service) createPVC(clientset *kubernetes.Clientset, instance string, owner string) (*apiv1.PersistentVolumeClaim, error) {
	log.Debugf("Kube: create pvc %s/%s for %s", s.Namespace, s.Name, owner)
	pvc := &apiv1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-", s.Name),
			Namespace:    s.Namespace,
			Labels:       map[string]string{"service": instance, "oidc.argovue.io/id": owner, "argovue.io/service": s.Name},
		},
		Spec: apiv1.PersistentVolumeClaimSpec{
			AccessModes: []apiv1.PersistentVolumeAccessMode{apiv1.ReadWriteOnce},
			Resources: apiv1.ResourceRequirements{
				Requests: apiv1.ResourceList{
					apiv1.ResourceName(apiv1.ResourceStorage): resource.MustParse(s.PrivateVolumeSize),
				},
			},
		},
	}
	return clientset.CoreV1().PersistentVolumeClaims(s.Namespace).Create(pvc)
}

func (s *Service) createService(clientset *kubernetes.Clientset, name, instance, owner string) (*apiv1.Service, error) {
	log.Debugf("Kube: create service %s/%s for %s", s.Namespace, s.Name, owner)
	svc := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: s.Namespace,
			Labels:    map[string]string{"service": instance, "oidc.argovue.io/id": owner, "argovue.io/service": s.Name},
		},
		Spec: apiv1.ServiceSpec{
			Ports:    []apiv1.ServicePort{{Port: 80, Protocol: apiv1.ProtocolTCP}},
			Type:     apiv1.ServiceTypeClusterIP,
			Selector: map[string]string{"service": instance, "oidc.argovue.io/id": owner},
		},
	}
	return clientset.CoreV1().Services(s.Namespace).Create(svc)
}

func (s *Service) Deploy(clientset *kubernetes.Clientset, owner string) (*appsv1.Deployment, error) {
	instance := fmt.Sprintf("%s-%s", s.Name, xid.New().String())
	baseUrl := fmt.Sprintf("/proxy/%s/%s/%d", s.Namespace, instance, 80)
	log.Debugf("Kube: create deployment %s/%s service:%s instance:%s %s", s.Namespace, s.Name, owner, instance, s.Args)
	for i, arg := range s.Args {
		if arg == "BASE_URL" {
			s.Args[i] = baseUrl
		}
	}
	pvc, err := s.createPVC(clientset, instance, owner)
	if err != nil {
		log.Errorf("Kube: can't create pvc, error:%s", err)
		return nil, err
	}
	_, err = s.createService(clientset, instance, instance, owner)
	if err != nil {
		log.Errorf("Kube: can't create service, error:%s", err)
		return nil, err
	}
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName:    fmt.Sprintf("%s-", s.Name),
			Namespace:       s.Namespace,
			Labels:          map[string]string{"service": instance, "oidc.argovue.io/id": owner, "argovue.io/service": s.Name},
			OwnerReferences: []metav1.OwnerReference{{APIVersion: "argovue.io/v1", Kind: "Service", Name: s.Name, UID: s.UID}},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"service": instance, "oidc.argovue.io/id": owner, "argovue.io/service": s.Name}},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"service": instance, "oidc.argovue.io/id": owner, "argovue.io/service": s.Name}},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:         s.Name,
							Image:        s.Image,
							Args:         s.Args,
							VolumeMounts: []apiv1.VolumeMount{{Name: "work", MountPath: "/work"}},
							Env:          []apiv1.EnvVar{{Name: "BASE_URL", Value: baseUrl}},
							Ports:        []apiv1.ContainerPort{{Name: "http", Protocol: apiv1.ProtocolTCP, ContainerPort: 80}},
						},
					},
					Volumes: []apiv1.Volume{
						{
							Name:         "work",
							VolumeSource: apiv1.VolumeSource{PersistentVolumeClaim: &apiv1.PersistentVolumeClaimVolumeSource{ClaimName: pvc.GetName()}},
						},
					},
				},
			},
		},
	}
	return clientset.AppsV1().Deployments(s.Namespace).Create(deployment)
}
