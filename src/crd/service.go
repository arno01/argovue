package crd

import (
	"fmt"

	"argovue/kube"

	v1 "argovue/apis/argovue.io/v1"
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

func createPVC(s *v1.ServiceType, clientset *kubernetes.Clientset, instance, owner string) (*apiv1.PersistentVolumeClaim, error) {
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

func createService(s *v1.ServiceType, clientset *kubernetes.Clientset, instance, owner string) (*apiv1.Service, error) {
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
				"oidc.argovue.io/id":          owner
			},
		},
	}
	return clientset.CoreV1().Services(s.Namespace).Create(svc)
}

func createDeployment(s *v1.ServiceType, clientset *kubernetes.Clientset, instance, owner, baseUrl string, pvc *apiv1.PersistentVolumeClaim) (*appsv1.Deployment, error) {
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
							VolumeMounts: []apiv1.VolumeMount{{Name: "work", MountPath: "/work"}},
							Env:          []apiv1.EnvVar{{Name: "BASE_URL", Value: baseUrl}},
							Ports:        []apiv1.ContainerPort{{Name: "http", Protocol: apiv1.ProtocolTCP, ContainerPort: s.Spec.Port}},
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

func Deploy(s *v1.ServiceType, owner string) error {
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
	pvc, err := createPVC(s, clientset, instance, owner)
	if err != nil {
		log.Errorf("Kube: can't create pvc, error:%s", err)
		return err
	}
	_, err = createService(s, clientset, instance, owner)
	if err != nil {
		log.Errorf("Kube: can't create service, error:%s", err)
		return err
	}
	_, err = createDeployment(s, clientset, instance, owner, baseUrl, pvc)
	if err != nil {
		log.Errorf("Kube: can't create deployment, error:%s", err)
		return err
	}
	return nil
}

func Delete(s *v1.ServiceType, instance string) error {
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
