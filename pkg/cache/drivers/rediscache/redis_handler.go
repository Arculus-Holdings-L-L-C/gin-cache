package rediscache

import (
	"context"
	"encoding/json"
	"math"
	"time"

	"github.com/Arculus-Holdings-L-L-C/gin-cache/pkg/define"
	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	cacheStore *redis.Client
	cacheTime  time.Duration
}

// NewRedisHandler do new Redis startup object
func NewRedisHandler(client *redis.Client, cacheTime time.Duration) *redisCache {
	return &redisCache{cacheStore: client, cacheTime: cacheTime}
}

func (r *redisCache) Load(ctx context.Context, key string) define.Response {
	val := r.cacheStore.Get(ctx, key).Val()
	if val == "" {
		return define.Response{}
	}

	var response define.Response
	if err := json.Unmarshal([]byte(val), &response); err != nil {
		return define.Response{}
	}
	return response
}

func (r *redisCache) Set(ctx context.Context, key string, data define.Response, timeout time.Duration) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	if timeout > 0 {
		r.cacheStore.Set(ctx, key, string(jsonData), timeout)
	} else {
		r.cacheStore.Set(ctx, key, string(jsonData), r.cacheTime)
	}
}

func (r *redisCache) DoEvict(ctx context.Context, keys []string) {
	var evictKeys []string
	for _, key := range keys {
		var cursor uint64
		deleteKeys, _, err := r.cacheStore.Scan(ctx, cursor, key, math.MaxUint16).Result()

		if err == nil {
			evictKeys = append(evictKeys, deleteKeys...)
		}
	}

	if len(evictKeys) > 0 {
		r.cacheStore.Del(ctx, evictKeys...)
	}
}
