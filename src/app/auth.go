package app

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc"
	log "github.com/sirupsen/logrus"
)

func (a *App) onLogout(sessionId string) {
	sessionData, ok := a.brokers[sessionId]
	if !ok {
		return
	}
	for name, _ := range sessionData {
		log.Debugf("Delete broker: %s", name)
		if broker, ok := sessionData[name]; ok {
			broker.Stop()
			delete(sessionData, name)
		}
	}
}

func (a *App) onLogin(sessionId string, profile map[string]interface{}) {
	groups := profile["groups"].([]interface{})
	wfBroker := a.newBroker(sessionId, "workflows")
	svcBroker := a.newBroker(sessionId, "services")
	catBroker := a.newBroker(sessionId, "catalogue")
	a.newBroker(sessionId, "pods")
	if len(groups) > 0 {
		strGroups := []string{}
		for _, group := range groups {
			if strGroup, ok := group.(string); ok {
				strGroups = append(strGroups, strGroup)
			}
		}
		selector := fmt.Sprintf("oidc.argovue.io/group in (%s)", strings.Join(strGroups, ","))
		wfBroker.AddCrd("argoproj.io", "v1alpha1", "workflows", selector)
		svcBroker.AddCrd("", "v1", "services", selector)
		catBroker.AddCrd("argovue.io", "v1", "services", selector)
	}
	if sub, ok := profile["sub"].(string); ok {
		selector := fmt.Sprintf("oidc.argovue.io/id in (%s)", sub)
		wfBroker.AddCrd("argoproj.io", "v1alpha1", "workflows", selector)
		svcBroker.AddCrd("", "v1", "services", selector)
		catBroker.AddCrd("argovue.io", "v1", "services", selector)
	}
}

// Profile returns user's profile if autorised
func (a *App) Profile(w http.ResponseWriter, r *http.Request) {
	var profile interface{}
	session, err := a.Store().Get(r, "auth-session")
	if err != nil {
		profile = new(map[string]interface{})
	} else {
		profile = session.Values["profile"]
	}
	obj, err := json.Marshal(profile)
	w.Write([]byte(obj))
}

// Logout clears authentication
func (a *App) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := a.Store().Get(r, "auth-session")
	if err != nil {
		http.Redirect(w, r, a.Args().UIRootURL(), http.StatusFound)
		return
	}
	a.onLogout(session.ID)
	delete(session.Values, "state")
	delete(session.Values, "auth-session")
	delete(session.Values, "profile")
	session.Options.MaxAge = -1
	if err = session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, a.Args().UIRootURL(), http.StatusFound)
}

// AuthInitiate initialises OIDC auth sequence by redirecting browser to OIDC provider
func (a *App) AuthInitiate(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	state := base64.StdEncoding.EncodeToString(b)
	session, _ := a.Store().Get(r, "auth-session")
	session.Values["state"] = state
	if err = session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, a.Auth().Config.AuthCodeURL(state), http.StatusTemporaryRedirect)
}

// AuthCallback processes OIDC provider response with state and code parameters
func (a *App) AuthCallback(w http.ResponseWriter, r *http.Request) {
	session, err := a.Store().Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.URL.Query().Get("state") != session.Values["state"] {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	token, err := a.Auth().Config.Exchange(context.TODO(), r.URL.Query().Get("code"))
	if err != nil {
		log.Errorf("no token found: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
		return
	}

	idToken, err := a.Auth().Provider.Verifier(&oidc.Config{ClientID: a.Auth().Config.ClientID}).Verify(context.TODO(), rawIDToken)
	if err != nil {
		http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Debugf("OIDC: auth name:%v", profile["name"])
	session.Values["profile"] = profile
	a.onLogin(session.ID, profile)

	if err = session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, a.Args().UIRootURL(), http.StatusSeeOther)
}
