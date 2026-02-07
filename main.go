package main

import (
	"evessotest/backend/auth"
	"evessotest/backend/session"
	"evessotest/backend/users"
	"net/http"

	// Load environment variables from .env file automatically.
	_ "github.com/joho/godotenv/autoload"
)

var userDB = users.Init()
var sessionStore = session.Init()
var authService = auth.Init(sessionStore, userDB)

func main() {
	http.HandleFunc("GET /auth/start", authService.Start)
	http.HandleFunc("GET /auth/callback", authService.Callback)

	http.ListenAndServe(":8080", nil)
}
