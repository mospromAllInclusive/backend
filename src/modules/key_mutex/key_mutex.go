package key_mutex

import (
	"sync"
)

type keyMutex struct {
	m  map[string]*sync.RWMutex
	mu sync.RWMutex
}

func NewKeyMutex() IKeyMutex {
	return &keyMutex{
		m: make(map[string]*sync.RWMutex),
	}
}

func (km *keyMutex) RLock(key string) func() {
	mu, ok := km.tryGetMutex(key)
	if !ok {
		km.registerForKey(key)
	}
	mu, ok = km.tryGetMutex(key)
	if !ok {
		panic("missing mutex for key at keyMutex")
	}

	mu.RLock()
	return mu.RUnlock
}

func (km *keyMutex) Lock(key string) func() {
	mu, ok := km.tryGetMutex(key)
	if !ok {
		km.registerForKey(key)
	}
	mu, ok = km.tryGetMutex(key)
	if !ok {
		panic("missing mutex for key at keyMutex")
	}

	mu.Lock()
	return mu.Unlock
}

func (km *keyMutex) tryGetMutex(key string) (*sync.RWMutex, bool) {
	km.mu.RLock()
	defer km.mu.RUnlock()
	mu, ok := km.m[key]
	return mu, ok
}

func (km *keyMutex) registerForKey(key string) {
	km.mu.Lock()
	defer km.mu.Unlock()
	if _, ok := km.m[key]; ok {
		return
	}
	km.m[key] = &sync.RWMutex{}
}
