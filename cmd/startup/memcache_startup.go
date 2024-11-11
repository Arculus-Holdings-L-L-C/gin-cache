package startup

import (
	"time"

	"github.com/Arculus-Holdings-L-L-C/gin-cache/pkg/cache"
	"github.com/Arculus-Holdings-L-L-C/gin-cache/pkg/cache/drivers/memcache"
	"github.com/Arculus-Holdings-L-L-C/gin-cache/pkg/cache/drivers/rediscache"
	"github.com/Arculus-Holdings-L-L-C/gin-cache/pkg/define"
	"github.com/redis/go-redis/v9"
)

func MemCache(onCacheHit ...define.OnCacheHit) (*cache.CacheHandler, error) {
	return cache.New(memcache.NewMemoryHandler(), onCacheHit...), nil
}

func RedisCache(client *redis.Client, cacheTime time.Duration, onCacheHit ...define.OnCacheHit) (*cache.CacheHandler, error) {
	return cache.New(rediscache.NewRedisHandler(client, cacheTime), onCacheHit...), nil
}
