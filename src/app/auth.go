package app

import (
	"argovue/crd"
	"argovue/util"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/coreos/go-oidc"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
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
	wfBroker := a.newBroker(sessionId, "workflows")
	catBroker := a.newBroker(sessionId, "catalogue")
	if groups, ok := profile["effective_groups"].([]string); ok && len(groups) > 0 {
		selector := fmt.Sprintf("oidc.argovue.io/group in (%s)", strings.Join(groups, ","))
		wfBroker.AddCrd(crd.New("argoproj.io", "v1alpha1", "workflows").SetLabelSelector(selector))
		catBroker.AddCrd(crd.New("argovue.io", "v1", "services").SetLabelSelector(selector))
	}
	if userId, ok := profile["effective_id"]; ok {
		label := util.EncodeLabel(util.I2s(userId))
		log.Debugf("App: using user label:%s", label)
		selector := fmt.Sprintf("oidc.argovue.io/id in (%s)", label)
		wfBroker.AddCrd(crd.New("argoproj.io", "v1alpha1", "workflows").SetLabelSelector(selector))
		catBroker.AddCrd(crd.New("argovue.io", "v1", "services").SetLabelSelector(selector))
	}
}

// Profile returns user's profile if autorised
func (a *App) Profile(w http.ResponseWriter, r *http.Request) {
	var profile interface{}
	session, err := a.Store().Get(r, "auth-session")
	if err == nil {
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
		log.Errorf("Can't delete session, error:%s", err)
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

	redirect := r.URL.Query().Get("redirect")

	state := base64.StdEncoding.EncodeToString(b)
	session, _ := a.Store().Get(r, "auth-session")
	session.Values["state"] = state
	if len(redirect) > 0 {
		unescape, err := url.PathUnescape(redirect)
		if err == nil {
			log.Debugf("AUTH: keep redirect value:%s", unescape)
			session.Values["redirect"] = unescape
		} else {
			log.Debugf("AUTH: error unescape path:%s, error:%s", redirect, err)
		}
	}

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

	var idTokenClaims map[string]interface{}
	if err := idToken.Claims(&idTokenClaims); err != nil {
		log.Errorf("Can't decode id_token claims, error:%s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Debugf("OIDC: id token claims: %s", idTokenClaims)
	profile := idTokenClaims
	log.Debugf("OIDC: auth name:%s", idTokenClaims["name"])
	if _, ok := idTokenClaims["groups"]; !ok {
		if userInfoclaims, err := a.userInfo(token); err == nil {
			log.Debugf("OIDC: user info claims: %s", userInfoclaims)
			profile = userInfoclaims
		}
	}

	if email, ok := idTokenClaims["email"]; ok {
		profile["email"] = email
	}

	effGroups := []string{}
	for _, group := range util.Li2s(profile["groups"]) {
		if k8sGroup, ok := a.groups[group]; ok {
			effGroups = append(effGroups, k8sGroup)
		}
	}
	profile["effective_groups"] = effGroups

	userIdKey := a.Args().OidcUserId()
	if userId, ok := profile[userIdKey]; ok {
		profile["effective_id"] = userId
	} else {
		profile["effective_id"] = profile["sub"]
		log.Warnf("OIDC: can't find map effective user id by name:%s, using sub value:%s", userIdKey, profile["sub"])
	}

	session.Values["profile"] = profile
	a.onLogin(session.ID, profile)

	redirectUrl := session.Values["redirect"]
	delete(session.Values, "redirect")

	if err = session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// this is to set cookie for the api domain
	redirect := `<html><head><script type="text/javascript">window.location.href="%s"</script></head><body></body></html>`
	if re, ok := redirectUrl.(string); ok {
		log.Debugf("AUTH: redirecting to:%s", re)
		fmt.Fprintf(w, redirect, re)
	} else {
		fmt.Fprintf(w, redirect, a.Args().UIRootURL())
	}
}

func (a *App) userInfo(token *oauth2.Token) (map[string]interface{}, error) {
	var claims map[string]interface{}
	userinfo, err := a.Auth().Provider.UserInfo(context.TODO(), a.Auth().Config.TokenSource(context.TODO(), token))
	if err != nil {
		log.Errorf("Can't request user info, error:%s", err)
		return nil, err
	}
	err = userinfo.Claims(&claims)
	if err != nil {
		log.Errorf("Can't decode claims, error:%s", err)
		return nil, err
	}
	return claims, nil
}
