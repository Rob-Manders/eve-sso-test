package auth

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func (a *Auth) Refresh() error {
	request, err := a.buildRefreshRequest("refresh_token")
	if err != nil {
		return err
	}

	response, err := a.client.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return err
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	// TODO: Save new token to store.
	fmt.Println(string(data))

	return nil
}

func (a *Auth) buildRefreshRequest(refreshToken string) (*http.Request, error) {
	form := url.Values{}
	form.Add("grant_type", "refresh_token")
	form.Add("refresh_token", refreshToken)
	form.Add("redirect_uri", a.config.RedirectURI)

	request, err := http.NewRequest("POST", tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", a.basicAuthEncodedHeader())
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return request, nil
}
