package cache

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	RWlock sync.RWMutex
	Data   map[string][]byte
}

func New() *Cache {
	return &Cache{
		RWlock: sync.RWMutex{},
		Data:   make(map[string][]byte),
	}
}

func (c *Cache) Get(key []byte) ([]byte, error) {
	keyStr := string(key)
	val, ok := c.Data[keyStr]
	if !ok {
		return val, fmt.Errorf("key not found %v", keyStr)
	}
	return val, nil
}

func (c *Cache) Set(key, val []byte, ttl time.Duration) error {
	keyStr := string(key)
	if _, ok := c.Data[keyStr]; ok {
		c.Delete(key)
	}
	c.Data[keyStr] = val
	return nil
}
func (c *Cache) Has(key []byte) bool {
	if _, ok := c.Data[string(key)]; !ok {
		return false
	}
	return true
}
func (c *Cache) Delete(key []byte) error {
	keyStr := string(key)
	if _, ok := c.Data[string(key)]; !ok {
		return fmt.Errorf("key does not exist %v", keyStr)
	}
	delete(c.Data, keyStr)
	return nil
}
