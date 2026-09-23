package utils

import "sync"

// KeyedMutex hands out a lock per key (e.g. per account ID) so callers can
// guard a critical section without a single global lock serializing everything.
type KeyedMutex struct {
	mu    sync.Mutex
	locks map[string]*sync.Mutex
}

func NewKeyedMutex() *KeyedMutex {
	return &KeyedMutex{locks: make(map[string]*sync.Mutex)}
}

func (k *KeyedMutex) Lock(key string) (unlock func()) {
	k.mu.Lock()
	l, ok := k.locks[key]
	if !ok {
		l = &sync.Mutex{}
		k.locks[key] = l
	}
	k.mu.Unlock()

	l.Lock()
	return l.Unlock
}
