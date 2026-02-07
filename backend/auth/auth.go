package auth

import (
	"encoding/base64"
	"errors"
	"evessotest/backend/session"
	"evessotest/backend/users"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const stateCookieName = "oauth_state"

type Auth struct {
	client       *http.Client
	credentials  *Credentials
	scopes       Scopes
	sessionStore *session.Store
	userDB       *users.DB
}

func Init(sessionStore *session.Store, userDB *users.DB) *Auth {
	return &Auth{
		client:       &http.Client{},
		credentials:  LoadAuthCredentials(),
		scopes:       ScopeList,
		sessionStore: sessionStore,
		userDB:       userDB,
	}
}

func (a *Auth) GetAccessToken(w http.ResponseWriter, sessionID string) (string, error) {
	userSession, ok := a.sessionStore.Get(sessionID)
	if !ok {
		return "", errors.New("session not found")
	}

	var expired bool
	if userSession.TokenExpiry < time.Now().Unix()-10 {
		expired = true
	}

	if !expired {
		return userSession.AccessToken, nil
	}

	token, err := a.Refresh(w, userSession.UserID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *Auth) createSession(userID uuid.UUID, accessToken Token) *http.Cookie {
	sessionID, sessionExpiry := a.sessionStore.Add(userID, accessToken.AccessToken, accessToken.ExpiresIn)

	sessionCookie := &http.Cookie{
		Name:     "session",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  sessionExpiry,
	}

	return sessionCookie
}

func (a *Auth) basicAuthEncodedHeader() string {
	basicAuth := fmt.Sprintf("%s:%s", a.credentials.ClientID, a.credentials.ClientSecret)
	encoded := base64.StdEncoding.EncodeToString([]byte(basicAuth))

	return fmt.Sprintf("Basic %s", encoded)
}
