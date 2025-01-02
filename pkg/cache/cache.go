package cache

import (
	"errors"
	"sync"

	redis "github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

var cache *Cache

var hdlOnce sync.Once

func NewOrGetSingletonCache(url string) (*Cache, error) {
	if url == "" {
		return nil, errors.New("url is empty")
	}

	var er error
	hdlOnce.Do(func() {
		c, err := newCache(url)
		if err != nil {
			er = err
		}

		cache = c
	})
	return cache, er
}

func newCache(url string) (*Cache, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	return &Cache{
		client: redis.NewClient(opts),
	}, nil
}
