package registry

import (
	"errors"
	"sync"
)

// Any Any
type Any = interface{}
type Registry struct {
	mu   *sync.Mutex
	data map[string]Any
}

// Delete Delete
func (registry *Registry) Delete(key string) {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	if _, ok := registry.data[key]; ok {
		delete(registry.data, key)
	}
}

// clear registry. registry size is 0
func (registry *Registry) Clear() {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	registry.data = make(map[string]Any)
}

// set <key,value> , return error
func (registry *Registry) Set(key string, value Any) error {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	if value == nil {
		return errors.New("registry: value can not be nil")
	}
	registry.data[key] = value
	return nil
}

// Get Get
func (registry *Registry) Get(key string) (value Any) {
	if _, ok := registry.data[key]; ok {
		return registry.data[key]
	}
	return nil
}

// GetRegistry GetRegistry
func GetRegistry() *Registry {
	return &Registry{
		mu:   &sync.Mutex{},
		data: map[string]Any{},
	}
}
