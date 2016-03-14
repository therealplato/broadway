package store

import "sync"

type memoryStore struct {
	sync.Mutex
	store map[string]string
}

func NewMemory() Store {
	return &memoryStore{store: map[string]string{}}
}

func (s *memoryStore) SetValue(path, value string) error {
	s.Lock()
	s.store[path] = value
	s.Unlock()
	return nil
}
func (s *memoryStore) Value(path string) string {
	s.Lock()
	defer s.Unlock()
	return s.store[path]
}
func (s *memoryStore) Values(path string) map[string]string {
	return s.store
}
func (s *memoryStore) Delete(path string) error {
	return nil
}
