package store

import (
	"context"
	"sync"
	"time"
)

type InMemoryStore struct {
	counts  map[string]int
	blocked map[string]time.Time // Armazena o tempo de expiração do bloqueio
	mu      sync.Mutex
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		counts:  make(map[string]int),
		blocked: make(map[string]time.Time),
	}
}

func (s *InMemoryStore) Increment(ctx context.Context, key string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.counts[key]++
	return s.counts[key], nil
}

func (s *InMemoryStore) IsBlocked(ctx context.Context, key string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if blockedTime, exists := s.blocked[key]; exists {
		if time.Now().Before(blockedTime) {
			// Ainda está bloqueado
			return true, nil
		}
		// Remove o bloqueio e reseta o contador após o tempo expirar
		delete(s.blocked, key)
		delete(s.counts, key) // Reseta o contador de requisições
	}
	return false, nil
}

func (s *InMemoryStore) Block(ctx context.Context, key string, blockTime time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.blocked[key] = time.Now().Add(blockTime)
	return nil
}
