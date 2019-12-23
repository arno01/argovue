package app

import (
	"argovue/broker"
	"argovue/crd"
)

type CrdBroker struct {
	crd    []*crd.Crd
	broker *broker.Broker
}

func NewCrdBroker(id string) *CrdBroker {
	cb := new(CrdBroker)
	cb.broker = broker.New(id)
	return cb
}

func (cb *CrdBroker) AddCrd(group, version, name, selector string) *CrdBroker {
	cb.crd = append(cb.crd, crd.New(group, version, name, selector))
	return cb
}

func (cb *CrdBroker) Stop() {
	for _, crd := range cb.crd {
		crd.Stop()
		close(cb.broker.Notifier)
	}
}

func (cb *CrdBroker) PassMessages() {
	for _, t := range cb.crd {
		go func(crd *crd.Crd) {
			for msg := range crd.Notifier() {
				cb.broker.Notifier <- msg
			}
		}(t)
	}
}
