package app

import (
	"argovue/broker"
	"argovue/crd"
	"fmt"
)

type CrdBroker struct {
	crd    []*crd.Crd
	broker *broker.Broker
}

func (a *App) addBroker(sessionId, name string, broker *CrdBroker) *CrdBroker {
	nameMap, ok := a.brokers[sessionId]
	if !ok {
		nameMap = make(map[string]*CrdBroker)
		a.brokers[sessionId] = nameMap
	}
	nameMap[name] = broker
	return broker
}

func (a *App) deleteBroker(sessionId, name string) {
	nameMap, ok := a.brokers[sessionId]
	if !ok {
		return
	}
	delete(nameMap, name)
}

func (a *App) newBroker(sessionId, name string) *CrdBroker {
	return a.addBroker(sessionId, name, NewCrdBroker(fmt.Sprintf("%s", name)))
}

func (a *App) getBroker(sessionId, name string) *CrdBroker {
	if nameMap, ok := a.brokers[sessionId]; ok {
		return nameMap[name]
	}
	return nil
}

func NewCrdBroker(id string) *CrdBroker {
	cb := new(CrdBroker)
	cb.broker = broker.New(id)
	return cb
}

func (cb *CrdBroker) AddCrd(group, version, name, selector string) *CrdBroker {
	id := crd.MakeId(group, version, name, selector)
	for _, i := range cb.crd {
		if id == i.Id() {
			return cb
		}
	}
	crd := crd.New(group, version, name, selector)
	cb.crd = append(cb.crd, crd)
	go func() {
		for msg := range crd.Notifier() {
			cb.broker.Notifier <- msg
		}
	}()
	return cb
}

func (cb *CrdBroker) Stop() {
	for _, crd := range cb.crd {
		crd.Stop()
	}
	close(cb.broker.Notifier)
}

func (cb *CrdBroker) Broker() *broker.Broker {
	return cb.broker
}
