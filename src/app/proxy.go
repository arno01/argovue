package app

import (
	"argovue/kube"
	"argovue/profile"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func (a *App) proxyDex(sid string, p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	vars := mux.Vars(r)
	r = mux.SetURLVars(r, map[string]string{"namespace": a.args.Namespace(), "name": a.Args().DexServiceName(), "port": "5556", "rest": vars["rest"]})
	return a.proxyService(sid, p, w, r)
}

func (a *App) proxyService(sid string, p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	v := mux.Vars(r)
	name, namespace, port, rest := v["name"], v["namespace"], v["port"], v["rest"]

	if name != a.Args().DexServiceName() {
		svc, err := kube.GetService(name, namespace)
		if err != nil {
			return makeError(http.StatusForbidden, "Proxy: no service %s/%s, access denied, error:%s", namespace, name, err)
		}
		if !p.Authorize(svc) {
			return makeError(http.StatusForbidden, "Proxy: %s/%s, no auth, access denied", namespace, name)
		}
	}

	schema := "http"
	if port == "443" {
		schema = "https"
	}

	target := fmt.Sprintf("%s://%s.%s.svc.cluster.local:%s", schema, name, namespace, port)
	if regexp.MustCompile("^/dex.*").MatchString(r.RequestURI) {
		rest = fmt.Sprintf("%s/%s", "/dex", rest)
	} else {
		rest = fmt.Sprintf("/proxy/%s/%s/%s/%s", namespace, name, port, rest)
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
	return nil
}
