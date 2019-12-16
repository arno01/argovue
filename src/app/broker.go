package app

func (a *App) addBroker(name, namespace string, broker *CrdBroker) *CrdBroker {
	nameMap, ok := a.brokers[namespace]
	if !ok {
		nameMap = make(map[string]*CrdBroker)
		a.brokers[namespace] = nameMap
	}
	nameMap[name] = broker
	return broker
}

func (a *App) deleteBroker(name, namespace string) {
	nameMap, ok := a.brokers[namespace]
	if !ok {
		return
	}
	delete(nameMap, name)
}

func (a *App) newBroker(group, version, name, namespace string) *CrdBroker {
	return a.addBroker(name, namespace, NewBroker(group, version, name, namespace))
}

func (a *App) getBroker(name, namespace string) *CrdBroker {
	if nameMap, ok := a.brokers[namespace]; ok {
		return nameMap[name]
	} else {
		return nil
	}
}
