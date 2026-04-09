package storage 

import (
	"in-memory/internal/storage/engine"
	"errors"
)

var ErrNotExists = errors.New("The item not exists")

type Storage interface {
	Set(key, value string) error 
	Get(key string) (string, error) 
	Del(key string) error
}

type StorageMemory struct {
	hashTable engine.Engine
}

func NewStorage(hashTable *engine.HashTable) Storage{
	return &StorageMemory{
		hashTable: hashTable,
	}
}

func (i *StorageMemory) Set(key, value string) error{
	i.hashTable.Set(key, value)
	return nil
}

func (i *StorageMemory) Get(key string) (string, error) {
	value, flag := i.hashTable.Get(key)
	if !flag {
		return "", ErrNotExists
	}
	return value, nil
}


func (i *StorageMemory) Del(key string) error {
	i.hashTable.Del(key)
	return nil
}