package auth

import "os"

type Credentials struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

func LoadAuthCredentials() *Credentials {
	return &Credentials{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURI:  os.Getenv("REDIRECT_URI"),
	}
}
