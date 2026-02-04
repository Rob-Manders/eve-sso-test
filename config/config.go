package config

import "os"

type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

func LoadConfig() *Config {
	return &Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURI:  os.Getenv("REDIRECT_URI"),
	}
}
