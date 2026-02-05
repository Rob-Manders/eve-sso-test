package main

import (
	"evessotest/auth"
	"evessotest/config"
	"evessotest/scopes"
	"net/http"

	// Load environment variables from .env file automatically.
	_ "github.com/joho/godotenv/autoload"
)

var authService = auth.Init(
	&http.Client{},
	config.Load(),
	scopes.ScopeList,
)

func main() {
	http.HandleFunc("GET /auth/start", authService.Start)
	http.HandleFunc("GET /auth/callback", authService.Callback)
	http.HandleFunc("GET /api/me", me)
	http.HandleFunc("GET /api/esi", esi)

	http.ListenAndServe(":8080", nil)
}

func me(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Me"))
}

func esi(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ESI"))
}
