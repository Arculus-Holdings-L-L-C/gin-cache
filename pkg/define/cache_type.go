package define

import (
	"time"

	"github.com/gin-gonic/gin"
)

// CacheHitHook startup on hit hook
type CacheHitHook []func(c *gin.Context, cacheValue string)

// GenKeyFunc startup on hit hook
type GenKeyFunc func(c *gin.Context) string

// CacheEvict do Evict
type CacheEvict GenKeyFunc

// Cacheable do caching
type Cacheable struct {
	GenKey     GenKeyFunc
	CacheTime  time.Duration
	OnCacheHit CacheHitHook // 命中缓存钩子 优先级最高, 可覆盖Caching的OnCacheHitting
}

// Caching mixins Cacheable and CacheEvict
type Caching struct {
	Cacheable Cacheable
	Evict     []CacheEvict

	// CacheErrorCodes error codes allowed to cache
	CacheErrorCodes []int
}
