package tinybox

import "sync"

type Cache struct {
	mu   *sync.Mutex
	lru  *LRU[Payload]
	size int
}

func NewCache(size int) *Cache {
	return &Cache{
		mu:   &sync.Mutex{},
		lru:  NewLRU[Payload](size, nil),
		size: size,
	}
}

func (c *Cache) Set(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.lru.Add(key, Payload{buf: value}, len(value))
}

func (c *Cache) Get(key string) (value []byte, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if payload, ok := c.lru.Get(key); ok {
		return payload.buf, true
	}
	return nil, false
}
