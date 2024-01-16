package startup

import (
	"github.com/Arculus-Holdings-L-L-C/gin-cache/internal"
	"github.com/Arculus-Holdings-L-L-C/gin-cache/internal/drivers/memcache"
	"github.com/gin-gonic/gin"
)

// MemCache NewMemoryCache init memory support
func MemCache(onCacheHit ...func(c *gin.Context, cacheValue string)) (*internal.CacheHandler, error) {
	return internal.New(memcache.NewMemoryHandler(), onCacheHit...), nil
}
