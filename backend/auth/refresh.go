package auth

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
)

func (a *Auth) Refresh(w http.ResponseWriter, userID uuid.UUID) (string, error) {
	refreshToken, ok := a.userDB.Get(userID)
	if !ok {
		return "", errors.New("user not found")
	}

	request, err := a.buildRefreshRequest(refreshToken)
	if err != nil {
		return "", err
	}

	response, err := a.client.Do(request)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", err
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var token Token
	err = json.Unmarshal(data, &token)
	if err != nil {
		return "", err
	}

	sessionCookie := a.createSession(userID, token)
	http.SetCookie(w, sessionCookie)

	return token.AccessToken, nil
}

func (a *Auth) buildRefreshRequest(refreshToken string) (*http.Request, error) {
	form := url.Values{}
	form.Add("grant_type", "refresh_token")
	form.Add("refresh_token", refreshToken)
	form.Add("redirect_uri", a.credentials.RedirectURI)

	request, err := http.NewRequest("POST", tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", a.basicAuthEncodedHeader())
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return request, nil
}
