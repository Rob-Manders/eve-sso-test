package sso

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const tokenURL = "https://login.eveonline.com/v2/oauth/token"

func (s *SSO) Callback(w http.ResponseWriter, r *http.Request) {
	if !validateState(w, r) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	request, err := s.buildTokenRequest(code)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := s.client.Do(request)
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

	var token AccessToken
	err = json.Unmarshal(data, &token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tokenClaims, err := s.validateToken(token.AccessToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	subject, err := tokenClaims.GetSubject()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	subjectSegments := strings.Split(subject, ":")
	if len(subjectSegments) != 2 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	characterID := subjectSegments[2]

	fmt.Print(characterID)
}

func (s *SSO) buildTokenRequest(code string) (*http.Request, error) {
	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", code)
	form.Add("redirect_uri", s.credentials.RedirectURI)

	request, err := http.NewRequest("POST", tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", s.basicAuthEncodedHeader())
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return request, nil
}

func validateState(w http.ResponseWriter, r *http.Request) bool {
	defer removeStateCookie(w)

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

func removeStateCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:    stateCookieName,
		Expires: time.Now().Add(-1 * time.Minute),
	})
}
