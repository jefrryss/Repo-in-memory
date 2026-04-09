package engine

import(
	"sync"
)

type Engine interface {
	Set(key string, value string)
	Get(key string) (string, bool)
	Del(key string)
}

type HashTable struct {
	mu sync.RWMutex
	data map[string]string
}

func NewHashTable() *HashTable{
	return &HashTable{
		data: make(map[string]string),
	}
}

func (m *HashTable) Get(key string) (string, bool){
	m.mu.RLock()
	defer m.mu.RUnlock()
	if val, ok := m.data[key]; ok {
		return val, true
	}
	return "", false
}

func (m *HashTable) Set(key, val string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = val
}

func (m *HashTable) Del(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.data, key)
}