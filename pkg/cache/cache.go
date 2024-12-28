package cache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

type Cache struct {
	*cache.Cache
}

func New() *Cache {
	return &Cache{
		Cache: cache.New(5*time.Minute, 30*time.Minute),
	}
}
