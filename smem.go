package smem

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"sync"
	"time"
)

type (
	Session struct {
		Token   string
		UserID  string
		Expires time.Time
	}
	store struct {
		mutex sync.Mutex
		data  map[string]Session
	}

	SessionStore interface {
		NewSession(string, time.Time) (string, error)
		ReadSession(string) (Session, bool)
		RemoveSession(string) error
		RemoveExpiredSessions() (int, error)
	}
)

func NewStore() store {
	return store{
		mutex: sync.Mutex{},
		data:  map[string]Session{},
	}
}

func (s *store) NewSession(userID string, expire time.Time) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	token, err := NewSessionToken()
	if err != nil {
		return "", err
	}

	s.data[token] = Session{
		Token:   token,
		UserID:  userID,
		Expires: expire,
	}

	return token, nil
}

func (s *store) ReadSession(token string) (Session, bool) {
	s.mutex.Lock()

	v, ok := s.data[token]
	if !ok {
		s.mutex.Unlock()
		return Session{}, false
	}

	if v.Expires.Before(time.Now()) {
		s.mutex.Unlock()
		s.RemoveSession(token)
		return Session{}, false
	}

	s.mutex.Unlock()

	return v, true
}

func (s *store) RemoveSession(token string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.data, token)

	return nil
}

func (s *store) RemoveExpiredSessions() (int, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	c := 0
	for key, val := range s.data {
		if val.Expires.Before(time.Now()) {
			delete(s.data, key)
			c++
		}
	}

	return c, nil
}

func NewSessionToken() (string, error) {
	buf := make([]byte, 2)

	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	c := sha256.New()
	hash := fmt.Sprintf("%x", c.Sum(buf))

	return hash, nil
}
