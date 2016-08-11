package store

// Store declares an interface for a key/value store
type Store interface {
	SetValue(path, value string) error
	Value(path string) string
	Values(path string) map[string]string
	Delete(path string) error
}
