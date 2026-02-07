package auth

import (
	"encoding/json"
	"evessotest/backend/token"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
)

const tokenURL = "https://login.eveonline.com/v2/oauth/token"

func (a *Auth) Callback(w http.ResponseWriter, r *http.Request) {
	if !validateState(r) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	request, err := a.buildTokenRequest(code)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := a.client.Do(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var accessToken token.Token
	err = json.Unmarshal(data, &accessToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sessionCookie := a.createSession(accessToken)
	http.SetCookie(w, sessionCookie)

	// TODO: Store refresh token somewhere.

	http.Redirect(w, r, "http://localhost:8080/api/esi", http.StatusFound)
}

func (a *Auth) buildTokenRequest(code string) (*http.Request, error) {
	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", code)
	form.Add("redirect_uri", a.config.RedirectURI)

	request, err := http.NewRequest("POST", tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", a.basicAuthEncodedHeader())
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return request, nil
}

func (a *Auth) createSession(accessToken token.Token) *http.Cookie {
	userID := uuid.New()
	sessionID, sessionExpiry := a.sessionStore.Add(userID, accessToken.AccessToken, accessToken.ExpiresIn)

	sessionCookie := &http.Cookie{
		Name:     "session",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  sessionExpiry,
	}

	return sessionCookie
}

func validateState(r *http.Request) bool {
	state, err := r.Cookie(stateCookieName)
	if err != nil {
		return false
	}

	if state.Value == "" {
		return false
	}

	receivedState := r.URL.Query().Get("state")
	if state.Value != receivedState {
		return false
	}

	return true
}
