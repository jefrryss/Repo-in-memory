package engine 

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
	"sync"
)


func TestHashTableBasicLogic(t *testing.T) {
	ht := NewHashTable()

	ht.Set("key1", "value1")
	val, ok := ht.Get("key1")
	
	assert.True(t, ok, "Ключ 'key1' должен сущетсвовать")
	assert.Equal(t, "value1", val, "Значение ключа должно совпадать")

	val, ok = ht.Get("key2")
	assert.False(t, ok, "Ключ 'key2' не должен существовать")
	
	ht.Del("key1")
	_, ok = ht.Get("key2")
	assert.False(t, ok, "Ключ не должен сущетсвовать")
	

	ht.Set("key3", "123")
	ht.Set("key3", "1234")
	val, ok = ht.Get("key3")
	assert.True(t, ok, "Ключ key3 должен существовать")
	assert.Equal(t, "1234", val, "Значение должено быть перезаписано")
}



func TestHashTableConccurency(t *testing.T) {

	wg := sync.WaitGroup{}
	workers := 1000
	ht := NewHashTable()
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func(){
			defer wg.Done()
			key := fmt.Sprintf("key_%d", i % 10)
			ht.Set(key, "value1")
		}()
	}
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func(){
			defer wg.Done()
			key := fmt.Sprintf("key_%d", i % 10)
			ht.Get(key)
		}()
	}
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func(){
			defer wg.Done()
			key := fmt.Sprintf("key_%d", i % 10)
			ht.Del(key)
		}()
	}
	wg.Wait()
}