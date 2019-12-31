package app

import (
	"net/http"

	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func i2s(i interface{}) (re string) {
	if s, ok := i.(string); ok {
		re = s
	}
	return
}

func li2s(li interface{}) (re []string) {
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

func authorizeByGroup(groupLabel string, groups []string) bool {
	if len(groupLabel) == 0 {
		return false
	}
	for _, group := range groups {
		if group == groupLabel {
			return true
		}
	}
	return false
}

func authorizeById(idLabel, id string) bool {
	if len(idLabel) == 0 || len(id) == 0 {
		return false
	}
	return idLabel == id
}

func authorize(labels map[string]string, profile map[string]interface{}) bool {
	var auth bool
	if groupLabel, ok := labels["oidc.argovue.io/group"]; ok {
		if groups, ok := profile["groups"]; ok {
			auth = authorizeByGroup(groupLabel, li2s(groups))
			if auth {
				log.Debugf("authorize by group:%s", groupLabel)
			}
		}
	}
	if auth {
		return auth
	}
	if idLabel, ok := labels["oidc.argovue.io/id"]; ok {
		if id, ok := profile["sub"]; ok {
			auth = authorizeById(idLabel, i2s(id))
			if auth {
				log.Debugf("authorize by id:%s", id)
			}
		}
	}
	return auth
}

func authHTTP(obj metav1.Object, store sessions.Store, r *http.Request) bool {
	session, err := store.Get(r, "auth-session")
	if err != nil {
		log.Debugf("authHTTP: no session")
		return false
	}
	pI, ok := session.Values["profile"]
	if !ok {
		log.Debugf("authHTTP: no profile")
		return false
	}
	profile, ok := pI.(map[string]interface{})
	if !ok {
		log.Debugf("authHTTP: invalid profile")
		return false
	}
	return authorize(obj.GetLabels(), profile)
}
