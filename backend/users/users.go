package users

import (
	"sync"

	"github.com/google/uuid"
)

// Mock user database for demonstration purposes so we have somewhere to store the refresh tokens.

type DB struct {
	users map[uuid.UUID]string
	mutex sync.Mutex
}

func Init() *DB {
	return &DB{
		users: make(map[uuid.UUID]string),
	}
}

func (u *DB) Create(refreshToken string) uuid.UUID {
	userId := uuid.New()

	u.mutex.Lock()
	u.users[userId] = refreshToken
	u.mutex.Unlock()

	return userId
}

func (u *DB) Get(userID uuid.UUID) (string, bool) {
	u.mutex.Lock()
	token, ok := u.users[userID]
	u.mutex.Unlock()

	return token, ok
}
