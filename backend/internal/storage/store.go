package storage

import "sync"

type Store interface {
    Save(slug string, target string)
    Get(slug string) (string, bool)
}

type InMemoryStore struct {
    mu sync.RWMutex
    m  map[string]string
}

func NewInMemoryStore() *InMemoryStore {
    return &InMemoryStore{m: make(map[string]string)}
}

func (s *InMemoryStore) Save(slug string, target string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.m[slug] = target
}

func (s *InMemoryStore) Get(slug string) (string, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    t, ok := s.m[slug]
    return t, ok
}
