package handler

import (
	"github.com/redis/go-redis/v9"
)

type CacheHandler struct {
	redis redis.UniversalClient
}

func NewCacheService(redis redis.UniversalClient) *CacheHandler {
	return &CacheHandler{redis: redis}
}
