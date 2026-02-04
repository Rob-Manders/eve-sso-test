package main

import (
	"encoding/base64"
	"evessotest/config"
	"evessotest/scopes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"

	// Load environment variables from .env file automatically.
	_ "github.com/joho/godotenv/autoload"
)

const authURL = "https://login.eveonline.com/v2/oauth/authorize"
const tokenURL = "https://login.eveonline.com/v2/oauth/token"

var cfg = config.Load()

func main() {
	http.HandleFunc("GET /auth/start", authStart)
	http.HandleFunc("GET /auth/callback", authCallback)
	http.HandleFunc("GET /api/me", me)
	http.HandleFunc("GET /api/esi", esi)

	http.ListenAndServe(":8080", nil)
}

func authStart(w http.ResponseWriter, r *http.Request) {
	state := uuid.New().String()
	state = strings.Replace(state, "-", "", -1)

	esiScopes := scopes.Compile()

	query := r.URL.Query()
	query.Add("response_type", "code")
	query.Add("client_id", cfg.ClientID)
	query.Add("redirect_uri", cfg.RedirectURI)
	query.Add("scope", esiScopes)
	query.Add("state", state)

	queryString := query.Encode()
	redirectURL := fmt.Sprintf("%s?%s", authURL, queryString)

	http.SetCookie(w, &http.Cookie{Name: "eve_oauth_state", Value: state, HttpOnly: true, Secure: true})
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func authCallback(w http.ResponseWriter, r *http.Request) {
	state, err := r.Cookie("eve_oauth_state")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	receivedState := r.URL.Query().Get("state")

	if state.Value != receivedState {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	basicAuth := fmt.Sprintf("%s:%s", cfg.ClientID, cfg.ClientSecret)
	basicAuthEncoded := base64.StdEncoding.EncodeToString([]byte(basicAuth))

	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", code)
	form.Add("redirect_uri", cfg.RedirectURI)

	request, err := http.NewRequest("POST", tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	request.Header.Add("Authorization", fmt.Sprintf("Basic %s", basicAuthEncoded))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println(string(data))

	http.Redirect(w, r, "http://localhost:8080/api/esi", http.StatusFound)
}

func me(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Me"))
}

func esi(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ESI"))
}
