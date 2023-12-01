package cache

import (
	"errors"
	"sync"
)

var (
	ErrNotFoundInCache = errors.New("not found in cache")
)

type Cache struct {
	entries map[string]string

	mu sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		entries: make(map[string]string),
		mu:      sync.RWMutex{},
	}
}

func (c *Cache) Set(domain, address string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[domain] = address
}

func (c *Cache) Get(domain string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	add, ok := c.entries[domain]
	if !ok {
		return "", ErrNotFoundInCache
	}

	return add, nil
}
