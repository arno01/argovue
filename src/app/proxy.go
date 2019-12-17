package app

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func (a *App) ProxyService(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	namespace := mux.Vars(r)["namespace"]
	target := fmt.Sprintf("http://%s.%s.svc.cluster.local", name, namespace)
	a.Proxy(name, namespace, target, w, r)
}

func (a *App) Proxy(name, namespace, target string, w http.ResponseWriter, r *http.Request) {
	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)
	r.URL.Host = url.Host
	r.URL.Scheme = url.Scheme
	newPath := regexp.MustCompile(fmt.Sprintf("^/proxy/%s/%s", namespace, name)).ReplaceAllString(r.URL.Path, "")
	log.Debugf("Rewrote URL to:%s", newPath)
	r.URL.Path = newPath
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = url.Host
	proxy.ServeHTTP(w, r)
}
