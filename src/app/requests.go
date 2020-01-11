package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (a *App) Version(w http.ResponseWriter, r *http.Request) {
	obj, _ := json.Marshal(a.ver)
	w.Write([]byte(obj))
}

func (a *App) Objects(w http.ResponseWriter, r *http.Request) {
	session, err := a.Store().Get(r, "auth-session")
	if err != nil {
		log.Debugf("Can't get session, error:%s", err)
	}
	var re []string
	for name, _ := range a.brokers[session.ID] {
		re = append(re, fmt.Sprintf("%s", name))
	}
	json.NewEncoder(w).Encode(re)
}

func (a *App) Profile(w http.ResponseWriter, r *http.Request) {
	var profile interface{}
	session, err := a.Store().Get(r, "auth-session")
	if err != nil {
		log.Debugf("Can't get session, error:%s", err)
	}
	profile = session.Values["profile"]
	obj, err := json.Marshal(profile)
	w.Write([]byte(obj))
}
