package app

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func (a *App) ProxyService(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	namespace := mux.Vars(r)["namespace"]
	port := mux.Vars(r)["port"]
	rest := mux.Vars(r)["rest"]
	schema := "http"
	if port == "443" {
		schema = "https"
	}
	target := fmt.Sprintf("%s://%s.%s.svc.cluster.local:%s", schema, name, namespace, port)
	a.Proxy(name, namespace, port, rest, target, w, r)
}

func (a *App) Proxy(name, namespace, port, rest, target string, w http.ResponseWriter, r *http.Request) {
	log.Debugf("Proxy to service target:%s, path:%s", target, rest)
	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)
	r.URL.Host = url.Host
	r.URL.Scheme = url.Scheme
	r.URL.Path = rest
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = url.Host
	proxy.ServeHTTP(w, r)
}
