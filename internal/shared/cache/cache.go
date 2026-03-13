package cache

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) (interface{}, error)

	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error

	Delete(ctx context.Context, key string) error

	Exists(ctx context.Context, key string) (bool, error)
}

type InMemoryCache struct {
	data map[string]*cacheItem
}

type cacheItem struct {
	value      interface{}
	expiration time.Time
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		data: make(map[string]*cacheItem),
	}
}

func (c *InMemoryCache) Get(ctx context.Context, key string) (interface{}, error) {
	item, ok := c.data[key]
	if !ok {
		return nil, nil
	}

	if !item.expiration.IsZero() && time.Now().After(item.expiration) {
		delete(c.data, key)
		return nil, nil
	}

	return item.value, nil
}

func (c *InMemoryCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	var exp time.Time
	if expiration > 0 {
		exp = time.Now().Add(expiration)
	}

	c.data[key] = &cacheItem{
		value:      value,
		expiration: exp,
	}
	return nil
}

func (c *InMemoryCache) Delete(ctx context.Context, key string) error {
	delete(c.data, key)
	return nil
}

func (c *InMemoryCache) Exists(ctx context.Context, key string) (bool, error) {
	_, ok := c.data[key]
	return ok, nil
}
