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

func Init() *Store {
	return &Store{
		sessions: make(map[string]Session),
	}
}

func (s *Store) Get(id string) (Session, bool) {
	s.mutex.Lock()
	session, ok := s.sessions[id]
	s.mutex.Unlock()

	return session, ok
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
