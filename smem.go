package smem

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"sync"
	"time"

	"github.com/tochti/session-stores"
)

type (
	store struct {
		mutex sync.Mutex
		data  map[string]Session
	}

	Session struct {
		token   string
		userID  string
		expires time.Time
	}
)

func (s Session) Token() string {
	return s.token
}

func (s Session) UserID() string {
	return s.userID
}

func (s Session) Expires() time.Time {
	return s.expires
}

func NewStore() store {
	return store{
		mutex: sync.Mutex{},
		data:  map[string]Session{},
	}
}

func (s *store) NewSession(userID string, expire time.Time) (s2tore.Session, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	token, err := NewSessionToken()
	if err != nil {
		return nil, err
	}

	session := Session{
		token:   token,
		userID:  userID,
		expires: expire,
	}

	s.data[token] = session

	return session, nil
}

func (s *store) ReadSession(token string) (s2tore.Session, bool) {
	s.mutex.Lock()

	v, ok := s.data[token]
	if !ok {
		s.mutex.Unlock()
		return Session{}, false
	}

	if v.expires.Before(time.Now()) {
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
		if val.expires.Before(time.Now()) {
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
