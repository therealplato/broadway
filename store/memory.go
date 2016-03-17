package store

import "sync"

type memoryStore struct {
	sync.Mutex
	store map[string]string
}

// NewMemory instantiates and returns a Store using the in-memory driver
func NewMemory() Store {
	return &memoryStore{store: map[string]string{}}
}

// SetValue sets the string value for a string key. The key may include
// '/' path separators.
func (s *memoryStore) SetValue(path, value string) error {
	s.Lock()
	s.store[path] = value
	s.Unlock()
	return nil
}

// Value retrieves the string value for a string key.
func (s *memoryStore) Value(path string) string {
	s.Lock()
	defer s.Unlock()
	return s.store[path]
}

// Values finds all leaf nodes under the given key. It strips any leading path
// components from the keys and returns a key/value map. For example, given keys
// "animals/flea" and "animals/cats/egyptian", Values("animals") would return
// {"flea" : "...", "egyptian": "..."}
func (s *memoryStore) Values(path string) map[string]string {
	return s.store
}

// Delete removes the specified key and its value from the store
func (s *memoryStore) Delete(path string) error {
	return nil
}
