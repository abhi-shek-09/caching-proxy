package main

import "sync"

type Cache struct {
    Data  map[string][]byte
    mutex sync.RWMutex
}

var cache = Cache{
    Data: make(map[string][]byte),
}

func (c *Cache) Get(url string) ([]byte, bool) {
    c.mutex.RLock()
    defer c.mutex.RUnlock()

    data, exists := c.Data[url]
    return data, exists
}

func (c *Cache) Set(url string, data []byte) {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    c.Data[url] = data
}

func (c *Cache) Clear() {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    c.Data = make(map[string][]byte)
}