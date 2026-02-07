package esi

import (
	"evessotest/backend/session"
	"net/http"
)

type ESI struct {
	sessionStore session.Store
}

func (e ESI) Handler(w http.ResponseWriter, r *http.Request) {
	// Stubbed ESI endpoint for now.
	w.Write([]byte("ESI"))
}
