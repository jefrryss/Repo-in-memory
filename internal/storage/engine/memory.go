package engine

import(
	"errors"
)

var (
	ErrNotExistsKey = errors.New("Key not exists")
	ErrAlreadyExistsKey = errors.New("Key already exist")
)

type Memory struct {
	memory map[string]string
}
func NewMemory() *Memory{
	return &Memory{
		memory: make(map[string]string),
	}
}
func (m *Memory) Get(key string) (string, error){
	if val, ok := m.memory[key]; ok {
		return val, nil
	} 
	return "", ErrNotExistsKey
}

func (m *Memory) Set(key, val string) error {
	if _, ok := m.memory[key]; ok {
		return ErrAlreadyExistsKey
	}
	m.memory[key] = val
	return nil
}

func (m *Memory) Del(key string) error {
	if _, ok := m.memory[key]; !ok {
		return ErrNotExistsKey
	}
	delete(m.memory, key)
	return nil
}