package util

func I2s(i interface{}) (re string) {
	if i == nil {
		return
	}
	if s, ok := i.(string); ok {
		re = s
	}
	return
}

func Li2s(li interface{}) (re []string) {
	if li == nil {
		return
	}
	lii, ok := li.([]interface{})
	if !ok {
		return
	}
	for _, l := range lii {
		if s, ok := l.(string); ok {
			re = append(re, s)
		}
	}
	return
}
