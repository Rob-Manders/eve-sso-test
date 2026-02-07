package auth

import (
	"encoding/base64"
	"evessotest/backend/config"
	"evessotest/backend/scopes"
	"evessotest/backend/session"
	"fmt"
	"net/http"
)

const stateCookieName = "oauth_state"

type Auth struct {
	client       *http.Client
	config       *config.Config
	scopes       scopes.Scopes
	sessionStore *session.Store
}

func Init(
	client *http.Client,
	config *config.Config,
	scopes scopes.Scopes,
) *Auth {
	return &Auth{
		client: client,
		config: config,
		scopes: scopes,
	}
}

func (a *Auth) basicAuthEncodedHeader() string {
	basicAuth := fmt.Sprintf("%s:%s", a.config.ClientID, a.config.ClientSecret)
	encoded := base64.StdEncoding.EncodeToString([]byte(basicAuth))

	return fmt.Sprintf("Basic %s", encoded)
}
