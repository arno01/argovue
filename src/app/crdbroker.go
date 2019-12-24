package app

import (
	"argovue/broker"
	"argovue/crd"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
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
	return a.addBroker(sessionId, name, NewCrdBroker(name))
}

func (a *App) getBroker(sessionId, name string) *CrdBroker {
	if nameMap, ok := a.brokers[sessionId]; ok {
		return nameMap[name]
	}
	return nil
}

func (a *App) newSubsetBroker(sessionId, name string) *CrdBroker {
	nameMap, ok := a.subset[sessionId]
	if !ok {
		nameMap = make(map[string]*CrdBroker)
		a.subset[sessionId] = nameMap
	}
	nameMap[name] = NewCrdBroker(name)
	return nameMap[name]
}

func (a *App) getSubsetBroker(sessionId string, keys ...string) *CrdBroker {
	id := strings.Join(keys, "|")
	nameMap, ok := a.subset[sessionId]
	if ok {
		return nameMap[id]
	}
	nameMap = make(map[string]*CrdBroker)
	if len(keys) == 3 {
		nameMap[id] = a.createNameBroker(id, keys[0], keys[1], keys[2])
		return nameMap[id]
	}
	if len(keys) == 4 {
		nameMap[id] = a.createNameSubsetBroker(id, keys[0], keys[1], keys[2], keys[3])
		return nameMap[id]
	}
	return nil
}

func (a *App) createNameBroker(id, namespace, kind, name string) (broker *CrdBroker) {
	broker = NewCrdBroker(id)
	fieldSelector := fmt.Sprintf("metadata.name=%s", name)
	switch kind {
	case "pods":
		broker.AddCrd(crd.New("", "v1", "pods").SetFieldSelector(fieldSelector))
	case "services":
		broker.AddCrd(crd.New("", "v1", "services").SetFieldSelector(fieldSelector))
	case "catalogue":
		broker.AddCrd(crd.New("argovue.io", "v1", "services").SetFieldSelector(fieldSelector))
	case "workflows":
		broker.AddCrd(crd.New("argoproj.io", "v1alpha1", "workflows").SetFieldSelector(fieldSelector))
	default:
		log.Errorf("Unknown resource to watch by name, id:%s, kind:%s", id, kind)
	}
	return
}

func (a *App) createNameSubsetBroker(id, namespace, kind, name, subset string) (broker *CrdBroker) {
	broker = NewCrdBroker(id)
	switch kind {
	case "catalogue":
		broker.AddCrd(crd.New("", "v1", "services").SetLabelSelector(subset))
	default:
		log.Errorf("createNameSubsetBroker: unknown resource to watch by name, id:%s, kind:%s", id, kind)
	}
	return
}

func NewCrdBroker(id string) *CrdBroker {
	cb := new(CrdBroker)
	cb.broker = broker.New(id)
	return cb
}

func (cb *CrdBroker) AddCrd(crd *crd.Crd) *CrdBroker {
	id := crd.Id()
	for _, i := range cb.crd {
		if id == i.Id() {
			return cb
		}
	}
	cb.crd = append(cb.crd, crd)
	go func() {
		for msg := range crd.Notifier() {
			cb.broker.Notifier <- msg
		}
	}()
	crd.Watch()
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
