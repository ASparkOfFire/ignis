package cache

import "sync"

type Maper[K comparable, V any] interface {
	Get(key K) V
	Add(key K, v V)
	Remove(key K)
	Has(key K) bool
}
type SafeMap[K comparable, V any] struct {
	mu sync.RWMutex
	m  map[K]V
}

func NewSafeMap[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{m: make(map[K]V)}
}

func (s *SafeMap[K, V]) Get(key K) V {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.m[key]

	var zeroVal V
	if !ok {
		return zeroVal
	}
	return value
}

func (s *SafeMap[K, V]) Add(key K, value V) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[key] = value
}

func (s *SafeMap[K, V]) Remove(key K) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, key)
}

func (s *SafeMap[K, V]) Has(key K) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.m[key]
	return ok
}
