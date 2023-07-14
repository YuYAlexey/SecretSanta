package http

import (
	"crypto/md5"
	"encoding/hex"
	"sync"
	"time"
)

const sessionMinSize int = 60 * 10

type Session[V any] struct {
	sync.Mutex

	store     map[string]*sessionValue[V]
	secretKey string
}

type sessionValue[V any] struct {
	value    V
	lifetime time.Time
}

func NewSession[V any](secretKey string) *Session[V] {
	return &Session[V]{
		store:     make(map[string]*sessionValue[V], sessionMinSize),
		secretKey: secretKey,
	}
}

func (s *Session[V]) Set(value V, maxAgg int) string {
	var (
		key string
	)

	for {
		key = s.hash()
		if _, ok := s.Get(key); !ok {
			break
		}
	}

	s.Lock()
	defer s.Unlock()
	s.store[key] = &sessionValue[V]{
		value:    value,
		lifetime: time.Now().Add(time.Second * time.Duration(maxAgg)),
	}

	return key
}

func (s *Session[V]) Get(key string) (V, bool) {
	var value V

	s.Lock()
	sv, ok := s.store[key]
	if !ok {
		s.Unlock()
		return value, false

	}
	s.Unlock()

	if time.Now().After(sv.lifetime) {
		delete(s.store, key)
		return value, false
	}

	return sv.value, true
}

func (s *Session[_]) hash() string {
	algorithm := md5.New()
	algorithm.Write([]byte(s.secretKey))
	return hex.EncodeToString(algorithm.Sum(nil))
}
