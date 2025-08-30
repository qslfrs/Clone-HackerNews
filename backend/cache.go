package main

import (
	"time"

	"github.com/patrickmn/go-cache"
)

// hasil item disimpan/diambil dari cache bila tersedia
// jd gk selalu ngambil dr Hacker News

var gc *cache.Cache

func InitCache(defaultTTL time.Duration) {
	gc = cache.New(defaultTTL, 2*defaultTTL)
}

func CacheGet(key string) (interface{}, bool) {
	return gc.Get(key)
}

func CacheSet(key string, value interface{}, ttl time.Duration) {
	gc.Set(key, value, ttl)
}
