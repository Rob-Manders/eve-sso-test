package main

import (
	// Load environment variables from .env file automatically.

	"evessotest/internal/sso"
	"net/http"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

var httpClient = &http.Client{
	Timeout: time.Second * 10,
}

var ssoCredentials = sso.Credentials{
	ClientID:     os.Getenv("CLIENT_ID"),
	ClientSecret: os.Getenv("CLIENT_SECRET"),
	RedirectURI:  os.Getenv("REDIRECT_URI"),
}

var ssoScopes = []string{
	"esi-corporations.read_corporation_membership.v1",
	"esi-corporations.read_structures.v1",
	"esi-corporations.track_members.v1",
	"esi-corporations.read_divisions.v1",
	"esi-corporations.read_contacts.v1",
	"esi-corporations.read_titles.v1",
	"esi-corporations.read_blueprints.v1",
	"esi-corporations.read_standings.v1",
	"esi-corporations.read_starbases.v1",
	"esi-corporations.read_container_logs.v1",
	"esi-corporations.read_facilities.v1",
	"esi-corporations.read_medals.v1",
	"esi-alliances.read_contacts.v1",
	"esi-corporations.read_fw_stats.v1",
	"esi-corporations.read_projects.v1",
	"esi-corporations.read_freelance_jobs.v1",
}

var auth = sso.New(httpClient, ssoCredentials, ssoScopes)

func main() {
	http.HandleFunc("GET /auth/start", auth.Start)
	http.HandleFunc("GET /auth/callback", auth.Callback)

	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
