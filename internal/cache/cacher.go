package cache

import (
	"github.com/tetratelabs/wazero"
)

// Cacher interface for cache operations
type Cacher[K comparable] interface {
	Get(K) wazero.CompilationCache
	Add(K, wazero.CompilationCache)
	Has(K) bool
	Remove(K)
}

// ModCache wraps SafeMap for thread-safe caching
type ModCache[K comparable] struct {
	m *SafeMap[K, wazero.CompilationCache] // Use pointer for proper initialization
}

// NewModCache initializes a new ModCache
func NewModCache[K comparable]() ModCache[K] {
	return ModCache[K]{
		m: NewSafeMap[K, wazero.CompilationCache](), // Properly initialize SafeMap
	}
}

// Get retrieves a CompilationCache from SafeMap
func (mc *ModCache[K]) Get(key K) wazero.CompilationCache {
	return mc.m.Get(key)
}

// Add Set stores a CompilationCache in SafeMap
func (mc *ModCache[K]) Add(key K, cache wazero.CompilationCache) {
	mc.m.Add(key, cache)
}

func (mc *ModCache[K]) Has(key K) bool {
	return mc.m.Has(key)
}

// Remove removes an entry from SafeMap
func (mc *ModCache[K]) Remove(key K) {
	mc.m.Remove(key)
}
