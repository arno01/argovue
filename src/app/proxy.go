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

func (a *App) ProxyDex(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	r = mux.SetURLVars(r, map[string]string{"namespace": a.args.Namespace(), "name": a.Args().DexName(), "port": "5556", "rest": vars["rest"]})
	a.ProxyService(w, r)
}

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
	if regexp.MustCompile("^/dex.*").MatchString(r.RequestURI) {
		rest = fmt.Sprintf("%s/%s", "/dex", rest)
	}
	log.Debugf("Proxy: %s to %s%s", r.URL.Path, target, rest)

	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)
	r.URL.Host = url.Host
	r.URL.Scheme = url.Scheme
	r.URL.Path = rest
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = url.Host
	proxy.ServeHTTP(w, r)
}
