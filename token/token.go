package token

import (
	"encoding/json"
	"fmt"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	Expiry       int64  `json:"expires_in"`
}

func Parse(tokenJSON json.RawMessage) (Token, error) {
	var token Token
	err := json.Unmarshal(tokenJSON, &token)

	if err != nil {
		return Token{}, fmt.Errorf("error parsing token JSON: %s", err)
	}

	return token, nil
}
