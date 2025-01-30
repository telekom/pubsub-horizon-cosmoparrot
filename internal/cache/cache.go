package cache

import (
	go_cache "github.com/patrickmn/go-cache"
	"time"
)

var Current *go_cache.Cache

func init() {
	Current = go_cache.New(24*time.Hour, 10*time.Minute)
}
