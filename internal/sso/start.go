package sso

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

const authURL = "https://login.eveonline.com/v2/oauth/authorize"

func (s *SSO) Start(w http.ResponseWriter, r *http.Request) {
	state := uuid.New().String()
	state = strings.Replace(state, "-", "", -1)

	esiScopes := s.compileScopes()

	query := r.URL.Query()
	query.Add("response_type", "code")
	query.Add("client_id", s.credentials.ClientID)
	query.Add("redirect_uri", s.credentials.RedirectURI)
	query.Add("scope", esiScopes)
	query.Add("state", state)

	queryString := query.Encode()
	redirectURL := fmt.Sprintf("%s?%s", authURL, queryString)

	stateCookie := &http.Cookie{
		Name:     stateCookieName,
		Value:    state,
		MaxAge:   120,
		HttpOnly: true,
		Secure:   true,
	}

	http.SetCookie(w, stateCookie)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// compileScopes assembles an array of scopes into a space separated list.
func (s *SSO) compileScopes() string {
	compiled := ""

	for _, scope := range s.scopes {
		compiled += fmt.Sprintf(" %s", scope)
	}

	compiled = strings.Trim(compiled, " ")
	return compiled
}
