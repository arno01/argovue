package app

import "fmt"

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
	return a.addBroker(sessionId, name, NewCrdBroker(fmt.Sprintf("%s/%s", sessionId, name)))
}

func (a *App) getBroker(sessionId, name string) *CrdBroker {
	if nameMap, ok := a.brokers[sessionId]; ok {
		return nameMap[name]
	}
	return nil
}
