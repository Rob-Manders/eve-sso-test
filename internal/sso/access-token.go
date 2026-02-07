package sso

import "time"

type AccessToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func (a AccessToken) ExpiresOn() time.Time {
	return time.Now().Add(time.Duration(a.ExpiresIn) * time.Second)
}
