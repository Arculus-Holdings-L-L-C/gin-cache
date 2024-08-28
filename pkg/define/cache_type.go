package define

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type OnCacheHit func(c *gin.Context, r Response)

// GenKeyFunc startup on hit hook
type GenKeyFunc func(c *gin.Context) string

// CacheEvict do Evict
type CacheEvict GenKeyFunc

// Cacheable do caching
type Cacheable struct {
	GenKey     GenKeyFunc
	CacheTime  time.Duration
	OnCacheHit []OnCacheHit
}

// Caching mixins Cacheable and CacheEvict
type Caching struct {
	Cacheable Cacheable
	Evict     []CacheEvict

	// CacheErrorCodes error codes allowed to cache
	CacheErrorCodes []int
}

type Response struct {
	Status int
	Header http.Header
	Body   string
}
