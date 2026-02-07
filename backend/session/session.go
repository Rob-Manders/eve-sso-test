package session

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// Max session = 30 days
const maxSessionLifetime = time.Hour * 24 * 30

type Store struct {
	sessions map[string]Session
	mutex    sync.Mutex
}

type Session struct {
	UserID        uuid.UUID
	AccessToken   string
	TokenExpiry   int64
	SessionExpiry time.Time
}

func (s *Store) GetAccessToken(id string) (string, bool) {
	s.mutex.Lock()
	session, ok := s.sessions[id]
	s.mutex.Unlock()

	if !ok {
		return "", false
	}

	currentTime := time.Now().Unix()
	if session.TokenExpiry < currentTime-10 {
		s.mutex.Lock()
		delete(s.sessions, id)
		s.mutex.Unlock()

		return "", false
	}

	return session.AccessToken, ok
}

func (s *Store) Add(userID uuid.UUID, accessToken string, expiresIn int64) (string, time.Time) {
	id := uuid.New().String()
	tokenExpiry := time.Now().Add(time.Second * time.Duration(expiresIn)).Unix()
	sessionExpiry := time.Now().Add(time.Second * maxSessionLifetime)

	s.mutex.Lock()
	s.sessions[id] = Session{
		UserID:        userID,
		AccessToken:   accessToken,
		TokenExpiry:   tokenExpiry,
		SessionExpiry: sessionExpiry,
	}
	s.mutex.Unlock()

	return id, sessionExpiry
}

func (s *Store) Delete(id string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.sessions, id)
}

func (s *Store) ClearExpired() {
	for id, session := range s.sessions {
		currentTime := time.Now().Unix()
		if session.TokenExpiry < currentTime-10 {
			s.mutex.Lock()
			delete(s.sessions, id)
			s.mutex.Unlock()
		}
	}
}
