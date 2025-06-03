package storage

import (
	"context"
	"sync"

	"github.com/gotd/td/session"
)

type MemorySessionStorage struct {
	data []byte
	mu   sync.RWMutex
}

func NewMemorySessionStorage() *MemorySessionStorage {
	return &MemorySessionStorage{}
}

func (s *MemorySessionStorage) LoadSession(_ context.Context) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.data == nil {
		return nil, session.ErrNotFound
	}
	return s.data, nil
}

func (s *MemorySessionStorage) StoreSession(_ context.Context, data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = make([]byte, len(data))
	copy(s.data, data)
	return nil
}
