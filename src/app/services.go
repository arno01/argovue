package app

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

func serviceId(s corev1.Service) string {
	return fmt.Sprintf("%s/%s", s.Namespace, s.Name)
}

func (a *App) watchServices(cb *CrdBroker) {
	for msg := range cb.crd.Notifier() {
		m := msg.Content.(corev1.Service)
		cb.broker.Notifier <- msg
		switch msg.Action {
		case "add":
			log.Infof("adding service %s", serviceId(m))
		case "delete":
			log.Infof("deleting service %s", serviceId(m))
		case "update":
			log.Infof("updating service %s", serviceId(m))
		default:
		}
	}
}
