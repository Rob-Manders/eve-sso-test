package sso

import (
	"encoding/base64"
	"fmt"
	"net/http"
)

const stateCookieName = "oauth_state"

type Credentials struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

type SSO struct {
	client      *http.Client
	credentials Credentials
	scopes      []string
	jwkCache    *JWKCache
}

func New(client *http.Client, credentials Credentials, scopes []string) *SSO {
	return &SSO{
		client:      client,
		credentials: credentials,
		scopes:      scopes,
		jwkCache:    &JWKCache{},
	}
}

func (s *SSO) basicAuthEncodedHeader() string {
	basicAuth := fmt.Sprintf("%s:%s", s.credentials.ClientID, s.credentials.ClientSecret)
	encoded := base64.StdEncoding.EncodeToString([]byte(basicAuth))

	return fmt.Sprintf("Basic %s", encoded)
}
