package engine 

import (
    "fmt"
    "sync"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestHashTableBasicLogic(t *testing.T) {
    ht := NewHashTable()

    ht.Set("key1", "value1")
    val, ok := ht.Get("key1")
    assert.True(t, ok, "Метод Get должен возвращать true для существующего ключа 'key1'")
    assert.Equal(t, "value1", val, "Полученное значение должно совпадать с записанным")

    val, ok = ht.Get("key2")
    assert.False(t, ok, "Метод Get должен возвращать false при запросе несуществующего ключа 'key2'")
    assert.Empty(t, val, "При запросе несуществующего ключа значение должно быть пустым")
    
    ht.Del("key1")
    _, ok = ht.Get("key1")
    assert.False(t, ok, "После вызова Del ключ 'key1' должен быть полностью удален из таблицы")
    
    ht.Set("key3", "123")
    ht.Set("key3", "1234")
    val, ok = ht.Get("key3")
    assert.True(t, ok, "Ключ 'key3' должен существовать после перезаписи")
    assert.Equal(t, "1234", val, "Значение по ключу 'key3' должно быть успешно обновлено на '1234'")
}

func TestHashTableConcurrency(t *testing.T) {
    ht := NewHashTable()
    workers := 1000
    var wg sync.WaitGroup

    wg.Add(workers * 3)

    for i := 0; i < workers; i++ {
        i := i 
        
        go func() {
            defer wg.Done()
            key := fmt.Sprintf("key_%d", i%10)
            ht.Set(key, "value")
        }()

        go func() {
            defer wg.Done()
            key := fmt.Sprintf("key_%d", i%10)
            ht.Get(key)
        }()

        go func() {
            defer wg.Done()
            key := fmt.Sprintf("key_%d", i%10)
            ht.Del(key)
        }()
    }

    wg.Wait()
}