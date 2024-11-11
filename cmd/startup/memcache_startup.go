package startup

import (
	"time"

	"github.com/Arculus-Holdings-L-L-C/gin-cache/internal"
	"github.com/Arculus-Holdings-L-L-C/gin-cache/internal/drivers/memcache"
	"github.com/Arculus-Holdings-L-L-C/gin-cache/internal/drivers/rediscache"
	"github.com/Arculus-Holdings-L-L-C/gin-cache/pkg/define"
	"github.com/redis/go-redis/v9"
)

func MemCache(onCacheHit ...define.OnCacheHit) (*internal.CacheHandler, error) {
	return internal.New(memcache.NewMemoryHandler(), onCacheHit...), nil
}

func RedisCache(client *redis.Client, cacheTime time.Duration, onCacheHit ...define.OnCacheHit) (*internal.CacheHandler, error) {
	return internal.New(rediscache.NewRedisHandler(client, cacheTime), onCacheHit...), nil
}
