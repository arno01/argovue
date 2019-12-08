package auth

import (
	"context"
	"strings"

	"github.com/coreos/go-oidc"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// Auth collects OIDC data
type Auth struct {
	Provider *oidc.Provider
	Config   oauth2.Config
	Ctx      context.Context
}

// New creates Auth instance
func New(oidcProvider, oidcClientID, oidcClientSecret, oidcRedirectURL, oidcScopes string) *Auth {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, oidcProvider)
	if err != nil {
		log.Fatalf("Failed to get provider: %v", err)
		return nil
	}

	scopes := append(strings.Split(oidcScopes, " "), oidc.ScopeOpenID, "profile")

	conf := oauth2.Config{
		ClientID:     oidcClientID,
		ClientSecret: oidcClientSecret,
		RedirectURL:  oidcRedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       scopes,
	}

	return &Auth{
		Provider: provider,
		Config:   conf,
		Ctx:      ctx,
	}
}
