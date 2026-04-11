package storage 

import (
	"in-memory/internal/storage/engine"
	"errors"
	"context"
)

var ErrNotExists = errors.New("The item not exists")

type Storage interface {
	Set(ctx context.Context, key, value string) error 
	Get(ctx context.Context, key string) (string, error) 
	Del(ctx context.Context, key string) error
}

type StorageMemory struct {
	hashTable engine.Engine
}

func NewStorage(hashTable *engine.HashTable) Storage{
	return &StorageMemory{
		hashTable: hashTable,
	}
}

func (i *StorageMemory) Set(ctx context.Context, key, value string) error{
	if err := ctx.Err(); err != nil {
		return err
	}
	i.hashTable.Set(key, value)
	return nil
}

func (i *StorageMemory) Get(ctx context.Context, key string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	value, flag := i.hashTable.Get(key)
	if !flag {
		return "", ErrNotExists
	}
	return value, nil
}


func (i *StorageMemory) Del(ctx context.Context, key string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	i.hashTable.Del(key)
	return nil
}