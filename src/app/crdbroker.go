package app

import (
	"kubevue/broker"
	"kubevue/crd"
)

type CrdBroker struct {
	crd    *crd.Crd
	broker *broker.Broker
}

func NewBroker(group, version, name, namespace string) *CrdBroker {
	crd := crd.New(group, version, name, namespace)
	broker := broker.New()
	return &CrdBroker{crd, broker}
}

func (cb *CrdBroker) Stop() {
	cb.crd.Stop()
	close(cb.broker.Notifier)
}

func (cb *CrdBroker) PassMessages() {
	go func() {
		for msg := range cb.crd.Notifier() {
			cb.broker.Notifier <- msg
		}
	}()
}