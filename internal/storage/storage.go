package storage 

import "in-memory/internal/storage/engine"

type Storage interface {
	Set(key, value string) error 
	Get(key string) (string, error) 
	Del(key string) error
}

type InMemory struct {
	memory *engine.Memory
}

func NewStorage(m *enigne.Memory) Storage{
	return &InMemory{
		memory: m,
	}
}

func (i *InMemory) Set(key, value string) error{
	err := i.memory.Set(key, value)
	return err
}

func (i *InMemory) Get(key strign) (string, error) {
	value, err := i.memory.Get(key)
	if err != nil {
		return "", err
	}
	return value, nil
}


func (i *InMemory) Del(key string) error {
	err := i.memory.Del(key)
	return err
}