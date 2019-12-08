package args

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"
)

// Args type
type Args struct {
	verboseLevel     string
	port             int
	bindAddr         string
	dir              string
	args             []string
	oidcProvider     string
	oidcClientID     string
	oidcClientSecret string
	oidcRedirectURL  string
	oidcScopes       string
	uiRootURL        string
}

// New type
func New() *Args {
	return new(Args).Parse()
}

// Parse parameters
func (a *Args) Parse() *Args {
	flag.StringVar(&a.verboseLevel, "verbose", "info", "Set verbosity level")
	flag.IntVar(&a.port, "port", 8080, "Listen port")
	flag.StringVar(&a.bindAddr, "bind", "", "Bind address")
	flag.StringVar(&a.dir, "dir", "ui/dist", "Static files folder")
	flag.StringVar(&a.oidcProvider, "oidc-provider", os.Getenv("OIDC_PROVIDER"), "OIDC provider")
	flag.StringVar(&a.oidcClientID, "oidc-client-id", os.Getenv("OIDC_CLIENT_ID"), "OIDC client id")
	flag.StringVar(&a.oidcClientSecret, "oidc-client-secret", os.Getenv("OIDC_CLIENT_SECRET"), "OIDC client secret")
	flag.StringVar(&a.oidcRedirectURL, "oidc-redirect-url", os.Getenv("OIDC_REDIRECT_URL"), "OIDC redirect url")
	flag.StringVar(&a.oidcScopes, "oidc-scopes", os.Getenv("OIDC_SCOPES"), "OIDC scopes")
	flag.StringVar(&a.uiRootURL, "ui-root-url", os.Getenv("UI_ROOT_URL"), "UI root url for redirects")

	flag.Parse()
	a.args = flag.Args()
	return a
}

// Dir to serve files from
func (a *Args) Dir() string {
	return a.dir
}

// BindAddr to bind to Web Server
func (a *Args) BindAddr() string {
	return a.bindAddr
}

// Port to bind to Web Server
func (a *Args) Port() int {
	return a.port
}

// UIRootURL returns UI root url
func (a *Args) UIRootURL() string {
	return a.uiRootURL
}

// OIDC returns OIDC parameters
func (a *Args) OIDC() (string, string, string, string, string) {
	return a.oidcProvider, a.oidcClientID, a.oidcClientSecret, a.oidcRedirectURL, a.oidcScopes
}

// LogLevel set loglevel
func (a *Args) LogLevel() *Args {
	switch a.verboseLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
	return a
}
