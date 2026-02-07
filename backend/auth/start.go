package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

const authURL = "https://login.eveonline.com/v2/oauth/authorize"

func (a *Auth) Start(w http.ResponseWriter, r *http.Request) {
	state := uuid.New().String()
	state = strings.Replace(state, "-", "", -1)

	esiScopes := a.scopes.Compile()

	query := r.URL.Query()
	query.Add("response_type", "code")
	query.Add("client_id", a.credentials.ClientID)
	query.Add("redirect_uri", a.credentials.RedirectURI)
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
