package startup

import (
	"github.com/Arculus-Holdings-L-L-C/gin-cache/internal"
	"github.com/Arculus-Holdings-L-L-C/gin-cache/internal/drivers/memcache"
	"github.com/Arculus-Holdings-L-L-C/gin-cache/pkg/define"
)

// MemCache NewMemoryCache init memory support
func MemCache(onCacheHit ...define.OnCacheHit) (*internal.CacheHandler, error) {
	return internal.New(memcache.NewMemoryHandler(), onCacheHit...), nil
}
