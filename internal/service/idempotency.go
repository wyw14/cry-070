package service

import "sync"

type Keys struct {
	mu     sync.Mutex
	values map[string]string
}

func NewKeys() *Keys { return &Keys{values: map[string]string{}} }
func (k *Keys) Get(key string) (string, bool) {
	k.mu.Lock()
	defer k.mu.Unlock()
	v, ok := k.values[key]
	return v, ok
}
func (k *Keys) Put(key, value string) bool {
	k.mu.Lock()
	defer k.mu.Unlock()
	if _, ok := k.values[key]; ok {
		return false
	}
	k.values[key] = value
	return true
}
